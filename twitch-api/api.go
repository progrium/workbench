package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/machinebox/graphql"
)

type TwitchAPI struct {
	Client    *graphql.Client
	ChannelID string

	clientID string
	token    string
}

func (api *TwitchAPI) Authenticate(authTool string) error {
	cmd := exec.Command(authTool)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	var auth map[string]string
	if err := json.Unmarshal(out, &auth); err != nil {
		return err
	}
	api.token = auth["token"]
	api.clientID = auth["client-id"]
	return nil
}

func (api *TwitchAPI) makeRequest(query string) *graphql.Request {
	req := graphql.NewRequest(query)
	req.Header.Set("Client-Id", api.clientID)
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", api.token))
	return req
}

// func (api *TwitchAPI) FutureEvents() ([]Event, error) {
// 	return FetchEvents(api.ChannelID)
// }
func (api *TwitchAPI) EventAt(t time.Time) (*Event, error) {
	events, err := api.FutureEvents()
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if event.StartTime.Equal(t.UTC()) {
			return &event, nil
		}
	}
	return nil, nil
}
func (api *TwitchAPI) EventByID(id string) (*Event, error) {
	events, err := api.FutureEvents()
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if event.ID == id {
			return &event, nil
		}
	}
	return nil, nil
}

func (api *TwitchAPI) CreateEvent(e *Event) error {
	req := api.makeRequest(`
		mutation($input: CreateSegmentEventInput!) {
			createSegmentEvent(input: $input) {
				event {
					id
				}
			}
		}
	`)
	req.Var("input", gqlEvent{
		ChannelID:   api.ChannelID,
		Title:       e.Title,
		Description: e.Description,
		StartAt:     e.StartTime.UTC().Format(time.RFC3339),
		EndAt:       e.EndTime.UTC().Format(time.RFC3339),
		GameID:      e.GameID,
		OwnerID:     api.ChannelID,
		ParentID:    e.SeriesID,
	})
	var resp map[string]map[string]map[string]interface{}
	err := api.Client.Run(context.Background(), req, &resp)
	if err != nil {
		return err
	}
	e.ID = resp["createSegmentEvent"]["event"]["id"].(string)
	return nil
}

func (api *TwitchAPI) DeleteEvent(id string) error {
	req := api.makeRequest(`
		mutation($input: DeleteEventLeafInput!) {
			deleteEventLeaf(input: $input) {
				event {
					id
				}
			}
	  	}
	`)
	req.Var("input", map[string]string{"eventID": id})
	var resp map[string]interface{}
	err := api.Client.Run(context.Background(), req, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (api *TwitchAPI) FutureEvents() ([]Event, error) {
	req := api.makeRequest(`query Events($criteria: ManagedEventLeavesCriteriaInput) {
		currentUser {
			managedEventLeaves(first: 100, criteria: $criteria) {
				edges {
					node {
						channel {
							id
						},
						createdAt,
						defaultTimeZone,
						description,
						endAt,
						game {
							id
						},
						imageURL,
						id,
						owner {
							id
						},
						parent{
							id
						},
						startAt,
						title,
						type
					}
				}
			}
		}
	}`)
	req.Var("criteria", criteria{
		StartsAfter: time.Now().UTC().Format(time.RFC3339),
	})
	var resp eventsResponse
	err := api.Client.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}
	var events []Event
	for _, edge := range resp.CurrentUser.ManagedEventLeaves.Edges {
		events = append(events, Event{
			Title:       edge.Node.Title,
			ID:          edge.Node.ID,
			Description: edge.Node.Description,
			StartTime:   edge.Node.StartAt,
			EndTime:     edge.Node.EndAt,
			SeriesID:    edge.Node.Parent.ID,
			GameID:      edge.Node.Game.ID,
		})
	}
	return events, nil
}

type criteria struct {
	StartsAfter string `json:"startsAfter"`
}

type eventsResponse struct {
	CurrentUser struct {
		ManagedEventLeaves struct {
			Edges []struct {
				Node struct {
					Channel struct {
						ID string `json:"id"`
					} `json:"channel"`
					CreatedAt       time.Time `json:"createdAt"`
					DefaultTimeZone string    `json:"defaultTimeZone"`
					Description     string    `json:"description"`
					EndAt           time.Time `json:"endAt"`
					Game            struct {
						ID string `json:"id"`
					} `json:"game"`
					ImageURL string `json:"imageURL"`
					ID       string `json:"id"`
					Owner    struct {
						ID string `json:"id"`
					} `json:"owner"`
					Parent struct {
						ID string `json:"id"`
					} `json:"parent"`
					StartAt time.Time `json:"startAt"`
					Title   string    `json:"title"`
					Type    string    `json:"type"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"managedEventLeaves"`
	} `json:"currentUser"`
}

func (api *TwitchAPI) UpdateEvent(id string, e *Event) error {
	req := api.makeRequest(`
		mutation($input: UpdateSegmentEventInput!) {
			updateSegmentEvent(input: $input) {
				event {
					id
				}
			}
	  	}
	`)
	req.Var("input", gqlEvent{
		ID:          id,
		ChannelID:   api.ChannelID,
		Title:       e.Title,
		Description: e.Description,
		GameID:      e.GameID,
		ParentID:    e.SeriesID,
		StartAt:     e.StartTime.UTC().Format(time.RFC3339),
		EndAt:       e.EndTime.UTC().Format(time.RFC3339),
	})
	var resp map[string]interface{}
	err := api.Client.Run(context.Background(), req, &resp)
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

type EventsResponse struct {
	Events []Event `json:"events"`
}

type Event struct {
	Title       string    `json:"title"`
	ID          string    `json:"_id"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	SeriesID    string
	GameID      string
}

func (e *Event) Tag() string {
	parts := strings.Split(e.Description, "#")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

type gqlEvent struct {
	ID          string `json:"id,omitempty"`
	ChannelID   string `json:"channelID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EndAt       string `json:"endAt"`
	StartAt     string `json:"startAt"`
	GameID      string `json:"gameID"`
	OwnerID     string `json:"ownerID"`
	ParentID    string `json:"parentID"`
}

func FetchEvents(channelID string) ([]Event, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/channels/%s/events", channelID), nil)
	if err != nil {
		return nil, err
	}
	var payload EventsResponse
	_, err = Do(req, &payload)
	if err != nil {
		return nil, err
	}
	return payload.Events, err
}

func Do(req *http.Request, r interface{}) (*http.Response, error) {
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Set("User-Agent", "progrium")
	req.Header.Set("Client-ID", "tf99t1lw9dcsxpprca8h4uwew53yos")
	req.Header.Set("Authorization", "OAuth "+os.Getenv("TWITCH_OAUTH_TOKEN"))
	if req.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkResponse(resp); err != nil {
		return resp, err
	}

	if r != nil {
		err = json.NewDecoder(resp.Body).Decode(r)
		if err == io.EOF {
			err = nil
		}
	}
	return resp, err
}

type ErrorResponse struct {
	// HTTP response that cause this error.
	Response *http.Response

	// Error message.
	Message string `json:"message,omitempty"`
}

func checkResponse(r *http.Response) error {
	if 200 <= r.StatusCode && r.StatusCode <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

func (e *ErrorResponse) Error() string {
	r := e.Response

	return fmt.Sprintf("%v %v: %d %v",
		r.Request.Method, r.Request.URL, r.StatusCode, e.Message)
}

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/machinebox/graphql"
)

// type TwitchAPI interface {
// 	FutureEvents() ([]TwitchEvent, error)
// 	EventAt(t time.Time) (*TwitchEvent, error)
// 	EventByID(id string) (*TwitchEvent, error)
// 	Post(e TwitchEvent) (*TwitchEvent, error)
// 	Delete(id string) error
// 	Put(id string, e TwitchEvent) (*TwitchEvent, error)
// }

// type TwitchEvent struct {
// 	ID          string
// 	Title       string
// 	Description string
// 	StartTime   time.Time
// 	EndTime     time.Time
// }

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	cmd := exec.Command("./auth/gql-auth")
	out, err := cmd.CombinedOutput()
	fatal(err)
	var auth map[string]string
	fatal(json.Unmarshal(out, &auth))

	api := &TwitchAPI{
		client:    graphql.NewClient("https://gql.twitch.tv/gql"),
		clientID:  auth["client-id"],
		token:     auth["token"],
		channelID: "5031651",
	}

	events, _ := api.FutureEvents()
	fmt.Println(events)

	// e := Event{
	// 	Title:       "New test event",
	// 	Description: "Here is a description",
	// 	StartTime:   time.Now().Add(48 * time.Hour),
	// 	EndTime:     time.Now().Add(50 * time.Hour),
	// }
	// _, err = api.Post(e)
	// fatal(err)

	fatal(api.Delete("23nl853ITbSXw2VniEND_w"))
	fatal(api.Delete("kUmt6nwqQyOsIPLS64zmwA"))
	fatal(api.Delete("485Ij4deROeuXwEdOMwb6w"))

	// events, _ = api.FutureEvents()
	// fmt.Println(events)
}

type TwitchAPI struct {
	client    *graphql.Client
	clientID  string
	token     string
	channelID string
}

func (api *TwitchAPI) makeRequest(query string) *graphql.Request {
	req := graphql.NewRequest(query)
	req.Header.Set("Client-Id", api.clientID)
	req.Header.Set("Authorization", fmt.Sprintf("OAuth %s", api.token))
	return req
}

func (api *TwitchAPI) FutureEvents() ([]Event, error) {
	return FetchEvents(api.channelID)
}
func (api *TwitchAPI) EventAt(t time.Time) (*Event, error) {
	events, err := api.FutureEvents()
	if err != nil {
		return nil, err
	}
	for _, event := range events {
		if event.StartTime.Equal(t) {
			return &event, nil
		}
	}
	return nil, fmt.Errorf("event not found")
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
	return nil, fmt.Errorf("event not found")
}

func (api *TwitchAPI) Post(e Event) (*Event, error) {
	req := api.makeRequest(`
		mutation($input: CreateSingleEventInput!) {
			createSingleEvent(input: $input) {
				event {
					id
				}
			}
		}
	`)
	req.Var("input", GQLEvent{
		ChannelID:   api.channelID,
		Title:       e.Title,
		Description: e.Description,
		StartAt:     e.StartTime.UTC().Format(time.RFC3339),
		EndAt:       e.EndTime.UTC().Format(time.RFC3339),
		GameID:      "488191",
		OwnerID:     "5031651",
	})
	var resp map[string]interface{}
	err := api.client.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (api *TwitchAPI) Delete(id string) error {
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
	err := api.client.Run(context.Background(), req, &resp)
	if err != nil {
		return err
	}
	return nil
}

func (api *TwitchAPI) Put(id string, e Event) (*Event, error) {
	req := api.makeRequest(`
		mutation($input: UpdateSingleEventInput!) {
			updateSingleEvent(input: $input) {
				event {
					id
				}
			}
	  	}
	`)
	req.Var("input", GQLEvent{
		ID:          id,
		ChannelID:   api.channelID,
		Title:       e.Title,
		Description: e.Description,
		GameID:      "488191",
	})
	var resp map[string]interface{}
	err := api.client.Run(context.Background(), req, &resp)
	if err != nil {
		return nil, err
	}
	return &e, nil
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
}

type GQLEvent struct {
	ID          string `json:"id,omitempty"`
	ChannelID   string `json:"channelID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EndAt       string `json:"endAt"`
	StartAt     string `json:"startAt"`
	GameID      string `json:"gameID"`
	OwnerID     string `json:"ownerID"`
}

func FetchEvents(channelID string) ([]Event, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/kraken/channels/%s/events", channelID), nil)
	if err != nil {
		return nil, err
	}
	var resp EventsResponse
	_, err = Do(req, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Events, err
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

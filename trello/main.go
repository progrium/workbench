package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/now"
	"github.com/machinebox/graphql"
	"github.com/mitchellh/hashstructure"
	twitch "github.com/progrium/workbench/twitch-api"

	"github.com/adlio/trello"
)

const CreativeGameID = "488191"

var (
	zeroTime = time.Date(0, 0, 0, 0, 0, 0, 0, tz)
	tz       = time.Now().Location()
	series   = map[string]string{
		"tigl3d":    "CruZjy5IS7e6JKgQxW2SEQ",
		"workbench": "Fms5J2-TRVWr7Cs-wOcdGQ",
	}
	schedule = map[string][]Slot{
		"workbench": []Slot{
			{
				time.Tuesday,
				time.Date(0, 0, 0, 17, 0, 0, 0, tz),
				time.Hour * 2,
			},
			{
				time.Thursday,
				time.Date(0, 0, 0, 17, 0, 0, 0, tz),
				time.Hour * 2,
			},
		},
		"tigl3d": []Slot{
			{
				time.Monday,
				time.Date(0, 0, 0, 15, 0, 0, 0, tz),
				time.Hour * 4,
			},
			{
				time.Tuesday,
				time.Date(0, 0, 0, 14, 0, 0, 0, tz),
				time.Hour * 2,
			},
			{
				time.Wednesday,
				time.Date(0, 0, 0, 15, 0, 0, 0, tz),
				time.Hour * 4,
			},
			{
				time.Thursday,
				time.Date(0, 0, 0, 14, 0, 0, 0, tz),
				time.Hour * 2,
			},
		},
	}
)

func main() {
	resync := false

	state, err := LoadState()
	fatal(err)

	client := trello.NewClient(os.Getenv("TRELLO_KEY"), os.Getenv("TRELLO_TOKEN"))
	board, err := client.GetBoard("EsjNNP3c", trello.Defaults())
	fatal(err)
	lists, err := board.GetLists(trello.Defaults())
	fatal(err)

	twitchapi := &twitch.TwitchAPI{
		Client:    graphql.NewClient("https://gql.twitch.tv/gql"),
		ChannelID: "5031651",
	}
	fatal(twitchapi.Authenticate("../twitch-api/auth/gql-auth"))

	futureEvents, err := twitchapi.FutureEvents()
	fatal(err)

	var pastEventCardIDs []string
	for cardID, twitchID := range state.Mapping {
		found := false
		for _, event := range futureEvents {
			if event.ID == twitchID {
				found = true
			}
		}
		if !found {
			pastEventCardIDs = append(pastEventCardIDs, cardID)
			resync = true
		}
	}

	for id, hash := range state.TwitchHashes {
		e, _ := twitchapi.EventByID(id)
		if hash != HashTwitchEvent(e) {
			resync = true
			break
		}
	}

	var events []Event
	for _, l := range lists {
		if l.Name == "Schedule" {
			cards, err := l.GetCards(trello.Defaults())
			fatal(err)
			for _, c := range cards {
				events = append(events, NewEvent(c))
			}
			break
		}
	}

	trelloHash := HashEvents(events)
	var deletedCardIDs []string
	if state.TrelloHash != trelloHash {
		resync = true
		for k, _ := range state.Mapping {
			found := false
			for _, e := range events {
				if e.ID == k {
					found = true
					break
				}
			}
			if !found {
				deletedCardIDs = append(deletedCardIDs, k)
			}
		}
	}

	// tevents, err := twitchapi.NewEvents()
	// fatal(err)
	// fmt.Println(tevents)
	// for _, e := range tevents {
	// 	twitchapi.DeleteEvent(e.ID)
	// }
	// return

	now.WeekStartDay = time.Monday
	thisWeek := now.BeginningOfWeek()
	nextWeek := now.New(now.EndOfWeek().Add(time.Hour)).BeginningOfWeek()

	if resync {
		fmt.Println("Resyncing...")
		state.TrelloHash = trelloHash

		if len(deletedCardIDs) > 0 {
			fmt.Println(" - Deleting unmapped events")
		}
		for _, id := range deletedCardIDs {
			for _, event := range futureEvents {
				if event.ID == state.Mapping[id] {
					fmt.Printf("   x %s\n", event.Title)
					fatal(twitchapi.DeleteEvent(state.Mapping[id]))
				}
			}
			delete(state.Mapping, id)
		}

		if len(pastEventCardIDs) > 0 {
			fmt.Println(" - Deleting past event cards")
			for _, id := range pastEventCardIDs {
				card, err := client.GetCard(id, trello.Defaults())
				fatal(err)
				fmt.Printf("   x %s\n", card.Name)
				fatal(card.Update(trello.Arguments{
					"closed": "true",
				}))
			}
		}
		var event *Event
		for _, weekBegin := range []time.Time{thisWeek, nextWeek} {
			for tag, slots := range schedule {
				for _, slot := range slots {
					if slot.StartTime(weekBegin).Before(time.Now()) {
						continue
					}
					event, events = ShiftEvent(events, tag)
					if event == nil {
						continue
					}
					findEvent, err := twitchapi.EventAt(slot.StartTime(weekBegin))
					fatal(err)
					replaceEvent := &twitch.Event{
						Title:       event.Name,
						Description: event.Description + "\n#" + tag,
						StartTime:   slot.StartTime(weekBegin),
						EndTime:     slot.EndTime(weekBegin),
						SeriesID:    series[tag],
						GameID:      CreativeGameID,
					}
					if findEvent == nil {
						fmt.Println("NEW:", event.Name)
						fatal(twitchapi.CreateEvent(replaceEvent))
					} else {
						fmt.Println("MOD:", event.Name)
						fatal(twitchapi.UpdateEvent(findEvent.ID, replaceEvent))
					}
					state.TwitchHashes[replaceEvent.ID] = HashTwitchEvent(replaceEvent)
					state.Mapping[event.ID] = replaceEvent.ID
				}
			}
		}

	}

	EnsureSlotEvents(twitchapi, thisWeek)
	EnsureSlotEvents(twitchapi, nextWeek)

	fatal(SaveState(state))
}

func EnsureSlotEvents(api TwitchAPI, weekBegin time.Time) {
	for tag, slots := range schedule {
		for _, slot := range slots {
			if slot.StartTime(weekBegin).Before(time.Now()) {
				continue
			}
			e, err := api.EventAt(slot.StartTime(weekBegin))
			fatal(err)
			if e == nil {
				api.CreateEvent(&twitch.Event{
					Title:       "TBD",
					Description: "#" + tag,
					StartTime:   slot.StartTime(weekBegin),
					EndTime:     slot.EndTime(weekBegin),
					SeriesID:    series[tag],
					GameID:      CreativeGameID,
				})
			}
		}
	}
}

func ShiftEvent(events []Event, tag string) (*Event, []Event) {
	var idx *int
	for i, e := range events {
		if e.Tag == tag {
			idx = &i
			break
		}
	}
	if idx == nil {
		return nil, events
	}
	e := events[*idx]
	return &e, append(events[:*idx], events[(*idx)+1:]...)
}

func HashEvents(events []Event) string {
	hash, err := hashstructure.Hash(events, nil)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", hash)
}

func HashTwitchEvent(event *twitch.Event) string {
	hash, err := hashstructure.Hash(event, nil)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", hash)
}

func LabelNames(labels []*trello.Label) []string {
	var names []string
	for _, l := range labels {
		names = append(names, l.Name)
	}
	return names
}

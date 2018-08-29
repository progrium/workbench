package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"time"

	"github.com/progrium/workbench/twitch-api"
)

type TwitchAPI interface {
	FutureEvents() ([]twitch.Event, error)
	EventAt(t time.Time) (*twitch.Event, error)
	EventByID(id string) (*twitch.Event, error)
	CreateEvent(e *twitch.Event) error
	DeleteEvent(id string) error
	UpdateEvent(id string, e *twitch.Event) error
}

type TwitchAPIMock struct {
	Events []twitch.Event
}

func (ft *TwitchAPIMock) FutureEvents() ([]twitch.Event, error) {
	var events []twitch.Event
	for _, event := range ft.Events {
		if event.StartTime.After(time.Now()) {
			events = append(events, event)
		}
	}
	sort.Slice(events, func(i, j int) bool { return events[i].StartTime.Before(events[j].StartTime) })
	return events, nil
}

func (ft *TwitchAPIMock) EventAt(t time.Time) (*twitch.Event, error) {
	for _, event := range ft.Events {
		if event.StartTime.Equal(t) {
			return &event, nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (ft *TwitchAPIMock) EventByID(id string) (*twitch.Event, error) {
	for _, event := range ft.Events {
		if event.ID == id {
			return &event, nil
		}
	}
	return nil, fmt.Errorf("not found: %s", id)
}

func (ft *TwitchAPIMock) CreateEvent(e *twitch.Event) error {
	e.ID = RandomString(10)
	ft.Events = append(ft.Events, *e)
	return nil
}

func (ft *TwitchAPIMock) DeleteEvent(id string) error {
	var idx *int
	for i, e := range ft.Events {
		if e.ID == id {
			idx = &i
			break
		}
	}
	if idx == nil {
		return fmt.Errorf("not found: %s", id)
	}
	ft.Events = append(ft.Events[:*idx], ft.Events[(*idx)+1:]...)
	return nil
}

func (ft *TwitchAPIMock) UpdateEvent(id string, e *twitch.Event) error {
	var idx *int
	for i, ee := range ft.Events {
		if ee.ID == id {
			idx = &i
			break
		}
	}
	if idx == nil {
		return fmt.Errorf("not found: %s", id)
	}
	e.ID = id
	ft.Events[*idx] = *e
	return nil
}

func LoadTwitchMock() (*TwitchAPIMock, error) {
	b, err := ioutil.ReadFile("twitch.json")
	if err != nil {
		return &TwitchAPIMock{}, nil
	}
	if len(b) == 0 {
		return &TwitchAPIMock{}, nil
	}
	var api TwitchAPIMock
	err = json.Unmarshal(b, &api)
	return &api, err
}

func SaveTwitchMock(t *TwitchAPIMock) error {
	buf, err := json.MarshalIndent(*t, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("twitch.json", buf, 0644)
}

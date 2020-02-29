package main

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/nlopes/slack/slackevents"
)

type fakeSeparator struct {
	message string
}

func (f *fakeSeparator) parseEvent(rawEvent json.RawMessage, opts slackevents.Option) (slackevents.EventsAPIEvent, error) {
	e := &slackevents.EventsAPIEvent{}
	e.Token = "foobar"
	e.Type = slackevents.CallbackEvent
	e.TeamID = "foobar"
	innerEvent := e.InnerEvent
	innerEvent.Type = ""

	data := new(slackevents.MessageEvent)
	data.Text = f.message
	data.User = "someone"
	data.Channel = "some channel"
	innerEvent.Data = data

	e.InnerEvent = innerEvent
	return *e, nil
}

func TestServe(t *testing.T) {
	s := "please input a message you want to test"

	fakebmo := new(BMO)
	fakebmo.bridge = &fakeSeparator{s}

	req := httptest.NewRequest("POST", "/events-endpoint", nil)
	rec := httptest.NewRecorder()
	fakebmo.ServeHTTP(rec, req)
}

func TestParse(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput map[string]*votes
	}{
		{
			input:          "bmo++ increment 😊",
			expectedOutput: map[string]*votes{"bmo": &votes{1, 1}},
		},
		{
			input:          "bmo++ ykhr53++ great job!",
			expectedOutput: map[string]*votes{"bmo": &votes{1, 1}, "ykhr53": &votes{1, 1}},
		},
		{
			input:          "bmo-- decrement 😢",
			expectedOutput: map[string]*votes{"bmo": &votes{-1, 1}},
		},
		{
			input:          "bmo++ bmo-- neutral 😐",
			expectedOutput: map[string]*votes{"bmo": &votes{0, 2}},
		},
	}

	for i, tc := range tests {
		output := parseVote(tc.input)
		if !reflect.DeepEqual(output, tc.expectedOutput) {
			t.Errorf("Test %d: test fails", i)
		}
	}
}

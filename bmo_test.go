package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

type fakeBMO struct {
	child   *BMO
	message string
	token   string
	uname   string
	api     *slack.Client
	client  *dynamodb.DynamoDB
}

func (f *fakeBMO) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.child.ServeHTTP(w, r)
}

func (f *fakeBMO) ParseEvent(rawEvent json.RawMessage, opts slackevents.Option) (slackevents.EventsAPIEvent, error) {
	e := &slackevents.EventsAPIEvent{}
	e.Token = "foobar"
	e.Type = slackevents.CallbackEvent
	e.TeamID = "foobar"
	innerEvent := e.InnerEvent
	innerEvent.Type = ""

	var data slackevents.MessageEvent
	data.Text = f.message
	data.User = "someone"
	data.Channel = "some channel"
	innerEvent.Data = data

	e.InnerEvent = innerEvent
	return *e, nil
}

func TestServe(t *testing.T) {
	fbmo := new(fakeBMO)
	fbmo.child = new(BMO)
	fbmo.message = "!!!!!!!!!! yokohei test !!!!!!!!!!"
	req := httptest.NewRequest("POST", "/events-endpoint", nil)
	rec := httptest.NewRecorder()
	fbmo.ServeHTTP(rec, req)
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
		output := parse(tc.input)
		if !reflect.DeepEqual(output, tc.expectedOutput) {
			t.Errorf("Test %d: test fails", i)
		}
	}
}

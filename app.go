package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

var token = getenv("SLACKTOKEN")
var vtoken = getenv("VTOKEN")
var botname = getenv("BOTUNAME")
var api = slack.New(token)

func main() {
	http.HandleFunc("/events-endpoint", func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		body := buf.String()
		eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: vtoken}))
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}
		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			case *slackevents.MessageEvent:
				if ev.User != botname && votable(ev.Text) {
					for _, name := range parse(ev.Text) {
						text := strings.TrimRight(name, "++ ") + " voted!"
						api.PostMessage(ev.Channel, slack.MsgOptionText(text, false))
					}
				}
			}
		}
	})
	fmt.Println("[INFO] Server listening")
	http.ListenAndServe(":3000", nil)
}

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func votable(text string) bool {
	r := regexp.MustCompile(`\S+\+\+\s`)
	return r.MatchString(text)
}

func parse(text string) []string {
	r := regexp.MustCompile(`\S+\+\+\s`)
	names := r.FindAllString(text, -1)
	if names == nil {
		return nil
	}
	return names
}

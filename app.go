package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
)

// BMO handles your task using following information.
type BMO struct {
	token string
	uname string
	api   *slack.Client
}

// NewBMO is BMO constructor.
func NewBMO() *BMO {
	bmo := new(BMO)
	oauthToken := getenv("SLACKTOKEN")

	bmo.token = getenv("VTOKEN")
	bmo.uname = getenv("BOTUNAME")
	bmo.api = slack.New(oauthToken)

	return bmo
}

func (b *BMO) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Event catch")
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: b.token}))

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
			b.api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
		case *slackevents.MessageEvent:
			if ev.User != b.uname && parse(ev.Text) != nil {
				for _, name := range parse(ev.Text) {
					text := strings.TrimRight(name, "++ ") + " : <ここに数字が入る> voted!"
					b.api.PostMessage(ev.Channel, slack.MsgOptionText(text, false))
				}
			}
		}
	}
}

func main() {
	bmo := NewBMO()
	mux := http.NewServeMux()
	mux.Handle("/events-endpoint", bmo)
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func parse(text string) []string {
	r := regexp.MustCompile(`\S+\+\+\s`)
	names := r.FindAllString(text, -1)
	return names
}

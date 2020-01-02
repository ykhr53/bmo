package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"

	"github.com/ykhr53/bmo/ddbfunc"
)

type funcSet interface {
	parseEvent(json.RawMessage, slackevents.Option) (slackevents.EventsAPIEvent, error)
}

type separator struct{}

type votes struct {
	sum   int
	count int
}

// BMO handles your task using following information.
type BMO struct {
	token  string
	uname  string
	api    *slack.Client
	client *dynamodb.DynamoDB
	bridge funcSet
}

// NewBMO is BMO constructor.
func NewBMO() *BMO {
	bmo := new(BMO)

	ddbClient := dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-west-1"))
	oauthToken := getenv("SLACKTOKEN")

	bmo.token = getenv("VTOKEN")
	bmo.uname = getenv("BOTUNAME")
	bmo.api = slack.New(oauthToken)
	bmo.client = ddbClient
	bmo.bridge = &separator{}

	return bmo
}

func (s *separator) parseEvent(rawEvent json.RawMessage, opts slackevents.Option) (slackevents.EventsAPIEvent, error) {
	return slackevents.ParseEvent(rawEvent, opts)
}

// ServeHTTP is hoge
func (b *BMO) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	eventsAPIEvent, e := b.bridge.parseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: b.token}))

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
			if ev.User != b.uname && votable(ev.Text) {
				m := parse(ev.Text)
				for name, votes := range m {
					if votes.sum == 0 {
						continue
					}
					// curr must be string typed number, "unvoted" or "err".
					curr, _ := ddbfunc.GetVal(b.client, name)
					if curr != "err" {
						var text string
						var iv int
						var sv string
						if curr != "unvoted" {
							iv, _ = strconv.Atoi(curr)
						}
						iv = iv + votes.sum
						sv = strconv.Itoa(iv)
						if votes.count > 1 {
							text = name + ": " + sv + " voted! (you got " + strconv.Itoa(votes.count) + " votes)"
						} else {
							text = name + ": " + sv + " voted!"
						}
						b.api.PostMessage(ev.Channel, slack.MsgOptionText(text, false))
						ddbfunc.SetVal(b.client, name, sv)
					}
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

func votable(text string) bool {
	r := regexp.MustCompile(`\S+(\+\+|--)\s`)
	return r.MatchString(text)
}

func parse(text string) map[string]*votes {
	var isPositive bool
	m := make(map[string]*votes)

	r := regexp.MustCompile(`\S+(\+\+|--)\s`)
	names := r.FindAllString(text, -1)

	for _, name := range names {
		isPositive = strings.HasSuffix(name, "+ ")
		name = strings.TrimRight(name, "+- ")
		if v, ok := m[name]; ok {
			v.count++
			if isPositive {
				v.sum++
			} else {
				v.sum--
			}
		} else {
			m[name] = new(votes)
			m[name].count = 1
			if isPositive {
				m[name].sum = 1
			} else {
				m[name].sum = -1
			}
		}
	}
	return m
}

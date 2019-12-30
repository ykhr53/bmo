package bmo

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

// BMO handles your task using following information.
type BMO struct {
	token  string
	uname  string
	api    *slack.Client
	client *dynamodb.DynamoDB
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

	return bmo
}

func (b *BMO) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
					name = strings.TrimRight(name, "+ ")
					vote, _ := ddbfunc.GetVal(b.client, name)
					var voteStr string
					if vote < 0 {
						voteStr = "1"
					} else {
						voteStr = strconv.Itoa(vote + 1)
					}
					text := name + ": " + voteStr + " voted!"
					b.api.PostMessage(ev.Channel, slack.MsgOptionText(text, false))
					ddbfunc.SetVal(b.client, name, voteStr)
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

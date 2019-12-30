package bmo

import (
	"bytes"
	"encoding/json"
	"fmt"
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
)

// BMO handles your task using following information.
type BMO struct {
	token string
	uname string
	api   *slack.Client
	ddb   *dynamodb.DynamoDB
}

// NewBMO is BMO constructor.
func NewBMO() *BMO {
	bmo := new(BMO)

	ddb := dynamodb.New(session.New(), aws.NewConfig().WithRegion("eu-west-1"))
	oauthToken := getenv("SLACKTOKEN")

	bmo.token = getenv("VTOKEN")
	bmo.uname = getenv("BOTUNAME")
	bmo.api = slack.New(oauthToken)
	bmo.ddb = ddb

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
					vote, _ := getVal(b.ddb, name)
					var voteStr string
					if vote < 0 {
						voteStr = "1"
					} else {
						voteStr = strconv.Itoa(vote + 1)
					}
					text := name + ": " + voteStr + " voted!"
					b.api.PostMessage(ev.Channel, slack.MsgOptionText(text, false))
					setVal(b.ddb, name, voteStr)
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

func getVal(ddb *dynamodb.DynamoDB, key string) (int, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String("bmo"),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(key),
			},
		},
		AttributesToGet: []*string{
			aws.String("votes"),
		},
		ConsistentRead:         aws.Bool(true),
		ReturnConsumedCapacity: aws.String("NONE"),
	}

	resp, err := ddb.GetItem(params)
	if err != nil {
		fmt.Println(err.Error())
		return -1, nil
	}
	if len(resp.Item) == 0 {
		return -1, nil
	}
	return strconv.Atoi(*resp.Item["votes"].N)
}

func setVal(ddb *dynamodb.DynamoDB, key string, n string) {
	param := &dynamodb.UpdateItemInput{
		TableName: aws.String("bmo"),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(key),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#votes": aws.String("votes"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":vote_val": {
				N: aws.String(n),
			},
		},
		UpdateExpression:            aws.String("set #votes = :vote_val"),
		ReturnConsumedCapacity:      aws.String("NONE"),
		ReturnItemCollectionMetrics: aws.String("NONE"),
		ReturnValues:                aws.String("NONE"),
	}

	_, err := ddb.UpdateItem(param)
	if err != nil {
		fmt.Println(err.Error())
	}
}

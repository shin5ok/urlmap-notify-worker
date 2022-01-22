package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	// "github.com/gin-gonic/gin"
)

func init() {
	log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

var projectId = os.Getenv("PROJECT")
var subscription = os.Getenv("SUBSCRIPTION")

/*
  Not use because slackUrl will be pulled from message
	var slackUrl = os.Getenv("SLACK_URL")
*/

// to use slack auth parameter
type SlackStruct struct {
	SlackUrl     string
	SlackChannel string
}

type dataJson struct {
	Message  string `json:"message"`
	SlackUrl string `json:"slack_url"`
	Email    string `json:"email"`
	NotifyTo string `json:"notify_to"`
}

type notifyInterface interface {
	Send(message string) error
}

func NotifyDo(n notifyInterface, message string) {
	n.Send(message)
}

func (s *SlackStruct) Send(message string) error {
	channel := s.SlackChannel
	url := s.SlackUrl
	body, _ := json.Marshal(map[string]string{"text": message, "channel": channel})
	res, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Error().Msgf("status code error: %d %s", res.StatusCode, res.Status)
		return errors.New(res.Status)
	} else {
		content, _ := ioutil.ReadAll(res.Body)
		log.Info().Msg(string(content))
		return nil
	}
}

func main() {
	client, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		log.Fatal().Err(err)
	}
	sub := client.Subscription(subscription)
	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {

		d := dataJson{}
		json.Unmarshal([]byte(m.Data), &d)
		x := strings.Split("#", d.NotifyTo)

		var s notifyInterface
		switch x[0] {
		case "slack":
			s = &SlackStruct{SlackUrl: d.SlackUrl, SlackChannel: x[1]}
		default:
			log.Error().Msgf("%+v", x)

		}
		NotifyDo(s, string(d.Message))
		log.Info().Msgf("%+v", string(m.Data))
		m.Ack()
	})
	if err != nil {
		log.Error().Err(err)
	}
}

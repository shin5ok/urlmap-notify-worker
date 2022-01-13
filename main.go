package main

import (
	"context"
	"os"
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

var ProjectId = os.Getenv("PROJECT")
var Subscription = os.Getenv("SUBSCRIPTION")
var SlackURL = os.Getenv("SLACK_URL")

func main() {
	client, err := pubsub.NewClient(context.Background(), ProjectId)
	if err != nil {
		log.Fatal().Err(err)
	}
	sub := client.Subscription(Subscription)
	err = sub.Receive(context.Background(), func(ctx context.Context, m *pubsub.Message) {
		// http.PostForm(SlackURL, )
		log.Info().Msgf("%+v", m)
		m.Ack()
	})
	if err != nil {
		log.Error().Err(err)
	}
}

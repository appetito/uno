package handlers

import (
	"fmt"

	"github.com/appetito/uno"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	analyticsapi "github.com/appetito/uno/examples/greetanalytics/api"
	"github.com/appetito/uno/examples/greeter/api"
	"github.com/rs/zerolog/log"
)

var client *analyticsapi.GreetAnalyticsClient
var js jetstream.JetStream

func InitClient(nc *nats.Conn){
	client = analyticsapi.NewGreetAnalyticsClient(nc, &uno.UnoClientConfig{})
}

func InitJetStream(nc *nats.Conn){
	s, err := jetstream.New(nc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to JetStream")
	}
	js = s
}

//Greet a user, with some additional information
func GreetHandler(r uno.Request, request api.GreetRequest){

	stats, err := client.GetUsersStats(r.Context(), analyticsapi.GetUsersStatsRequest{Name: request.Name})
	if err != nil {
		r.Error("INTERNAL_ERROR", "Fail to get stats", []byte(err.Error()))
		return
	}
	response := api.Greeting{
		Message: fmt.Sprintf("Hello, %s! This is greeting number %d", request.Name, stats.GreetCount),
	}
	
	ack, err := js.Publish(r.Context(), "greets", []byte(request.Name))
	if err != nil {
		log.Error().Err(err).Msg("Failed to publish greet")
		r.Error("INTERNAL_ERROR", "Failed to publish greet", []byte(err.Error()))
	}else{
		log.Info().Uint64("Ack", ack.Sequence).Msg("Greet published")
		r.RespondJSON(response)
	}
}
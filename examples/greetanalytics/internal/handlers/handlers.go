package handlers

import (
	"context"
	"sort"
	"sync"

	"github.com/appetito/uno"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"github.com/appetito/uno/examples/greetanalytics/api"
	"github.com/rs/zerolog/log"
)

var stats map[string]int64
var mx sync.Mutex

// var c jetstream.Consumer

func StartConsumer(nc *nats.Conn) jetstream.ConsumeContext {
	stats = make(map[string]int64)

	s, err := jetstream.New(nc)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to JetStream")
	}
	// js = s
	c, err := s.Consumer(context.TODO(), "GREETS", "ga")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create consumer")
	}

	cc, err := c.Consume(func(m jetstream.Msg){
		mx.Lock()
		defer mx.Unlock()
		log.Info().Msgf("Received message: %s", string(m.Data()))
		stats[string(m.Data())]++
		m.Ack()
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to consume")
	}
	log.Info().Msgf("Consumer created: %s", c.CachedInfo().Name)
	return cc
}


//Get user's greet stats
func GetUsersStatsHandler(r uno.Request, request api.GetUsersStatsRequest){
	count := stats[request.Name]
	
	response := api.UserStats{
		Name: request.Name,
		GreetCount: count,
	}
	r.RespondJSON(response)
}

//Get top greeted users
func TopGreetedUsersHandler(r uno.Request, request api.TopGreetedUsersRequest){
	response := []api.UserStats{}
	mx.Lock()
	defer mx.Unlock()
	for name, count := range stats {
		response = append(response, api.UserStats{
			Name: name,
			GreetCount: count,
		})
	}
	//sort by greet count
	sort.Slice(response, func(i, j int) bool {
		return response[i].GreetCount > response[j].GreetCount
	})
	if int(request.Count) < len(response) {
		response = response[:request.Count]
	}
	r.RespondJSON(response)
}


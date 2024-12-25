package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	analyticsapi "github.com/appetito/uno/examples/greetanalytics/api"
	//greet api
	greetapi "github.com/appetito/uno/examples/greeter/api"
	//nats
	"github.com/nats-io/nats.go"
	//uno
	"github.com/appetito/uno"
)

func main(){

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}

	ac := analyticsapi.NewGreetAnalyticsClient(nc, &uno.UnoClientConfig{})
	gc := greetapi.NewGreeterClient(nc, &uno.UnoClientConfig{})

	gc.Greet(context.Background(), greetapi.GreetRequest{Name: "Alice"})
	gc.Greet(context.Background(), greetapi.GreetRequest{Name: "Bob"})
	gc.Greet(context.Background(), greetapi.GreetRequest{Name: "Alice"})
	gc.Greet(context.Background(), greetapi.GreetRequest{Name: "Alice"})
	resp, err := ac.TopGreetedUsers(context.Background(), analyticsapi.TopGreetedUsersRequest{Count: 10})
	if err != nil {
		log.Error().Err(err).Msg("Failed to get top greeted users")
	}
	fmt.Printf("Top greeted users: %v\n", resp)
}	
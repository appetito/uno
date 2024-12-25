package service

import (
	"github.com/appetito/uno"
	"github.com/nats-io/nats.go"

	"github.com/rs/zerolog/log"

	"github.com/appetito/uno/examples/greetanalytics/api"
	"github.com/appetito/uno/examples/greetanalytics/internal/config"
	"github.com/appetito/uno/examples/greetanalytics/internal/handlers"
)

func New(cfg *config.Config) uno.Service {

	log.Info().Str("URL", cfg.NatsServers).Msg("Connecting to NATS")
	nc, err := nats.Connect(cfg.NatsServers)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}

	svc, err := uno.AddService(nc, uno.Config{
		Name:       "example" + "_" +  "GreetAnalytics",
		Version:     "0.0.1",
		Description: "GreetAnalytics",
		Interceptors: []uno.InterceptorFunc{
			uno.NewPanicInterceptor,
			uno.NewMetricsInterceptor,
			uno.NewTracingInterceptor, 
			uno.NewLoggingInterceptor,   
		},
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to add service")
	}

	root := svc.AddGroup(api.SERVICE_NAME)


	root.AddEndpoint(api.GET_USERS_STATS, uno.AsStructHandler[api.GetUsersStatsRequest](handlers.GetUsersStatsHandler))

	root.AddEndpoint(api.TOP_GREETED_USERS, uno.AsStructHandler[api.TopGreetedUsersRequest](handlers.TopGreetedUsersHandler))

	handlers.StartConsumer(nc)
	return svc
}


package service

import (

    "github.com/appetito/uno"
    "github.com/nats-io/nats.go"

	"github.com/rs/zerolog/log"

	"github.com/appetito/examples/greeter/api"
	"github.com/appetito/examples/greeter/internal/config"
	"github.com/appetito/examples/greeter/internal/handlers"

)

func New(cfg *config.Config) uno.Service {

	log.Info().Str("URL", cfg.NatsServers).Msg("Connecting to NATS")
	nc, err := nats.Connect(cfg.NatsServers)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}

	svc, err := uno.AddService(nc, uno.Config{
		Name:       "example" + "_" +  "Greeter",
		Version:     "",
		Description: "Greeter",
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


	root.AddEndpoint(api.GREET, uno.AsStructHandler[api.GreetRequest](handlers.GreetHandler))

	
	return svc
}
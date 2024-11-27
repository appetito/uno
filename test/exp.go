package main

import (
	"fmt"

	"github.com/appetito/uno"
	"github.com/nats-io/nats.go"

	// "github.com/nats-io/nats.go/uno"
	"github.com/rs/zerolog/log"
)


type Handler struct{
	wrapped uno.Handler
}

func (h Handler) Handle(req uno.Request) {
	log.Info().Str("data", string(req.Data())).Msg("Got request")
	// req.RespondJSON("kuku")
	err := req.RespondJSON("zzz")
	log.Info().Msgf("Reply %v",err)
}


func wrapper (h uno.HandlerFunc) uno.HandlerFunc {
	return func(req uno.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Msgf("Panic: %v", r)
				req.Error("500", "Internal Server Error", nil)
			}
		}()
		log.Info().Str("data", string(req.Data())).Msg("Handling request")
		h(req)

		log.Info().Str("data", string(req.Data())).Msg("Request handled successfully")
	}
 
}


func NewLoggingInterceptor (h uno.HandlerFunc) uno.HandlerFunc {
	return func(req uno.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Msgf("Panic: %v", r)
				req.Error("500", "Internal Server Error", nil)
			}
		}()
		log.Info().Str("data", string(req.Data())).Msg("Handling request")
		h(req)

		if req.HasError(){
			log.Info().Str("error", req.GetServiceError().Error()).Msg("Request handled with error")
		}else{
			log.Info().Str("data", string(req.Data())).Msg("Request handled successfully")
		}
	}
 
}

func NewHdrInterceptor (h uno.HandlerFunc) uno.HandlerFunc {
	return func(req uno.Request) {
		
		log.Info().Str("headers", fmt.Sprintf("%v", req.Headers())).Msg("Hdr")
		h(req)
	}
 
}


func main() {
	log.Info().Str("URL", nats.DefaultURL).Msg("Connecting to NATS")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}

	// database.Init(cfg)

	svc, err := uno.AddService(nc, uno.Config{
		Name:        "TestService",
		Version:     "0.0.1",
		Description: "TestService Controller",
		Interceptors: []uno.InterceptorFunc{
			NewHdrInterceptor,
			NewLoggingInterceptor,
		},
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to add service")
	}


	root := svc.AddGroup("test")

	root.AddEndpoint("get", Handler{})

	root.AddEndpoint("foo", wrapper(func(r uno.Request){
		panic("foo")
		r.RespondJSON("foo")
	}))
	
	root.AddEndpoint("bar", wrapper(func(r uno.Request){
		r.Error("404", "NotFound", []byte("bar not found"))
	}))
	
	root.AddEndpoint("int", uno.HandlerFunc(func(r uno.Request){
		r.Error("404", "NotFound", []byte("bar not found"))
	}))


	svc.Serve()
	// return nil
}


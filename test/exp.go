package main

import (
	"fmt"
	"os"
	"time"

	"github.com/appetito/uno"
	"github.com/nats-io/nats.go"

	// "github.com/nats-io/nats.go/uno"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)


type Handler struct{
}

type Foo struct {
	A string `json:"a"`
}

type Bar struct {
	B string `json:"b"`
	T uint   `json:"t"`
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

func NewHdrInterceptor (h uno.HandlerFunc) uno.HandlerFunc {
	return func(req uno.Request) {
		
		log.Info().Str("headers", fmt.Sprintf("%v", req.Headers())).Msg("Hdr")
		h(req)
	}
 
}


func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	log.Info().Str("URL", nats.DefaultURL).Msg("Connecting to NATS")
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to NATS")
	}
	log.With().Str("URL", nats.DefaultURL).Logger()
	// database.Init(cfg)


	svc, err := uno.AddService(nc, uno.Config{
		Name:        "TestService",
		Version:     "0.0.1",
		Description: "TestService Controller",
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

	root.AddEndpoint("struct_foo", uno.AsStructHandler[Foo](
		func(r uno.Request, f Foo) {
			r.Logger().Info().Str("data", string(r.Data())).Msgf("Handling Struct: %v", f)
			panic("azaza!")
			r.RespondJSON(f)
	
	}))

	root.AddEndpoint("struct_bar", uno.AsStructHandler[Bar](BarHandler))


	svc.ServeForever()
	// return nil
}

func BarHandler(r uno.Request, b Bar) {
	r.Logger().Info().Str("data", string(r.Data())).Msgf("Handling Struct: %v", b)
	time.Sleep(time.Duration(b.T) * time.Second)
	r.RespondJSON(b)

}
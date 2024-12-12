package uno

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	otelcodes "go.opentelemetry.io/otel/codes"
)

func NewPanicInterceptor (h HandlerFunc) HandlerFunc {
	return func(req Request) {
		defer func() {
			if r := recover(); r != nil {
				req.Logger().Error().Msgf("Panic: %v", r)
				req.Error("INTERNAL", "Internal Server Error", nil)
			}
		}()		
		h(req)
	}
 
}


func NewLoggingInterceptor (h HandlerFunc) HandlerFunc {
	return func(req Request) {
	
		// req.Logger().Debug().Msg("Log Int")
		req.Logger().
			Info().
			Str("data", string(req.Data())).
			// Str("headers", fmt.Sprintf("%v", req.Headers())).
			Msg("Handling request")

		h(req)

		if req.HasError(){
			e := req.GetServiceError()
			var level zerolog.Level = zerolog.InfoLevel
			if e.Code == "INTERNAL" {
				level = zerolog.ErrorLevel
			}
			req.Logger().WithLevel(level).
				Str("error", req.GetServiceError().Error()).
				Msg("Request handled with error")
		}else{
			req.Logger().Info().
				Dur("duration", time.Since(req.StartTime())).
				Msg("Request handled successfully")
		}
		// req.Logger().Debug().Msg("Log Int after")
	}
 
}


func NewTracingInterceptor (h HandlerFunc) HandlerFunc {
	return func(req Request) {

		// req.Logger().Debug().Msg("Trace Int")	
		ctx, span := Tracer.Start(req.Context(), req.Endpoint().Name)
		req.SetContext(ctx)
		defer span.End()

		h(req)

		if req.HasError() {
			e := req.GetServiceError()
			span.RecordError(e)
			span.SetStatus(otelcodes.Error, e.Error())
			// span.SetAttributes()
		}else{
			span.SetStatus(otelcodes.Ok, "")
		}
	
		// req.Logger().Debug().Msg("Trace Int after")
	}
 
}


func NewMetricsInterceptor (h HandlerFunc) HandlerFunc {
	return func(req Request) {

		// req.Logger().Debug().Msg("Prom Int")	

		h(req)

		requestsTotal.With(prometheus.Labels{"service": req.Endpoint().service.Name, "endpoint": req.Endpoint().Name, "status": req.Status()}).Inc()
		requestsDuration.With(prometheus.Labels{"service": req.Endpoint().service.Name, "endpoint": req.Endpoint().Name}).Observe(time.Since(req.StartTime()).Seconds())
		// req.Logger().Debug().Msg("Prom Int after")
	}
 
}
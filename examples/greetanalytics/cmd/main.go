package main

import (
	"github.com/appetito/uno/examples/greetanalytics/internal/config"
	"github.com/appetito/uno/examples/greetanalytics/internal/service"
)

func main(){
	cfg := config.GetConfig()
	svc := service.New(cfg)
	svc.ServeForever()
}

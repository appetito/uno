package main

import (
	"github.com/appetito/uno/examples/greeter/internal/config"
	"github.com/appetito/uno/examples/greeter/internal/service"
)

func main(){
	cfg := config.GetConfig()
	svc := service.New(cfg)
	svc.ServeForever()
}

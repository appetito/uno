
package main

import (

	"github.com/appetito/examples/greeter/internal/config"
	"github.com/appetito/examples/greeter/internal/service"
)

func main(){
	cfg := config.GetConfig()
	svc := service.New(cfg)
	svc.ServeForever()
}

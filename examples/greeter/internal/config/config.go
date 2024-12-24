
package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	NatsServers string `env:"NATS_SERVERS" env-default:"localhost:4222"`
}



func GetConfig() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}
	return &cfg
}

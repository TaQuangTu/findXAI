package config

import (
	"github.com/caarlos0/env"
)

type Config struct {
	//service
	PORT string `env:"PORT" default:"50051"`

	//postgres
	POSTGRES_DSN string `env:"POSTGRES_DSN"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}

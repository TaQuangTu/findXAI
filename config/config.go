package config

import "github.com/caarlos0/env"

type Config struct {
	//service
	PORT string `env:"PORT" default:"50051"`

	//postgres
	POSTGRES_DSN           string `env:"POSTGRES_DSN"`
	REDIS_LOCKDB_DNS       string `env:"REDIS_LOCKDB_DNS"`
	REDIS_RATE_LIMITDB_DNS string `env:"REDIS_RATE_LIMITDB_DNS"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}

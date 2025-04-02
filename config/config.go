package config

import "github.com/caarlos0/env"

type Config struct {
	//service
	PORT string `env:"PORT" default:"50051"`

	//postgres
	POSTGRES_DSN           string `env:"POSTGRES_DSN"`
	REDIS_LOCKDB_DSN       string `env:"REDIS_LOCKDB_DSN"`
	REDIS_RATE_LIMITDB_DSN string `env:"REDIS_RATE_LIMITDB_DSN"`

	// default requests per day is 100
	MAX_REQUEST_PER_DAY int `env:"MAX_REQUEST_PER_DAY" default:"100"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	//always 100
	cfg.MAX_REQUEST_PER_DAY = 100
	if err != nil {
		panic(err)
	}
	return &cfg
}

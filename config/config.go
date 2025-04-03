package config

import "github.com/caarlos0/env"

type Config struct {
	//service
	HOST string `env:"HOST" default:"localhost"`
	PORT string `env:"PORT" default:"50051"`

	//postgres
	POSTGRES_DSN           string `env:"POSTGRES_DSN"`
	REDIS_LOCKDB_DSN       string `env:"REDIS_LOCKDB_DSN"`
	REDIS_RATE_LIMITDB_DSN string `env:"REDIS_RATE_LIMITDB_DSN"`

	// default requests per day is 100
	MAX_REQUEST_PER_DAY int `env:"MAX_REQUEST_PER_DAY" envDefault:"100"`
	APP_KEY_BUCKET      int `env:"APP_KEY_BUCKET" envDefault:"3"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}

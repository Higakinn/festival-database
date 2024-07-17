package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env  string `env:"TODO_ENV" env_Default:"dev"`
	Port int    `env:"PORT" env_Default:"80"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

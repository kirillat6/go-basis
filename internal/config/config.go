package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        string `envconfig:"PORT" 		 default:"8080" `
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
}

func Load() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
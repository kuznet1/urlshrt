package config

import (
	"flag"
	"github.com/caarlos0/env"
)

type Config struct {
	ListenAddr      string `env:"SERVER_ADDRESS"`
	ShortenerPrefix string `env:"BASE_URL"`
}

func ParseArgs() Config {
	cfg := Config{}
	flag.StringVar(&cfg.ListenAddr, "a", ":8080", "address to listen on")
	flag.StringVar(&cfg.ShortenerPrefix, "b", "http://localhost:8080", "prefix for url shortening")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}

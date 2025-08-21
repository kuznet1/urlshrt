package config

import (
	"flag"
	"github.com/caarlos0/env"
)

type Config struct {
	ListenAddr      string `env:"SERVER_ADDRESS"`
	ShortenerPrefix string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func ParseArgs() (Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.ListenAddr, "a", ":8080", "address to listen on")
	flag.StringVar(&cfg.ShortenerPrefix, "b", "http://localhost:8080", "prefix for url shortening")
	flag.StringVar(&cfg.FileStoragePath, "f", "repo.json", "file storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", "postgres://postgres@localhost:5432/urlshrt", "database connection string")
	flag.Parse()

	return cfg, env.Parse(&cfg)
}

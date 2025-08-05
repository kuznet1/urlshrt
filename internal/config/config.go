package config

import "flag"

type Config struct {
	ListenAddr      string
	ShortenerPrefix string
}

func ParseArgs() Config {
	cfg := Config{}
	flag.StringVar(&cfg.ListenAddr, "a", ":8080", "address to listen on")
	flag.StringVar(&cfg.ShortenerPrefix, "b", "http://localhost:8080", "prefix for url shortening")
	flag.Parse()
	return cfg
}

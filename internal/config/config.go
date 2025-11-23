package config

import (
	"flag"
	"github.com/caarlos0/env"
	"time"
)

// Config contains runtime configuration loaded from flags and environment variables.
// See struct tags for env variable names; command-line flags mirror these fields.
type Config struct {
	ListenAddr         string        `env:"SERVER_ADDRESS"`
	ShortenerPrefix    string        `env:"BASE_URL"`
	FileStoragePath    string        `env:"FILE_STORAGE_PATH"`
	DatabaseDSN        string        `env:"DATABASE_DSN"`
	SecretKey          string        `env:"SECRET_KEY"`
	DeleteBatchSize    int           `env:"DELETE_BATCH_SIZE"`
	DeleteBatchTimeout time.Duration `env:"DELETE_BATCH_TIMEOUT"`
	AuditFile          string        `env:"AUDIT_FILE"`
	AuditURL           string        `env:"AUDIT_URL"`
	AuditURLTimeout    time.Duration `env:"AUDIT_URL_REQ_TIMEOUT"`
}

// ParseArgs populates Config from command-line flags and environment variables.
// Environment variables take the form described by struct tags (e.g., SERVER_ADDRESS, BASE_URL).
func ParseArgs() (Config, error) {
	cfg := Config{}
	flag.StringVar(&cfg.ListenAddr, "a", ":8080", "address to listen on")
	flag.StringVar(&cfg.ShortenerPrefix, "b", "http://localhost:8080", "prefix for url shortening")
	flag.StringVar(&cfg.FileStoragePath, "f", "", "file storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", "", "database connection string")
	flag.StringVar(&cfg.SecretKey, "k", "", "secret key for cookie signing")
	flag.IntVar(&cfg.DeleteBatchSize, "bs", 1, "delete batch size")
	flag.DurationVar(&cfg.DeleteBatchTimeout, "t", time.Second, "delete timeout")
	flag.StringVar(&cfg.AuditFile, "af", "", "file to save audit logs")
	flag.StringVar(&cfg.AuditURL, "au", "", "url to send audit logs to")
	flag.DurationVar(&cfg.AuditURLTimeout, "aut", 10*time.Second, "audit request timeout")
	flag.Parse()

	return cfg, env.Parse(&cfg)
}

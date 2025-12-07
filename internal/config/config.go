package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env"
	"os"
	"time"
)

// Config contains runtime configuration loaded from flags and environment variables.
// See struct tags for env variable names; command-line flags mirror these fields.
type Config struct {
	ConfigPath         string        `env:"CONFIG"`
	ListenAddr         string        `env:"SERVER_ADDRESS" json:"server_address"`
	HTTPReadTimeout    time.Duration `env:"HTTP_READ_TIMEOUT"  json:"http_read_timeout"`
	HTTPWriteTimeout   time.Duration `env:"HTTP_WRITE_TIMEOUT" json:"http_write_timeout"`
	HTTPIdleTimeout    time.Duration `env:"HTTP_IDLE_TIMEOUT"  json:"http_idle_timeout"`
	ShortenerPrefix    string        `env:"BASE_URL" json:"base_url"`
	FileStoragePath    string        `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DatabaseDSN        string        `env:"DATABASE_DSN" json:"database_dsn"`
	SecretKey          string        `env:"SECRET_KEY" json:"secret_key"`
	DeleteBatchSize    int           `env:"DELETE_BATCH_SIZE" json:"delete_batch_size"`
	DeleteBatchTimeout time.Duration `env:"DELETE_BATCH_TIMEOUT" json:"delete_batch_timeout"`
	AuditFile          string        `env:"AUDIT_FILE" json:"audit_file"`
	AuditURL           string        `env:"AUDIT_URL" json:"audit_url"`
	AuditURLTimeout    time.Duration `env:"AUDIT_URL_REQ_TIMEOUT" json:"audit_url_req_timeout"`
	EnableHTTPS        bool          `env:"ENABLE_HTTPS" json:"enable_https"`
	HTTPSCertFile      string        `env:"HTTPS_CERT_FILE" json:"https_cert_file"`
	HTTPSCertKey       string        `env:"HTTPS_CERT_KEY" json:"https_cert_key"`
	ShutdownTimeout    time.Duration `env:"SHUTDOWN_TIMEOUT"  json:"shutdown_timeout"`
}

func parseConfig(cfg *Config) error {
	err := env.Parse(cfg)
	if err != nil {
		return fmt.Errorf("unable to get env `CONFIG`: %w", err)
	}

	if cfg.ConfigPath == "" {
		return nil
	}

	f, err := os.Open(cfg.ConfigPath)
	if err != nil {
		return fmt.Errorf("unable to open config: %w", err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return fmt.Errorf("unable to parse config: %w", err)
	}

	return nil
}

// ParseArgs populates Config from command-line flags and environment variables.
// Environment variables take the form described by struct tags (e.g., SERVER_ADDRESS, BASE_URL).
func ParseArgs() (Config, error) {
	cfg := Config{
		ConfigPath:         "",
		ListenAddr:         ":8080",
		HTTPReadTimeout:    10 * time.Second,
		HTTPWriteTimeout:   10 * time.Second,
		HTTPIdleTimeout:    60 * time.Second,
		ShortenerPrefix:    "http://localhost:8080",
		FileStoragePath:    "",
		DatabaseDSN:        "",
		SecretKey:          "",
		DeleteBatchSize:    1,
		DeleteBatchTimeout: time.Second,
		AuditFile:          "",
		AuditURL:           "",
		AuditURLTimeout:    10 * time.Second,
		EnableHTTPS:        false,
		HTTPSCertFile:      "cert.pem",
		HTTPSCertKey:       "key.pem",
		ShutdownTimeout:    15 * time.Second,
	}

	flag.StringVar(&cfg.ConfigPath, "config", cfg.ConfigPath, "path to config")
	flag.StringVar(&cfg.ConfigPath, "c", cfg.ConfigPath, "path to config (shorthand)")
	flag.StringVar(&cfg.ListenAddr, "a", cfg.ListenAddr, "address to listen on")
	flag.DurationVar(&cfg.HTTPReadTimeout, "rt", cfg.HTTPReadTimeout, "http read timeout")
	flag.DurationVar(&cfg.HTTPWriteTimeout, "wt", cfg.HTTPWriteTimeout, "http write timeout")
	flag.DurationVar(&cfg.HTTPIdleTimeout, "it", cfg.HTTPIdleTimeout, "http idle timeout")
	flag.StringVar(&cfg.ShortenerPrefix, "b", cfg.ShortenerPrefix, "prefix for url shortening")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database connection string")
	flag.StringVar(&cfg.SecretKey, "k", cfg.SecretKey, "secret key for cookie signing")
	flag.IntVar(&cfg.DeleteBatchSize, "bs", cfg.DeleteBatchSize, "delete batch size")
	flag.DurationVar(&cfg.DeleteBatchTimeout, "t", cfg.DeleteBatchTimeout, "delete timeout")
	flag.StringVar(&cfg.AuditFile, "af", cfg.AuditFile, "file to save audit logs")
	flag.StringVar(&cfg.AuditURL, "au", cfg.AuditURL, "url to send audit logs to")
	flag.DurationVar(&cfg.AuditURLTimeout, "aut", cfg.AuditURLTimeout, "audit request timeout")
	flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "enable https")
	flag.StringVar(&cfg.HTTPSCertFile, "sc", cfg.HTTPSCertFile, "x509 certificate file")
	flag.StringVar(&cfg.HTTPSCertKey, "sk", cfg.HTTPSCertKey, "x509 certificate key")
	flag.DurationVar(&cfg.ShutdownTimeout, "st", cfg.ShutdownTimeout, "shutdown timeout")
	flag.Parse()

	err := parseConfig(&cfg)
	if err != nil {
		return cfg, err
	}

	flag.Parse()
	return cfg, env.Parse(&cfg)
}

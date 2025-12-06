// Package config
package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Address   string
	Mode      string
	Interval  time.Duration
	LogLevel  string
	LogFormat string
}

const (
	ModeServe    = "serve"
	ModeStream   = "stream"
	ModeSnapshot = "snapshot"
)

func Load() *Config {
	godotenv.Load()

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":3000"
	}

	interval := time.Second
	if raw := os.Getenv("SCRAPE_INTERVAL"); raw != "" {
		if parsed, err := time.ParseDuration(raw); err == nil && parsed > 0 {
			interval = parsed
		}
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	logFormat := os.Getenv("LOG_FORMAT")
	if logFormat == "" {
		logFormat = "text"
	}

	return &Config{
		Address:   addr,
		Mode:      ModeServe,
		Interval:  interval,
		LogLevel:  logLevel,
		LogFormat: logFormat,
	}
}

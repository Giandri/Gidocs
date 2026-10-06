package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port               string
	MaxUploadMB        int
	MaxFiles           int
	RequestTimeout     time.Duration
	CORSAllowedOrigins string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	maxUploadMB := 50
	if v := os.Getenv("MAX_UPLOAD_MB"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			maxUploadMB = parsed
		}
	}

	maxFiles := 20
	if v := os.Getenv("MAX_FILES"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			maxFiles = parsed
		}
	}

	timeout := 120 * time.Second
	if v := os.Getenv("REQUEST_TIMEOUT"); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			timeout = parsed
		}
	}

	cors := os.Getenv("CORS_ALLOWED_ORIGINS")

	return Config{
		Port:               port,
		MaxUploadMB:        maxUploadMB,
		MaxFiles:           maxFiles,
		RequestTimeout:     timeout,
		CORSAllowedOrigins: cors,
	}
}
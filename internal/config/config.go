package config

import (
	"fmt"
	"os"
)

type Config struct {
	Addr string
	DB   DBconfig
}

type DBconfig struct {
	User, Password, Database, Addr string
}

func envOr(key, fallback string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	return fallback
}

func Load() (*Config, error) {
	for _, key := range []string{"POSTGRES_USER", "POSTGRES_PASSWORD", "POSTGRES_DB"} {
		if (os.Getenv(key)) == "" {
			return nil, fmt.Errorf("required env var %s is not set", key)
		}
	}
	host := envOr("POSTGRES_HOST", "localhost")
	port := envOr("POSTGRES_PORT", "5432")
	cfg := &Config{
		Addr: ":" + envOr("APP_PORT", "8080"),
		DB: DBconfig{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
			Addr:     host + ":" + port,
		},
	}
	return cfg, nil
}

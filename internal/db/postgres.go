package db

import (
	"os"

	"github.com/go-pg/pg/v10"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	return value
}

func Connect() *pg.DB {
	database := pg.Connect(&pg.Options{
		User:     getEnv("POSTGRES_USER"),
		Password: getEnv("POSTGRES_PASSWORD"),
		Database: getEnv("POSTGRES_DB"),
		Addr:     getEnv("POSTGRES_HOST") + ":" + getEnv("POSTGRES_PORT"),
	})

	return database
}

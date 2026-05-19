package db

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/v3lichko/student-distribution/internal/config"
)

func Connect(cfg config.DBconfig) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Database,
		Addr:     cfg.Addr,
	})
	if _, err := db.Exec("SELECT 1"); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}

package db

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"auth/service/internal/config"
	"auth/service/internal/logger"
)

func Connect(cfg *config.Config) (*sql.DB, error) {
	logger.Init()

	log := logger.L()
	db, err := sql.Open("postgres", cfg.DBConnString())
	if err != nil {
		log.Errorf("Failed to open DB connection: %v\n", err)
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		log.Errorf("Failed to ping DB: %v\n", err)
	}

	log.Info("Connected to PostgreSql successfully")

	return db, nil
}

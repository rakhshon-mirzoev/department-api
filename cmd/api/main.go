package main

import (
	"net/http"
	"os"

	"github.com/pressly/goose/v3"
	"github.com/rakhshon-mirzoev/department-api/internal/config"
	"github.com/rakhshon-mirzoev/department-api/internal/logger"
	"github.com/rakhshon-mirzoev/department-api/internal/repository"
	"github.com/rakhshon-mirzoev/department-api/pkg/db"
)

func main() {
	log := logger.New()

	cfg := config.Load()

	gormdb, err := db.NewPostgres(cfg.DB)
	if err != nil {
		log.Error("failed to init db", "err", err)
		os.Exit(1)
	}
	log.Info("db connected")

	sqlDb, err := gormdb.DB()
	if err != nil {
		log.Error("failed to get db", "err", err)
		os.Exit(1)
	}

	if err := goose.Up(sqlDb, "migrations"); err != nil {
		log.Error("failed to migrate", "err", err)
		os.Exit(1)
	}

	_ = repository.NewRepository(gormdb)

	log.Info("started", "port", 8080)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Error("failed to start", "err", err)
		os.Exit(1)
	}
}

package main

import (
	"log"
	"net/http"

	"github.com/pressly/goose/v3"
	"github.com/rakhshon-mirzoev/department-api/internal/config"
	"github.com/rakhshon-mirzoev/department-api/internal/repository"
	"github.com/rakhshon-mirzoev/department-api/pkg/db"
)

func main() {
	cfg := config.Load()
	gormdb, err := db.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	sqlDb, err := gormdb.DB()
	if err != nil {
		log.Fatalf("failed to get db: %v", err)
	}

	if err := goose.Up(sqlDb, "migrations"); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	_ = repository.NewRepository(gormdb)

	log.Println("server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

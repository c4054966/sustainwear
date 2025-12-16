package main

import (
	"log"
	"net/http"

	"sustainwear/internal/api"
	"sustainwear/internal/config"
	"sustainwear/internal/storage"
)

func main() {
	// LOAD CONFIG
	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// GET DB CONNECTION STRING
	connStr, err := cfg.Database.GetConnectionString()
	if err != nil {
		log.Fatalf("Failed to get database connection string: %v", err)
	}

	// INITIALISE DB
	db, err := storage.NewDB(cfg.Database.Driver, connStr)
	if err != nil {
		log.Fatalf("Failed to initialise database: %v", err)
	}
	defer db.Close()

	// INITIALISE ROUTER
	router := api.NewRouter(cfg, db.DB)

	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("SustainWear Server listening on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

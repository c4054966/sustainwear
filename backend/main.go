package main

import (
	"log"
	"net/http"

	"sustainwear/internal/api"
	"sustainwear/internal/config"
	"sustainwear/internal/storage"
)

func main() {
	// LOAD CONFIGURATION
	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting SustainWear API server...")

	// GET DATABASE CONNECTION STRING
	connStr, err := cfg.Database.GetConnectionString()
	if err != nil {
		log.Fatalf("Failed to get database connection string: %v", err)
	}

	// INITIALIZE DATABASE
	db, err := storage.NewDB(cfg.Database.Driver, connStr)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// INITIALIZE ROUTER
	router := api.NewRouter(cfg, db.DB)

	// START SERVER
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server listening on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

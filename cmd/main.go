package main

import (
	"log"

	"spotsync/internal/config"
	"spotsync/internal/server"
)

func main() {
	cfg := config.Load()

	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	app := server.NewHTTPServer(cfg, db)

	if err := app.Start(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

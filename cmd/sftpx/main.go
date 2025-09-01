// cmd/sftpx/main.go
package main

import (
	"log"

	"sftpx/internal/config"
	"sftpx/internal/logger"
	"sftpx/internal/watcher"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("configs/config.json")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Setup logger
	logger.Setup(cfg.LogDir, cfg.LogFile)

	log.Println("Starting SFTPX...")

	// Start folder watcher
	if err := watcher.Start(cfg); err != nil {
		log.Fatalf("watcher failed: %v", err)
	}
}

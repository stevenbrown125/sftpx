package logger

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Setup configures logging to a timestamped file in a given dir
func Setup(logDir, logFile string) {
	if strings.HasPrefix(logDir, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			logDir = filepath.Join(home, logDir[2:])
		}
	}

	// Ensure directory exists
	if logDir != "" {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("failed to create log dir %s: %v", logDir, err)
		}
	}

	// Split filename and add timestamp
	ext := filepath.Ext(logFile) // ".log"
	base := logFile[:len(logFile)-len(ext)]
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileWithTs := base + "-" + timestamp + ext

	if logDir != "" {
		logFileWithTs = filepath.Join(logDir, logFileWithTs)
	}

	// Open log file
	f, err := os.OpenFile(logFileWithTs, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("Logging initialized → %s", logFileWithTs)
}

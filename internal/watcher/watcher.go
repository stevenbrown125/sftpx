package watcher

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"

	"sftpx/internal/config"
	"sftpx/internal/sftp"
)

func Start(cfg *config.Config) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	if err := watcher.Add(cfg.WatchDir); err != nil {
		return err
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				info, err := os.Stat(event.Name)
				if err != nil {
					log.Println("Stat error:", err)
					continue
				}

				if info.IsDir() {
					log.Printf("Detected new folder: %s", event.Name)
					go queueFolder(cfg, event.Name)
				} else {
					log.Printf("Detected new file: %s", event.Name)
					go queueFile(cfg, event.Name)
				}
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Println("Watcher error:", err)
		}
	}
}

func queueFolder(cfg *config.Config, path string) {
	time.Sleep(time.Duration(cfg.DelaySeconds) * time.Second)

	client, err := sftp.NewClient(cfg)
	if err != nil {
		log.Println("SFTP connection failed:", err)
		return
	}
	defer client.Close()

	jobs := make(chan string)
	done := make(chan bool)

	workerCount := cfg.Workers
	if workerCount < 1 {
		workerCount = 1
	}

	// Start worker pool
	for i := 0; i < workerCount; i++ {
		go uploadWorker(cfg, client, jobs, done)
	}

	// Walk folder and send jobs
	filepath.Walk(path, func(localPath string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("Walk error:", err)
			return nil
		}
		if !info.IsDir() {
			jobs <- localPath
		}
		return nil
	})
	close(jobs)

	// Wait for all workers
	for i := 0; i < workerCount; i++ {
		<-done
	}
}

func queueFile(cfg *config.Config, localPath string) {
	time.Sleep(time.Duration(cfg.DelaySeconds) * time.Second)

	client, err := sftp.NewClient(cfg)
	if err != nil {
		log.Println("SFTP connection failed:", err)
		return
	}
	defer client.Close()

	jobs := make(chan string, 1)
	done := make(chan bool)

	workerCount := cfg.Workers
	if workerCount < 1 {
		workerCount = 1
	}

	// Start worker pool
	for i := 0; i < workerCount; i++ {
		go uploadWorker(cfg, client, jobs, done)
	}

	// Send single file as a job
	jobs <- localPath
	close(jobs)

	// Wait for all workers
	for i := 0; i < workerCount; i++ {
		<-done
	}
}

func uploadWorker(cfg *config.Config, client *sftp.Client, jobs <-chan string, done chan<- bool) {
	for localPath := range jobs {
		rel, _ := filepath.Rel(cfg.WatchDir, localPath)
		remotePath := cfg.RemoteDir + "/" + filepath.ToSlash(rel)

		log.Printf("Uploading %s → %s", localPath, remotePath)
		if err := sftp.UploadFile(client, localPath, remotePath); err != nil {
			log.Printf("Upload failed for %s: %v", localPath, err)
		} else {
			log.Printf("Upload complete: %s → %s", localPath, remotePath)
		}
	}
	done <- true
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"video-transcoding-service/internal/config"
	"video-transcoding-service/internal/scheduler"
)

func main() {
	cfg := config.Load()

	sched, err := scheduler.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize scheduler: %v", err)
	}

	// Setup graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Run scheduler in a goroutine
	go func() {
		if err := sched.Run(); err != nil {
			log.Fatalf("Scheduler failed: %v", err)
		}
	}()

	log.Println("Service started successfully")
	<-shutdown
	log.Println("Shutting down...")

	// Add any cleanup logic here
	log.Println("Service stopped")
}

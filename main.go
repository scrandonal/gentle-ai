package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gentle-ai/gentle-ai/internal/config"
	"github.com/gentle-ai/gentle-ai/internal/server"
)

// Version is set at build time via ldflags
var Version = "dev"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration from environment / config file
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("starting gentle-ai %s on %s", Version, cfg.Addr())

	// Initialise and start the HTTP server
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	// Run server in background goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := srv.Start(ctx); err != nil {
			errCh <- fmt.Errorf("server error: %w", err)
		}
	}()

	// Wait for OS signal or server error
	// Listening for SIGINT (Ctrl+C) and SIGTERM (e.g. from Docker/systemd)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("received signal %s, shutting down gracefully", sig)
	case err := <-errCh:
		// If the server dies on its own, log it clearly and attempt clean shutdown
		log.Printf("server stopped unexpectedly: %v", err)
	}

	// Trigger graceful shutdown
	cancel()
	if err := srv.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
		os.Exit(1)
	}

	log.Println("gentle-ai stopped")
}

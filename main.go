package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/little-secrets/amnezia-wg-easy/internal/api"
	"github.com/little-secrets/amnezia-wg-easy/internal/config"
	"github.com/little-secrets/amnezia-wg-easy/internal/wireguard"
)

func main() {
	// Load configuration from environment variables
	cfg := config.Load()

	// Print startup info
	log.Printf("AmneziaWG Easy v%s", cfg.Release)
	log.Printf("Web UI: http://%s:%s", cfg.WebUIHost, cfg.Port)

	if cfg.NoWebUI {
		log.Println("Web UI is disabled (NO_WEB_UI=true)")
	}

	if cfg.WGHost == "" {
		log.Fatal("Error: WG_HOST environment variable is required")
	}

	// Initialize WireGuard service
	wg := wireguard.New(cfg)
	if err := wg.Init(); err != nil {
		log.Fatalf("Failed to initialize WireGuard: %v", err)
	}

	// Setup HTTP router
	router := api.SetupRouter(cfg, wg)

	// Create HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.WebUIHost, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start cron job for periodic tasks
	go startCronJob(wg)

	// Start server in goroutine
	go func() {
		log.Printf("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Shutdown WireGuard
	if err := wg.Shutdown(); err != nil {
		log.Printf("Error shutting down WireGuard: %v", err)
	}

	log.Println("Server exited")
}

// startCronJob runs periodic tasks every minute
func startCronJob(wg *wireguard.WireGuard) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := wg.CronJob(); err != nil {
			log.Printf("Cron job error: %v", err)
		}
	}
}

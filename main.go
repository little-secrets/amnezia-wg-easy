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

	// Audit C-2: with PASSWORD_HASH empty, the auth middleware used to
	// fall through (`!cfg.RequiresPassword() => c.Next()`) and expose
	// every protected route -- including peer creation, deletion, key
	// retrieval, and full-config backup -- to any caller on the network.
	// Refuse to start unless the operator explicitly opted in via
	// NO_AUTH=true; in that mode force the listener to 127.0.0.1 so a
	// passwordless deployment can never reach the public Internet by
	// mistake.
	if cfg.PasswordHash == "" {
		if !cfg.NoAuth {
			log.Fatal("Error: PASSWORD_HASH is empty. Set a bcrypt hash " +
				"(see cmd/wgpw) or, for a strictly local development run, " +
				"set NO_AUTH=true to acknowledge that the API will accept " +
				"unauthenticated requests.")
		}
		log.Println("WARNING: NO_AUTH=true -- API is unauthenticated. " +
			"This is intended for local development only.")
		// Hard-pin the bind address so a misconfigured docker-compose or
		// systemd unit cannot accidentally expose the no-auth API.
		if cfg.WebUIHost != "127.0.0.1" && cfg.WebUIHost != "::1" && cfg.WebUIHost != "localhost" {
			log.Printf("WARNING: NO_AUTH mode forces WEBUI_HOST from %q to 127.0.0.1", cfg.WebUIHost)
			cfg.WebUIHost = "127.0.0.1"
		}
	}

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

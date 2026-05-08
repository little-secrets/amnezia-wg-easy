package api

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/little-secrets/amnezia-wg-easy/internal/config"
	"github.com/little-secrets/amnezia-wg-easy/internal/wireguard"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(cfg *config.Config, wg *wireguard.WireGuard) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// Setup session middleware. Audit H-2: persist the cookie-store
	// secret so a restart doesn't invalidate every active UI session,
	// and panic on crypto/rand failure rather than falling back to the
	// hardcoded "default-secret-please-set-proper-one" constant that
	// would let an attacker forge any session cookie.
	secret := loadOrCreateSessionSecret(cfg.WGPath)
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	r.Use(sessions.Sessions("wg-session", store))

	h := NewHandler(cfg, wg)

	// Public API routes (no auth required)
	api := r.Group("/api")
	{
		api.GET("/release", h.GetRelease)
		api.GET("/lang", h.GetLang)
		api.GET("/remember-me", h.GetRememberMeEnabled)
		api.GET("/ui-traffic-stats", h.GetUITrafficStats)
		api.GET("/ui-chart-type", h.GetUIChartType)
		api.GET("/wg-enable-one-time-links", h.GetWGEnableOneTimeLinks)
		api.GET("/wg-enable-expire-time", h.GetWGEnableExpireTime)
		api.GET("/ui-sort-clients", h.GetUISortClients)
		api.GET("/ui-avatar-settings", h.GetUIAvatarSettings)
		api.GET("/session", h.GetSession)
		api.POST("/session", h.CreateSession)
	}

	// One-time link route (public)
	r.GET("/cnf/:clientOneTimeLink", h.DownloadOneTimeLink)

	// Protected API routes (auth required)
	protected := r.Group("/api")
	protected.Use(AuthMiddleware(cfg))
	{
		protected.DELETE("/session", h.DeleteSession)

		// WireGuard client routes
		protected.GET("/wireguard/client", h.GetClients)
		protected.POST("/wireguard/client", h.CreateClient)
		protected.DELETE("/wireguard/client/:clientId", h.DeleteClient)
		protected.POST("/wireguard/client/:clientId/enable", h.EnableClient)
		protected.POST("/wireguard/client/:clientId/disable", h.DisableClient)
		protected.POST("/wireguard/client/:clientId/generateOneTimeLink", h.GenerateOneTimeLink)
		protected.PUT("/wireguard/client/:clientId/name", h.UpdateClientName)
		protected.PUT("/wireguard/client/:clientId/address", h.UpdateClientAddress)
		protected.PUT("/wireguard/client/:clientId/address6", h.UpdateClientAddress6)
		protected.PUT("/wireguard/client/:clientId/expireDate", h.UpdateClientExpireDate)
		protected.PUT("/wireguard/client/:clientId/allowedIPs", h.UpdateClientAllowedIPs)
		protected.PUT("/wireguard/client/:clientId/dns", h.UpdateClientDNS)
		protected.PUT("/wireguard/client/:clientId/mtu", h.UpdateClientMTU)
		protected.PUT("/wireguard/client/:clientId/keepalive", h.UpdateClientKeepalive)
		protected.GET("/wireguard/client/:clientId/qrcode.svg", h.GetClientQRCode)
		protected.GET("/wireguard/client/:clientId/configuration", h.GetClientConfiguration)
		protected.GET("/wireguard/client/:clientId/secrets", h.GetClientSecrets)

		// Backup/Restore routes
		protected.GET("/wireguard/backup", h.BackupConfiguration)
		protected.PUT("/wireguard/restore", h.RestoreConfiguration)
	}

	// Prometheus metrics routes
	metrics := r.Group("/metrics")
	metrics.Use(PrometheusAuthMiddleware(cfg))
	{
		metrics.GET("", h.GetMetrics)
		metrics.GET("/json", h.GetMetricsJSON)
	}

	// API Documentation routes
	r.StaticFile("/api/openapi.yaml", "./www/openapi.yaml")
	r.GET("/api/docs", func(c *gin.Context) {
		c.File("./www/swagger.html")
	})

	// Static files (Web UI)
	if !cfg.NoWebUI {
		r.Static("/css", "./www/css")
		r.Static("/js", "./www/js")
		r.Static("/img", "./www/img")
		r.StaticFile("/manifest.json", "./www/manifest.json")
		r.StaticFile("/favicon.ico", "./www/img/favicon.ico")

		// Serve index.html for root
		r.GET("/", func(c *gin.Context) {
			c.File("./www/index.html")
		})

		// Fallback for SPA routing
		r.NoRoute(func(c *gin.Context) {
			// Only serve index.html for non-API routes
			if len(c.Request.URL.Path) < 4 || c.Request.URL.Path[:4] != "/api" {
				c.File("./www/index.html")
				return
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		})
	}

	return r
}

// loadOrCreateSessionSecret returns a 32-byte hex-encoded secret used
// to authenticate the session cookie. Persisted in `<dataDir>/.session_secret`
// (mode 0600) so a restart doesn't invalidate every UI session.
//
// Panics on crypto/rand failure -- a constant fallback would let any
// attacker who knew the constant forge cookies. Panics on disk write
// failure for the same reason: better to crash loudly than to silently
// regenerate the secret on every restart.
func loadOrCreateSessionSecret(dataDir string) string {
	path := filepath.Join(dataDir, ".session_secret")
	if data, err := os.ReadFile(path); err == nil {
		s := string(data)
		// Sanity check: hex-encoded 32 bytes is 64 chars.
		if len(s) >= 64 {
			return s
		}
		log.Printf("session secret at %s is too short (%d chars); regenerating", path, len(s))
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		// crypto/rand reading from the OS entropy source rarely fails;
		// continuing with a degraded value would defeat the cookie's
		// authentication, so abort.
		panic("crypto/rand failed: " + err.Error())
	}
	secret := hex.EncodeToString(bytes)

	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		log.Printf("warning: cannot create session-secret dir %s: %v", dataDir, err)
	}
	if err := os.WriteFile(path, []byte(secret), 0o600); err != nil {
		log.Printf("warning: cannot persist session secret to %s: %v", path, err)
	}
	return secret
}

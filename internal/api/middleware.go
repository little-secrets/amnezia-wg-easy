package api

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/little-secrets/amnezia-wg-easy/internal/config"
	"github.com/little-secrets/amnezia-wg-easy/internal/models"
)

// AuthMiddleware checks if the user is authenticated
func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if no password required
		if !cfg.RequiresPassword() {
			c.Next()
			return
		}

		// Check session
		session := sessions.Default(c)
		if session.Get("authenticated") == true {
			c.Next()
			return
		}

		// Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && checkPassword(authHeader, cfg.PasswordHash) {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Not Logged In"})
		c.Abort()
	}
}

// PrometheusAuthMiddleware checks Basic Auth for prometheus metrics
func PrometheusAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.RequiresPrometheusPassword() {
			c.Next()
			return
		}

		_, password, ok := c.Request.BasicAuth()
		if !ok {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Not Logged In"})
			c.Abort()
			return
		}

		if !checkPassword(password, cfg.PrometheusMetricsPassword) {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Incorrect Password"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// checkPassword compares a password with a bcrypt hash
func checkPassword(password, hash string) bool {
	if hash == "" || password == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


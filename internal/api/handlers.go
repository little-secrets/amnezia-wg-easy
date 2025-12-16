package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/little-secrets/amnezia-wg-easy/internal/config"
	"github.com/little-secrets/amnezia-wg-easy/internal/models"
	"github.com/little-secrets/amnezia-wg-easy/internal/wireguard"
)

// Handler holds all HTTP handlers
type Handler struct {
	cfg *config.Config
	wg  *wireguard.WireGuard
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config, wg *wireguard.WireGuard) *Handler {
	return &Handler{
		cfg: cfg,
		wg:  wg,
	}
}

// GetRelease returns the application version
func (h *Handler) GetRelease(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfg.Release)
}

// GetLang returns the configured language
func (h *Handler) GetLang(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfg.Lang)
}

// GetRememberMeEnabled returns whether remember me is enabled
func (h *Handler) GetRememberMeEnabled(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfg.MaxAge > 0)
}

// GetUITrafficStats returns traffic stats setting
func (h *Handler) GetUITrafficStats(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfg.EnableTrafficStats)
}

// GetUIChartType returns chart type setting
func (h *Handler) GetUIChartType(c *gin.Context) {
	c.JSON(http.StatusOK, fmt.Sprintf("%d", h.cfg.ChartType))
}

// GetWGEnableOneTimeLinks returns one-time links setting
func (h *Handler) GetWGEnableOneTimeLinks(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfg.EnableOneTimeLinks)
}

// GetWGEnableExpireTime returns expire time setting
func (h *Handler) GetWGEnableExpireTime(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfg.EnableExpiresTime)
}

// GetUISortClients returns sort clients setting
func (h *Handler) GetUISortClients(c *gin.Context) {
	c.JSON(http.StatusOK, h.cfg.EnableSortClients)
}

// GetUIAvatarSettings returns avatar settings
func (h *Handler) GetUIAvatarSettings(c *gin.Context) {
	c.JSON(http.StatusOK, models.AvatarSettings{
		Dicebear: h.cfg.DicebearType,
		Gravatar: h.cfg.UseGravatar,
	})
}

// GetSession returns current session status
func (h *Handler) GetSession(c *gin.Context) {
	session := sessions.Default(c)
	authenticated := !h.cfg.RequiresPassword() || session.Get("authenticated") == true

	c.JSON(http.StatusOK, models.SessionResponse{
		RequiresPassword: h.cfg.RequiresPassword(),
		Authenticated:    authenticated,
	})
}

// CreateSession handles login
func (h *Handler) CreateSession(c *gin.Context) {
	if !h.cfg.RequiresPassword() {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Invalid state"})
		return
	}

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}

	if !checkPassword(req.Password, h.cfg.PasswordHash) {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{Error: "Incorrect Password"})
		return
	}

	session := sessions.Default(c)
	session.Set("authenticated", true)
	if h.cfg.MaxAge > 0 && req.Remember {
		session.Options(sessions.Options{
			MaxAge: h.cfg.MaxAge / 1000, // Convert ms to seconds
		})
	}
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// DeleteSession handles logout
func (h *Handler) DeleteSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// GetClients returns all clients
func (h *Handler) GetClients(c *gin.Context) {
	clients, err := h.wg.GetClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, clients)
}

// CreateClient creates a new client
func (h *Handler) CreateClient(c *gin.Context) {
	var req models.CreateClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}

	var expiredDate *time.Time
	if req.ExpiredDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpiredDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid expiry date format"})
			return
		}
		expiredDate = &t
	}

	_, err := h.wg.CreateClient(req.Name, expiredDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// DeleteClient deletes a client
func (h *Handler) DeleteClient(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	if err := h.wg.DeleteClient(clientID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// EnableClient enables a client
func (h *Handler) EnableClient(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	if err := h.wg.EnableClient(clientID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// DisableClient disables a client
func (h *Handler) DisableClient(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	if err := h.wg.DisableClient(clientID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// UpdateClientName updates client name
func (h *Handler) UpdateClientName(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	var req models.UpdateClientNameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}

	if err := h.wg.UpdateClientName(clientID, req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// UpdateClientAddress updates client IP address
func (h *Handler) UpdateClientAddress(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	var req models.UpdateClientAddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}

	if err := h.wg.UpdateClientAddress(clientID, req.Address); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// UpdateClientExpireDate updates client expiry date
func (h *Handler) UpdateClientExpireDate(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	var req models.UpdateClientExpireDateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}

	var expireDate *time.Time
	if req.ExpireDate != "" {
		t, err := time.Parse("2006-01-02", req.ExpireDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid date format"})
			return
		}
		expireDate = &t
	}

	if err := h.wg.UpdateClientExpireDate(clientID, expireDate); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// GenerateOneTimeLink generates a one-time download link
func (h *Handler) GenerateOneTimeLink(c *gin.Context) {
	if !h.cfg.EnableOneTimeLinks {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Feature disabled"})
		return
	}

	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	if err := h.wg.GenerateOneTimeLink(clientID); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// GetClientQRCode returns QR code SVG for client configuration
func (h *Handler) GetClientQRCode(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	config, err := h.wg.GetClientConfiguration(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Generate SVG QR code using inline SVG generation
	svg := generateQRCodeSVG(config)
	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svg)
}

// GetClientConfiguration returns client WireGuard configuration file
func (h *Handler) GetClientConfiguration(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	client, err := h.wg.GetClient(clientID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	config, err := h.wg.GetClientConfiguration(clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Sanitize filename
	configName := sanitizeFilename(client.Name)
	if configName == "" {
		configName = clientID
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.conf"`, configName))
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, config)
}

// DownloadOneTimeLink handles one-time link downloads
func (h *Handler) DownloadOneTimeLink(c *gin.Context) {
	if !h.cfg.EnableOneTimeLinks {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Feature disabled"})
		return
	}

	oneTimeLink := c.Param("clientOneTimeLink")

	client, err := h.wg.GetClientByOneTimeLink(oneTimeLink)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Invalid link"})
		return
	}

	config, err := h.wg.GetClientConfiguration(client.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Mark link as used
	_ = h.wg.EraseOneTimeLink(client.ID)

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.conf"`, oneTimeLink))
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, config)
}

// BackupConfiguration returns configuration backup
func (h *Handler) BackupConfiguration(c *gin.Context) {
	backup, err := h.wg.BackupConfiguration()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="wg0.json"`)
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, backup)
}

// RestoreConfiguration restores configuration from backup
func (h *Handler) RestoreConfiguration(c *gin.Context) {
	var req models.RestoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}

	if err := h.wg.RestoreConfiguration(req.File); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// GetMetrics returns Prometheus metrics
func (h *Handler) GetMetrics(c *gin.Context) {
	if !h.cfg.EnablePrometheusMetrics {
		c.String(http.StatusOK, "")
		return
	}

	metrics, err := h.wg.GetMetrics()
	if err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusOK, metrics)
}

// GetMetricsJSON returns metrics in JSON format
func (h *Handler) GetMetricsJSON(c *gin.Context) {
	if !h.cfg.EnablePrometheusMetrics {
		c.JSON(http.StatusOK, gin.H{})
		return
	}

	metrics, err := h.wg.GetMetricsJSON()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// Helper functions

func isValidClientID(id string) bool {
	// Protect against prototype pollution
	if id == "__proto__" || id == "constructor" || id == "prototype" {
		return false
	}
	return id != ""
}

func sanitizeFilename(name string) string {
	// Replace invalid characters
	re := regexp.MustCompile(`[^a-zA-Z0-9_=+.-]`)
	name = re.ReplaceAllString(name, "-")

	// Replace multiple dashes
	re = regexp.MustCompile(`(-{2,}|-$)`)
	name = re.ReplaceAllString(name, "-")

	// Remove trailing dash
	name = strings.TrimSuffix(name, "-")

	// Limit length
	if len(name) > 32 {
		name = name[:32]
	}

	return name
}

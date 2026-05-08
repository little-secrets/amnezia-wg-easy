package api

import (
	"fmt"
	"log"
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

	// Audit H-5 / C-1 sibling: reject any caller-supplied ID that is
	// not a valid UUID-v4. The previous code accepted arbitrary strings
	// (e.g. one containing "\nPostUp=...") which then reached wg0.conf
	// comment lines and were executed by wg-quick.
	if req.ID != nil && *req.ID != "" && !isValidClientID(*req.ID) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}
	if rejectInjection(c,
		struct{ label, value string }{"name", req.Name},
		struct{ label, value string }{"allowed_ips", derefStr(req.AllowedIPs)},
		struct{ label, value string }{"dns", derefStr(req.DNS)},
		struct{ label, value string }{"mtu", derefStr(req.MTU)},
		struct{ label, value string }{"persistent_keepalive", derefStr(req.PersistentKeepalive)},
		struct{ label, value string }{"address", derefStr(req.Address)},
		struct{ label, value string }{"address6", derefStr(req.Address6)},
	) {
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

	// Build create params with optional parameters
	params := &models.CreateClientParams{
		ID:          req.ID,
		Name:        req.Name,
		ExpiredDate: expiredDate,

		// Network configuration
		Address:    req.Address,
		Address6:   req.Address6,
		AllowedIPs: req.AllowedIPs,

		// Keys
		PrivateKey:   req.PrivateKey,
		PublicKey:    req.PublicKey,
		PreSharedKey: req.PreSharedKey,

		// WireGuard parameters
		DNS:                 req.DNS,
		MTU:                 req.MTU,
		PersistentKeepalive: req.PersistentKeepalive,

		// AmneziaWG obfuscation parameters
		Jc:   req.Jc,
		Jmin: req.Jmin,
		Jmax: req.Jmax,
		S1:   req.S1,
		S2:   req.S2,
		H1:   req.H1,
		H2:   req.H2,
		H3:   req.H3,
		H4:   req.H4,
	}

	_, err := h.wg.CreateClient(params)
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
	if rejectInjection(c, struct{ label, value string }{"name", req.Name}) {
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
	if rejectInjection(c, struct{ label, value string }{"address", req.Address}) {
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

// UpdateClientAllowedIPs updates client allowed IPs
func (h *Handler) UpdateClientAllowedIPs(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	var req models.UpdateClientAllowedIPsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}
	if rejectInjection(c, struct{ label, value string }{"allowed_ips", req.AllowedIPs}) {
		return
	}

	if err := h.wg.UpdateClientAllowedIPs(clientID, req.AllowedIPs); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// UpdateClientDNS updates client DNS servers
func (h *Handler) UpdateClientDNS(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	var req models.UpdateClientDNSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}
	if rejectInjection(c, struct{ label, value string }{"dns", req.DNS}) {
		return
	}

	if err := h.wg.UpdateClientDNS(clientID, req.DNS); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// UpdateClientMTU updates client MTU
func (h *Handler) UpdateClientMTU(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	var req models.UpdateClientMTURequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}
	if rejectInjection(c, struct{ label, value string }{"mtu", req.MTU}) {
		return
	}

	if err := h.wg.UpdateClientMTU(clientID, req.MTU); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// UpdateClientKeepalive updates client persistent keepalive
func (h *Handler) UpdateClientKeepalive(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	var req models.UpdateClientKeepaliveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}
	if rejectInjection(c, struct{ label, value string }{"persistent_keepalive", req.PersistentKeepalive}) {
		return
	}

	if err := h.wg.UpdateClientKeepalive(clientID, req.PersistentKeepalive); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{Success: true})
}

// UpdateClientAddress6 updates client IPv6 address
func (h *Handler) UpdateClientAddress6(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	var req models.UpdateClientAddress6Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request"})
		return
	}
	if rejectInjection(c, struct{ label, value string }{"address6", req.Address6}) {
		return
	}

	if err := h.wg.UpdateClientAddress6(clientID, req.Address6); err != nil {
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

// GetClientSecrets returns client keys (privateKey, publicKey, preSharedKey)
func (h *Handler) GetClientSecrets(c *gin.Context) {
	clientID := c.Param("clientId")
	if !isValidClientID(clientID) {
		c.JSON(http.StatusForbidden, models.ErrorResponse{Error: "Invalid client ID"})
		return
	}

	secrets, err := h.wg.GetClientSecrets(clientID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, secrets)
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

// BackupConfiguration returns configuration backup.
//
// Audit H-4: the response is the raw wg0.json -- every peer's private
// key and PSK in cleartext, plus the server private key. A single
// authenticated GET extracts every secret managed by the service.
// Require an explicit `?confirm=true` query parameter so the endpoint
// cannot be hit accidentally (e.g. by a careless tab in a browser
// extension that prefetches links). Log every successful backup so
// operators have an audit trail.
//
// The deeper fix -- encrypting the backup with a passphrase before
// returning it, or requiring step-up auth -- is a follow-up.
func (h *Handler) BackupConfiguration(c *gin.Context) {
	if c.Query("confirm") != "true" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Backup contains every peer's private key and PSK. " +
				"Re-issue the request with ?confirm=true to acknowledge.",
		})
		return
	}

	backup, err := h.wg.BackupConfiguration()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	log.Printf("[backup] configuration exported from %s", c.ClientIP())

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

// uuidV4Pattern is the strict UUID-v4 form used by every clientID we
// generate (uuid.New()). Audit H-5: the previous validator only blocked
// three JavaScript prototype-pollution literals -- a Go-irrelevant
// concern that gave a false sense of security while letting through
// arbitrary strings, including ones containing newline characters that
// flowed into wg0.conf comment lines and turned into wg-quick PostUp
// injection opportunities.
var uuidV4Pattern = regexp.MustCompile(
	`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
)

func isValidClientID(id string) bool {
	return uuidV4Pattern.MatchString(id)
}

// containsConfigInjection rejects any byte that, when written into
// wg0.conf, would let a peer field break out of its line and either
// declare a new section ([Peer], [Interface]) or inject a wg-quick
// hook (PostUp = curl ...). Audit C-1 sibling.
//
// The legitimate values for the affected fields (name, AllowedIPs,
// DNS, MTU, keepalive) never contain newlines, NUL, or shell control
// characters in any deployment we ship, so a strict allow-policy is
// safe; a future addition that needs richer characters can scope a
// dedicated validator to that field.
func containsConfigInjection(s string) bool {
	for _, r := range s {
		switch r {
		case '\n', '\r', '\x00':
			return true
		}
		// Reject any other ASCII control char (< 0x20 except tab).
		if r < 0x20 && r != '\t' {
			return true
		}
		// 0x7F is DEL.
		if r == 0x7F {
			return true
		}
	}
	return false
}

// derefStr returns the string the pointer references, or "" if nil.
// Used to feed optional request fields into rejectInjection without
// peppering the call sites with nil checks.
func derefStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// rejectInjection writes a 400 response and returns true if any of the
// supplied (label, value) pairs contains a forbidden byte. Use it at
// the top of every handler that forwards caller-controlled strings into
// wg0.conf-bound fields.
func rejectInjection(c *gin.Context, pairs ...struct{ label, value string }) bool {
	for _, p := range pairs {
		if containsConfigInjection(p.value) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "Invalid characters in " + p.label,
			})
			return true
		}
	}
	return false
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

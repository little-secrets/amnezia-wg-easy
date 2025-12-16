package wireguard

import (
	"encoding/json"
	"fmt"
	"hash/crc32"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/little-secrets/amnezia-wg-easy/internal/config"
	"github.com/little-secrets/amnezia-wg-easy/internal/models"
)

// WireGuard manages WireGuard configuration and operations
type WireGuard struct {
	cfg    *config.Config
	config *models.WGConfig
	mu     sync.RWMutex
}

// New creates a new WireGuard instance
func New(cfg *config.Config) *WireGuard {
	return &WireGuard{
		cfg: cfg,
	}
}

// Init initializes WireGuard configuration
func (wg *WireGuard) Init() error {
	if wg.cfg.WGHost == "" {
		return fmt.Errorf("WG_HOST environment variable not set")
	}

	log.Println("Loading WireGuard configuration...")

	wgConfig, err := wg.loadConfig()
	if err != nil {
		log.Println("No existing configuration found, generating new one...")
		wgConfig, err = wg.generateNewConfig()
		if err != nil {
			return fmt.Errorf("failed to generate config: %w", err)
		}
	}

	wg.config = wgConfig

	// Save and sync configuration
	if err := wg.saveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Bring up WireGuard interface
	_ = wg.exec("wg-quick down wg0") // Ignore error if not running
	if err := wg.exec("wg-quick up wg0"); err != nil {
		if strings.Contains(err.Error(), "Cannot find device") {
			return fmt.Errorf("WireGuard exited with error: Cannot find device \"wg0\"\nThis usually means that your host's kernel does not support WireGuard")
		}
		return err
	}

	if err := wg.syncConfig(); err != nil {
		return fmt.Errorf("failed to sync config: %w", err)
	}

	log.Println("WireGuard initialized successfully")
	return nil
}

// Shutdown stops the WireGuard interface
func (wg *WireGuard) Shutdown() error {
	log.Println("Shutting down WireGuard...")
	return wg.exec("wg-quick down wg0")
}

// GetClients returns all clients with their runtime statistics
func (wg *WireGuard) GetClients() ([]*models.ClientResponse, error) {
	wg.mu.RLock()
	defer wg.mu.RUnlock()

	clients := make([]*models.ClientResponse, 0, len(wg.config.Clients))
	for _, client := range wg.config.Clients {
		clients = append(clients, client.ToResponse())
	}

	// Get runtime statistics from WireGuard
	dump, err := wg.getDump()
	if err != nil {
		log.Printf("Warning: failed to get wg dump: %v", err)
		return clients, nil
	}

	// Merge runtime stats
	for _, client := range clients {
		if entry, ok := dump[client.PublicKey]; ok {
			if entry.LatestHandshakeAt > 0 {
				t := time.Unix(entry.LatestHandshakeAt, 0)
				client.LatestHandshakeAt = &t
			}
			if entry.Endpoint != "(none)" {
				client.Endpoint = &entry.Endpoint
			}
			client.TransferRx = &entry.TransferRx
			client.TransferTx = &entry.TransferTx
			client.PersistentKeepalive = &entry.PersistentKeepalive
		}
	}

	return clients, nil
}

// GetClient returns a single client by ID
func (wg *WireGuard) GetClient(clientID string) (*models.Client, error) {
	wg.mu.RLock()
	defer wg.mu.RUnlock()

	client, ok := wg.config.Clients[clientID]
	if !ok {
		return nil, fmt.Errorf("client not found: %s", clientID)
	}
	return client, nil
}

// CreateClient creates a new WireGuard client
func (wg *WireGuard) CreateClient(params *models.CreateClientParams) (*models.Client, error) {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	if params.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Generate keys
	privateKey, err := wg.execOutput("wg genkey")
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	publicKey, err := wg.execPipe(privateKey, "wg pubkey")
	if err != nil {
		return nil, fmt.Errorf("failed to generate public key: %w", err)
	}

	preSharedKey, err := wg.execOutput("wg genpsk")
	if err != nil {
		return nil, fmt.Errorf("failed to generate preshared key: %w", err)
	}

	// Find next available IP address
	address, err := wg.getNextAvailableIP()
	if err != nil {
		return nil, err
	}

	// Create client
	now := time.Now()
	client := &models.Client{
		ID:           uuid.New().String(),
		Name:         params.Name,
		Address:      address,
		PrivateKey:   strings.TrimSpace(privateKey),
		PublicKey:    strings.TrimSpace(publicKey),
		PreSharedKey: strings.TrimSpace(preSharedKey),
		Enabled:      true,
		CreatedAt:    now,
		UpdatedAt:    now,

		// Set custom AmneziaWG parameters if provided (nil = use server defaults)
		Jc:   params.Jc,
		Jmin: params.Jmin,
		Jmax: params.Jmax,
		S1:   params.S1,
		S2:   params.S2,
		H1:   params.H1,
		H2:   params.H2,
		H3:   params.H3,
		H4:   params.H4,
	}

	if params.ExpiredDate != nil {
		// Set expiry to end of day
		exp := time.Date(params.ExpiredDate.Year(), params.ExpiredDate.Month(), params.ExpiredDate.Day(), 23, 59, 59, 0, params.ExpiredDate.Location())
		client.ExpiredAt = &exp
	}

	wg.config.Clients[client.ID] = client

	if err := wg.saveAndSync(); err != nil {
		delete(wg.config.Clients, client.ID)
		return nil, err
	}

	return client, nil
}

// DeleteClient removes a client
func (wg *WireGuard) DeleteClient(clientID string) error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	if _, ok := wg.config.Clients[clientID]; !ok {
		return nil // Already deleted
	}

	delete(wg.config.Clients, clientID)
	return wg.saveAndSync()
}

// EnableClient enables a client
func (wg *WireGuard) EnableClient(clientID string) error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	client, ok := wg.config.Clients[clientID]
	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}

	client.Enabled = true
	client.UpdatedAt = time.Now()

	return wg.saveAndSync()
}

// DisableClient disables a client
func (wg *WireGuard) DisableClient(clientID string) error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	client, ok := wg.config.Clients[clientID]
	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}

	client.Enabled = false
	client.UpdatedAt = time.Now()

	return wg.saveAndSync()
}

// UpdateClientName updates client name
func (wg *WireGuard) UpdateClientName(clientID, name string) error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	client, ok := wg.config.Clients[clientID]
	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}

	client.Name = name
	client.UpdatedAt = time.Now()

	return wg.saveAndSync()
}

// UpdateClientAddress updates client IP address
func (wg *WireGuard) UpdateClientAddress(clientID, address string) error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	if !isValidIPv4(address) {
		return fmt.Errorf("invalid IPv4 address: %s", address)
	}

	client, ok := wg.config.Clients[clientID]
	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}

	client.Address = address
	client.UpdatedAt = time.Now()

	return wg.saveAndSync()
}

// UpdateClientExpireDate updates client expiration date
func (wg *WireGuard) UpdateClientExpireDate(clientID string, expireDate *time.Time) error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	client, ok := wg.config.Clients[clientID]
	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}

	if expireDate != nil {
		exp := time.Date(expireDate.Year(), expireDate.Month(), expireDate.Day(), 23, 59, 59, 0, expireDate.Location())
		client.ExpiredAt = &exp
	} else {
		client.ExpiredAt = nil
	}
	client.UpdatedAt = time.Now()

	return wg.saveAndSync()
}

// GenerateOneTimeLink generates a one-time download link for client
func (wg *WireGuard) GenerateOneTimeLink(clientID string) error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	client, ok := wg.config.Clients[clientID]
	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}

	key := fmt.Sprintf("%s-%d", clientID, time.Now().UnixNano()%1000)
	hash := crc32.ChecksumIEEE([]byte(key))
	link := fmt.Sprintf("%x", hash)
	expiresAt := time.Now().Add(5 * time.Minute)

	client.OneTimeLink = &link
	client.OneTimeLinkExpiresAt = &expiresAt
	client.UpdatedAt = time.Now()

	return wg.saveAndSync()
}

// EraseOneTimeLink marks one-time link for expiration
func (wg *WireGuard) EraseOneTimeLink(clientID string) error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	client, ok := wg.config.Clients[clientID]
	if !ok {
		return fmt.Errorf("client not found: %s", clientID)
	}

	expiresAt := time.Now().Add(10 * time.Second)
	client.OneTimeLinkExpiresAt = &expiresAt
	client.UpdatedAt = time.Now()

	return wg.saveAndSync()
}

// GetClientByOneTimeLink finds a client by their one-time link
func (wg *WireGuard) GetClientByOneTimeLink(link string) (*models.Client, error) {
	wg.mu.RLock()
	defer wg.mu.RUnlock()

	for _, client := range wg.config.Clients {
		if client.OneTimeLink != nil && *client.OneTimeLink == link {
			return client, nil
		}
	}
	return nil, fmt.Errorf("client not found for one-time link")
}

// GetClientConfiguration returns the WireGuard configuration for a client
func (wg *WireGuard) GetClientConfiguration(clientID string) (string, error) {
	wg.mu.RLock()
	defer wg.mu.RUnlock()

	client, ok := wg.config.Clients[clientID]
	if !ok {
		return "", fmt.Errorf("client not found: %s", clientID)
	}

	privateKey := client.PrivateKey
	if privateKey == "" {
		privateKey = "REPLACE_ME"
	}

	// Use client-specific parameters if set, otherwise use server defaults
	jc := wg.getParamOrDefault(client.Jc, wg.config.Server.Jc)
	jmin := wg.getParamOrDefault(client.Jmin, wg.config.Server.Jmin)
	jmax := wg.getParamOrDefault(client.Jmax, wg.config.Server.Jmax)
	s1 := wg.getParamOrDefault(client.S1, wg.config.Server.S1)
	s2 := wg.getParamOrDefault(client.S2, wg.config.Server.S2)
	h1 := wg.getParamOrDefault(client.H1, wg.config.Server.H1)
	h2 := wg.getParamOrDefault(client.H2, wg.config.Server.H2)
	h3 := wg.getParamOrDefault(client.H3, wg.config.Server.H3)
	h4 := wg.getParamOrDefault(client.H4, wg.config.Server.H4)

	var sb strings.Builder
	sb.WriteString("[Interface]\n")
	sb.WriteString(fmt.Sprintf("PrivateKey = %s\n", privateKey))
	sb.WriteString(fmt.Sprintf("Address = %s/24\n", client.Address))
	if wg.cfg.WGDefaultDNS != "" {
		sb.WriteString(fmt.Sprintf("DNS = %s\n", wg.cfg.WGDefaultDNS))
	}
	if wg.cfg.WGMTU != "" {
		sb.WriteString(fmt.Sprintf("MTU = %s\n", wg.cfg.WGMTU))
	}
	sb.WriteString(fmt.Sprintf("Jc = %s\n", jc))
	sb.WriteString(fmt.Sprintf("Jmin = %s\n", jmin))
	sb.WriteString(fmt.Sprintf("Jmax = %s\n", jmax))
	sb.WriteString(fmt.Sprintf("S1 = %s\n", s1))
	sb.WriteString(fmt.Sprintf("S2 = %s\n", s2))
	sb.WriteString(fmt.Sprintf("H1 = %s\n", h1))
	sb.WriteString(fmt.Sprintf("H2 = %s\n", h2))
	sb.WriteString(fmt.Sprintf("H3 = %s\n", h3))
	sb.WriteString(fmt.Sprintf("H4 = %s\n", h4))
	sb.WriteString("\n[Peer]\n")
	sb.WriteString(fmt.Sprintf("PublicKey = %s\n", wg.config.Server.PublicKey))
	if client.PreSharedKey != "" {
		sb.WriteString(fmt.Sprintf("PresharedKey = %s\n", client.PreSharedKey))
	}
	sb.WriteString(fmt.Sprintf("AllowedIPs = %s\n", wg.cfg.WGAllowedIPs))
	sb.WriteString(fmt.Sprintf("PersistentKeepalive = %s\n", wg.cfg.WGPersistentKeepalive))
	sb.WriteString(fmt.Sprintf("Endpoint = %s:%s", wg.cfg.WGHost, wg.cfg.WGConfigPort))

	return sb.String(), nil
}

// getParamOrDefault returns client parameter if set, otherwise server default
func (wg *WireGuard) getParamOrDefault(clientParam *string, serverDefault string) string {
	if clientParam != nil && *clientParam != "" {
		return *clientParam
	}
	return serverDefault
}

// GetClientQRCode returns QR code SVG for client configuration
func (wg *WireGuard) GetClientQRCode(clientID string) ([]byte, error) {
	config, err := wg.GetClientConfiguration(clientID)
	if err != nil {
		return nil, err
	}

	qr, err := qrcode.New(config, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return []byte(qr.ToSmallString(false)), nil
}

// BackupConfiguration returns the current configuration as JSON
func (wg *WireGuard) BackupConfiguration() (string, error) {
	wg.mu.RLock()
	defer wg.mu.RUnlock()

	data, err := json.MarshalIndent(wg.config, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// RestoreConfiguration restores configuration from JSON
func (wg *WireGuard) RestoreConfiguration(data string) error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	var newConfig models.WGConfig
	if err := json.Unmarshal([]byte(data), &newConfig); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	wg.config = &newConfig

	if err := wg.saveConfig(); err != nil {
		return err
	}

	return wg.syncConfig()
}

// GetMetrics returns Prometheus-formatted metrics
func (wg *WireGuard) GetMetrics() (string, error) {
	clients, err := wg.GetClients()
	if err != nil {
		return "", err
	}

	var peerCount, enabledCount, connectedCount int
	var sentBytes, receivedBytes, latestHandshake strings.Builder

	for _, client := range clients {
		peerCount++
		if client.Enabled {
			enabledCount++
		}
		if client.Endpoint != nil {
			connectedCount++
		}

		var tx, rx uint64
		if client.TransferTx != nil {
			tx = *client.TransferTx
		}
		if client.TransferRx != nil {
			rx = *client.TransferRx
		}

		sentBytes.WriteString(fmt.Sprintf(`wireguard_sent_bytes{interface="wg0",enabled="%t",address="%s",name="%s"} %d`+"\n",
			client.Enabled, client.Address, client.Name, tx))
		receivedBytes.WriteString(fmt.Sprintf(`wireguard_received_bytes{interface="wg0",enabled="%t",address="%s",name="%s"} %d`+"\n",
			client.Enabled, client.Address, client.Name, rx))

		var handshakeSeconds float64
		if client.LatestHandshakeAt != nil {
			handshakeSeconds = time.Since(*client.LatestHandshakeAt).Seconds()
		}
		latestHandshake.WriteString(fmt.Sprintf(`wireguard_latest_handshake_seconds{interface="wg0",enabled="%t",address="%s",name="%s"} %.0f`+"\n",
			client.Enabled, client.Address, client.Name, handshakeSeconds))
	}

	var result strings.Builder
	result.WriteString("# HELP wg-easy and wireguard metrics\n")
	result.WriteString("\n# HELP wireguard_configured_peers\n")
	result.WriteString("# TYPE wireguard_configured_peers gauge\n")
	result.WriteString(fmt.Sprintf(`wireguard_configured_peers{interface="wg0"} %d`+"\n", peerCount))
	result.WriteString("\n# HELP wireguard_enabled_peers\n")
	result.WriteString("# TYPE wireguard_enabled_peers gauge\n")
	result.WriteString(fmt.Sprintf(`wireguard_enabled_peers{interface="wg0"} %d`+"\n", enabledCount))
	result.WriteString("\n# HELP wireguard_connected_peers\n")
	result.WriteString("# TYPE wireguard_connected_peers gauge\n")
	result.WriteString(fmt.Sprintf(`wireguard_connected_peers{interface="wg0"} %d`+"\n", connectedCount))
	result.WriteString("\n# HELP wireguard_sent_bytes Bytes sent to the peer\n")
	result.WriteString("# TYPE wireguard_sent_bytes counter\n")
	result.WriteString(sentBytes.String())
	result.WriteString("\n# HELP wireguard_received_bytes Bytes received from the peer\n")
	result.WriteString("# TYPE wireguard_received_bytes counter\n")
	result.WriteString(receivedBytes.String())
	result.WriteString("\n# HELP wireguard_latest_handshake_seconds UNIX timestamp seconds of the last handshake\n")
	result.WriteString("# TYPE wireguard_latest_handshake_seconds gauge\n")
	result.WriteString(latestHandshake.String())

	return result.String(), nil
}

// GetMetricsJSON returns metrics in JSON format
func (wg *WireGuard) GetMetricsJSON() (*models.MetricsJSON, error) {
	clients, err := wg.GetClients()
	if err != nil {
		return nil, err
	}

	var peerCount, enabledCount, connectedCount int
	for _, client := range clients {
		peerCount++
		if client.Enabled {
			enabledCount++
		}
		if client.Endpoint != nil {
			connectedCount++
		}
	}

	return &models.MetricsJSON{
		ConfiguredPeers: peerCount,
		EnabledPeers:    enabledCount,
		ConnectedPeers:  connectedCount,
	}, nil
}

// CronJob runs periodic tasks (expire clients, expire one-time links)
func (wg *WireGuard) CronJob() error {
	wg.mu.Lock()
	defer wg.mu.Unlock()

	needSave := false
	now := time.Now()

	for _, client := range wg.config.Clients {
		// Check client expiration
		if wg.cfg.EnableExpiresTime && client.Enabled && client.ExpiredAt != nil {
			if now.After(*client.ExpiredAt) {
				log.Printf("Client %s (%s) expired", client.Name, client.ID)
				client.Enabled = false
				client.UpdatedAt = now
				needSave = true
			}
		}

		// Check one-time link expiration
		if wg.cfg.EnableOneTimeLinks && client.OneTimeLink != nil && client.OneTimeLinkExpiresAt != nil {
			if now.After(*client.OneTimeLinkExpiresAt) {
				log.Printf("Client %s (%s) one-time link expired", client.Name, client.ID)
				client.OneTimeLink = nil
				client.OneTimeLinkExpiresAt = nil
				client.UpdatedAt = now
				needSave = true
			}
		}
	}

	if needSave {
		return wg.saveAndSync()
	}
	return nil
}

// Helper methods

func (wg *WireGuard) loadConfig() (*models.WGConfig, error) {
	path := filepath.Join(wg.cfg.WGPath, "wg0.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config models.WGConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Clients == nil {
		config.Clients = make(map[string]*models.Client)
	}

	return &config, nil
}

func (wg *WireGuard) generateNewConfig() (*models.WGConfig, error) {
	privateKey, err := wg.execOutput("wg genkey")
	if err != nil {
		return nil, err
	}

	publicKey, err := wg.execPipe(privateKey, "wg pubkey")
	if err != nil {
		return nil, err
	}

	address := strings.Replace(wg.cfg.WGDefaultAddress, "x", "1", 1)

	return &models.WGConfig{
		Server: models.ServerConfig{
			PrivateKey: strings.TrimSpace(privateKey),
			PublicKey:  strings.TrimSpace(publicKey),
			Address:    address,
			Jc:         wg.cfg.Jc,
			Jmin:       wg.cfg.Jmin,
			Jmax:       wg.cfg.Jmax,
			S1:         wg.cfg.S1,
			S2:         wg.cfg.S2,
			H1:         wg.cfg.H1,
			H2:         wg.cfg.H2,
			H3:         wg.cfg.H3,
			H4:         wg.cfg.H4,
		},
		Clients: make(map[string]*models.Client),
	}, nil
}

func (wg *WireGuard) saveAndSync() error {
	if err := wg.saveConfig(); err != nil {
		return err
	}
	return wg.syncConfig()
}

func (wg *WireGuard) saveConfig() error {
	// Save JSON config
	jsonPath := filepath.Join(wg.cfg.WGPath, "wg0.json")
	jsonData, err := json.MarshalIndent(wg.config, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(jsonPath, jsonData, 0660); err != nil {
		return err
	}

	// Generate and save wg0.conf
	confPath := filepath.Join(wg.cfg.WGPath, "wg0.conf")
	confData := wg.generateWGConf()
	return os.WriteFile(confPath, []byte(confData), 0600)
}

func (wg *WireGuard) generateWGConf() string {
	var sb strings.Builder

	sb.WriteString("# Note: Do not edit this file directly.\n")
	sb.WriteString("# Your changes will be overwritten!\n\n")
	sb.WriteString("# Server\n")
	sb.WriteString("[Interface]\n")
	sb.WriteString(fmt.Sprintf("PrivateKey = %s\n", wg.config.Server.PrivateKey))
	sb.WriteString(fmt.Sprintf("Address = %s/24\n", wg.config.Server.Address))
	sb.WriteString(fmt.Sprintf("ListenPort = %s\n", wg.cfg.WGPort))
	if wg.cfg.WGPreUp != "" {
		sb.WriteString(fmt.Sprintf("PreUp = %s\n", wg.cfg.WGPreUp))
	}
	if wg.cfg.WGPostUp != "" {
		sb.WriteString(fmt.Sprintf("PostUp = %s\n", wg.cfg.WGPostUp))
	}
	if wg.cfg.WGPreDown != "" {
		sb.WriteString(fmt.Sprintf("PreDown = %s\n", wg.cfg.WGPreDown))
	}
	if wg.cfg.WGPostDown != "" {
		sb.WriteString(fmt.Sprintf("PostDown = %s\n", wg.cfg.WGPostDown))
	}
	sb.WriteString(fmt.Sprintf("Jc = %s\n", wg.config.Server.Jc))
	sb.WriteString(fmt.Sprintf("Jmin = %s\n", wg.config.Server.Jmin))
	sb.WriteString(fmt.Sprintf("Jmax = %s\n", wg.config.Server.Jmax))
	sb.WriteString(fmt.Sprintf("S1 = %s\n", wg.config.Server.S1))
	sb.WriteString(fmt.Sprintf("S2 = %s\n", wg.config.Server.S2))
	sb.WriteString(fmt.Sprintf("H1 = %s\n", wg.config.Server.H1))
	sb.WriteString(fmt.Sprintf("H2 = %s\n", wg.config.Server.H2))
	sb.WriteString(fmt.Sprintf("H3 = %s\n", wg.config.Server.H3))
	sb.WriteString(fmt.Sprintf("H4 = %s\n", wg.config.Server.H4))

	for clientID, client := range wg.config.Clients {
		if !client.Enabled {
			continue
		}

		sb.WriteString(fmt.Sprintf("\n# Client: %s (%s)\n", client.Name, clientID))
		sb.WriteString("[Peer]\n")
		sb.WriteString(fmt.Sprintf("PublicKey = %s\n", client.PublicKey))
		if client.PreSharedKey != "" {
			sb.WriteString(fmt.Sprintf("PresharedKey = %s\n", client.PreSharedKey))
		}
		sb.WriteString(fmt.Sprintf("AllowedIPs = %s/32\n", client.Address))
	}

	return sb.String()
}

func (wg *WireGuard) syncConfig() error {
	log.Println("Syncing WireGuard configuration...")
	return wg.exec("wg syncconf wg0 <(wg-quick strip wg0)")
}

func (wg *WireGuard) getDump() (map[string]*models.WGDumpEntry, error) {
	output, err := wg.execOutput("wg show wg0 dump")
	if err != nil {
		return nil, err
	}

	result := make(map[string]*models.WGDumpEntry)
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Skip first line (server info)
	for i := 1; i < len(lines); i++ {
		parts := strings.Split(lines[i], "\t")
		if len(parts) < 8 {
			continue
		}

		latestHandshake, _ := strconv.ParseInt(parts[4], 10, 64)
		transferRx, _ := strconv.ParseUint(parts[5], 10, 64)
		transferTx, _ := strconv.ParseUint(parts[6], 10, 64)

		entry := &models.WGDumpEntry{
			PublicKey:           parts[0],
			PreSharedKey:        parts[1],
			Endpoint:            parts[2],
			AllowedIPs:          parts[3],
			LatestHandshakeAt:   latestHandshake,
			TransferRx:          transferRx,
			TransferTx:          transferTx,
			PersistentKeepalive: parts[7],
		}
		result[entry.PublicKey] = entry
	}

	return result, nil
}

func (wg *WireGuard) getNextAvailableIP() (string, error) {
	usedIPs := make(map[string]bool)
	for _, client := range wg.config.Clients {
		usedIPs[client.Address] = true
	}

	for i := 2; i < 255; i++ {
		ip := strings.Replace(wg.cfg.WGDefaultAddress, "x", strconv.Itoa(i), 1)
		if !usedIPs[ip] {
			return ip, nil
		}
	}

	return "", fmt.Errorf("maximum number of clients reached")
}

func isValidIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}

// exec runs a command without capturing output
func (wg *WireGuard) exec(cmd string) error {
	if runtime.GOOS != "linux" {
		log.Printf("[DRY RUN] $ %s", cmd)
		return nil
	}
	log.Printf("$ %s", cmd)
	_, err := wg.execOutput(cmd)
	return err
}

// execOutput runs a command and returns its output
func (wg *WireGuard) execOutput(cmd string) (string, error) {
	if runtime.GOOS != "linux" {
		return "", nil
	}
	return execCommand(cmd)
}

// execPipe runs a command with input piped to stdin
func (wg *WireGuard) execPipe(input, cmd string) (string, error) {
	if runtime.GOOS != "linux" {
		return "", nil
	}
	return execCommandWithInput(input, cmd)
}

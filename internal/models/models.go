package models

import "time"

// WGConfig represents the main WireGuard configuration stored in wg0.json
type WGConfig struct {
	Server  ServerConfig       `json:"server"`
	Clients map[string]*Client `json:"clients"`
}

// ServerConfig represents WireGuard server configuration
type ServerConfig struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Address    string `json:"address"`
	Jc         string `json:"jc"`
	Jmin       string `json:"jmin"`
	Jmax       string `json:"jmax"`
	S1         string `json:"s1"`
	S2         string `json:"s2"`
	H1         string `json:"h1"`
	H2         string `json:"h2"`
	H3         string `json:"h3"`
	H4         string `json:"h4"`
}

// Client represents a WireGuard client/peer
type Client struct {
	ID                   string     `json:"id"`
	Name                 string     `json:"name"`
	Address              string     `json:"address"`
	Address6             string     `json:"address6,omitempty"`
	PublicKey            string     `json:"publicKey"`
	PrivateKey           string     `json:"privateKey"`
	PreSharedKey         string     `json:"preSharedKey"`
	Enabled              bool       `json:"enabled"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	ExpiredAt            *time.Time `json:"expiredAt"`
	AllowedIPs           string     `json:"allowedIPs,omitempty"`
	OneTimeLink          *string    `json:"oneTimeLink,omitempty"`
	OneTimeLinkExpiresAt *time.Time `json:"oneTimeLinkExpiresAt,omitempty"`

	// WireGuard custom parameters (optional, uses server defaults if not set)
	DNS                 *string `json:"dns,omitempty"`
	MTU                 *string `json:"mtu,omitempty"`
	PersistentKeepalive *string `json:"persistentKeepalive,omitempty"`

	// AmneziaWG obfuscation parameters (optional, uses server defaults if not set)
	Jc   *string `json:"jc,omitempty"`
	Jmin *string `json:"jmin,omitempty"`
	Jmax *string `json:"jmax,omitempty"`
	S1   *string `json:"s1,omitempty"`
	S2   *string `json:"s2,omitempty"`
	H1   *string `json:"h1,omitempty"`
	H2   *string `json:"h2,omitempty"`
	H3   *string `json:"h3,omitempty"`
	H4   *string `json:"h4,omitempty"`
}

// ClientResponse is the API response for client queries (includes runtime stats)
type ClientResponse struct {
	ID                   string     `json:"id"`
	Name                 string     `json:"name"`
	Enabled              bool       `json:"enabled"`
	Address              string     `json:"address"`
	Address6             string     `json:"address6,omitempty"`
	PublicKey            string     `json:"publicKey"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            time.Time  `json:"updatedAt"`
	ExpiredAt            *time.Time `json:"expiredAt,omitempty"`
	AllowedIPs           string     `json:"allowedIPs,omitempty"`
	DNS                  *string    `json:"dns,omitempty"`
	MTU                  *string    `json:"mtu,omitempty"`
	OneTimeLink          *string    `json:"oneTimeLink,omitempty"`
	OneTimeLinkExpiresAt *time.Time `json:"oneTimeLinkExpiresAt,omitempty"`
	DownloadableConfig   bool       `json:"downloadableConfig"`
	PersistentKeepalive  *string    `json:"persistentKeepalive"`
	LatestHandshakeAt    *time.Time `json:"latestHandshakeAt"`
	TransferRx           *uint64    `json:"transferRx"`
	TransferTx           *uint64    `json:"transferTx"`
	Endpoint             *string    `json:"endpoint"`
}

// ToResponse converts a Client to ClientResponse
func (c *Client) ToResponse() *ClientResponse {
	return &ClientResponse{
		ID:                   c.ID,
		Name:                 c.Name,
		Enabled:              c.Enabled,
		Address:              c.Address,
		Address6:             c.Address6,
		PublicKey:            c.PublicKey,
		CreatedAt:            c.CreatedAt,
		UpdatedAt:            c.UpdatedAt,
		ExpiredAt:            c.ExpiredAt,
		AllowedIPs:           c.AllowedIPs,
		DNS:                  c.DNS,
		MTU:                  c.MTU,
		OneTimeLink:          c.OneTimeLink,
		OneTimeLinkExpiresAt: c.OneTimeLinkExpiresAt,
		DownloadableConfig:   c.PrivateKey != "",
	}
}

// WGDumpEntry represents parsed output from `wg show wg0 dump`
type WGDumpEntry struct {
	PublicKey           string
	PreSharedKey        string
	Endpoint            string
	AllowedIPs          string
	LatestHandshakeAt   int64
	TransferRx          uint64
	TransferTx          uint64
	PersistentKeepalive string
}

// CreateClientRequest is the request body for creating a client
type CreateClientRequest struct {
	Name        string `json:"name" binding:"required"`
	ExpiredDate string `json:"expiredDate,omitempty"`

	// Network configuration (optional)
	Address    *string `json:"address,omitempty"`    // IPv4 address (auto-generated if not provided)
	Address6   *string `json:"address6,omitempty"`   // IPv6 address (optional)
	AllowedIPs *string `json:"allowedIPs,omitempty"` // Override default allowed IPs

	// Keys (optional, auto-generated if not provided)
	PrivateKey   *string `json:"privateKey,omitempty"`   // Client private key (if provided, publicKey will be derived)
	PublicKey    *string `json:"publicKey,omitempty"`    // Client public key (alternative to privateKey for import)
	PreSharedKey *string `json:"preSharedKey,omitempty"` // Pre-shared key for extra security

	// WireGuard parameters (optional, uses server defaults if not set)
	DNS                 *string `json:"dns,omitempty"`                 // DNS servers
	MTU                 *string `json:"mtu,omitempty"`                 // MTU size
	PersistentKeepalive *string `json:"persistentKeepalive,omitempty"` // Keepalive interval

	// AmneziaWG obfuscation parameters (optional)
	Jc   *string `json:"jc,omitempty"`
	Jmin *string `json:"jmin,omitempty"`
	Jmax *string `json:"jmax,omitempty"`
	S1   *string `json:"s1,omitempty"`
	S2   *string `json:"s2,omitempty"`
	H1   *string `json:"h1,omitempty"`
	H2   *string `json:"h2,omitempty"`
	H3   *string `json:"h3,omitempty"`
	H4   *string `json:"h4,omitempty"`
}

// CreateClientParams holds parameters for creating a new client
type CreateClientParams struct {
	Name        string
	ExpiredDate *time.Time

	// Network configuration (optional)
	Address    *string
	Address6   *string
	AllowedIPs *string

	// Keys (optional, auto-generated if not provided)
	PrivateKey   *string
	PublicKey    *string
	PreSharedKey *string

	// WireGuard parameters (optional, nil = use server defaults)
	DNS                 *string
	MTU                 *string
	PersistentKeepalive *string

	// AmneziaWG obfuscation parameters (optional, nil = use server defaults)
	Jc   *string
	Jmin *string
	Jmax *string
	S1   *string
	S2   *string
	H1   *string
	H2   *string
	H3   *string
	H4   *string
}

// UpdateClientNameRequest is the request body for updating client name
type UpdateClientNameRequest struct {
	Name string `json:"name" binding:"required"`
}

// UpdateClientAddressRequest is the request body for updating client address
type UpdateClientAddressRequest struct {
	Address string `json:"address" binding:"required"`
}

// UpdateClientExpireDateRequest is the request body for updating expiry date
type UpdateClientExpireDateRequest struct {
	ExpireDate string `json:"expireDate,omitempty"`
}

// UpdateClientAllowedIPsRequest is the request body for updating allowed IPs
type UpdateClientAllowedIPsRequest struct {
	AllowedIPs string `json:"allowedIPs" binding:"required"`
}

// UpdateClientDNSRequest is the request body for updating DNS
type UpdateClientDNSRequest struct {
	DNS string `json:"dns"`
}

// UpdateClientMTURequest is the request body for updating MTU
type UpdateClientMTURequest struct {
	MTU string `json:"mtu"`
}

// UpdateClientKeepaliveRequest is the request body for updating persistent keepalive
type UpdateClientKeepaliveRequest struct {
	PersistentKeepalive string `json:"persistentKeepalive"`
}

// UpdateClientAddress6Request is the request body for updating IPv6 address
type UpdateClientAddress6Request struct {
	Address6 string `json:"address6"`
}

// LoginRequest is the login request body
type LoginRequest struct {
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember"`
}

// SessionResponse is the session status response
type SessionResponse struct {
	RequiresPassword bool `json:"requiresPassword"`
	Authenticated    bool `json:"authenticated"`
}

// SuccessResponse is a generic success response
type SuccessResponse struct {
	Success bool `json:"success"`
}

// ErrorResponse is a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// RestoreRequest is the request body for restoring configuration
type RestoreRequest struct {
	File string `json:"file" binding:"required"`
}

// AvatarSettings represents UI avatar configuration
type AvatarSettings struct {
	Dicebear string `json:"dicebear"`
	Gravatar bool   `json:"gravatar"`
}

// MetricsJSON represents JSON format of prometheus metrics
type MetricsJSON struct {
	ConfiguredPeers int `json:"wireguard_configured_peers"`
	EnabledPeers    int `json:"wireguard_enabled_peers"`
	ConnectedPeers  int `json:"wireguard_connected_peers"`
}

// ClientSecretsResponse returns sensitive client data
type ClientSecretsResponse struct {
	ID           string `json:"id"`
	PrivateKey   string `json:"privateKey"`
	PublicKey    string `json:"publicKey"`
	PreSharedKey string `json:"preSharedKey"`
}

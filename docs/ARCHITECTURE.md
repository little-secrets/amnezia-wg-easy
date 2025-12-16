# Architecture Guide

This document describes the architecture and structure of AmneziaWG Easy (Go Edition).

## 📁 Project Structure

```
amnezia-wg-easy/
├── main.go                     # Application entry point
├── go.mod                      # Go module dependencies
├── go.sum                      # Dependency checksums
│
├── cmd/                        # Command-line tools
│   └── wgpw/                  # Password generation utility
│       └── main.go            # wgpw implementation
│
├── internal/                   # Private application code
│   ├── config/                # Configuration management
│   │   └── config.go          # ENV loading & defaults
│   │
│   ├── models/                # Data structures
│   │   └── models.go          # All structs & types
│   │
│   ├── wireguard/             # WireGuard operations
│   │   ├── wireguard.go       # Core WG logic
│   │   └── exec.go            # Shell command execution
│   │
│   └── api/                   # HTTP layer
│       ├── routes.go          # Router setup & middleware
│       ├── handlers.go        # HTTP request handlers
│       ├── middleware.go      # Authentication & auth
│       └── qrcode.go          # QR code generation
│
├── www/                        # Web UI (Vue.js)
│   ├── index.html             # SPA entry point
│   ├── js/                    # JavaScript files
│   ├── css/                   # Stylesheets
│   └── img/                   # Images & icons
│
├── docs/                       # Documentation
│   ├── ARCHITECTURE.md        # This file
│   ├── ENVIRONMENT_VARIABLES.md
│   ├── AMNEZIAWG_PARAMETERS.md
│   ├── PASSWORD_GENERATION.md
│   └── API_REFERENCE.md
│
├── Dockerfile                  # Container build
├── docker-compose.yml          # Compose configuration
├── Makefile                    # Build automation
└── README.md                   # Main documentation
```

## 🏗️ Core Components

### 1. Main Entry Point (`main.go`)

**Responsibilities:**
- Load configuration from environment
- Initialize WireGuard service
- Setup HTTP router
- Start cron jobs
- Handle graceful shutdown

**Flow:**
```go
func main() {
    config := config.Load()              // 1. Load ENV
    wg := wireguard.New(config)          // 2. Create WG service
    wg.Init()                            // 3. Initialize WG
    router := api.SetupRouter(cfg, wg)   // 4. Setup HTTP
    server := http.Server{...}           // 5. Create server
    go startCronJob(wg)                  // 6. Start cron
    server.ListenAndServe()              // 7. Start server
    // Wait for signal
    // Graceful shutdown
}
```

### 2. Configuration Layer (`internal/config/`)

**Purpose:** Centralized environment variable management

**Key Features:**
- Default values for all settings
- Type conversion (string → int, bool)
- Random generation for AmneziaWG parameters
- Build iptables rules dynamically

**Example:**
```go
cfg := config.Load()  // Loads all ENV with defaults

fmt.Println(cfg.Port)      // "51821" (default)
fmt.Println(cfg.WGHost)    // from ENV (required)
fmt.Println(cfg.Jc)        // random 3-10
```

### 3. Models Layer (`internal/models/`)

**Purpose:** Data structures for the entire application

**Main Types:**
- `WGConfig` - Server + clients configuration
- `ServerConfig` - WireGuard server settings
- `Client` - Individual client/peer
- `ClientResponse` - API response with stats
- `WGDumpEntry` - Parsed wg dump output
- Request/Response DTOs

**Example:**
```go
type Client struct {
    ID           string    `json:"id"`
    Name         string    `json:"name"`
    PublicKey    string    `json:"publicKey"`
    Enabled      bool      `json:"enabled"`
    CreatedAt    time.Time `json:"createdAt"`
    ExpiredAt    *time.Time `json:"expiredAt"`
}
```

### 4. WireGuard Layer (`internal/wireguard/`)

**Purpose:** All WireGuard operations and business logic

**Key Operations:**
```go
// Initialization
wg.Init()  // Load/generate config, start WG interface

// Client management
wg.CreateClient(name, expiry)
wg.DeleteClient(id)
wg.EnableClient(id)
wg.DisableClient(id)

// Configuration
wg.GetClientConfiguration(id)
wg.GetClientQRCode(id)
wg.BackupConfiguration()
wg.RestoreConfiguration(data)

// Statistics
wg.GetClients()  // With runtime stats
wg.GetMetrics()  // Prometheus format
wg.GetMetricsJSON()

// Maintenance
wg.CronJob()  // Expire clients, clean one-time links
wg.Shutdown()
```

**File Structure:**

#### `wireguard.go`
- Main WireGuard service
- Config load/save/sync
- Client CRUD operations
- Metrics generation

#### `exec.go`
- Shell command execution
- Safe command handling
- Error propagation

**Command Execution:**
```go
// Generate keys
privateKey := wg.execOutput("wg genkey")
publicKey := wg.execPipe(privateKey, "wg pubkey")

// Manage interface
wg.exec("wg-quick up wg0")
wg.exec("wg syncconf wg0 <(wg-quick strip wg0)")

// Get stats
dump := wg.execOutput("wg show wg0 dump")
```

### 5. API Layer (`internal/api/`)

**Purpose:** HTTP server, routing, and authentication

#### `routes.go`
**Router setup with Gin framework:**
```go
func SetupRouter(cfg, wg) *gin.Engine {
    r := gin.New()
    
    // Session middleware
    r.Use(sessions.Sessions("wg-session", store))
    
    // Public routes
    api.GET("/api/session")
    api.POST("/api/session")
    
    // Protected routes (with auth)
    protected.Use(AuthMiddleware)
    protected.GET("/api/wireguard/client")
    protected.POST("/api/wireguard/client")
    
    // Metrics (optional auth)
    metrics.Use(PrometheusAuthMiddleware)
    metrics.GET("/metrics")
    
    // Static files (Web UI)
    if !cfg.NoWebUI {
        r.Static("/js", "./www/js")
        r.Static("/css", "./www/css")
    }
    
    return r
}
```

#### `handlers.go`
**HTTP request handlers:**
- Session management (login/logout)
- Client CRUD
- Configuration download/QR codes
- Backup/restore
- Metrics
- UI settings endpoints

**Example:**
```go
func (h *Handler) CreateClient(c *gin.Context) {
    var req models.CreateClientRequest
    c.ShouldBindJSON(&req)  // Parse JSON
    
    _, err := h.wg.CreateClient(req.Name, expiry)
    if err != nil {
        c.JSON(500, ErrorResponse{err.Error()})
        return
    }
    
    c.JSON(200, SuccessResponse{true})
}
```

#### `middleware.go`
**Authentication:**
- Session-based auth
- Header-based auth (Authorization)
- Basic Auth for Prometheus

```go
func AuthMiddleware(cfg) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check session
        if session.Get("authenticated") == true {
            c.Next()
            return
        }
        
        // Check Authorization header
        if checkPassword(header, cfg.PasswordHash) {
            c.Next()
            return
        }
        
        c.JSON(401, ErrorResponse{"Not Logged In"})
        c.Abort()
    }
}
```

#### `qrcode.go`
**QR code generation:**
- SVG format
- 512x512 size
- Embedded in API response

## 🔄 Data Flow

### 1. Creating a Client

```
User → HTTP POST /api/wireguard/client
  ↓
Handler.CreateClient
  ↓
WireGuard.CreateClient
  ├─→ Generate keys (wg genkey, wg pubkey, wg genpsk)
  ├─→ Find next available IP
  ├─→ Create Client struct
  ├─→ Add to config.Clients map
  └─→ SaveAndSync
      ├─→ Save JSON (wg0.json)
      ├─→ Generate conf (wg0.conf)
      └─→ Sync interface (wg syncconf)
  ↓
Response: {success: true}
```

### 2. Getting Client List with Stats

```
User → HTTP GET /api/wireguard/client
  ↓
Handler.GetClients
  ↓
WireGuard.GetClients
  ├─→ Read config.Clients
  ├─→ Convert to ClientResponse
  └─→ Get runtime stats
      ├─→ Execute: wg show wg0 dump
      ├─→ Parse dump output
      └─→ Merge stats (handshake, transfer, endpoint)
  ↓
Response: [ClientResponse, ...]
```

### 3. Client Configuration Download

```
User → HTTP GET /api/wireguard/client/:id/configuration
  ↓
Handler.GetClientConfiguration
  ↓
WireGuard.GetClientConfiguration
  ├─→ Load client from config
  ├─→ Build WireGuard config text
  │   [Interface]
  │   PrivateKey = ...
  │   Jc = ...
  │   [Peer]
  │   PublicKey = ...
  │   Endpoint = ...
  └─→ Return string
  ↓
Response: text/plain (.conf file)
```

## 🧩 Dependency Management

### External Dependencies

```go
require (
    github.com/gin-gonic/gin              // HTTP framework
    github.com/gin-contrib/sessions       // Session management
    github.com/google/uuid                // UUID generation
    github.com/skip2/go-qrcode           // QR code generation
    golang.org/x/crypto/bcrypt           // Password hashing
    golang.org/x/term                    // Terminal input
)
```

### Standard Library Usage
- `encoding/json` - JSON serialization
- `os/exec` - Shell command execution
- `net/http` - HTTP server
- `time` - Timestamps & intervals
- `sync` - Mutex for concurrent access
- `path/filepath` - File path handling

## 🔒 Concurrency & Safety

### Mutex Protection

```go
type WireGuard struct {
    config *models.WGConfig
    mu     sync.RWMutex  // Protects config access
}

// Read lock for queries
func (wg *WireGuard) GetClients() {
    wg.mu.RLock()
    defer wg.mu.RUnlock()
    // Read config safely
}

// Write lock for modifications
func (wg *WireGuard) CreateClient() {
    wg.mu.Lock()
    defer wg.mu.Unlock()
    // Modify config safely
}
```

### Goroutines

**Cron Job:**
```go
func startCronJob(wg *WireGuard) {
    ticker := time.NewTicker(1 * time.Minute)
    for range ticker.C {
        wg.CronJob()  // Expire clients, clean links
    }
}
```

**HTTP Server:**
- Gin handles each request in a goroutine
- Mutex ensures safe concurrent access to config

## 📝 Configuration Storage

### Files

**wg0.json** (Primary storage):
```json
{
  "server": {
    "privateKey": "...",
    "publicKey": "...",
    "jc": "7",
    "h1": "1234567891",
    ...
  },
  "clients": {
    "uuid-1": {
      "id": "uuid-1",
      "name": "client1",
      "publicKey": "...",
      "enabled": true,
      ...
    }
  }
}
```

**wg0.conf** (Generated for WireGuard):
```ini
[Interface]
PrivateKey = ...
Jc = 7
H1 = 1234567891

[Peer]
PublicKey = ...
AllowedIPs = 10.8.0.2/32
```

### Storage Location
- Default: `/etc/wireguard/`
- Customizable via `WG_PATH` ENV
- Mounted as Docker volume

## 🛠️ Build Process

### Development Build
```bash
go build -o amnezia-wg-easy .
go build -o wgpw ./cmd/wgpw
```

### Production Build (Docker)
```dockerfile
# Stage 1: Build
FROM golang:1.24-alpine AS builder
COPY . .
RUN go build -ldflags="-s -w" -o amnezia-wg-easy .

# Stage 2: Runtime
FROM amneziavpn/amnezia-wg:latest
COPY --from=builder /build/amnezia-wg-easy /app/
COPY www /app/www
CMD ["/app/amnezia-wg-easy"]
```

### Binary Size
- With debug symbols: ~25MB
- Stripped (`-ldflags="-s -w"`): ~15MB
- Compressed in container: ~8MB

## 🔌 Extension Points

### Adding New API Endpoints

**1. Define handler in `handlers.go`:**
```go
func (h *Handler) MyNewEndpoint(c *gin.Context) {
    result := h.wg.SomeOperation()
    c.JSON(200, result)
}
```

**2. Register route in `routes.go`:**
```go
protected.GET("/api/my-endpoint", h.MyNewEndpoint)
```

### Adding New WireGuard Operations

**1. Implement in `wireguard.go`:**
```go
func (wg *WireGuard) NewOperation() error {
    wg.mu.Lock()
    defer wg.mu.Unlock()
    
    // Your logic here
    
    return wg.saveAndSync()
}
```

**2. Add handler to expose via API**

### Adding New Configuration Options

**1. Add to `config.go`:**
```go
type Config struct {
    MyNewOption string
}

func Load() *Config {
    cfg := &Config{
        MyNewOption: getEnv("MY_NEW_OPTION", "default"),
    }
    return cfg
}
```

**2. Use in application:**
```go
if cfg.MyNewOption == "something" {
    // Do something
}
```

## 📊 Performance Characteristics

### Memory Usage
- Base: ~10-15MB
- Per client: ~1KB
- 1000 clients: ~20-25MB total

### CPU Usage
- Idle: <1%
- During client creation: 5-10% (key generation)
- Per request: <1ms

### Disk I/O
- Config save: 1-10KB per operation
- Batch operations: Minimal impact

### Network
- REST API: ~1-5ms response time
- Static files: Served directly by Gin
- WireGuard: Kernel module (no overhead)

## 🧪 Testing Strategy

### Unit Tests
```bash
go test ./...
```

### Integration Tests
```bash
# Start container
docker compose up -d

# Test API
curl http://localhost:51821/api/session

# Test metrics
curl http://localhost:51821/metrics
```

### Load Testing
```bash
# Using Apache Bench
ab -n 1000 -c 10 http://localhost:51821/api/wireguard/client
```

## 🐛 Debugging

### Enable Debug Logging
```go
// In main.go
gin.SetMode(gin.DebugMode)
```

### View WireGuard State
```bash
docker compose exec amnezia-wg-easy wg show
docker compose exec amnezia-wg-easy cat /etc/wireguard/wg0.json
```

### Check Logs
```bash
docker compose logs -f
```

---

**Next Steps:**
- [Environment Variables](./ENVIRONMENT_VARIABLES.md)
- [API Reference](./API_REFERENCE.md)
- [AmneziaWG Parameters](./AMNEZIAWG_PARAMETERS.md)


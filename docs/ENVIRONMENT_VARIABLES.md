# Environment Variables Reference

Complete reference of all environment variables supported by AmneziaWG Easy.

## ⚠️ Required Variables

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `WG_HOST` | string | **REQUIRED** | Public IP address or domain name of your server |

**This is the ONLY required variable.** Everything else has sensible defaults.

---

## 🌐 Server Settings

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `PORT` | string | `51821` | TCP port for Web UI and API |
| `WEBUI_HOST` | string | `0.0.0.0` | IP address to bind the web server to |
| `RELEASE` | string | `1.0.0` | Application version (auto-set) |

### Examples

```bash
# Custom port
PORT=8080

# Bind to localhost only
WEBUI_HOST=127.0.0.1

# Default (listen on all interfaces)
WEBUI_HOST=0.0.0.0
```

---

## 🔐 Authentication

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `PASSWORD_HASH` | string | *(empty)* | Bcrypt hash for Web UI login |
| `MAX_AGE` | int | `0` | Session max age in minutes (0 = browser session) |
| `PROMETHEUS_METRICS_PASSWORD` | string | *(empty)* | Bcrypt hash for Prometheus metrics Basic Auth |

### Important Notes

- **No `PASSWORD_HASH`** = No authentication (⚠️ not recommended for production)
- Generate hashes with: `docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw mypassword`
- See [Password Generation Guide](./PASSWORD_GENERATION.md)

### Examples

```bash
# Enable Web UI password
PASSWORD_HASH='$2a$12$xELb112CO5ZgDqydj4SET.bxuHr2hcMb2SWgTlBU/XKSt8NEGjUge'

# Session expires after 24 hours
MAX_AGE=1440

# Protect metrics with password
ENABLE_PROMETHEUS_METRICS=true
PROMETHEUS_METRICS_PASSWORD='$2a$12$...'
```

---

## 🔧 WireGuard Settings

### Basic Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `WG_PATH` | string | `/etc/wireguard/` | Directory for WireGuard configuration files |
| `WG_DEVICE` | string | `eth0` | Network interface for traffic forwarding |
| `WG_PORT` | string | `51820` | UDP port for WireGuard |
| `WG_CONFIG_PORT` | string | `$WG_PORT` | Port in client configurations (for port forwarding) |

### Network Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `WG_DEFAULT_ADDRESS` | string | `10.8.0.x` | Client IP address range (`x` is replaced with number) |
| `WG_DEFAULT_DNS` | string | `1.1.1.1` | DNS server for clients (empty = no DNS) |
| `WG_ALLOWED_IPS` | string | `0.0.0.0/0, ::/0` | Allowed IPs for clients (routes all traffic) |
| `WG_MTU` | string | *(empty)* | MTU for clients (empty = default) |
| `WG_PERSISTENT_KEEPALIVE` | string | `0` | Persistent keepalive in seconds (0 = disabled) |

### iptables Hooks

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `WG_PRE_UP` | string | *(empty)* | Commands before interface goes up |
| `WG_POST_UP` | string | *(auto-generated)* | Commands after interface goes up |
| `WG_PRE_DOWN` | string | *(empty)* | Commands before interface goes down |
| `WG_POST_DOWN` | string | *(auto-generated)* | Commands after interface goes down |

**Default `WG_POST_UP`:**
```bash
iptables -t nat -A POSTROUTING -s 10.8.0.0/24 -o eth0 -j MASQUERADE;
iptables -A INPUT -p udp -m udp --dport 51820 -j ACCEPT;
iptables -A FORWARD -i wg0 -j ACCEPT;
iptables -A FORWARD -o wg0 -j ACCEPT;
```

**Default `WG_POST_DOWN`:**
```bash
iptables -t nat -D POSTROUTING -s 10.8.0.0/24 -o eth0 -j MASQUERADE;
iptables -D INPUT -p udp -m udp --dport 51820 -j ACCEPT;
iptables -D FORWARD -i wg0 -j ACCEPT;
iptables -D FORWARD -o wg0 -j ACCEPT;
```

### Examples

```bash
# Custom subnet
WG_DEFAULT_ADDRESS=192.168.100.x

# Multiple DNS servers
WG_DEFAULT_DNS=8.8.8.8,1.1.1.1

# Only route specific network
WG_ALLOWED_IPS=10.0.0.0/8,192.168.0.0/16

# Enable keepalive for NAT traversal
WG_PERSISTENT_KEEPALIVE=25

# Lower MTU for problematic networks
WG_MTU=1420

# Custom iptables
WG_POST_UP="iptables -A FORWARD -i wg0 -j ACCEPT"
```

---

## 🎛️ Feature Flags

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `NO_WEB_UI` | bool | `false` | Disable Web UI (API-only mode) |
| `WG_ENABLE_ONE_TIME_LINKS` | bool | `false` | Enable one-time download links (expire after 5 min) |
| `WG_ENABLE_EXPIRES_TIME` | bool | `false` | Enable automatic client expiration |
| `ENABLE_PROMETHEUS_METRICS` | bool | `false` | Enable Prometheus metrics at `/metrics` |
| `UI_TRAFFIC_STATS` | bool | `false` | Enable detailed RX/TX stats in Web UI |
| `UI_ENABLE_SORT_CLIENTS` | bool | `false` | Enable client sorting by name in UI |

### Examples

```bash
# API-only mode (no Web UI)
NO_WEB_UI=true

# Enable all features
WG_ENABLE_ONE_TIME_LINKS=true
WG_ENABLE_EXPIRES_TIME=true
ENABLE_PROMETHEUS_METRICS=true
UI_TRAFFIC_STATS=true
UI_ENABLE_SORT_CLIENTS=true
```

---

## 🎨 UI Settings

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `LANG` | string | `en` | Interface language (en, ru, de, fr, es, etc.) |
| `UI_CHART_TYPE` | int | `0` | Chart type: 0=disabled, 1=line, 2=area, 3=bar |
| `DICEBEAR_TYPE` | string | *(empty)* | Avatar style (bottts, avataaars, etc.) |
| `USE_GRAVATAR` | bool | `false` | Use Gravatar for avatars |

### Supported Languages

`en`, `ua`, `ru`, `tr`, `no`, `pl`, `fr`, `de`, `ca`, `es`, `ko`, `vi`, `nl`, `is`, `pt`, `chs`, `cht`, `it`, `th`, `hi`

### Examples

```bash
# Russian interface with line charts
LANG=ru
UI_CHART_TYPE=1

# Use Gravatar avatars
USE_GRAVATAR=true

# Use Dicebear bottts style
DICEBEAR_TYPE=bottts
```

---

## 🎲 AmneziaWG Obfuscation

**These parameters are automatically randomized if not set!**

See [AmneziaWG Parameters Guide](./AMNEZIAWG_PARAMETERS.md) for detailed explanations.

### Junk Packet Settings

| Variable | Type | Auto-Default | Range | Description |
|----------|------|--------------|-------|-------------|
| `JC` | int | Random 3-10 | 0-128 | Junk packet count |
| `JMIN` | int | `50` | 0-1500 | Minimum junk packet size (bytes) |
| `JMAX` | int | `1000` | JMIN-1500 | Maximum junk packet size (bytes) |

### Packet Junk Size

| Variable | Type | Auto-Default | Range | Description |
|----------|------|--------------|-------|-------------|
| `S1` | int | Random 15-150 | 0-1500 | Init packet junk size (bytes) |
| `S2` | int | Random 15-150 | 0-1500 | Response packet junk size (bytes) |

### Magic Headers

| Variable | Type | Auto-Default | Range | Description |
|----------|------|--------------|-------|-------------|
| `H1` | uint32 | Random | 1-4294967295 | Init packet magic header |
| `H2` | uint32 | Random | 1-4294967295 | Response packet magic header |
| `H3` | uint32 | Random | 1-4294967295 | Underload packet magic header |
| `H4` | uint32 | Random | 1-4294967295 | Transport packet magic header |

### Examples

```bash
# Let all values randomize (recommended)
# Don't set JC, S1, S2, H1-H4

# Or set custom values
JC=7
JMIN=50
JMAX=1000
S1=100
S2=100
H1=1234567891
H2=1234567892
H3=1234567893
H4=1234567894
```

---

## 📋 Configuration Examples

### 1. Minimal (Development)

```bash
WG_HOST=192.168.1.100
```

**Result:**
- No password (open access)
- All defaults applied
- Random AmneziaWG parameters

### 2. Production (Secure)

```bash
WG_HOST=vpn.example.com
PASSWORD_HASH='$2a$12$xELb112CO5ZgDqydj4SET.bxuHr2hcMb2SWgTlBU/XKSt8NEGjUge'
PORT=51821
WG_PORT=51820
LANG=en
WG_DEFAULT_DNS=8.8.8.8,1.1.1.1
WG_PERSISTENT_KEEPALIVE=25
```

### 3. API-Only Mode

```bash
WG_HOST=api.vpn.example.com
PASSWORD_HASH='$2a$12$...'
NO_WEB_UI=true
ENABLE_PROMETHEUS_METRICS=true
```

### 4. Full Features

```bash
# Required
WG_HOST=vpn.example.com
PASSWORD_HASH='$2a$12$...'

# Server
PORT=51821
WG_PORT=51820

# Features
WG_ENABLE_ONE_TIME_LINKS=true
WG_ENABLE_EXPIRES_TIME=true
ENABLE_PROMETHEUS_METRICS=true
PROMETHEUS_METRICS_PASSWORD='$2a$12$...'
UI_TRAFFIC_STATS=true
UI_ENABLE_SORT_CLIENTS=true

# UI
LANG=en
UI_CHART_TYPE=1
USE_GRAVATAR=true

# WireGuard
WG_DEFAULT_DNS=8.8.8.8,1.1.1.1
WG_PERSISTENT_KEEPALIVE=25
WG_MTU=1420

# Session
MAX_AGE=1440

# AmneziaWG (custom)
JC=7
JMIN=50
JMAX=1000
S1=100
S2=100
H1=1234567891
H2=1234567892
H3=1234567893
H4=1234567894
```

### 5. Split Tunnel (Specific Routes)

```bash
WG_HOST=vpn.example.com
PASSWORD_HASH='$2a$12$...'

# Only route internal networks through VPN
WG_ALLOWED_IPS=10.0.0.0/8,192.168.0.0/16

# Use local DNS
WG_DEFAULT_DNS=
```

---

## 🔍 Variable Resolution

### Priority Order

1. Environment variables
2. Default values
3. Auto-generated values

### Example Resolution

```bash
# User sets
WG_PORT=51820

# Application resolves
WG_CONFIG_PORT=${WG_PORT}  # → 51820
WG_DEFAULT_ADDRESS="10.8.0.x"  # → default
Jc=random(3, 10)  # → e.g., 7
```

---

## 🔗 Environment from File

### .env File

```bash
# Create .env file
cat > .env << 'EOF'
WG_HOST=vpn.example.com
PASSWORD_HASH='$2a$12$...'
PORT=51821
EOF

# Use with Docker Compose
docker compose up -d
```

### docker-compose.yml

```yaml
services:
  amnezia-wg-easy:
    environment:
      - WG_HOST=${WG_HOST}
      - PASSWORD_HASH=${PASSWORD_HASH}
    # Or use env_file
    env_file:
      - .env
```

---

## 🐛 Troubleshooting

### Check Current Values

```bash
# View environment in container
docker compose exec amnezia-wg-easy env | grep WG_

# View generated config
docker compose exec amnezia-wg-easy cat /etc/wireguard/wg0.json | jq .server
```

### Common Issues

**Problem:** `WG_HOST` not set
```bash
Error: WG_HOST environment variable is required
```
**Solution:** Set `WG_HOST=your.server.ip`

**Problem:** Password not working
```bash
# Check hash format
echo $PASSWORD_HASH
# Should be: $2a$12$...

# Regenerate
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw yourpassword
```

**Problem:** Metrics not appearing
```bash
# Enable metrics
ENABLE_PROMETHEUS_METRICS=true
```

---

## 📚 Related Documentation

- [AmneziaWG Parameters](./AMNEZIAWG_PARAMETERS.md)
- [Password Generation](./PASSWORD_GENERATION.md)
- [Architecture](./ARCHITECTURE.md)
- [API Reference](./API_REFERENCE.md)

---

**Tip:** Most variables have sensible defaults. Start with just `WG_HOST` and add others as needed!


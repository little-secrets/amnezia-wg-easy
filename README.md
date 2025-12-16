# AmneziaWG Easy

The easiest way to install & manage AmneziaWG with a Web UI, rewritten in Go.

<p align="center">
  <img src="./assets/screenshot.png" width="802" />
</p>

## ✨ Features

- **All-in-one**: AmneziaWG + Web UI + REST API
- **Lightweight**: Single Go binary (~15MB) instead of Node.js
- **Easy setup**: Only one required environment variable (`WG_HOST`)
- **API-only mode**: Disable Web UI with `NO_WEB_UI=true`
- **Built-in password tool**: Generate bcrypt hashes with `wgpw`
- **Full AmneziaWG support**: Traffic obfuscation out of the box
- **Prometheus metrics**: Monitor your VPN with `/metrics` endpoint
- **Modern UI**: Automatic Dark/Light mode, multilingual support

### Why Go Edition?

| Feature | Node.js | Go |
|---------|---------|-----|
| Binary size | ~300MB | ~15MB |
| Dependencies | node_modules | None |
| Memory usage | ~100MB | ~30MB |
| Startup time | ~2s | <100ms |

## 🚀 Quick Start

### Prerequisites

- Docker with NET_ADMIN capability
- Linux host with WireGuard support

### 1. Generate Password Hash

```bash
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw mypassword
# Copy the output: PASSWORD_HASH='$2a$12$...'
```

See [Password Generation Guide](./docs/PASSWORD_GENERATION.md) for more options.

### 2. Create `.env` File

```bash
# Required
WG_HOST=your.server.ip.or.domain
PASSWORD_HASH='$2a$12$...'

# Optional (with defaults)
PORT=51821
WG_PORT=51820
LANG=en
```

See [Environment Variables Reference](./docs/ENVIRONMENT_VARIABLES.md) for all options.

### 3. Start with Docker Compose

```bash
docker compose up -d
```

**Or with docker run:**

```bash
docker run -d \
  --name=amnezia-wg-easy \
  -e WG_HOST=your.server.ip \
  -e PASSWORD_HASH='$2a$12$...' \
  -v ~/.amnezia-wg:/etc/wireguard \
  -p 51820:51820/udp \
  -p 51821:51821/tcp \
  --cap-add=NET_ADMIN \
  --cap-add=SYS_MODULE \
  --sysctl="net.ipv4.ip_forward=1" \
  --sysctl="net.ipv4.conf.all.src_valid_mark=1" \
  --device=/dev/net/tun:/dev/net/tun \
  --restart unless-stopped \
  ghcr.io/w0rng/amnezia-wg-easy:latest
```

### 4. Access Web UI

Open `http://your.server.ip:51821` and login with your password.

## 📖 Documentation

- **[Environment Variables](./docs/ENVIRONMENT_VARIABLES.md)** - Complete list of all ENV options
- **[AmneziaWG Parameters](./docs/AMNEZIAWG_PARAMETERS.md)** - Deep dive into obfuscation settings
- **[Password Generation](./docs/PASSWORD_GENERATION.md)** - Using the `wgpw` utility
- **[Architecture](./docs/ARCHITECTURE.md)** - Project structure and components
- **[API Reference](./docs/API_REFERENCE.md)** - REST API endpoints

## 🔧 Important Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `WG_HOST` | **YES** | - | Your server's public IP or domain |
| `PASSWORD_HASH` | No | - | Bcrypt hash for Web UI login |
| `PORT` | No | `51821` | Web UI TCP port |
| `WG_PORT` | No | `51820` | WireGuard UDP port |
| `NO_WEB_UI` | No | `false` | Disable Web UI (API only) |

### AmneziaWG Obfuscation (Auto-Random)

These parameters are **automatically randomized** if not set:

- `JC` (Junk packet count): Random 3-10
- `S1`, `S2` (Junk sizes): Random 15-150
- `H1`, `H2`, `H3`, `H4` (Magic headers): Random uint32

See [AmneziaWG Parameters Guide](./docs/AMNEZIAWG_PARAMETERS.md) for detailed explanation.

## 🛠️ Advanced Usage

### API-Only Mode (No Web UI)

```bash
NO_WEB_UI=true
docker compose up -d
```

All REST API endpoints remain available at `http://localhost:51821/api/*`

### Prometheus Metrics

```bash
ENABLE_PROMETHEUS_METRICS=true
docker compose up -d

# Access metrics
curl http://localhost:51821/metrics
```

### Custom AmneziaWG Settings

```bash
# Set specific obfuscation parameters
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

## 🏗️ Development

### Build from Source

```bash
# Clone repository
git clone https://github.com/your-repo/amnezia-wg-easy.git
cd amnezia-wg-easy

# Build binaries
go build -o amnezia-wg-easy .
go build -o wgpw ./cmd/wgpw

# Or use Makefile
make build
make build-wgpw
```

### Project Structure

See [Architecture Guide](./docs/ARCHITECTURE.md) for detailed project structure.

```
amnezia-wg-easy/
├── main.go                 # Application entry point
├── cmd/wgpw/              # Password generation utility
├── internal/
│   ├── config/            # Environment configuration
│   ├── models/            # Data structures
│   ├── wireguard/         # WireGuard operations
│   └── api/               # HTTP handlers & routes
└── www/                   # Web UI (Vue.js)
```

### Docker Build

```bash
docker build -t amnezia-wg-easy:dev .

# Or with compose
docker compose -f docker-compose.yml up --build
```

## 🔐 Security

- **Always set `PASSWORD_HASH`** in production
- Use strong passwords (12+ characters)
- Keep WireGuard private keys secure
- Use firewall rules to restrict access
- Enable `PROMETHEUS_METRICS_PASSWORD` if exposing metrics

## 🔄 Migration from Node.js Version

1. Stop the old container
2. Backup your `/etc/wireguard/wg0.json`
3. Start the Go version with the same volume
4. Configuration is automatically compatible

## 📊 REST API

### Authentication

```bash
# Login
curl -X POST http://localhost:51821/api/session \
  -H "Content-Type: application/json" \
  -d '{"password":"your_password"}'

# Or use Authorization header
curl -H "Authorization: your_password" \
  http://localhost:51821/api/wireguard/client
```

### Client Management

```bash
# List clients
GET /api/wireguard/client

# Create client
POST /api/wireguard/client
{"name":"client1"}

# Delete client
DELETE /api/wireguard/client/:id

# Download config
GET /api/wireguard/client/:id/configuration

# Get QR code
GET /api/wireguard/client/:id/qrcode.svg
```

See [API Reference](./docs/API_REFERENCE.md) for complete documentation.

## 🐛 Troubleshooting

### Error: "WG_HOST environment variable is required"

```bash
# Solution: Set WG_HOST
docker run -e WG_HOST=192.168.1.1 ...
```

### Error: "Cannot find device wg0"

Your kernel doesn't support WireGuard. Use `amneziavpn/amnezia-wg` base image.

### Error: "no such file or directory: /etc/wireguard"

```bash
# Solution: Mount volume
docker run -v ~/.amnezia-wg:/etc/wireguard ...
```

### Metrics not working

```bash
# Enable in .env
ENABLE_PROMETHEUS_METRICS=true
```

## 🙏 Credits

- [wg-easy](https://github.com/wg-easy/wg-easy) - Original project
- [amnezia-wg-easy](https://github.com/spcfox/amnezia-wg-easy) - AmneziaWG integration
- [AmneziaVPN](https://github.com/amnezia-vpn) - AmneziaWG protocol

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🌟 Support

- ⭐ Star this repository
- 🐛 Report issues
- 💡 Suggest features
- 🤝 Contribute code

---

**Made with ❤️ using Go**

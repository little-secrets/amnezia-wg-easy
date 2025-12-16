# Documentation

Complete documentation for AmneziaWG Easy (Go Edition).

## 📚 Available Guides

### Getting Started
- **[Main README](../README.md)** - Quick start and overview
- **[Password Generation](./PASSWORD_GENERATION.md)** - Using the `wgpw` utility

### Configuration
- **[Environment Variables](./ENVIRONMENT_VARIABLES.md)** - Complete ENV reference
- **[AmneziaWG Parameters](./AMNEZIAWG_PARAMETERS.md)** - Traffic obfuscation settings

### Development
- **[Architecture](./ARCHITECTURE.md)** - Project structure and components
- **[API Reference](./API_REFERENCE.md)** - REST API endpoints

### Other
- **[Go Edition Changelog](./CHANGELOG_GO.md)** - What's new in the Go rewrite

## 🎯 Quick Links

### Most Common Tasks

**Generate password hash:**
```bash
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw mypassword
```
→ [Password Generation Guide](./PASSWORD_GENERATION.md)

**Set environment variables:**
```bash
WG_HOST=your.server.ip
PASSWORD_HASH='$2a$12$...'
```
→ [Environment Variables Reference](./ENVIRONMENT_VARIABLES.md)

**Understand obfuscation:**
```bash
JC=7  # Junk packet count
H1=1234567891  # Magic header
```
→ [AmneziaWG Parameters](./AMNEZIAWG_PARAMETERS.md)

**Use the API:**
```bash
curl http://localhost:51821/api/wireguard/client
```
→ [API Reference](./API_REFERENCE.md)

## 📖 Reading Order

### For First-Time Users:
1. [Main README](../README.md) - Start here
2. [Password Generation](./PASSWORD_GENERATION.md) - Set up authentication
3. [Environment Variables](./ENVIRONMENT_VARIABLES.md) - Configure your setup

### For Developers:
1. [Architecture](./ARCHITECTURE.md) - Understand the codebase
2. [API Reference](./API_REFERENCE.md) - Integrate with the API
3. [Environment Variables](./ENVIRONMENT_VARIABLES.md) - Available options

### For Advanced Users:
1. [AmneziaWG Parameters](./AMNEZIAWG_PARAMETERS.md) - Fine-tune obfuscation
2. [API Reference](./API_REFERENCE.md) - Automate operations
3. [Architecture](./ARCHITECTURE.md) - Extend functionality

## 🔍 Search by Topic

### Authentication
- [Password Generation](./PASSWORD_GENERATION.md)
- [Environment Variables - Authentication](./ENVIRONMENT_VARIABLES.md#authentication)
- [API Reference - Authentication](./API_REFERENCE.md#authentication)

### AmneziaWG Obfuscation
- [AmneziaWG Parameters](./AMNEZIAWG_PARAMETERS.md)
- [Environment Variables - AmneziaWG](./ENVIRONMENT_VARIABLES.md#amneziawg-obfuscation)

### API Integration
- [API Reference](./API_REFERENCE.md)
- [Architecture - API Layer](./ARCHITECTURE.md#5-api-layer)

### Configuration
- [Environment Variables](./ENVIRONMENT_VARIABLES.md)
- [Architecture - Configuration Layer](./ARCHITECTURE.md#2-configuration-layer)

### Metrics
- [API Reference - Metrics](./API_REFERENCE.md#prometheus-metrics)
- [Environment Variables - Feature Flags](./ENVIRONMENT_VARIABLES.md#feature-flags)

## 💡 Tips

- Most environment variables have sensible defaults
- Only `WG_HOST` is required
- AmneziaWG parameters auto-randomize for security
- Check [CHANGELOG_GO.md](./CHANGELOG_GO.md) for what's new in Go edition

## 🐛 Found an Issue?

- Check the specific guide for troubleshooting sections
- See [Environment Variables](./ENVIRONMENT_VARIABLES.md#troubleshooting)
- Open an issue on GitHub

---

**Need help?** Start with the [Main README](../README.md) or jump to the guide you need above!


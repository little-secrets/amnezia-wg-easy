# Per-Client AmneziaWG Parameters

**New in v1.1:** You can now set individual AmneziaWG obfuscation parameters for each client.

## 🎯 Why Per-Client Parameters?

### Use Cases

1. **Testing**: Test different obfuscation settings without affecting all clients
2. **Diversity**: Different clients use different obfuscation patterns
3. **Specialized**: Some clients need stronger obfuscation than others
4. **Migration**: Gradually migrate clients to new parameters

### Example Scenarios

**Scenario 1**: Your VPN is getting blocked in Country A but works fine in Country B
- Give Country A clients stronger obfuscation (higher JC, custom headers)
- Keep Country B clients on default settings

**Scenario 2**: Testing new obfuscation parameters
- Create test client with new parameters
- Verify it works before rolling out to all clients

**Scenario 3**: High-security users
- Some users need maximum obfuscation
- Other users prefer performance over stealth

## 🚀 Quick Start

### Default Client (uses server settings)

```bash
curl -X POST http://localhost:51821/api/wireguard/client \
  -H "Content-Type: application/json" \
  -H "Authorization: your_password" \
  -d '{
    "name": "normal-client"
  }'
```

**Result:** Client inherits all server AmneziaWG parameters.

### Custom Client (specific parameters)

```bash
curl -X POST http://localhost:51821/api/wireguard/client \
  -H "Content-Type: application/json" \
  -H "Authorization: your_password" \
  -d '{
    "name": "stealth-client",
    "jc": "10",
    "h1": "9999999",
    "h2": "8888888",
    "h3": "7777777",
    "h4": "6666666"
  }'
```

**Result:** Client uses:
- Custom JC, H1-H4 (as specified)
- Server defaults for JMIN, JMAX, S1, S2 (not specified)

## 📊 Parameter Inheritance

```
Server ENV:
  JC=5, JMIN=50, JMAX=1000, S1=75, S2=75, H1=random, H2=random, ...

Client Request:
  {
    "name": "client1",
    "jc": "10",     ← Override
    "h1": "9999"    ← Override
                    ← Others inherit from server
  }

Final Client Config:
  JC=10           ← From request
  JMIN=50         ← From server
  JMAX=1000       ← From server
  S1=75           ← From server
  S2=75           ← From server
  H1=9999         ← From request
  H2=random       ← From server
  ...
```

## 🔧 Complete Examples

### Example 1: Maximum Stealth Client

```bash
curl -X POST http://localhost:51821/api/wireguard/client \
  -H "Content-Type: application/json" \
  -H "Authorization: your_password" \
  -d '{
    "name": "max-stealth",
    "jc": "10",
    "jmin": "100",
    "jmax": "1400",
    "s1": "150",
    "s2": "150",
    "h1": "2147483647",
    "h2": "2147483646",
    "h3": "2147483645",
    "h4": "2147483644"
  }'
```

**When to use:**
- User in highly censored network
- Maximum DPI evasion needed
- Performance is not critical

### Example 2: Performance-Optimized Client

```bash
curl -X POST http://localhost:51821/api/wireguard/client \
  -H "Content-Type: application/json" \
  -H "Authorization: your_password" \
  -d '{
    "name": "fast-client",
    "jc": "3",
    "jmin": "20",
    "jmax": "200",
    "s1": "30",
    "s2": "30"
  }'
```

**When to use:**
- Low latency required
- Network has no DPI
- Client on limited bandwidth

### Example 3: Custom Headers Only

```bash
curl -X POST http://localhost:51821/api/wireguard/client \
  -H "Content-Type: application/json" \
  -H "Authorization: your_password" \
  -d '{
    "name": "custom-headers",
    "h1": "1111111",
    "h2": "2222222",
    "h3": "3333333",
    "h4": "4444444"
  }'
```

**When to use:**
- Want unique headers per client
- Keep other parameters at server defaults
- Create fingerprint diversity

## 📱 Web UI Support

Currently, per-client parameters are **API-only**. The Web UI creates clients with server defaults.

**To create client with custom parameters:**
1. Use API endpoint (see examples above)
2. Or use automation tools (Python, Bash, etc.)

**Web UI shows:**
- Client name, IP, status
- Configuration download (includes client's custom parameters)
- QR code (includes client's custom parameters)

## 🔍 Viewing Client Parameters

### Via API

```bash
# Get all clients
curl http://localhost:51821/api/wireguard/client \
  -H "Authorization: your_password"
```

**Response includes client parameters:**
```json
[
  {
    "id": "uuid",
    "name": "client1",
    "address": "10.8.0.2",
    // ... other fields ...
    // Custom parameters shown in client configuration
  }
]
```

### Via Configuration File

```bash
# Download client configuration
curl http://localhost:51821/api/wireguard/client/:id/configuration \
  -H "Authorization: your_password"
```

**Output shows parameters:**
```ini
[Interface]
PrivateKey = ...
Address = 10.8.0.2/24
Jc = 10        ← Client-specific or server default
Jmin = 50
Jmax = 1000
S1 = 150       ← Client-specific
S2 = 150       ← Client-specific
H1 = 9999999   ← Client-specific
...
```

### Via Backend JSON

```bash
# SSH into server
docker compose exec amnezia-wg-easy cat /etc/wireguard/wg0.json | jq
```

**Shows stored parameters:**
```json
{
  "clients": {
    "uuid": {
      "name": "client1",
      "jc": "10",       ← Stored if custom
      "s1": "150",      ← Stored if custom
      // null if using server defaults
    }
  }
}
```

## 🛡️ Security Considerations

### Pros of Per-Client Parameters

✅ **Diversity**: Different obfuscation per client makes pattern detection harder
✅ **Flexibility**: Adapt to specific network conditions
✅ **Isolation**: If one client's parameters are detected, others remain safe

### Cons of Per-Client Parameters

⚠️ **Complexity**: More parameters to manage
⚠️ **Debugging**: Harder to troubleshoot issues
⚠️ **Coordination**: Must track which client has which parameters

### Best Practices

1. **Document**: Keep track of which clients have custom parameters
2. **Test**: Always test custom parameters before deploying
3. **Monitor**: Watch for connection issues
4. **Consistency**: Use similar parameters for clients in same network
5. **Backup**: Regular backups include all custom parameters

## 🔄 Updating Client Parameters

Currently, you **cannot update** parameters for existing clients via API.

**To change parameters:**
1. Delete the old client
2. Create new client with updated parameters
3. Redistribute configuration to user

**Alternative (Advanced):**
1. Edit `/etc/wireguard/wg0.json` directly
2. Restart container

```bash
# Manual edit (advanced users only)
docker compose exec amnezia-wg-easy vi /etc/wireguard/wg0.json
docker compose restart
```

## 📊 Comparison Table

| Feature | Server Parameters | Per-Client Parameters |
|---------|------------------|----------------------|
| **Set via** | Environment variables | API request |
| **Scope** | All clients (default) | Specific client |
| **When applied** | At client creation | At client creation |
| **Can update** | Yes (restart required) | No (recreate client) |
| **Web UI support** | Yes | No (API only) |
| **Complexity** | Low | Medium |
| **Use case** | Standard setup | Advanced scenarios |

## 🧪 Testing Examples

### Python Script

```python
import requests

api_url = "http://localhost:51821/api"
password = "your_password"

# Login
session = requests.Session()
session.post(f"{api_url}/session", json={"password": password})

# Create client with custom parameters
response = session.post(f"{api_url}/wireguard/client", json={
    "name": "test-client",
    "jc": "10",
    "s1": "150",
    "s2": "150",
    "h1": "9999999",
    "h2": "8888888",
    "h3": "7777777",
    "h4": "6666666"
})

print(response.json())  # {"success": true}

# Download config
client_id = "..." # Get from GET /api/wireguard/client
config = session.get(f"{api_url}/wireguard/client/{client_id}/configuration")
print(config.text)
```

### Bash Script

```bash
#!/bin/bash

API_URL="http://localhost:51821/api"
PASSWORD="your_password"

# Login and get cookie
curl -c cookies.txt -X POST "$API_URL/session" \
  -H "Content-Type: application/json" \
  -d "{\"password\":\"$PASSWORD\"}"

# Create client with custom params
curl -b cookies.txt -X POST "$API_URL/wireguard/client" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-client",
    "jc": "10",
    "h1": "9999999"
  }'
```

## 📚 Related Documentation

- [AmneziaWG Parameters Guide](./AMNEZIAWG_PARAMETERS.md) - Detailed parameter explanations
- [API Reference](./API_REFERENCE.md) - Complete API documentation
- [Environment Variables](./ENVIRONMENT_VARIABLES.md) - Server-wide configuration

---

**Questions?** Check the [main README](../README.md) or open an issue on GitHub.


# API Reference

Complete REST API documentation for AmneziaWG Easy.

## 🌐 Base URL

```
http://your-server:51821/api
```

## 🔐 Authentication

### Session-Based (Web UI)

**Login:**
```bash
POST /api/session
Content-Type: application/json

{
  "password": "your_password",
  "remember": true
}
```

**Response:**
```json
{
  "success": true
}
```

**Logout:**
```bash
DELETE /api/session
```

### Header-Based (API Clients)

```bash
curl -H "Authorization: your_password" \
  http://localhost:51821/api/wireguard/client
```

---

## 📊 System Endpoints

### GET `/api/release`

Get application version.

```bash
curl http://localhost:51821/api/release
```

**Response:**
```json
"1.0.0"
```

### GET `/api/lang`

Get configured language.

```bash
curl http://localhost:51821/api/lang
```

**Response:**
```json
"en"
```

### GET `/api/session`

Check authentication status.

```bash
curl http://localhost:51821/api/session
```

**Response:**
```json
{
  "requiresPassword": true,
  "authenticated": false
}
```

---

## 👥 Client Management

### GET `/api/wireguard/client`

List all clients with their stats.

```bash
curl http://localhost:51821/api/wireguard/client
```

**Response:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "client1",
    "enabled": true,
    "address": "10.8.0.2",
    "publicKey": "ABC123...",
    "createdAt": "2024-12-16T10:00:00Z",
    "updatedAt": "2024-12-16T10:00:00Z",
    "expiredAt": null,
    "downloadableConfig": true,
    "persistentKeepalive": "0",
    "latestHandshakeAt": "2024-12-16T12:30:00Z",
    "transferRx": 1024000,
    "transferTx": 2048000,
    "endpoint": "192.168.1.100:51820"
  }
]
```

### POST `/api/wireguard/client`

Create a new client.

```bash
curl -X POST http://localhost:51821/api/wireguard/client \
  -H "Content-Type: application/json" \
  -d '{
    "name": "client1",
    "expiredDate": "2025-12-31"
  }'
```

**Request:**
```json
{
  "name": "client1",           // Required
  "expiredDate": "2025-12-31"  // Optional (YYYY-MM-DD)
}
```

**Response:**
```json
{
  "success": true
}
```

### DELETE `/api/wireguard/client/:clientId`

Delete a client.

```bash
curl -X DELETE http://localhost:51821/api/wireguard/client/550e8400-e29b-41d4-a716-446655440000
```

**Response:**
```json
{
  "success": true
}
```

### POST `/api/wireguard/client/:clientId/enable`

Enable a client.

```bash
curl -X POST http://localhost:51821/api/wireguard/client/550e8400-e29b-41d4-a716-446655440000/enable
```

**Response:**
```json
{
  "success": true
}
```

### POST `/api/wireguard/client/:clientId/disable`

Disable a client.

```bash
curl -X POST http://localhost:51821/api/wireguard/client/550e8400-e29b-41d4-a716-446655440000/disable
```

**Response:**
```json
{
  "success": true
}
```

### PUT `/api/wireguard/client/:clientId/name`

Update client name.

```bash
curl -X PUT http://localhost:51821/api/wireguard/client/550e8400-e29b-41d4-a716-446655440000/name \
  -H "Content-Type: application/json" \
  -d '{"name": "new-name"}'
```

**Request:**
```json
{
  "name": "new-name"
}
```

**Response:**
```json
{
  "success": true
}
```

### PUT `/api/wireguard/client/:clientId/address`

Update client IP address.

```bash
curl -X PUT http://localhost:51821/api/wireguard/client/550e8400-e29b-41d4-a716-446655440000/address \
  -H "Content-Type: application/json" \
  -d '{"address": "10.8.0.100"}'
```

**Request:**
```json
{
  "address": "10.8.0.100"
}
```

**Response:**
```json
{
  "success": true
}
```

### PUT `/api/wireguard/client/:clientId/expireDate`

Update client expiration date.

```bash
curl -X PUT http://localhost:51821/api/wireguard/client/550e8400-e29b-41d4-a716-446655440000/expireDate \
  -H "Content-Type: application/json" \
  -d '{"expireDate": "2025-12-31"}'
```

**Request:**
```json
{
  "expireDate": "2025-12-31"  // Or empty string to remove expiry
}
```

**Response:**
```json
{
  "success": true
}
```

---

## 📥 Configuration Download

### GET `/api/wireguard/client/:clientId/configuration`

Download client WireGuard configuration file.

```bash
curl http://localhost:51821/api/wireguard/client/550e8400-e29b-41d4-a716-446655440000/configuration \
  -O -J
```

**Response:**
```ini
Content-Disposition: attachment; filename="client1.conf"
Content-Type: text/plain

[Interface]
PrivateKey = ABC123...
Address = 10.8.0.2/24
DNS = 1.1.1.1
Jc = 7
Jmin = 50
Jmax = 1000
S1 = 89
S2 = 134
H1 = 1847362945
H2 = 492817364
H3 = 1638274562
H4 = 2014738291

[Peer]
PublicKey = SERVER_PUBLIC_KEY
PresharedKey = PRESHARED_KEY
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 0
Endpoint = vpn.example.com:51820
```

### GET `/api/wireguard/client/:clientId/qrcode.svg`

Get QR code for mobile apps.

```bash
curl http://localhost:51821/api/wireguard/client/550e8400-e29b-41d4-a716-446655440000/qrcode.svg
```

**Response:**
```xml
Content-Type: image/svg+xml

<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512">
  <!-- QR code SVG data -->
</svg>
```

---

## 🔗 One-Time Links

**Note:** Requires `WG_ENABLE_ONE_TIME_LINKS=true`

### POST `/api/wireguard/client/:clientId/generateOneTimeLink`

Generate a one-time download link (expires in 5 minutes).

```bash
curl -X POST http://localhost:51821/api/wireguard/client/550e8400-e29b-41d4-a716-446655440000/generateOneTimeLink
```

**Response:**
```json
{
  "success": true
}
```

After generation, the client object will have:
```json
{
  "oneTimeLink": "a1b2c3d4",
  "oneTimeLinkExpiresAt": "2024-12-16T12:35:00Z"
}
```

### GET `/cnf/:oneTimeLink`

Public endpoint to download config via one-time link.

```bash
curl http://localhost:51821/cnf/a1b2c3d4 -O -J
```

**Note:** Link expires after use or after 5 minutes.

---

## 💾 Backup & Restore

### GET `/api/wireguard/backup`

Download complete configuration backup.

```bash
curl http://localhost:51821/api/wireguard/backup -O -J
```

**Response:**
```json
Content-Disposition: attachment; filename="wg0.json"

{
  "server": {
    "privateKey": "...",
    "publicKey": "...",
    "address": "10.8.0.1",
    "jc": "7",
    ...
  },
  "clients": {
    "uuid-1": {...},
    "uuid-2": {...}
  }
}
```

### PUT `/api/wireguard/restore`

Restore configuration from backup.

```bash
curl -X PUT http://localhost:51821/api/wireguard/restore \
  -H "Content-Type: application/json" \
  -d '{
    "file": "{\"server\":{...},\"clients\":{...}}"
  }'
```

**Request:**
```json
{
  "file": "JSON_STRING_OF_CONFIG"
}
```

**Response:**
```json
{
  "success": true
}
```

---

## 📊 Prometheus Metrics

**Note:** Requires `ENABLE_PROMETHEUS_METRICS=true`

### GET `/metrics`

Get Prometheus metrics.

```bash
# Without password
curl http://localhost:51821/metrics

# With Basic Auth
curl -u ":your_metrics_password" http://localhost:51821/metrics
```

**Response:**
```prometheus
# HELP wireguard_configured_peers
# TYPE wireguard_configured_peers gauge
wireguard_configured_peers{interface="wg0"} 3

# HELP wireguard_enabled_peers
# TYPE wireguard_enabled_peers gauge
wireguard_enabled_peers{interface="wg0"} 2

# HELP wireguard_connected_peers
# TYPE wireguard_connected_peers gauge
wireguard_connected_peers{interface="wg0"} 1

# HELP wireguard_sent_bytes
# TYPE wireguard_sent_bytes counter
wireguard_sent_bytes{interface="wg0",enabled="true",address="10.8.0.2",name="client1"} 1024000

# HELP wireguard_received_bytes
# TYPE wireguard_received_bytes counter
wireguard_received_bytes{interface="wg0",enabled="true",address="10.8.0.2",name="client1"} 2048000

# HELP wireguard_latest_handshake_seconds
# TYPE wireguard_latest_handshake_seconds gauge
wireguard_latest_handshake_seconds{interface="wg0",enabled="true",address="10.8.0.2",name="client1"} 30
```

### GET `/metrics/json`

Get metrics in JSON format.

```bash
curl http://localhost:51821/metrics/json
```

**Response:**
```json
{
  "wireguard_configured_peers": 3,
  "wireguard_enabled_peers": 2,
  "wireguard_connected_peers": 1
}
```

---

## 🎛️ UI Settings Endpoints

### GET `/api/remember-me`

Check if "Remember Me" is enabled.

```bash
curl http://localhost:51821/api/remember-me
```

**Response:**
```json
true
```

### GET `/api/ui-traffic-stats`

Check if traffic stats are enabled.

```bash
curl http://localhost:51821/api/ui-traffic-stats
```

**Response:**
```json
true
```

### GET `/api/ui-chart-type`

Get chart type setting.

```bash
curl http://localhost:51821/api/ui-chart-type
```

**Response:**
```json
"1"
```

### GET `/api/wg-enable-one-time-links`

Check if one-time links are enabled.

```bash
curl http://localhost:51821/api/wg-enable-one-time-links
```

**Response:**
```json
true
```

### GET `/api/wg-enable-expire-time`

Check if client expiry is enabled.

```bash
curl http://localhost:51821/api/wg-enable-expire-time
```

**Response:**
```json
true
```

### GET `/api/ui-sort-clients`

Check if client sorting is enabled.

```bash
curl http://localhost:51821/api/ui-sort-clients
```

**Response:**
```json
true
```

### GET `/api/ui-avatar-settings`

Get avatar settings.

```bash
curl http://localhost:51821/api/ui-avatar-settings
```

**Response:**
```json
{
  "dicebear": "bottts",
  "gravatar": true
}
```

---

## ⚠️ Error Responses

All errors follow this format:

```json
{
  "error": "Error message here"
}
```

### Common HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| `200` | Success | Request completed successfully |
| `204` | No Content | Operation completed, no response body |
| `400` | Bad Request | Invalid request format |
| `401` | Unauthorized | Not logged in or wrong password |
| `403` | Forbidden | Invalid client ID (security) |
| `404` | Not Found | Client or resource not found |
| `500` | Internal Server Error | Server-side error |

### Example Errors

**Not authenticated:**
```json
Status: 401
{
  "error": "Not Logged In"
}
```

**Client not found:**
```json
Status: 404
{
  "error": "client not found: 550e8400-e29b-41d4-a716-446655440000"
}
```

**Invalid IP address:**
```json
Status: 400
{
  "error": "invalid IPv4 address: 10.8.0.999"
}
```

---

## 🔧 Complete API Examples

### Python

```python
import requests

base_url = "http://localhost:51821/api"
password = "mypassword"

# Login
session = requests.Session()
response = session.post(f"{base_url}/session", json={
    "password": password,
    "remember": True
})

# Create client
response = session.post(f"{base_url}/wireguard/client", json={
    "name": "python-client"
})

# List clients
response = session.get(f"{base_url}/wireguard/client")
clients = response.json()

for client in clients:
    print(f"{client['name']}: {client['address']}")

# Download config
client_id = clients[0]['id']
response = session.get(f"{base_url}/wireguard/client/{client_id}/configuration")
with open("client.conf", "w") as f:
    f.write(response.text)
```

### JavaScript/Node.js

```javascript
const axios = require('axios');

const baseURL = 'http://localhost:51821/api';
const password = 'mypassword';

// Create axios instance with session
const api = axios.create({
  baseURL,
  withCredentials: true
});

async function main() {
  // Login
  await api.post('/session', { password, remember: true });
  
  // Create client
  await api.post('/wireguard/client', { name: 'js-client' });
  
  // List clients
  const { data: clients } = await api.get('/wireguard/client');
  
  clients.forEach(client => {
    console.log(`${client.name}: ${client.address}`);
  });
  
  // Download config
  const clientId = clients[0].id;
  const { data: config } = await api.get(
    `/wireguard/client/${clientId}/configuration`
  );
  
  require('fs').writeFileSync('client.conf', config);
}

main();
```

### Bash/cURL

```bash
#!/bin/bash

BASE_URL="http://localhost:51821/api"
PASSWORD="mypassword"
COOKIE_FILE="/tmp/wg-cookies.txt"

# Login
curl -c $COOKIE_FILE -X POST "$BASE_URL/session" \
  -H "Content-Type: application/json" \
  -d "{\"password\":\"$PASSWORD\",\"remember\":true}"

# Create client
curl -b $COOKIE_FILE -X POST "$BASE_URL/wireguard/client" \
  -H "Content-Type: application/json" \
  -d '{"name":"bash-client"}'

# List clients
curl -b $COOKIE_FILE "$BASE_URL/wireguard/client" | jq .

# Download config
CLIENT_ID=$(curl -b $COOKIE_FILE "$BASE_URL/wireguard/client" | jq -r '.[0].id')
curl -b $COOKIE_FILE "$BASE_URL/wireguard/client/$CLIENT_ID/configuration" \
  -o client.conf
```

---

## 📚 Related Documentation

- [Environment Variables](./ENVIRONMENT_VARIABLES.md)
- [Architecture](./ARCHITECTURE.md)
- [Password Generation](./PASSWORD_GENERATION.md)

---

**Need help?** Open an issue on GitHub or check the [main README](../README.md).


# API Client Examples

This directory contains example code for interacting with the AmneziaWG Easy API.

## 📁 Available Examples

### Python Client (`openapi_client.py`)

Full-featured Python client with all API methods.

**Installation:**
```bash
pip install requests
```

**Usage:**
```bash
# Edit the script to set your password
python openapi_client.py
```

**Features:**
- Session management (login/logout)
- Client CRUD operations
- Configuration download
- QR code generation
- Backup/restore

### JavaScript/Node.js Client (`openapi_client.js`)

Complete Node.js client using axios.

**Installation:**
```bash
npm install axios
```

**Usage:**
```bash
# Edit the script to set your password
node openapi_client.js
```

**Features:**
- All API endpoints
- Promise-based async/await
- Error handling
- File downloads
- Metrics support

## 🚀 Quick Start

### 1. Start AmneziaWG Easy

```bash
docker compose up -d
```

### 2. Choose Your Language

**Python:**
```bash
cd docs/examples
pip install requests
python openapi_client.py
```

**JavaScript:**
```bash
cd docs/examples
npm install axios
node openapi_client.js
```

### 3. Customize

Edit the example files to:
- Set your password
- Change base URL
- Add custom logic
- Integrate into your application

## 📚 Using Generated Clients

For production applications, consider generating a typed client from the OpenAPI spec:

### Generate Python Client

```bash
npm install -g @openapitools/openapi-generator-cli

openapi-generator-cli generate \
  -i http://localhost:51821/api/openapi.yaml \
  -g python \
  -o ./amnezia-client-python
```

**Usage:**
```python
import amnezia_client
from amnezia_client.api import clients_api
from amnezia_client.model.create_client_request import CreateClientRequest

# Configure client
configuration = amnezia_client.Configuration(
    host = "http://localhost:51821"
)

# Create API instance
with amnezia_client.ApiClient(configuration) as api_client:
    api_instance = clients_api.ClientsApi(api_client)
    
    # Create client
    create_request = CreateClientRequest(
        name="my-client",
        expired_date="2025-12-31"
    )
    
    api_response = api_instance.create_client(create_request)
```

### Generate TypeScript Client

```bash
openapi-generator-cli generate \
  -i http://localhost:51821/api/openapi.yaml \
  -g typescript-axios \
  -o ./amnezia-client-ts
```

**Usage:**
```typescript
import { 
    Configuration, 
    ClientsApi, 
    AuthenticationApi 
} from './amnezia-client-ts';

const config = new Configuration({
    basePath: 'http://localhost:51821'
});

const authApi = new AuthenticationApi(config);
const clientsApi = new ClientsApi(config);

// Login
await authApi.createSession({
    password: 'your_password',
    remember: true
});

// Create client
await clientsApi.createClient({
    name: 'my-client',
    expiredDate: '2025-12-31'
});
```

### Generate Go Client

```bash
openapi-generator-cli generate \
  -i http://localhost:51821/api/openapi.yaml \
  -g go \
  -o ./amnezia-client-go
```

**Usage:**
```go
package main

import (
    "context"
    "fmt"
    amnezia "github.com/your-org/amnezia-client-go"
)

func main() {
    config := amnezia.NewConfiguration()
    config.Servers = amnezia.ServerConfigurations{
        {
            URL: "http://localhost:51821",
        },
    }
    
    client := amnezia.NewAPIClient(config)
    ctx := context.Background()
    
    // Login
    loginReq := amnezia.LoginRequest{
        Password: "your_password",
        Remember: true,
    }
    
    _, _, err := client.AuthenticationApi.CreateSession(ctx).
        LoginRequest(loginReq).Execute()
    
    if err != nil {
        panic(err)
    }
    
    // Create client
    createReq := amnezia.CreateClientRequest{
        Name: "my-client",
    }
    
    _, _, err = client.ClientsApi.CreateClient(ctx).
        CreateClientRequest(createReq).Execute()
    
    if err != nil {
        panic(err)
    }
}
```

## 🔧 Common Tasks

### Create Multiple Clients

**Python:**
```python
client = AmneziaWGClient()
client.login("password")

for i in range(10):
    client.create_client(
        name=f"client-{i+1}",
        expired_date="2025-12-31"
    )
```

**JavaScript:**
```javascript
const client = new AmneziaWGClient();
await client.login('password');

for (let i = 0; i < 10; i++) {
    await client.createClient({
        name: `client-${i+1}`,
        expiredDate: '2025-12-31'
    });
}
```

### Bulk Download Configurations

**Python:**
```python
clients = client.get_clients()

for c in clients:
    filename = f"{c['name']}.conf"
    client.download_config(c['id'], filename)
    print(f"Downloaded {filename}")
```

**JavaScript:**
```javascript
const clients = await client.getClients();

for (const c of clients) {
    const filename = `${c.name}.conf`;
    await client.downloadConfig(c.id, filename);
    console.log(`Downloaded ${filename}`);
}
```

### Monitor Connected Clients

**Python:**
```python
import time

while True:
    clients = client.get_clients()
    connected = [c for c in clients if c.get('latestHandshakeAt')]
    
    print(f"Connected: {len(connected)}/{len(clients)}")
    
    time.sleep(10)
```

**JavaScript:**
```javascript
setInterval(async () => {
    const clients = await client.getClients();
    const connected = clients.filter(c => c.latestHandshakeAt);
    
    console.log(`Connected: ${connected.length}/${clients.length}`);
}, 10000);
```

### Custom AmneziaWG Parameters

**Python:**
```python
client.create_client(
    name="stealth-client",
    jc="10",
    jmin="30",
    jmax="1200",
    s1="150",
    s2="150",
    h1="9876543",
    h2="8765432",
    h3="7654321",
    h4="6543210"
)
```

**JavaScript:**
```javascript
await client.createClient({
    name: 'stealth-client',
    jc: '10',
    jmin: '30',
    jmax: '1200',
    s1: '150',
    s2: '150',
    h1: '9876543',
    h2: '8765432',
    h3: '7654321',
    h4: '6543210'
});
```

## 🧪 Testing

### Test with cURL

```bash
# Login
curl -c cookies.txt -X POST http://localhost:51821/api/session \
  -H "Content-Type: application/json" \
  -d '{"password":"your_password","remember":true}'

# List clients
curl -b cookies.txt http://localhost:51821/api/wireguard/client

# Create client
curl -b cookies.txt -X POST http://localhost:51821/api/wireguard/client \
  -H "Content-Type: application/json" \
  -d '{"name":"test-client"}'
```

### Test with HTTPie

```bash
# Login and save session
http --session=wg POST http://localhost:51821/api/session \
  password=your_password remember:=true

# List clients
http --session=wg http://localhost:51821/api/wireguard/client

# Create client
http --session=wg POST http://localhost:51821/api/wireguard/client \
  name=test-client expiredDate=2025-12-31
```

## 📖 Documentation

- [API Reference](../API_REFERENCE.md)
- [OpenAPI Guide](../OPENAPI_GUIDE.md)
- [Swagger UI](http://localhost:51821/api/docs)

## 🐛 Troubleshooting

### Connection Refused

Make sure AmneziaWG Easy is running:
```bash
docker compose ps
curl http://localhost:51821/api/release
```

### Authentication Failed

Check your password:
```bash
# Regenerate password hash
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw mypassword
```

### Module Not Found

Install dependencies:
```bash
# Python
pip install requests

# JavaScript
npm install axios
```

## 💡 Tips

1. **Use environment variables** for passwords:
   ```bash
   export WG_PASSWORD="your_password"
   python openapi_client.py
   ```

2. **Enable metrics** for monitoring:
   ```bash
   ENABLE_PROMETHEUS_METRICS=true docker compose up -d
   ```

3. **Use one-time links** for secure config distribution:
   ```bash
   WG_ENABLE_ONE_TIME_LINKS=true docker compose up -d
   ```

## 🤝 Contributing

Have a useful example? Submit a PR!

Possible contributions:
- Bash script examples
- PowerShell examples
- Ruby client
- PHP client
- Rust client
- Integration examples (Ansible, Terraform, etc.)

---

**Need help?** Check the [API documentation](../API_REFERENCE.md) or [open an issue](https://github.com/little-secrets/amnezia-wg-easy/issues)


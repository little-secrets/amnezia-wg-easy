# OpenAPI Documentation Guide

This guide explains how to use the OpenAPI documentation for AmneziaWG Easy API.

## 📚 What is OpenAPI?

OpenAPI (formerly Swagger) is a standard specification for describing REST APIs. It provides:

- **Machine-readable** format (YAML/JSON)
- **Interactive documentation** (Swagger UI)
- **Code generation** capabilities
- **API testing** tools

## 🌐 Accessing the Documentation

### Swagger UI (Interactive)

The easiest way to explore the API:

```
http://localhost:51821/api/docs
```

**Features:**
- 📖 Browse all endpoints with descriptions
- 🧪 Test API calls directly from your browser
- 📝 View request/response examples
- 🔐 Authenticate and save credentials
- 📋 Copy code snippets in multiple languages

### OpenAPI Specification File

Download or view the raw specification:

```
http://localhost:51821/api/openapi.yaml
```

**Use cases:**
- Import into Postman
- Generate client SDKs
- Validate API responses
- Generate documentation

## 🚀 Using Swagger UI

### 1. Open the UI

Navigate to `http://localhost:51821/api/docs`

### 2. Authenticate

For protected endpoints, click the **Authorize** button:

**Option A: Cookie Auth (Session)**
1. Click "Authorize" on `cookieAuth`
2. Login via `POST /api/session` first
3. Session cookie is automatically used

**Option B: Header Auth (Direct)**
1. Click "Authorize" on `headerAuth`
2. Enter your password
3. All requests will include `Authorization` header

### 3. Try an Endpoint

1. Select an endpoint (e.g., `GET /api/wireguard/client`)
2. Click **"Try it out"**
3. Fill in parameters (if any)
4. Click **"Execute"**
5. View the response

### 4. Copy Code Examples

After executing a request:
- Click on the **request snippet** section
- Choose your language (curl, Python, JavaScript, etc.)
- Copy and use in your code

## 🛠️ Using OpenAPI Spec with Tools

### Postman

1. Open Postman
2. Click **Import**
3. Enter URL: `http://localhost:51821/api/openapi.yaml`
4. All endpoints are imported as a collection

### Insomnia

1. Open Insomnia
2. Create/Open workspace
3. Click **Import/Export** → **Import Data**
4. Select **From URL**
5. Enter: `http://localhost:51821/api/openapi.yaml`

### Code Generation

Generate client libraries using OpenAPI Generator:

```bash
# Install openapi-generator
npm install -g @openapitools/openapi-generator-cli

# Generate Python client
openapi-generator-cli generate \
  -i http://localhost:51821/api/openapi.yaml \
  -g python \
  -o ./amnezia-client-python

# Generate Go client
openapi-generator-cli generate \
  -i http://localhost:51821/api/openapi.yaml \
  -g go \
  -o ./amnezia-client-go

# Generate TypeScript/JavaScript client
openapi-generator-cli generate \
  -i http://localhost:51821/api/openapi.yaml \
  -g typescript-axios \
  -o ./amnezia-client-ts
```

### API Testing

Use the spec for automated testing:

**With Dredd:**
```bash
npm install -g dredd
dredd http://localhost:51821/api/openapi.yaml http://localhost:51821
```

**With Schemathesis:**
```bash
pip install schemathesis
schemathesis run http://localhost:51821/api/openapi.yaml \
  --base-url http://localhost:51821 \
  --hypothesis-max-examples=100
```

## 📖 Specification Structure

The OpenAPI spec is organized into:

### Tags (Categories)

- **System**: Version, language, settings
- **Authentication**: Login/logout
- **Clients**: WireGuard client CRUD operations
- **Configuration**: Config files and QR codes
- **One-Time Links**: Temporary download links
- **Backup**: Configuration backup/restore
- **Metrics**: Prometheus metrics

### Security Schemes

1. **cookieAuth**: Session-based (Web UI)
2. **headerAuth**: Header-based (API clients)

### Schemas (Models)

All request/response models are defined with:
- Field types
- Required fields
- Examples
- Descriptions

## 💡 Tips and Best Practices

### Testing in Swagger UI

1. **Start with public endpoints** (no auth required):
   - `GET /api/release`
   - `GET /api/session`

2. **Then authenticate**:
   - Use `POST /api/session` to login
   - Or use Authorization header

3. **Test protected endpoints**:
   - `GET /api/wireguard/client`
   - `POST /api/wireguard/client`

### Using with CI/CD

```yaml
# GitHub Actions example
- name: Validate API responses
  run: |
    docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli validate \
      -i /local/docs/openapi.yaml
```

### Custom Modifications

If you need to customize the spec:

1. Edit `docs/openapi.yaml`
2. Restart the application
3. Refresh Swagger UI

**Common modifications:**
- Add custom examples
- Update descriptions
- Add new endpoints (if extending the API)

## 🔍 Endpoint Reference

### Quick Lookup

| Category | Endpoints Count | Authentication |
|----------|----------------|----------------|
| System | 8 | None |
| Authentication | 3 | Varies |
| Clients | 11 | Required |
| Configuration | 2 | Required |
| One-Time Links | 2 | Varies |
| Backup | 2 | Required |
| Metrics | 2 | Optional |

### Most Used Endpoints

1. **Create Client**: `POST /api/wireguard/client`
2. **List Clients**: `GET /api/wireguard/client`
3. **Download Config**: `GET /api/wireguard/client/{id}/configuration`
4. **Get QR Code**: `GET /api/wireguard/client/{id}/qrcode.svg`
5. **Login**: `POST /api/session`

## 🐛 Troubleshooting

### Swagger UI not loading

**Check:**
- Application is running
- Navigate to `http://localhost:51821/api/docs`
- Check browser console for errors
- Verify `docs/openapi.yaml` exists

### Can't authenticate

**Solutions:**
1. Use `POST /api/session` first
2. Or use Authorization header with password
3. Check PASSWORD_HASH is set correctly

### OpenAPI spec validation errors

```bash
# Validate the spec
docker run --rm -v ${PWD}:/local openapitools/openapi-generator-cli validate \
  -i /local/docs/openapi.yaml
```

### Missing endpoints

The spec only includes `/api/*` routes. Static files and Web UI routes are not included as they're not REST API endpoints.

## 📚 Related Documentation

- [API Reference (Markdown)](./API_REFERENCE.md)
- [Environment Variables](./ENVIRONMENT_VARIABLES.md)
- [Architecture](./ARCHITECTURE.md)

## 🔗 Useful Links

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI Documentation](https://swagger.io/tools/swagger-ui/)
- [OpenAPI Generator](https://openapi-generator.tech/)
- [Postman OpenAPI Import](https://learning.postman.com/docs/designing-and-developing-your-api/importing-an-api/)

---

**Need help?** Open an issue on [GitHub](https://github.com/little-secrets/amnezia-wg-easy/issues)


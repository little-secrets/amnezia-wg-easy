# 📚 Documentation Index

Complete documentation for AmneziaWG Easy Go Edition.

## 🚀 Quick Links

- **[Main README](../README.md)** - Project overview and quick start
- **[Swagger UI](http://localhost:51821/api/docs)** - Interactive API documentation 🆕

## 📖 Core Documentation

### Getting Started

- **[Environment Variables](./ENVIRONMENT_VARIABLES.md)** - Configuration options
- **[Password Generation](./PASSWORD_GENERATION.md)** - Using the `wgpw` utility

### Features

- **[AmneziaWG Parameters](./AMNEZIAWG_PARAMETERS.md)** - Traffic obfuscation explained
- **[Per-Client Parameters](./PER_CLIENT_PARAMETERS.md)** - Custom obfuscation per client

### API Documentation 🆕

- **[API Reference (Markdown)](./API_REFERENCE.md)** - Complete REST API reference
- **[OpenAPI Specification (YAML)](./openapi.yaml)** - Machine-readable API spec 🆕
- **[OpenAPI Guide](./OPENAPI_GUIDE.md)** - How to use OpenAPI documentation 🆕
- **[Swagger UI](http://localhost:51821/api/docs)** - Interactive API testing 🆕

### Code Examples 🆕

- **[Python Client](./examples/openapi_client.py)** - Full Python API client 🆕
- **[JavaScript Client](./examples/openapi_client.js)** - Node.js API client 🆕
- **[Examples README](./examples/README.md)** - Usage guide and examples 🆕

### Technical

- **[Architecture](./ARCHITECTURE.md)** - Project structure and components
- **[Changelog (Go)](./CHANGELOG_GO.md)** - Changes from Node.js version

## 🎯 By Use Case

### I want to...

#### **Use the Web UI**
→ See [Main README](../README.md#-quick-start)

#### **Use the REST API**
→ Start with [Swagger UI](http://localhost:51821/api/docs) 🆕
→ Or read [API Reference](./API_REFERENCE.md)

#### **Generate API clients**
→ Read [OpenAPI Guide](./OPENAPI_GUIDE.md) 🆕
→ Use [OpenAPI Spec](./openapi.yaml) 🆕

#### **Write integration scripts**
→ See [Python Example](./examples/openapi_client.py) 🆕
→ Or [JavaScript Example](./examples/openapi_client.js) 🆕

#### **Configure obfuscation**
→ Read [AmneziaWG Parameters](./AMNEZIAWG_PARAMETERS.md)
→ And [Per-Client Parameters](./PER_CLIENT_PARAMETERS.md)

#### **Deploy to production**
→ Read [Environment Variables](./ENVIRONMENT_VARIABLES.md)
→ And [Architecture](./ARCHITECTURE.md)

#### **Understand the codebase**
→ Read [Architecture](./ARCHITECTURE.md)
→ And [Changelog](./CHANGELOG_GO.md)

## 📊 Documentation Features

### Interactive API Testing 🆕

Access the **Swagger UI** for interactive API exploration:

```
http://localhost:51821/api/docs
```

**Features:**
- 📖 Browse all endpoints
- 🧪 Test API calls in browser
- 🔐 Authenticate directly
- 📝 View schemas and examples
- 📋 Generate code snippets

### OpenAPI Specification 🆕

Download or use the machine-readable spec:

```
http://localhost:51821/api/openapi.yaml
```

**Use cases:**
- Import into Postman/Insomnia
- Generate typed clients (Python, TypeScript, Go, etc.)
- Automated API testing
- Documentation generation

### Code Examples 🆕

Ready-to-use API clients:

**Python:**
```bash
pip install requests
python docs/examples/openapi_client.py
```

**JavaScript:**
```bash
npm install axios
node docs/examples/openapi_client.js
```

## 🗂️ File Structure

```
docs/
├── README.md                    # This file
├── ARCHITECTURE.md              # Project structure
├── AMNEZIAWG_PARAMETERS.md      # Obfuscation guide
├── API_REFERENCE.md             # REST API reference
├── CHANGELOG_GO.md              # Go version changes
├── ENVIRONMENT_VARIABLES.md     # Configuration
├── PASSWORD_GENERATION.md       # Password utility
├── PER_CLIENT_PARAMETERS.md     # Per-client settings
├── openapi.yaml                 # 🆕 OpenAPI 3.0 spec
├── OPENAPI_GUIDE.md             # 🆕 OpenAPI usage guide
└── examples/                    # 🆕 Code examples
    ├── README.md                # Examples guide
    ├── openapi_client.py        # Python client
    └── openapi_client.js        # JavaScript client
```

## 📈 What's New in This Release

### OpenAPI Documentation 🆕

Complete OpenAPI 3.0 specification with:
- All endpoints documented
- Request/response schemas
- Authentication methods
- Interactive Swagger UI
- Code generation support

### API Client Examples 🆕

Production-ready client implementations:
- **Python** - Full-featured async client
- **JavaScript** - Promise-based Node.js client
- Examples for all common operations
- Error handling and best practices

### Enhanced Guides 🆕

- [OpenAPI Guide](./OPENAPI_GUIDE.md) - Complete guide to using the OpenAPI spec
- [Examples README](./examples/README.md) - How to use and customize clients

## 🔗 External Resources

### OpenAPI Ecosystem

- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger Editor](https://editor.swagger.io/)
- [OpenAPI Generator](https://openapi-generator.tech/)

### API Testing Tools

- [Postman](https://www.postman.com/)
- [Insomnia](https://insomnia.rest/)
- [HTTPie](https://httpie.io/)
- [cURL](https://curl.se/)

### Client Generators

- [Python](https://openapi-generator.tech/docs/generators/python/)
- [TypeScript/Axios](https://openapi-generator.tech/docs/generators/typescript-axios/)
- [Go](https://openapi-generator.tech/docs/generators/go/)
- [All Generators](https://openapi-generator.tech/docs/generators/)

## 🤝 Contributing

Found an issue in the docs? Have a suggestion?

1. Check existing [issues](https://github.com/little-secrets/amnezia-wg-easy/issues)
2. Open a new issue with the `documentation` label
3. Or submit a PR with improvements

## 📝 Documentation Style

Our docs follow these principles:

- ✅ **Clear examples** for every feature
- ✅ **Copy-paste ready** code snippets
- ✅ **Multiple formats** (Interactive, Markdown, YAML)
- ✅ **Practical use cases** over theory
- ✅ **Up-to-date** with the codebase

## 🆘 Need Help?

- **Quick questions**: Check the [Main README](../README.md)
- **API questions**: Use [Swagger UI](http://localhost:51821/api/docs)
- **Bug reports**: Open an [issue](https://github.com/little-secrets/amnezia-wg-easy/issues)
- **Discussions**: Start a [discussion](https://github.com/little-secrets/amnezia-wg-easy/discussions)

---

**Documentation generated for AmneziaWG Easy v1.0.0**

🔗 [GitHub Repository](https://github.com/little-secrets/amnezia-wg-easy) | 📖 [API Docs](http://localhost:51821/api/docs) | ⭐ Star us on GitHub!

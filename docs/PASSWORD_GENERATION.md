# Password Generation with wgpw

The `wgpw` utility generates bcrypt password hashes for securing your AmneziaWG Easy installation.

## 🔐 Why Bcrypt?

- **Secure**: Industry-standard password hashing
- **Slow by design**: Resistant to brute-force attacks
- **Salted**: Each hash is unique even for same password
- **Cost factor 12**: Good balance of security and performance

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Generate hash for a password
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw mypassword

# Output:
PASSWORD_HASH='$2a$12$xELb112CO5ZgDqydj4SET.bxuHr2hcMb2SWgTlBU/XKSt8NEGjUge'
```

Copy the entire `PASSWORD_HASH='...'` line to your `.env` file.

### Interactive Mode

```bash
# Prompt for password (hidden input)
docker run --rm -it ghcr.io/w0rng/amnezia-wg-easy wgpw

# Enter your password: ****
# PASSWORD_HASH='$2a$12$...'
```

### Verify Password

```bash
# Check if password matches hash
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw mypassword '$2a$12$...'

# Output:
# Password matches the hash! ✅
# OR
# Password does not match the hash. ❌
```

## 💻 Using Locally (Go)

### Build wgpw

```bash
go build -o wgpw ./cmd/wgpw
```

### Generate Hash

```bash
# Direct password
./wgpw mypassword

# Interactive
./wgpw

# Verify
./wgpw mypassword '$2a$12$...'
```

### Windows

```bash
go build -o wgpw.exe ./cmd/wgpw
wgpw.exe mypassword
```

## 📖 Usage Examples

### 1. Generate Hash for Web UI

```bash
# Generate
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw MySecurePassword123

# Add to .env
echo "PASSWORD_HASH='$2a$12$xELb112CO5ZgDqydj4SET.bxuHr2hcMb2SWgTlBU/XKSt8NEGjUge'" >> .env
```

### 2. Generate Hash for Prometheus Metrics

```bash
# Generate different password for metrics
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw MetricsPassword456

# Add to .env
echo "PROMETHEUS_METRICS_PASSWORD='$2a$12$...' >> .env
```

### 3. Update Existing Password

```bash
# Generate new hash
NEW_HASH=$(docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw NewPassword789)

# Update .env
sed -i "s/^PASSWORD_HASH=.*/PASSWORD_HASH='$NEW_HASH'/" .env

# Restart container
docker compose restart
```

### 4. Batch Generate Multiple Passwords

```bash
#!/bin/bash
# generate-passwords.sh

echo "Generating passwords..."
echo ""

read -p "Enter Web UI password: " -s WEB_PASSWORD
echo ""
WEB_HASH=$(docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw "$WEB_PASSWORD")

read -p "Enter Metrics password: " -s METRICS_PASSWORD
echo ""
METRICS_HASH=$(docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw "$METRICS_PASSWORD")

cat > .env << EOF
WG_HOST=your.server.ip
$WEB_HASH
ENABLE_PROMETHEUS_METRICS=true
$METRICS_HASH
EOF

echo ""
echo "✅ Passwords generated and saved to .env"
```

## 🔧 Command Line Options

### Syntax

```bash
wgpw [PASSWORD] [HASH]
```

### Arguments

| Args | Mode | Description |
|------|------|-------------|
| None | Interactive | Prompts for password (hidden) |
| 1 arg | Generate | Generates hash for PASSWORD |
| 2 args | Verify | Compares PASSWORD with HASH |

### Examples

```bash
# Interactive mode
wgpw
# Enter your password: ****

# Generate mode
wgpw MyPassword123
# PASSWORD_HASH='$2a$12$...'

# Verify mode
wgpw MyPassword123 '$2a$12$...'
# Password matches the hash!
```

## 🛡️ Security Best Practices

### Password Strength

✅ **Good passwords:**
- `MySecure_VPN_Pass2024!` (20+ chars, mixed case, numbers, symbols)
- `correct-horse-battery-staple-2024` (passphrase)
- `Tr0ub4dor&3-Extended` (modified dictionary)

❌ **Bad passwords:**
- `password` (too common)
- `12345678` (sequential)
- `mypassword` (no complexity)
- `qwerty` (keyboard pattern)

### Password Management

**DO:**
- Use unique passwords for each service
- Store hashes securely (encrypted .env file)
- Use password manager for storing original passwords
- Change passwords periodically

**DON'T:**
- Reuse passwords across services
- Share passwords in plain text
- Store passwords in git repositories
- Use default/example passwords

### Hash Storage

```bash
# ✅ Good: Encrypted .env file
ansible-vault encrypt .env

# ✅ Good: Environment from secure vault
docker run -e PASSWORD_HASH=$(vault read secret/wg-password) ...

# ❌ Bad: Plain text in repository
git add .env  # DON'T DO THIS

# ❌ Bad: Hard-coded in docker-compose
environment:
  PASSWORD_HASH: '$2a$12$...'  # DON'T DO THIS
```

## 🧪 Testing

### Test Password Generation

```bash
# Generate hash
HASH=$(docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw testpass | grep PASSWORD_HASH | cut -d"'" -f2)

# Verify it works
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw testpass "$HASH"
# Should output: Password matches the hash!
```

### Test in Container

```bash
# Start container with password
docker compose up -d

# Try to login
curl -X POST http://localhost:51821/api/session \
  -H "Content-Type: application/json" \
  -d '{"password":"MyPassword123"}'

# Should return: {"success":true}
```

## 📊 Bcrypt Hash Format

### Structure

```
$2a$12$xELb112CO5ZgDqydj4SET.bxuHr2hcMb2SWgTlBU/XKSt8NEGjUge
│  │ │  └──────────────────────────────────────┬────────────────┘
│  │ │                                          │
│  │ └─ Cost factor (12 = 2^12 iterations)     └─ Hash (31 chars)
│  └─── Algorithm variant (2a)
└────── Bcrypt identifier

```

### Cost Factor

This implementation uses cost factor **12**:
- **2^12 = 4,096 iterations**
- **~0.3 seconds** to hash
- **Good balance** of security and UX

## 🔄 Migration from Node.js Version

If migrating from the Node.js version, passwords remain compatible:

```bash
# Old Node.js password hash
PASSWORD_HASH='$2y$05$Ci...'

# New Go password hash
PASSWORD_HASH='$2a$12$...'

# Both work! Bcrypt is cross-compatible
```

**Note:** Go version uses cost 12 vs Node.js cost 5. Consider regenerating for better security.

## 🐛 Troubleshooting

### Error: "Password does not match"

```bash
# Check hash format (should start with $2a$12$)
echo $PASSWORD_HASH

# Verify manually
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw yourpassword "$PASSWORD_HASH"

# Regenerate if needed
NEW_HASH=$(docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw yourpassword)
```

### Error: "Incorrect Password" in Web UI

1. Verify hash in .env file:
   ```bash
   cat .env | grep PASSWORD_HASH
   ```

2. Test hash with wgpw:
   ```bash
   docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw yourpassword "$HASH"
   ```

3. Check for special characters in .env:
   ```bash
   # Make sure hash is quoted
   PASSWORD_HASH='$2a$12$...'  # ✅ Good
   PASSWORD_HASH=$2a$12$...     # ❌ Bad (shell interprets $)
   ```

### Hash looks corrupted

```bash
# Regenerate clean hash
docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw newpassword > hash.txt

# View hash
cat hash.txt
# PASSWORD_HASH='$2a$12$...'

# Add to .env
cat hash.txt >> .env
```

## 🔗 Integration Examples

### Shell Script

```bash
#!/bin/bash
PASSWORD=${1:-$(openssl rand -base64 32)}
HASH=$(docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw "$PASSWORD")
echo "Generated password: $PASSWORD"
echo "$HASH"
```

### Ansible Playbook

```yaml
- name: Generate WireGuard password
  shell: docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw "{{ wg_password }}"
  register: wg_hash

- name: Create .env file
  template:
    src: env.j2
    dest: /opt/wg-easy/.env
  vars:
    password_hash: "{{ wg_hash.stdout }}"
```

### Terraform

```hcl
resource "random_password" "wg_password" {
  length  = 32
  special = true
}

resource "null_resource" "generate_hash" {
  provisioner "local-exec" {
    command = "docker run --rm ghcr.io/w0rng/amnezia-wg-easy wgpw '${random_password.wg_password.result}' > hash.txt"
  }
}
```

## 📚 Further Reading

- [Bcrypt Algorithm](https://en.wikipedia.org/wiki/Bcrypt)
- [Password Security Best Practices](https://www.owasp.org/index.php/Authentication_Cheat_Sheet)
- [Go Bcrypt Package](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

---

**Remember:** Always use strong, unique passwords and store hashes securely!


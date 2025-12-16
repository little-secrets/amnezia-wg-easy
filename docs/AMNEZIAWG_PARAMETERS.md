# AmneziaWG Obfuscation Parameters

AmneziaWG adds traffic obfuscation to WireGuard by injecting junk packets and modifying packet headers. This makes VPN traffic harder to detect and block.

## 🎯 Overview

When you don't set these parameters, they are **automatically randomized** for maximum security. You only need to set them manually if you want specific values.

## 📊 Parameters Reference

### Junk Packet Parameters

#### `JC` - Junk Packet Count
- **Default**: Random between 3-10
- **Range**: 0-128
- **Description**: Number of junk packets sent before the actual handshake
- **Purpose**: Makes DPI harder by adding noise to the connection start

**Example:**
```bash
JC=5
```

#### `JMIN` - Junk Packet Minimum Size
- **Default**: 50
- **Range**: 0-1500
- **Description**: Minimum size (in bytes) of junk packets
- **Purpose**: Ensures junk packets have a minimum size

**Example:**
```bash
JMIN=50
```

#### `JMAX` - Junk Packet Maximum Size
- **Default**: 1000
- **Range**: JMIN-1500
- **Description**: Maximum size (in bytes) of junk packets
- **Purpose**: Limits the size of junk packets to avoid MTU issues

**Example:**
```bash
JMAX=1000
```

### Packet Junk Size Parameters

#### `S1` - Init Packet Junk Size
- **Default**: Random between 15-150
- **Range**: 0-1500
- **Description**: Size of random data added to the init packet
- **Purpose**: Obfuscates the handshake initialization

**Example:**
```bash
S1=75
```

#### `S2` - Response Packet Junk Size
- **Default**: Random between 15-150
- **Range**: 0-1500
- **Description**: Size of random data added to the response packet
- **Purpose**: Obfuscates the handshake response

**Example:**
```bash
S2=75
```

### Magic Header Parameters

These headers are prepended to packets to make them look like other protocols.

#### `H1` - Init Packet Magic Header
- **Default**: Random uint32 (1 to 2,147,483,647)
- **Range**: 1 to 4,294,967,295
- **Description**: Magic header for the first handshake byte
- **Purpose**: Makes the init packet look like a different protocol

**Example:**
```bash
H1=1234567891
```

#### `H2` - Response Packet Magic Header
- **Default**: Random uint32
- **Range**: 1 to 4,294,967,295
- **Description**: Magic header for the handshake response
- **Purpose**: Makes the response packet look like a different protocol

**Example:**
```bash
H2=1234567892
```

#### `H3` - Underload Packet Magic Header
- **Default**: Random uint32
- **Range**: 1 to 4,294,967,295
- **Description**: Magic header for underload packets
- **Purpose**: Obfuscates underload signaling

**Example:**
```bash
H3=1234567893
```

#### `H4` - Transport Packet Magic Header
- **Default**: Random uint32
- **Range**: 1 to 4,294,967,295
- **Description**: Magic header for data transport packets
- **Purpose**: Makes data packets look like a different protocol

**Example:**
```bash
H4=1234567894
```

## 🔐 Security Recommendations

### 1. Use Auto-Random (Default)

**Best for most users:**
```bash
# Don't set any parameters - let them randomize
WG_HOST=192.168.1.1
PASSWORD_HASH='$2a$12$...'
```

**Why?**
- Maximum unpredictability
- Each server has unique fingerprint
- Harder for DPI to create signatures

### 2. Set Custom Values (Advanced)

**For specific requirements:**
```bash
# Custom obfuscation
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

**When to use?**
- You need consistent values across multiple servers
- Debugging connection issues
- Compliance requirements

### 3. Coordinate with Clients

**Important:** All clients connecting to your server must use the **exact same** obfuscation parameters. These are automatically included in the client configuration.

## 📈 Performance Impact

| Parameter | Impact on Performance | Impact on Obfuscation |
|-----------|----------------------|----------------------|
| `JC` | Low (only at connection start) | High |
| `JMIN/JMAX` | Low | Medium |
| `S1/S2` | Very Low | High |
| `H1-H4` | None | High |

### Bandwidth Overhead

Approximate additional traffic per connection:
- Handshake overhead: `JC * avg(JMIN, JMAX) + S1 + S2` bytes
- Ongoing overhead: Negligible (<1%)

**Example calculation:**
```
JC=5, JMIN=50, JMAX=1000, S1=75, S2=75
Handshake overhead ≈ 5 * 525 + 75 + 75 = 2,775 bytes
```

## 🔍 How It Works

### 1. Handshake Obfuscation

```
Client                    Server
  |                         |
  |--- [H1][S1][Junk]-----> | (Init packet with header + junk)
  |                         |
  |<-- [H2][S2][Junk] ------| (Response with header + junk)
  |                         |
  |--- [H3][Data] --------> | (Underload)
  |                         |
```

### 2. Transport Obfuscation

```
Every data packet is prepended with H4:
[H4][WireGuard Packet]
```

### 3. Junk Packet Injection

Before handshake:
```
[Junk1][Junk2][Junk3]...[JunkN][Real Handshake]
```

## 📝 Configuration Examples

### Maximum Stealth

```bash
# High junk count, variable sizes
JC=10
JMIN=100
JMAX=1400
S1=150
S2=150
H1=2147483647
H2=2147483646
H3=2147483645
H4=2147483644
```

### Performance Optimized

```bash
# Low overhead, still obfuscated
JC=3
JMIN=20
JMAX=200
S1=30
S2=30
H1=1000000
H2=1000001
H3=1000002
H4=1000003
```

### Balanced (Default-like)

```bash
# Good balance of stealth and performance
JC=5
JMIN=50
JMAX=1000
S1=75
S2=75
# Leave H1-H4 empty for random
```

## 🧪 Testing Your Configuration

### 1. Check Generated Values

```bash
# View server configuration
docker compose exec amnezia-wg-easy cat /etc/wireguard/wg0.json | jq .server

# Output:
{
  "privateKey": "...",
  "publicKey": "...",
  "address": "10.8.0.1",
  "jc": "7",
  "jmin": "50",
  "jmax": "1000",
  "s1": "89",
  "s2": "134",
  "h1": "1847362945",
  "h2": "492817364",
  "h3": "1638274562",
  "h4": "2014738291"
}
```

### 2. Verify Client Configuration

```bash
# Download a client config
curl http://localhost:51821/api/wireguard/client/:id/configuration
```

**Output should include:**
```ini
[Interface]
PrivateKey = ...
Address = 10.8.0.2/24
Jc = 7
Jmin = 50
Jmax = 1000
S1 = 89
S2 = 134
H1 = 1847362945
H2 = 492817364
H3 = 1638274562
H4 = 1638274562

[Peer]
PublicKey = ...
Endpoint = server.ip:51820
...
```

### 3. Monitor Connections

```bash
# Check WireGuard status
docker compose exec amnezia-wg-easy wg show

# View logs
docker compose logs -f
```

## ⚠️ Common Issues

### Issue: Client can't connect

**Cause:** Mismatched obfuscation parameters

**Solution:**
1. Regenerate client configuration
2. Ensure parameters match server
3. Check that client supports AmneziaWG

### Issue: Connection drops frequently

**Cause:** `JMAX` too large, exceeding MTU

**Solution:**
```bash
# Reduce JMAX
JMAX=800
# Or set MTU
WG_MTU=1420
```

### Issue: Slow connection speed

**Cause:** Too many junk packets

**Solution:**
```bash
# Reduce junk count
JC=3
JMIN=20
JMAX=200
```

## 📚 Further Reading

- [AmneziaVPN Documentation](https://docs.amnezia.org/)
- [WireGuard Protocol](https://www.wireguard.com/protocol/)
- [Deep Packet Inspection](https://en.wikipedia.org/wiki/Deep_packet_inspection)

## 🤔 FAQ

**Q: Do I need to set these parameters?**
A: No! They auto-randomize for maximum security.

**Q: Can I change them after clients are created?**
A: Yes, but you must regenerate all client configs.

**Q: Are random values secure?**
A: Yes! Random values are more secure than predictable ones.

**Q: Do parameters affect speed?**
A: Minimal impact. Only handshake has noticeable overhead.

**Q: Can DPI still detect my VPN?**
A: AmneziaWG makes detection much harder, but not impossible.

---

**Remember:** When in doubt, leave parameters unset for automatic randomization!


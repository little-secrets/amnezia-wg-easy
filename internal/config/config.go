package config

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// secureRandUint32 returns a uniform 32-bit unsigned int from the OS
// entropy source. Audit H-6: the previous implementation used
// math/rand seeded from time.Now().UnixNano(), so an attacker who
// knew the deployment's startup time (logs, TLS cert NotBefore, DNS
// cache, SSH banner) could reconstruct the AmneziaWG obfuscation
// parameters and undermine the traffic-classification resistance the
// scheme is supposed to provide.
//
// Panics on CSPRNG failure -- continuing with a known-bad value would
// silently weaken every connection negotiated after that point.
func secureRandUint32() uint32 {
	var b [4]byte
	if _, err := cryptorand.Read(b[:]); err != nil {
		panic("crypto/rand failed: " + err.Error())
	}
	return binary.BigEndian.Uint32(b[:])
}

// secureRandIntInRange returns a uniformly distributed int in [lo, hi].
func secureRandIntInRange(lo, hi int) int {
	if hi <= lo {
		return lo
	}
	span := uint32(hi - lo + 1)
	return lo + int(secureRandUint32()%span)
}

type Config struct {
	// Server
	Port      string
	WebUIHost string
	Release   string

	// Auth
	PasswordHash              string
	NoAuth                    bool // explicit opt-in for unauthenticated mode
	MaxAge                    int
	PrometheusMetricsPassword string

	// WireGuard
	WGHost               string
	WGPort               string
	WGConfigPort         string
	WGPath               string
	WGDevice             string
	WGDefaultAddress     string
	WGAddressCIDR        int
	WGDefaultDNS         string
	WGAllowedIPs         string
	WGMTU                string
	WGPersistentKeepalive string
	WGPreUp              string
	WGPostUp             string
	WGPreDown            string
	WGPostDown           string

	// Feature flags
	NoWebUI                  bool
	EnableOneTimeLinks       bool
	EnableExpiresTime        bool
	EnablePrometheusMetrics  bool
	EnableTrafficStats       bool
	EnableSortClients        bool

	// UI
	Lang         string
	ChartType    int
	DicebearType string
	UseGravatar  bool

	// AmneziaWG obfuscation
	Jc   int
	Jmin int
	Jmax int
	S1   int
	S2   int
	H1   uint32
	H2   uint32
	H3   uint32
	H4   uint32
}

func (c *Config) RequiresPassword() bool {
	return c.PasswordHash != ""
}

func (c *Config) RequiresPrometheusPassword() bool {
	return c.PrometheusMetricsPassword != ""
}

func Load() *Config {
	wgPort := envStr("WG_PORT", "51820")
	wgDevice := envStr("WG_DEVICE", "eth0")
	wgDefaultAddress := envStr("WG_DEFAULT_ADDRESS", "10.8.0.x")
	wgPort_ := envStr("WG_CONFIG_PORT", wgPort)

	// Parse CIDR from default address (e.g. "10.8.0.x" -> 24)
	cidr := parseCIDR(wgDefaultAddress)

	// Build default PostUp/PostDown using device and port
	subnet := strings.Replace(wgDefaultAddress, "x", "0", 1) + "/" + strconv.Itoa(cidr)
	defaultPostUp := fmt.Sprintf(
		"iptables -t nat -A POSTROUTING -s %s -o %s -j MASQUERADE; "+
			"iptables -A INPUT -p udp -m udp --dport %s -j ACCEPT; "+
			"iptables -A FORWARD -i wg0 -j ACCEPT; "+
			"iptables -A FORWARD -o wg0 -j ACCEPT;",
		subnet, wgDevice, wgPort,
	)
	defaultPostDown := fmt.Sprintf(
		"iptables -t nat -D POSTROUTING -s %s -o %s -j MASQUERADE; "+
			"iptables -D INPUT -p udp -m udp --dport %s -j ACCEPT; "+
			"iptables -D FORWARD -i wg0 -j ACCEPT; "+
			"iptables -D FORWARD -o wg0 -j ACCEPT;",
		subnet, wgDevice, wgPort,
	)

	return &Config{
		// Server
		Port:      envStr("PORT", "51821"),
		WebUIHost: envStr("WEBUI_HOST", "0.0.0.0"),
		Release:   envStr("RELEASE", "1.0.0"),

		// Auth
		PasswordHash:              envStr("PASSWORD_HASH", ""),
		NoAuth:                    envBool("NO_AUTH", false),
		MaxAge:                    envInt("MAX_AGE", 0),
		PrometheusMetricsPassword: envStr("PROMETHEUS_METRICS_PASSWORD", ""),

		// WireGuard
		WGHost:                wgHost(),
		WGPort:                wgPort,
		WGConfigPort:          wgPort_,
		WGPath:                envStr("WG_PATH", "/etc/wireguard/"),
		WGDevice:              wgDevice,
		WGDefaultAddress:      wgDefaultAddress,
		WGAddressCIDR:         cidr,
		WGDefaultDNS:          envStr("WG_DEFAULT_DNS", "1.1.1.1"),
		WGAllowedIPs:          envStr("WG_ALLOWED_IPS", "0.0.0.0/0, ::/0"),
		WGMTU:                 envStr("WG_MTU", ""),
		WGPersistentKeepalive: envStr("WG_PERSISTENT_KEEPALIVE", "0"),
		WGPreUp:               envStr("WG_PRE_UP", ""),
		WGPostUp:              envStr("WG_POST_UP", defaultPostUp),
		WGPreDown:             envStr("WG_PRE_DOWN", ""),
		WGPostDown:            envStr("WG_POST_DOWN", defaultPostDown),

		// Feature flags
		NoWebUI:                 envBool("NO_WEB_UI", false),
		EnableOneTimeLinks:      envBool("WG_ENABLE_ONE_TIME_LINKS", false),
		EnableExpiresTime:       envBool("WG_ENABLE_EXPIRES_TIME", false),
		EnablePrometheusMetrics: envBool("ENABLE_PROMETHEUS_METRICS", false),
		EnableTrafficStats:      envBool("UI_TRAFFIC_STATS", false),
		EnableSortClients:       envBool("UI_ENABLE_SORT_CLIENTS", false),

		// UI
		Lang:         envStr("LANG", "en"),
		ChartType:    envInt("UI_CHART_TYPE", 0),
		DicebearType: envStr("DICEBEAR_TYPE", ""),
		UseGravatar:  envBool("USE_GRAVATAR", false),

		// AmneziaWG obfuscation -- generated from crypto/rand so they
		// are not reconstructible from the deployment's startup time.
		Jc:   envInt("JC", secureRandIntInRange(3, 10)),
		Jmin: envInt("JMIN", 50),
		Jmax: envInt("JMAX", 1000),
		S1:   envInt("S1", secureRandIntInRange(15, 150)),
		S2:   envInt("S2", secureRandIntInRange(15, 150)),
		H1:   envUint32("H1", secureRandUint32()|1),
		H2:   envUint32("H2", secureRandUint32()|1),
		H3:   envUint32("H3", secureRandUint32()|1),
		H4:   envUint32("H4", secureRandUint32()|1),
	}
}

func wgHost() string {
	return os.Getenv("WG_HOST")
}

func envStr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func envUint32(key string, fallback uint32) uint32 {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.ParseUint(v, 10, 32); err == nil {
			return uint32(i)
		}
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		v = strings.ToLower(v)
		return v == "true" || v == "1" || v == "yes"
	}
	return fallback
}

// parseCIDR guesses CIDR from address pattern like "10.8.0.x"
func parseCIDR(address string) int {
	parts := strings.Split(strings.Split(address, "/")[0], ".")
	if len(parts) != 4 {
		return 24
	}
	// Count non-"x" octets for subnet mask
	fixed := 0
	for _, p := range parts {
		if p != "x" {
			fixed++
		}
	}
	return fixed * 8
}

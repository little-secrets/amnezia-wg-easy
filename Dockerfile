# Build stage for Go binary
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the main application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o amnezia-wg-easy .

# Build the wgpw tool
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o wgpw ./cmd/wgpw

# Final image
FROM amneziavpn/amnezia-wg:latest

HEALTHCHECK CMD /usr/bin/timeout 5s /bin/sh -c "/usr/bin/wg show | /bin/grep -q interface || exit 1" --interval=1m --timeout=5s --retries=3

# Install Linux packages
RUN apk add --no-cache \
    dpkg \
    dumb-init \
    iptables

# Use iptables-legacy
RUN update-alternatives --install /sbin/iptables iptables /sbin/iptables-legacy 10 \
    --slave /sbin/iptables-restore iptables-restore /sbin/iptables-legacy-restore \
    --slave /sbin/iptables-save iptables-save /sbin/iptables-legacy-save

# Copy binaries from builder
COPY --from=builder /build/amnezia-wg-easy /app/amnezia-wg-easy
COPY --from=builder /build/wgpw /bin/wgpw

# Copy Web UI files
COPY www /app/www

# Create WireGuard config directory
RUN mkdir -p /etc/wireguard

# Set working directory
WORKDIR /app

# Expose ports
EXPOSE 51820/udp
EXPOSE 51821/tcp

# Run application
CMD ["/usr/bin/dumb-init", "/app/amnezia-wg-easy"]

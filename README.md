# SQRL Redis Hoard

[![CI/CD Pipeline](https://github.com/sqrldev/server-go-ssp-redishoard/actions/workflows/ci.yml/badge.svg)](https://github.com/sqrldev/server-go-ssp-redishoard/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sqrldev/server-go-ssp-redishoard)](https://goreportcard.com/report/github.com/sqrldev/server-go-ssp-redishoard)
[![GoDoc](https://godoc.org/github.com/sqrldev/server-go-ssp-redishoard?status.svg)](https://godoc.org/github.com/sqrldev/server-go-ssp-redishoard)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A secure, production-ready Redis-backed implementation of the `ssp.Hoard` interface for the SQRL (Secure Quick Reliable Login) authentication protocol.

## Features

- **Secure Memory Handling**: Implements secure memory clearing to prevent sensitive data leakage (CWE-226 mitigation)
- **No Sensitive Data Logging**: Removes all sensitive information from logs (CWE-200 mitigation)
- **Input Validation**: Validates cryptographic nonces to prevent injection attacks
- **Platform-Aware**: Memory locking support for Unix/Linux, macOS, and Windows
- **Atomic Operations**: Uses Redis transactions for race-condition-free operations
- **Go 1.23+**: Built with modern Go features and security patches
- **Redis 7.2+**: Supports latest Redis features including RESP3

## Installation

```bash
go get github.com/sqrldev/server-go-ssp-redishoard
```

## Requirements

- Go 1.23 or later
- Redis 7.2 or later
- [github.com/sqrldev/server-go-ssp](https://github.com/sqrldev/server-go-ssp)

## Quick Start

```go
package main

import (
    "time"

    "github.com/redis/go-redis/v9"
    ssp "github.com/sqrldev/server-go-ssp"
    "github.com/sqrldev/server-go-ssp-redishoard"
)

func main() {
    // Create Redis client
    client := redis.NewUniversalClient(&redis.UniversalOptions{
        Addrs: []string{"localhost:6379"},
        // For production, enable TLS:
        // TLSConfig: &tls.Config{...},
    })
    defer client.Close()

    // Create the Hoard
    hoard := redishoard.NewHoard(client)

    // Use with SQRL SSP API
    nut := ssp.Nut("your-cryptographic-nonce")
    cache := &ssp.HoardCache{
        State: "authentication-state",
    }

    // Save with 5-minute TTL
    err := hoard.Save(nut, cache, 5*time.Minute)
    if err != nil {
        panic(err)
    }

    // Retrieve
    retrieved, err := hoard.Get(nut)
    if err != nil {
        panic(err)
    }

    // Atomic get-and-delete (for one-time use)
    value, err := hoard.GetAndDelete(nut)
    if err != nil {
        panic(err)
    }
}
```

## Security Features

### Secure Memory Clearing

All sensitive data (JSON serialized authentication state, cryptographic material) is securely cleared from memory after use:

```go
// Automatic clearing in all operations
data, err := client.Get(ctx, key).Bytes()
defer redishoard.ClearBytes(data) // Securely clears memory
```

### Platform-Aware Memory Locking

Prevents sensitive data from being swapped to disk:

```go
// Unix/Linux/macOS
unix.Mlock(sensitiveData)

// Windows
windows.VirtualLock(sensitiveData)
```

### Input Validation

All nut values are validated before use:

```go
err := redishoard.ValidateNut(nut)
// Checks:
// - Length: 8-64 characters
// - Characters: base64url safe only (A-Z, a-z, 0-9, -, _)
```

## API Reference

### Hoard

```go
type Hoard struct {
    client redis.UniversalClient
}

// NewHoard creates a Redis-backed Hoard
func NewHoard(client redis.UniversalClient) *Hoard

// Get retrieves cached state by nut
func (h *Hoard) Get(nut ssp.Nut) (*ssp.HoardCache, error)

// GetAndDelete atomically retrieves and deletes state
func (h *Hoard) GetAndDelete(nut ssp.Nut) (*ssp.HoardCache, error)

// Save stores state with TTL
func (h *Hoard) Save(nut ssp.Nut, value *ssp.HoardCache, expiration time.Duration) error
```

### Security Functions

```go
// ClearBytes securely zeros a byte slice
func ClearBytes(b []byte)

// ClearString securely clears a string
func ClearString(s *string)

// ValidateNut validates nut format
func ValidateNut(nut string) error

// SecureBuffer for sensitive operations
type SecureBuffer struct { ... }
func NewSecureBuffer(size int, lock bool) (*SecureBuffer, error)
```

## Testing

```bash
# Unit tests (no Redis required)
go test -short ./...

# Integration tests (requires local Redis)
go test -v ./...

# Race condition detection
go test -race ./...

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Running Redis for Tests

```bash
# Docker
docker run -d --name redis-test -p 6379:6379 redis:7.4

# Or install locally
# Ubuntu: sudo apt install redis-server
# macOS: brew install redis
```

## CI/CD Pipeline

This project uses GitHub Actions for continuous integration:

- **Build & Test**: Go 1.23/1.24 with Redis 7.2/7.4
- **Linting**: golangci-lint, go vet, staticcheck
- **Security Scanning**: gosec, govulncheck
- **Dependency Review**: Automatic PR dependency analysis
- **Coverage**: Codecov integration

## Configuration

### Environment Variables

```bash
REDIS_URL=redis://localhost:6379
REDIS_TLS=true
REDIS_PASSWORD=your-password
```

### Redis Client Options

```go
client := redis.NewUniversalClient(&redis.UniversalOptions{
    Addrs:    []string{"localhost:6379"},
    Password: os.Getenv("REDIS_PASSWORD"),
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
    },
    PoolSize:     10,
    MinIdleConns: 5,
    DialTimeout:  5 * time.Second,
    ReadTimeout:  3 * time.Second,
    WriteTimeout: 3 * time.Second,
})
```

## Security Considerations

1. **Enable TLS** for Redis connections in production
2. **Set strong passwords** for Redis authentication
3. **Configure memory limits** to prevent DoS
4. **Monitor logs** for security events (no sensitive data logged)
5. **Regular dependency updates** for security patches

## Compliance

- **CWE-200**: Exposure of Sensitive Information ✅ Fixed
- **CWE-226**: Sensitive Information in Resource Not Removed ✅ Fixed
- **CWE-312**: Cleartext Storage of Sensitive Information ⚠️ Documented
- **OWASP Top 10**: A02:2021 Cryptographic Failures ✅ Mitigated

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass: `go test -race ./...`
5. Run linters: `golangci-lint run`
6. Submit a pull request

## Documentation

- [Requirements](REQUIREMENTS.md)
- [Security Plan](SECURITY_PLAN.md)
- [Upgrade Path](UPGRADE_PATH.md)
- [SQRL SSP API](https://github.com/sqrldev/server-go-ssp)

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Projects

- [server-go-ssp](https://github.com/sqrldev/server-go-ssp) - SQRL Server Side Protocol
- [server-go-ssp-gormauthstore](https://github.com/sqrldev/server-go-ssp-gormauthstore) - GORM AuthStore
- [SQRL Specification](https://www.grc.com/sqrl/sqrl.htm) - GRC SQRL Documentation

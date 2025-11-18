# Dependency Upgrade Path

## Current State (Pre-Upgrade)

| Dependency | Current Import Path | Current Version | Status |
|------------|---------------------|-----------------|--------|
| go-redis | `github.com/go-redis/redis` | Unversioned (v6/v7) | DEPRECATED |
| server-go-ssp | `github.com/sqrldev/server-go-ssp` | Unversioned | OUTDATED |
| Go Module | Missing | N/A | MISSING |

## Target State (Post-Upgrade)

| Dependency | New Import Path | Target Version | Release Date | Notes |
|------------|-----------------|----------------|--------------|-------|
| go-redis | `github.com/redis/go-redis/v9` | **v9.16.0** | Oct 23, 2025 | Latest stable |
| server-go-ssp | `github.com/sqrldev/server-go-ssp` | **v0.0.0-20241212182118** | Dec 12, 2024 | Uses dxcSithLord fork via replace |
| golang.org/x/sys | `golang.org/x/sys` | **v0.28.0** | Nov 2025 | Platform APIs |
| Go | N/A | **1.23.x** | 2024 | Current stable |

> **Development Note**: Currently using `github.com/dxcSithLord/server-go-ssp` fork via go.mod replace directive. This allows development work to proceed while maintaining import path compatibility. The fork will be merged back to sqrldev/server-go-ssp when ready.

---

## 1. go-redis Migration

### Version History & Breaking Changes

| Version | Import Path | Key Changes |
|---------|-------------|-------------|
| v6.x | `github.com/go-redis/redis` | Original unversioned |
| v7.x | `github.com/go-redis/redis/v7` | First versioned release |
| v8.x | `github.com/go-redis/redis/v8` | Context-first APIs, OpenTelemetry |
| v9.0 | `github.com/go-redis/redis/v9` | RESP3 support, connection pooling improvements |
| v9.x | `github.com/redis/go-redis/v9` | **Repository moved to redis org** |

### Migration Steps

1. **Update Import Path**
   ```go
   // OLD (DEPRECATED)
   import "github.com/go-redis/redis"

   // NEW (CURRENT)
   import "github.com/redis/go-redis/v9"
   ```

2. **API Changes from v6/v7 to v9**

   | Operation | Old API (v6/v7) | New API (v9) |
   |-----------|-----------------|--------------|
   | Get | `client.Get(key).Bytes()` | `client.Get(ctx, key).Bytes()` |
   | Set | `client.Set(key, val, exp).Err()` | `client.Set(ctx, key, val, exp).Err()` |
   | Del | `pipe.Del(key)` | `pipe.Del(ctx, key)` |
   | TxPipelined | `client.TxPipelined(func(pipe Pipeliner))` | `client.TxPipelined(ctx, func(pipe Pipeliner))` |
   | Error Check | `err == redis.Nil` | `err == redis.Nil` (unchanged) |

3. **Context Requirement**
   - All operations now require `context.Context` as first parameter
   - Use `context.Background()` for simple cases
   - Use `context.WithTimeout()` for deadline enforcement

4. **Code Changes Required**

   ```go
   // OLD
   func (h *Hoard) Get(nut ssp.Nut) (*ssp.HoardCache, error) {
       data, err := h.client.Get(string(nut)).Bytes()
       // ...
   }

   // NEW
   func (h *Hoard) Get(nut ssp.Nut) (*ssp.HoardCache, error) {
       ctx := context.Background()
       data, err := h.client.Get(ctx, string(nut)).Bytes()
       // ...
   }
   ```

### v9.16.0 Features

- **RESP3 Push Notifications** - Real-time server events
- **Hitless Upgrades** - Zero-downtime cluster upgrades
- **Trace Filtering** - Selective tracing for observability
- **256KiB Buffers** - Improved default buffer sizes
- **Redis 8.2 Support** - Latest Redis server compatibility

### Compatibility Matrix

| go-redis Version | Redis Server | Go Version |
|------------------|--------------|------------|
| v9.16.0 | 7.2, 7.4, 8.0, 8.2 | 1.23, 1.24 |
| v9.0.0 | 6.2+ | 1.18+ |
| v8.x | 6.0+ | 1.17+ |
| v7.x | 5.0+ | 1.13+ |

---

## 2. server-go-ssp Version

### Current Package Information

- **Latest Version:** `v0.0.0-20241212182118-c8230b16b87d`
- **Published:** December 12, 2024
- **License:** MIT
- **Status:** Pre-release (v0.x)

### Key Types Used

```go
// Nut - Cryptographic nonce (string alias)
type Nut string

// Hoard - Interface this package implements
type Hoard interface {
    Get(nut Nut) (*HoardCache, error)
    GetAndDelete(nut Nut) (*HoardCache, error)
    Save(nut Nut, value *HoardCache, expiration time.Duration) error
}

// HoardCache - State object stored in Redis
type HoardCache struct {
    RemoteIP     string
    OriginalNut  Nut
    PagNut       Nut
    LastRequest  CliRequest
    Identity     interface{}
    LastResponse CliResponse
}

// ErrNotFound - Standard error for missing nut
var ErrNotFound = errors.New("nut not found")
```

### Upgrade Considerations

- **No breaking changes** expected (same interface)
- Pseudo-version based on git commit
- Monitor for future v1.0 release with potential API changes

---

## 3. golang.org/x/sys (New Dependency)

### Purpose
Platform-specific system calls for secure memory operations:
- `unix.Mlock()` - Lock memory pages (prevent swap)
- `unix.Munlock()` - Unlock memory pages
- `windows.VirtualLock()` - Windows equivalent

### Version
- **Target:** v0.27.0 (November 2025)
- **License:** BSD-3-Clause
- **Go Version:** 1.18+

---

## 4. Go Version Requirements

### Minimum Version: Go 1.23

**Rationale:**
- Security patches current
- Module system mature
- Context package improvements
- Better build tag support (`//go:build`)
- Enhanced race detector
- Official redis client support

### Upgrade Path
```bash
# Check current version
go version

# Download latest Go 1.23
wget https://go.dev/dl/go1.23.3.linux-amd64.tar.gz

# Or use package manager
sudo apt update && sudo apt install golang-go
```

---

## 5. Implementation Steps

### Step 1: Initialize Go Module
```bash
go mod init github.com/sqrldev/server-go-ssp-redishoard
```

### Step 2: Update Import Statements
```go
import (
    "context"  // NEW - for v9 API
    "encoding/json"
    "fmt"
    "time"

    "github.com/redis/go-redis/v9"  // UPDATED
    ssp "github.com/sqrldev/server-go-ssp"
)
```

### Step 3: Add Context to All Redis Operations
```go
func (h *Hoard) Get(nut ssp.Nut) (*ssp.HoardCache, error) {
    ctx := context.Background()
    data, err := h.client.Get(ctx, string(nut)).Bytes()
    // ...
}
```

### Step 4: Fetch Dependencies
```bash
go mod tidy
```

### Step 5: Verify Versions
```bash
go list -m all
```

### Expected Output:
```
github.com/sqrldev/server-go-ssp-redishoard
github.com/redis/go-redis/v9 v9.16.0
github.com/sqrldev/server-go-ssp v0.0.0-20241212182118-c8230b16b87d
golang.org/x/sys v0.27.0
```

---

## 6. Verification Checklist

- [ ] All imports updated to new paths
- [ ] Context added to all Redis operations
- [ ] `go.mod` file created with correct module path
- [ ] `go.sum` generated with dependency checksums
- [ ] `go mod tidy` runs without errors
- [ ] `go build` compiles successfully
- [ ] `go test` passes all tests
- [ ] `go vet` reports no issues
- [ ] No deprecated API warnings

---

## 7. Rollback Plan

If upgrade causes issues:

1. **Revert imports** to old paths
2. **Remove** `go.mod` and `go.sum`
3. **Pin** to older version temporarily
4. **Document** incompatibility

---

## 8. Post-Upgrade Validation

```bash
# Security scan
govulncheck ./...

# Static analysis
staticcheck ./...

# Test with race detector
go test -race ./...

# Dependency audit
go mod verify
```

---

## 9. Future Considerations

### Upcoming Releases to Monitor

1. **go-redis v10** - Expected major version with potential API changes
2. **server-go-ssp v1.0** - Stable release with finalized API
3. **Go 1.24** - New language features and performance improvements

### Deprecation Warnings

- `github.com/go-redis/redis` namespace is deprecated
- Move to `github.com/redis/go-redis` completed
- Old import paths may stop working in future Go versions

---

## 10. References

- **go-redis Migration:** https://github.com/redis/go-redis
- **go-redis Releases:** https://github.com/redis/go-redis/releases
- **SQRL SSP Package:** https://pkg.go.dev/github.com/sqrldev/server-go-ssp
- **Go Modules:** https://go.dev/ref/mod
- **Redis Server:** https://redis.io/downloads/

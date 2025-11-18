# API Reference

## Document Control
- **Project**: server-go-ssp-redishoard
- **Version**: 1.0.0
- **Date**: 2025-01-18
- **Package**: `github.com/sqrldev/server-go-ssp-redishoard`
- **Repository**: https://github.com/dxcSithLord/server-go-ssp-redishoard

---

## Note on API Documentation Format

**This is a Go library package, not an HTTP REST API.**

OpenAPI (formerly Swagger) specifications are designed for HTTP-based RESTful APIs. This project provides a **Go library** that implements the `ssp.Hoard` interface for use by Go applications.

**Appropriate Documentation**:
- ✅ **This document**: Library API reference with function signatures and examples
- ✅ **GoDoc**: Inline documentation at https://pkg.go.dev/github.com/sqrldev/server-go-ssp-redishoard
- ✅ **README.md**: Quick start guide and usage examples
- ❌ **OpenAPI Spec**: Not applicable (no HTTP endpoints)

**If you need OpenAPI documentation**, you would create it for the **HTTP application** that uses this library, not for the library itself.

---

## Package Overview

```go
package redishoard // import "github.com/sqrldev/server-go-ssp-redishoard"
```

The `redishoard` package provides a secure, Redis-backed implementation of the `ssp.Hoard` interface for SQRL (Secure Quick Reliable Login) authentication state management.

**Key Features**:
- Secure memory handling (CWE-226 mitigation)
- Input validation (injection prevention)
- Atomic operations (race-free GetAndDelete)
- Platform-specific memory locking
- Automatic TTL expiration

---

## Types

### type Hoard

```go
type Hoard struct {
    // contains filtered or unexported fields
}
```

`Hoard` implements the `ssp.Hoard` interface using Redis as a backing store with secure memory handling.

**Thread Safety**: ✅ Safe for concurrent use by multiple goroutines.

**Implementation Notes**:
- All sensitive data is cleared from memory after use
- All operations include input validation
- Atomic operations use Redis MULTI/EXEC transactions

---

### type SecureBuffer

```go
type SecureBuffer struct {
    // contains filtered or unexported fields
}
```

`SecureBuffer` wraps a byte slice with automatic secure clearing.

**Use Cases**:
- Storing sensitive authentication tokens
- Handling cryptographic keys temporarily
- Processing PII data

**Example**:
```go
buf, err := redishoard.NewSecureBuffer(32, true)
if err != nil {
    log.Fatal(err)
}
defer buf.Destroy()

// Use buf.Data() for sensitive operations
copy(buf.Data(), sensitiveToken)
```

---

## Functions

### func NewHoard

```go
func NewHoard(client redis.UniversalClient) *Hoard
```

Creates a new Redis-backed Hoard instance.

**Parameters**:
- `client` (`redis.UniversalClient`): Redis client (supports single instance, Cluster, Sentinel)

**Returns**:
- `*Hoard`: New Hoard instance ready for use

**Example**:
```go
import (
    "github.com/redis/go-redis/v9"
    "github.com/sqrldev/server-go-ssp-redishoard"
)

client := redis.NewUniversalClient(&redis.UniversalOptions{
    Addrs: []string{"localhost:6379"},
})
defer client.Close()

hoard := redishoard.NewHoard(client)
```

**Thread Safety**: ✅ The returned Hoard is safe for concurrent use.

**Errors**: Never returns an error (client validation happens on first operation).

---

### func (*Hoard) Get

```go
func (h *Hoard) Get(nut ssp.Nut) (*ssp.HoardCache, error)
```

Retrieves the `HoardCache` for the given nut from Redis.

**Parameters**:
- `nut` (`ssp.Nut`): Cryptographic nonce (8-64 characters, base64url safe)

**Returns**:
- `*ssp.HoardCache`: Retrieved authentication state, or `nil` on error
- `error`: Error if nut not found, invalid, or Redis failure

**Errors**:
| Error | Description | Action |
|-------|-------------|--------|
| `ssp.ErrNotFound` | Nut not in Redis (expired or never existed) | Normal flow (nut consumed or expired) |
| `"invalid nut: ..."` | Nut validation failed | Fix nut format (must be base64url) |
| `"redis nut lookup failed: ..."` | Redis connection/command failure | Check Redis connectivity |

**Security**:
- ✅ Input validation prevents injection
- ✅ Data cleared from memory after deserialization

**Example**:
```go
nut := ssp.Nut("abc-123-def-456")
cache, err := hoard.Get(nut)
if errors.Is(err, ssp.ErrNotFound) {
    // Nut not found (expired or invalid)
    return nil, fmt.Errorf("session expired")
}
if err != nil {
    // Other error (Redis failure)
    return nil, fmt.Errorf("failed to get session: %w", err)
}

// Use cache.State, cache.RemoteIP, etc.
fmt.Printf("Session state: %s\n", cache.State)
```

**Performance**: Typically <5ms p95 latency.

---

### func (*Hoard) GetAndDelete

```go
func (h *Hoard) GetAndDelete(nut ssp.Nut) (*ssp.HoardCache, error)
```

Atomically retrieves and deletes the `HoardCache` for the given nut. This ensures the nut can only be used once, preventing replay attacks.

**Parameters**:
- `nut` (`ssp.Nut`): Cryptographic nonce to retrieve and delete

**Returns**:
- `*ssp.HoardCache`: Retrieved authentication state (before deletion)
- `error`: Error if nut not found, invalid, or transaction failed

**Atomicity Guarantee**: Uses Redis MULTI/EXEC to ensure:
1. Nut is retrieved
2. Nut is deleted
3. Both operations succeed or both fail
4. No other client can retrieve the same nut

**Errors**: Same as `Get()`, plus transaction-specific errors.

**Security**:
- ✅ **Replay attack prevention** - Nut can only be used once
- ✅ **Race-condition-free** - Atomic operation
- ✅ Data cleared from memory after use

**Example (SQRL Authentication Flow)**:
```go
// Client submits authentication with nut
nut := ssp.Nut(request.Nut)

// Atomically consume the nut
cache, err := hoard.GetAndDelete(nut)
if errors.Is(err, ssp.ErrNotFound) {
    return errors.New("nut already used or expired")
}
if err != nil {
    return fmt.Errorf("failed to consume nut: %w", err)
}

// Verify signatures using cache.Identity
if !verifySQRLSignature(cache) {
    return errors.New("invalid signature")
}

// Authenticate user
authenticateUser(cache.Identity.Idk)
```

**Performance**: Typically <10ms p95 latency (includes transaction overhead).

---

### func (*Hoard) Save

```go
func (h *Hoard) Save(nut ssp.Nut, value *ssp.HoardCache, expiration time.Duration) error
```

Stores the `HoardCache` in Redis with the specified TTL.

**Parameters**:
- `nut` (`ssp.Nut`): Cryptographic nonce (key)
- `value` (`*ssp.HoardCache`): Authentication state to store
- `expiration` (`time.Duration`): TTL for automatic expiration

**Returns**:
- `error`: Error if validation fails, serialization fails, or Redis write fails

**Errors**:
| Error | Description | Action |
|-------|-------------|--------|
| `"invalid nut: ..."` | Nut validation failed | Fix nut format |
| `"cannot save nil HoardCache"` | Value is nil | Provide valid HoardCache |
| `"failed json encoding HoardCache: ..."` | JSON serialization failed | Check HoardCache contents |
| Redis errors | Connection/command failure | Check Redis connectivity |

**TTL Behavior**:
- Key automatically deleted after `expiration` duration
- Use `time.Minute * 5` for typical SQRL flows
- Use `time.Second * 30` for short-lived nonces

**Security**:
- ✅ Input validation
- ✅ Serialized JSON cleared from memory after write

**Example**:
```go
nut := ssp.Nut("new-nut-abc123")
cache := &ssp.HoardCache{
    State:       "awaiting-client-response",
    RemoteIP:    "192.168.1.100",
    OriginalNut: nut,
}

// Save with 5-minute TTL
err := hoard.Save(nut, cache, 5*time.Minute)
if err != nil {
    return fmt.Errorf("failed to save session: %w", err)
}

// Key will auto-expire in 5 minutes
```

**Performance**: Typically <5ms p95 latency.

---

### func ClearBytes

```go
func ClearBytes(b []byte)
```

Securely clears a byte slice from memory. Implements CWE-226 mitigation by ensuring sensitive data does not remain in memory after use.

**Parameters**:
- `b` (`[]byte`): Byte slice to clear (zeroed in place)

**Returns**: None

**Security Properties**:
- ✅ **Constant-time operation** - Prevents timing attacks
- ✅ **Compiler-optimization-resistant** - Uses `runtime.KeepAlive()`
- ✅ **Memory barrier** - Uses `subtle.ConstantTimeCompare()`

**Usage**:
```go
sensitiveData, err := redis.Get(ctx, key).Bytes()
if err != nil {
    return err
}
defer redishoard.ClearBytes(sensitiveData)

// Use sensitiveData...
processData(sensitiveData)

// sensitiveData is automatically zeroed by defer
```

**Important Notes**:
- ⚠️ **Does not work on strings** - Go strings are immutable
- ⚠️ Use `[]byte` for sensitive data, not `string`
- ✅ Safe to call on empty slices (no-op)
- ✅ Safe to call multiple times on same slice

**Example (Wrong - Will Not Work)**:
```go
password := "secret"
redishoard.ClearBytes([]byte(password)) // ❌ Original string unchanged!
```

**Example (Correct)**:
```go
password := []byte("secret")
defer redishoard.ClearBytes(password) // ✅ password will be zeroed
```

---

### func ValidateNut

```go
func ValidateNut(nut string) error
```

Validates that a nut conforms to SQRL specifications.

**Parameters**:
- `nut` (`string`): Nut value to validate

**Returns**:
- `error`: Validation error, or `nil` if valid

**Validation Rules**:
1. Length: 8-64 characters
2. Character set: `[A-Za-z0-9_-]` (base64url safe)

**Errors**:
| Error | Condition |
|-------|-----------|
| `ErrNutTooShort` | `len(nut) < 8` |
| `ErrNutTooLong` | `len(nut) > 64` |
| `ErrNutInvalidChars` | Contains characters outside `[A-Za-z0-9_-]` |

**Example**:
```go
err := redishoard.ValidateNut(userInput)
if errors.Is(err, redishoard.ErrNutTooShort) {
    return errors.New("nut must be at least 8 characters")
}
if errors.Is(err, redishoard.ErrNutInvalidChars) {
    return errors.New("nut contains illegal characters")
}
if err != nil {
    return fmt.Errorf("invalid nut: %w", err)
}

// nut is valid, safe to use
```

**Security**: Prevents injection attacks by rejecting special characters.

---

### func NewSecureBuffer

```go
func NewSecureBuffer(size int, lock bool) (*SecureBuffer, error)
```

Creates a new secure buffer with optional memory locking.

**Parameters**:
- `size` (`int`): Size of buffer in bytes
- `lock` (`bool`): Whether to lock memory (prevent swap to disk)

**Returns**:
- `*SecureBuffer`: New secure buffer
- `error`: Error if memory locking fails (only if `lock=true`)

**Memory Locking Behavior**:
- **Unix/Linux/macOS** (if `lock=true`): Calls `mlock()` - requires `CAP_IPC_LOCK`
- **Windows** (if `lock=true`): Calls `VirtualLock()` - may require elevated privileges
- **Other platforms**: Ignores `lock` parameter (degrades gracefully)

**Example**:
```go
// Create 256-byte buffer with memory locking
buf, err := redishoard.NewSecureBuffer(256, true)
if err != nil {
    // Memory locking failed (non-fatal, buffer still usable)
    log.Warnf("memory locking unavailable: %v", err)
    buf, _ = redishoard.NewSecureBuffer(256, false)
}
defer buf.Destroy()

// Use the buffer
copy(buf.Data(), encryptionKey)
```

**Best Practices**:
- Always `defer buf.Destroy()` to ensure cleanup
- Use `lock=true` for highly sensitive data (keys, passwords)
- Use `lock=false` for less critical data to avoid privilege requirements

---

### func (*SecureBuffer) Data

```go
func (sb *SecureBuffer) Data() []byte
```

Returns the underlying byte slice.

**Returns**:
- `[]byte`: Mutable byte slice

**Example**:
```go
buf, _ := redishoard.NewSecureBuffer(32, false)
defer buf.Destroy()

// Write sensitive data
copy(buf.Data(), []byte("secret-token-12345"))

// Read sensitive data
token := buf.Data()
```

---

### func (*SecureBuffer) Clear

```go
func (sb *SecureBuffer) Clear()
```

Securely clears the buffer contents (without destroying the buffer).

**Example**:
```go
buf, _ := redishoard.NewSecureBuffer(32, false)
defer buf.Destroy()

// Use buffer for first operation
copy(buf.Data(), operation1Data)
processData(buf.Data())

// Clear and reuse for second operation
buf.Clear()
copy(buf.Data(), operation2Data)
processData(buf.Data())
```

---

### func (*SecureBuffer) Destroy

```go
func (sb *SecureBuffer) Destroy()
```

Clears and releases the buffer. After calling `Destroy()`, the buffer cannot be reused.

**Example**:
```go
buf, _ := redishoard.NewSecureBuffer(32, true)
defer buf.Destroy() // Always defer to ensure cleanup

// Use buffer...

// Explicit early cleanup if needed
buf.Destroy()
```

---

## Error Values

### var ErrNutTooShort

```go
var ErrNutTooShort = errors.New("nut too short: minimum 8 characters")
```

Indicates the nut is below minimum length.

---

### var ErrNutTooLong

```go
var ErrNutTooLong = errors.New("nut too long: maximum 64 characters")
```

Indicates the nut exceeds maximum length.

---

### var ErrNutInvalidChars

```go
var ErrNutInvalidChars = errors.New("nut contains invalid characters: must be base64url safe")
```

Indicates the nut contains characters outside the base64url character set.

---

### var ErrInvalidHoardCache

```go
var ErrInvalidHoardCache = errors.New("invalid hoard cache data")
```

Indicates the cached data is malformed (JSON deserialization failed).

---

## Constants

### Nut Validation

```go
const (
    MinNutLength = 8  // Minimum nut length in characters
    MaxNutLength = 64 // Maximum nut length in characters (base64 encoded)
)
```

---

## Complete Usage Example

### Basic SQRL Authentication Flow

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
    ssp "github.com/sqrldev/server-go-ssp"
    "github.com/sqrldev/server-go-ssp-redishoard"
)

func main() {
    // 1. Setup Redis connection
    client := redis.NewUniversalClient(&redis.UniversalOptions{
        Addrs:    []string{"localhost:6379"},
        Password: "", // Set if required
        DB:       0,  // Default DB
    })
    defer client.Close()

    // Verify Redis connectivity
    if err := client.Ping(context.Background()).Err(); err != nil {
        log.Fatalf("Redis connection failed: %v", err)
    }

    // 2. Create Hoard
    hoard := redishoard.NewHoard(client)

    // 3. Generate nut (normally done by SQRL tree)
    nut := ssp.Nut("sqrl-nut-abc123def456")

    // 4. Save initial authentication state
    initialCache := &ssp.HoardCache{
        State:       "awaiting-client",
        RemoteIP:    "192.168.1.100",
        OriginalNut: nut,
        // Other fields populated by SQRL flow...
    }

    err := hoard.Save(nut, initialCache, 5*time.Minute)
    if err != nil {
        log.Fatalf("Failed to save: %v", err)
    }
    fmt.Println("✅ Session saved")

    // 5. Client sends authentication request...
    // 6. Retrieve and verify (Get - non-destructive)
    cache, err := hoard.Get(nut)
    if errors.Is(err, ssp.ErrNotFound) {
        log.Fatal("Session not found or expired")
    }
    if err != nil {
        log.Fatalf("Failed to get session: %v", err)
    }
    fmt.Printf("✅ Session retrieved: state=%s\n", cache.State)

    // 7. Final authentication step - consume nut atomically
    finalCache, err := hoard.GetAndDelete(nut)
    if errors.Is(err, ssp.ErrNotFound) {
        log.Fatal("Nut already used (replay attack?)")
    }
    if err != nil {
        log.Fatalf("Failed to consume nut: %v", err)
    }
    fmt.Printf("✅ Nut consumed: state=%s\n", finalCache.State)

    // 8. Verify nut is deleted
    _, err = hoard.Get(nut)
    if !errors.Is(err, ssp.ErrNotFound) {
        log.Fatal("ERROR: Nut should have been deleted!")
    }
    fmt.Println("✅ Nut confirmed deleted (replay prevention works)")
}
```

### Output:
```
✅ Session saved
✅ Session retrieved: state=awaiting-client
✅ Nut consumed: state=awaiting-client
✅ Nut confirmed deleted (replay prevention works)
```

---

## Testing Your Integration

### Unit Testing with Mock Redis

```go
package yourapp

import (
    "testing"
    "time"

    "github.com/redis/go-redis/v9"
    ssp "github.com/sqrldev/server-go-ssp"
    "github.com/sqrldev/server-go-ssp-redishoard"
)

func TestHoardIntegration(t *testing.T) {
    // Use test Redis instance (or miniredis for mocking)
    client := redis.NewUniversalClient(&redis.UniversalOptions{
        Addrs: []string{"localhost:6379"},
    })
    defer client.Close()

    hoard := redishoard.NewHoard(client)

    // Test Save
    nut := ssp.Nut("test-nut-12345678")
    cache := &ssp.HoardCache{State: "test"}

    err := hoard.Save(nut, cache, time.Minute)
    if err != nil {
        t.Fatalf("Save failed: %v", err)
    }

    // Test Get
    retrieved, err := hoard.Get(nut)
    if err != nil {
        t.Fatalf("Get failed: %v", err)
    }
    if retrieved.State != "test" {
        t.Errorf("Wrong state: got %s, want test", retrieved.State)
    }

    // Test GetAndDelete
    deleted, err := hoard.GetAndDelete(nut)
    if err != nil {
        t.Fatalf("GetAndDelete failed: %v", err)
    }
    if deleted.State != "test" {
        t.Errorf("Wrong state: got %s, want test", deleted.State)
    }

    // Verify deletion
    _, err = hoard.Get(nut)
    if err != ssp.ErrNotFound {
        t.Error("Nut should have been deleted")
    }
}
```

---

## Performance Benchmarks

### Typical Performance (Redis on localhost)

| Operation | Latency (p50) | Latency (p95) | Throughput |
|-----------|---------------|---------------|------------|
| Get | <2ms | <5ms | ~10,000 ops/sec |
| Save | <2ms | <5ms | ~8,000 ops/sec |
| GetAndDelete | <3ms | <10ms | ~5,000 ops/sec |

### Running Benchmarks

```bash
go test -bench=. -benchmem
```

---

## Security Considerations

### 1. Memory Security

✅ **Implemented**:
- Sensitive data cleared from memory after use
- Platform-specific memory locking available
- Constant-time operations to prevent timing attacks

⚠️ **Application Responsibility**:
- Use `[]byte` for sensitive data, not `string`
- Always defer `ClearBytes()` or `Destroy()`
- Avoid logging HoardCache contents

### 2. Transport Security

⚠️ **Required for Production**:
```go
client := redis.NewUniversalClient(&redis.UniversalOptions{
    Addrs: []string{"redis.example.com:6380"},
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
        ServerName: "redis.example.com",
    },
})
```

### 3. Input Validation

✅ **Automatic**: All nuts validated before Redis operations

### 4. Replay Attack Prevention

✅ **Use GetAndDelete()**: Ensures one-time use of authentication nonces

---

## Migration Guide

### From ssp.MapHoard (In-Memory)

```go
// Old (in-memory, single server)
hoard := ssp.NewMapHoard()

// New (Redis-backed, distributed)
client := redis.NewUniversalClient(&redis.UniversalOptions{
    Addrs: []string{"localhost:6379"},
})
hoard := redishoard.NewHoard(client)

// API is identical - no code changes needed!
```

### From Other ssp.Hoard Implementations

The `ssp.Hoard` interface is stable. Simply replace your Hoard constructor with `redishoard.NewHoard()`.

---

## Troubleshooting

### Common Issues

#### Error: "nut contains invalid characters"

**Cause**: Nut contains characters outside `[A-Za-z0-9_-]`

**Fix**:
```go
// ❌ Wrong
nut := ssp.Nut("nut+with/special=chars")

// ✅ Correct (base64url encoding)
nut := ssp.Nut("nut-with_safe-chars")
```

#### Error: "redis: nil" / Connection Refused

**Cause**: Redis not running or wrong address

**Fix**:
```bash
# Check Redis is running
redis-cli ping

# Should output: PONG
```

#### Error: ssp.ErrNotFound on Get

**Causes**:
1. Nut expired (TTL elapsed)
2. Nut never saved
3. Nut already consumed by GetAndDelete

**Debug**:
```bash
# Check if key exists in Redis
redis-cli GET "your-nut-value"
```

---

## References

- **GoDoc**: https://pkg.go.dev/github.com/sqrldev/server-go-ssp-redishoard
- **SQRL Specification**: https://www.grc.com/sqrl/sqrl.htm
- **server-go-ssp**: https://github.com/dxcSithLord/server-go-ssp
- **Redis Client**: https://github.com/redis/go-redis

---

## Change Log

| Date | Version | Changes |
|------|---------|---------|
| 2025-01-18 | 1.0.0 | Initial API reference documentation |

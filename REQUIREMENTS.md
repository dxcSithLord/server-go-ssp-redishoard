# SQRL Redis Hoard - Requirements Documentation

## Project Overview

**Package Name:** `redishoard`
**Purpose:** Redis-backed implementation of the `ssp.Hoard` interface for SQRL (Secure Quick Reliable Login) protocol state management.
**License:** MIT
**Go Module Path:** `github.com/sqrldev/server-go-ssp-redishoard`

## Functional Requirements

### FR-1: State Storage
- **FR-1.1:** Store SQRL authentication state (HoardCache) in Redis with configurable TTL
- **FR-1.2:** Support JSON serialization/deserialization of state objects
- **FR-1.3:** Use cryptographic nonces (Nut) as unique keys

### FR-2: State Retrieval
- **FR-2.1:** Retrieve cached state by nut identifier
- **FR-2.2:** Return `ssp.ErrNotFound` when nut does not exist
- **FR-2.3:** Properly differentiate between "not found" and actual errors

### FR-3: Atomic Operations
- **FR-3.1:** Implement atomic get-and-delete for one-time nut consumption
- **FR-3.2:** Use Redis transactions to prevent race conditions
- **FR-3.3:** Ensure state is deleted even if retrieval partially fails

### FR-4: Interface Compliance
- **FR-4.1:** Implement `ssp.Hoard` interface completely
- **FR-4.2:** Support horizontal scaling via Redis
- **FR-4.3:** Work with any `redis.UniversalClient` implementation

## Non-Functional Requirements

### NFR-1: Security
- **NFR-1.1:** Clear sensitive data from memory after use (CWE-226 mitigation)
- **NFR-1.2:** Avoid logging sensitive authentication data (CWE-200 mitigation)
- **NFR-1.3:** Validate input nut values before use
- **NFR-1.4:** Support secure Redis connections (TLS)
- **NFR-1.5:** Platform-aware memory clearing implementation

### NFR-2: Performance
- **NFR-2.1:** Minimize Redis round-trips
- **NFR-2.2:** Use efficient serialization (JSON)
- **NFR-2.3:** Support connection pooling via UniversalClient

### NFR-3: Reliability
- **NFR-3.1:** Handle Redis connection failures gracefully
- **NFR-3.2:** Return meaningful error messages
- **NFR-3.3:** Automatic TTL-based cache expiration

### NFR-4: Maintainability
- **NFR-4.1:** Go module with proper versioning
- **NFR-4.2:** Comprehensive test coverage (>80%)
- **NFR-4.3:** CI/CD pipeline for automated testing
- **NFR-4.4:** Static analysis and security scanning

## Architecture

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────┐
│   SQRL Client   │────>│  server-go-ssp   │────>│ Redis Hoard │
│                 │     │   (SSP API)      │     │  (this pkg) │
└─────────────────┘     └──────────────────┘     └─────────────┘
                                                        │
                                                        v
                                                 ┌─────────────┐
                                                 │    Redis    │
                                                 │   Server    │
                                                 └─────────────┘
```

### Data Flow
1. SQRL client sends authentication request
2. SSP API generates nut (cryptographic nonce)
3. Redis Hoard stores HoardCache with TTL
4. Client completes authentication
5. Redis Hoard retrieves and deletes state atomically

## Data Structures

### HoardCache (from ssp package)
```go
type HoardCache struct {
    RemoteIP    string      // Client IP address
    OriginalNut Nut         // Original nonce
    PagNut      Nut         // Page-specific nonce
    LastRequest CliRequest  // Last client request
    Identity    interface{} // User identity data
    LastResponse CliResponse // Server response
}
```

### Nut (from ssp package)
```go
type Nut string // Cryptographic nonce (8-20 random bytes)
```

## Security Considerations

### Sensitive Data Handled
1. **Cryptographic nonces (Nut)** - Used as Redis keys
2. **Client IP addresses** - PII stored in HoardCache
3. **Authentication state** - Session validation data
4. **Identity credentials** - Idk, Suk, Vuk keys (Ed25519)
5. **Signature data** - Cryptographic verification material

### Threat Model
- **T1:** Memory disclosure attacks (cold boot, memory dump)
- **T2:** Log injection/exfiltration of sensitive data
- **T3:** Redis data compromise (unencrypted storage)
- **T4:** Race conditions in authentication flow
- **T5:** Timing attacks on cryptographic operations

## Testing Requirements

### Unit Tests
- [ ] NewHoard initialization
- [ ] Save with valid data
- [ ] Save with invalid JSON (error handling)
- [ ] Get existing nut
- [ ] Get non-existent nut
- [ ] GetAndDelete atomic behavior
- [ ] Concurrent access patterns
- [ ] TTL expiration behavior

### Integration Tests
- [ ] Redis connection handling
- [ ] Transaction rollback scenarios
- [ ] Memory clearing verification
- [ ] Large payload handling

### Security Tests
- [ ] Sensitive data not leaked to logs
- [ ] Memory cleared after operations
- [ ] Input validation enforced

## Deployment Requirements

### Redis Server
- **Version:** 7.2+ (recommended for RESP3)
- **Configuration:**
  - TLS enabled for production
  - Authentication configured
  - Appropriate memory limits
  - Persistence optional (state is ephemeral)

### Go Runtime
- **Version:** 1.23+ (for latest security patches)
- **Platform:** Linux, macOS, Windows (platform-aware code)

### Environment Variables
- `REDIS_URL` - Redis connection string
- `REDIS_TLS` - Enable TLS (true/false)
- `REDIS_PASSWORD` - Authentication credential

## Compliance

This implementation should adhere to:
- **OWASP Top 10** - Especially A02:2021 Cryptographic Failures
- **CWE-226** - Sensitive Information in Resource Not Removed Before Reuse
- **CWE-200** - Exposure of Sensitive Information
- **CWE-312** - Cleartext Storage of Sensitive Information
- **SQRL Specification** - GRC SQRL SSP API requirements

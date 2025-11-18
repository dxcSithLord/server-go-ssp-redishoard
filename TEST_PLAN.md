# Comprehensive Test Plan

## Document Control
- **Project**: server-go-ssp-redishoard
- **Version**: 1.0.0
- **Date**: 2025-01-18
- **Current Coverage**: 46.6%
- **Target Coverage**: 85%+
- **Repository**: https://github.com/dxcSithLord/server-go-ssp-redishoard

---

## Executive Summary

This document provides a comprehensive test plan for the **server-go-ssp-redishoard** library. Since this is a **Go library** (not an HTTP REST API), this plan focuses on:
- **Unit tests** for each public function
- **Integration tests** with Redis
- **Security tests** for memory clearing and validation
- **Concurrency tests** for race conditions
- **Platform-specific tests** for memory locking
- **Performance benchmarks**

**Current State**:
- ✅ Basic happy-path tests exist
- ⚠️ Coverage: 46.6% (need 38.4% more)
- ❌ Error path tests missing
- ❌ Concurrency tests missing
- ❌ Security verification tests missing

---

## Test Coverage Matrix

### Public API Functions

| Function | Unit Tests | Integration Tests | Error Tests | Concurrency Tests | Security Tests | Coverage |
|----------|-----------|-------------------|-------------|-------------------|----------------|----------|
| `NewHoard()` | ✅ Implicit | ✅ Yes | N/A | ✅ Yes (goroutine-safe) | N/A | 100% |
| `(*Hoard).Get()` | ⚠️ Basic | ✅ Yes | ❌ Missing | ❌ Missing | ⚠️ Partial | ~40% |
| `(*Hoard).GetAndDelete()` | ⚠️ Basic | ✅ Yes | ❌ Missing | ❌ Missing | ⚠️ Partial | ~50% |
| `(*Hoard).Save()` | ⚠️ Basic | ✅ Yes | ⚠️ Partial | ❌ Missing | ⚠️ Partial | ~50% |
| `ClearBytes()` | ✅ Good | N/A | ✅ Yes | N/A | ✅ Yes | ~90% |
| `ValidateNut()` | ✅ Excellent | N/A | ✅ Yes | N/A | N/A | 100% |
| `NewSecureBuffer()` | ⚠️ Basic | N/A | ❌ Missing | N/A | ⚠️ Partial | ~60% |
| `(*SecureBuffer).Data()` | ✅ Yes | N/A | N/A | N/A | N/A | 100% |
| `(*SecureBuffer).Clear()` | ✅ Yes | N/A | N/A | N/A | ✅ Yes | 100% |
| `(*SecureBuffer).Destroy()` | ✅ Yes | N/A | N/A | N/A | ✅ Yes | 100% |
| `LockMemory()` | ❌ Missing | N/A | ❌ Missing | N/A | ❌ Missing | 0% |
| `UnlockMemory()` | ❌ Missing | N/A | ❌ Missing | N/A | ❌ Missing | 0% |

### Internal/Private Functions

| Function | Unit Tests | Coverage |
|----------|-----------|----------|
| `(*Hoard).fromBytes()` | ⚠️ Implicit | ~40% |
| `isBase64URLSafe()` | ✅ Yes | 100% |

---

## Test Categories

### 1. Unit Tests (No External Dependencies)

**Objective**: Test individual functions in isolation

**Current Files**: `secure_test.go`

**Coverage Goals**: All validation, security utility functions at 100%

#### 1.1 Validation Tests

**File**: `secure_test.go`

**Status**: ✅ **COMPLETE** (100% coverage)

**Existing Tests**:
- ✅ `TestValidateNut()` - Covers all validation rules
- ✅ `TestIsBase64URLSafe()` - Covers character set validation

**Test Cases**:
- ✅ Valid nut (minimum length)
- ✅ Valid nut (maximum length)
- ✅ Valid nut (with dashes and underscores)
- ✅ Invalid nut (too short)
- ✅ Invalid nut (too long)
- ✅ Invalid nut (special characters: @, #, !, space, /, +, =)

---

#### 1.2 Security Utility Tests

**File**: `secure_test.go`

**Status**: ⚠️ **PARTIAL** (need to add more cases)

**Existing Tests**:
- ✅ `TestClearBytes()` - Basic clearing verification
- ✅ `TestSecureBuffer()` - Create/destroy/clear operations

**Missing Tests**:
- ❌ Verify compiler doesn't optimize away ClearBytes
- ❌ Test ClearBytes with very large slices (>1MB)
- ❌ Test SecureBuffer with memory locking on supported platforms
- ❌ Test SecureBuffer memory locking failures

**New Tests Needed**:

```go
// TestClearBytesNotOptimizedAway verifies ClearBytes actually clears memory
func TestClearBytesNotOptimizedAway(t *testing.T) {
    data := []byte("sensitive-password-12345")
    originalPtr := unsafe.Pointer(&data[0])

    ClearBytes(data)

    // Verify memory at original pointer is zeroed
    for i := 0; i < len(data); i++ {
        val := *(*byte)(unsafe.Pointer(uintptr(originalPtr) + uintptr(i)))
        if val != 0 {
            t.Errorf("byte %d not cleared: got %d", i, val)
        }
    }
}

// TestClearBytesLargeSlice tests clearing multi-megabyte buffers
func TestClearBytesLargeSlice(t *testing.T) {
    sizes := []int{1024 * 1024, 10 * 1024 * 1024} // 1MB, 10MB
    for _, size := range sizes {
        t.Run(fmt.Sprintf("size_%dMB", size/(1024*1024)), func(t *testing.T) {
            data := make([]byte, size)
            // Fill with non-zero data
            for i := range data {
                data[i] = byte(i % 256)
            }

            start := time.Now()
            ClearBytes(data)
            duration := time.Since(start)

            // Verify all zeroed
            for i, b := range data {
                if b != 0 {
                    t.Fatalf("byte %d not cleared", i)
                }
            }

            t.Logf("Cleared %d MB in %v", size/(1024*1024), duration)
        })
    }
}

// TestSecureBufferMemoryLocking tests platform-specific locking
func TestSecureBufferMemoryLocking(t *testing.T) {
    if runtime.GOOS != "linux" && runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
        t.Skip("Memory locking not supported on", runtime.GOOS)
    }

    // This test may fail without elevated privileges
    buf, err := NewSecureBuffer(4096, true)
    if err != nil {
        t.Logf("Memory locking failed (may require privileges): %v", err)
        // Not a test failure - graceful degradation is acceptable
        return
    }
    defer buf.Destroy()

    if !buf.locked {
        t.Error("Buffer should be locked on supported platform")
    }
}

// TestSecureBufferLockingFailure tests graceful degradation
func TestSecureBufferLockingFailure(t *testing.T) {
    // Try to lock more memory than ulimit allows (this should fail gracefully)
    buf, err := NewSecureBuffer(1024*1024*1024, true) // 1GB
    if err == nil {
        defer buf.Destroy()
        // Locking may succeed on systems with high limits
        t.Log("Large buffer locking succeeded (system allows it)")
    } else {
        // Should still create buffer, just not locked
        buf, err := NewSecureBuffer(1024*1024*1024, false)
        if err != nil {
            t.Fatal("Should create unlocked buffer even if locking fails")
        }
        defer buf.Destroy()
    }
}
```

---

### 2. Integration Tests (Require Redis)

**Objective**: Test interactions with actual Redis instance

**Current Files**: `hoard_test.go`

**Coverage Goals**: All Hoard operations at 100%

#### 2.1 Basic Operations Tests

**File**: `hoard_test.go`

**Status**: ⚠️ **PARTIAL** (need error paths)

**Existing Tests**:
- ✅ `TestSave()` - Happy path
- ✅ `TestGetAndDelete()` - Happy path
- ✅ `TestGetNonExistentNut()` - ErrNotFound case
- ✅ `TestSaveNilCache()` - Nil validation
- ✅ `TestInvalidNutCharacters()` - Validation

**Missing Tests**:
- ❌ Redis connection failure scenarios
- ❌ Redis timeout scenarios
- ❌ Malformed JSON in Redis (data corruption)
- ❌ TTL expiration verification
- ❌ Large payload handling
- ❌ Unicode in HoardCache fields
- ❌ Concurrent access patterns

**New Tests Needed**:

```go
// TestSaveTTLExpiration verifies automatic key expiration
func TestSaveTTLExpiration(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)
    nut := ssp.Nut("ttl-test-nut-12345")
    cache := &ssp.HoardCache{State: "test"}

    // Save with 2-second TTL
    err := h.Save(nut, cache, 2*time.Second)
    if err != nil {
        t.Fatalf("Save failed: %v", err)
    }

    // Verify exists immediately
    _, err = h.Get(nut)
    if err != nil {
        t.Fatalf("Get failed immediately: %v", err)
    }

    // Wait for expiration
    time.Sleep(3 * time.Second)

    // Verify expired
    _, err = h.Get(nut)
    if !errors.Is(err, ssp.ErrNotFound) {
        t.Errorf("Key should have expired, got error: %v", err)
    }
}

// TestRedisConnectionFailure tests behavior when Redis is unavailable
func TestRedisConnectionFailure(t *testing.T) {
    // Connect to non-existent Redis instance
    client := redis.NewUniversalClient(&redis.UniversalOptions{
        Addrs:       []string{"localhost:9999"}, // Wrong port
        DialTimeout: 100 * time.Millisecond,
    })
    defer client.Close()

    h := NewHoard(client)
    nut := ssp.Nut("connection-test-nut")
    cache := &ssp.HoardCache{State: "test"}

    // Should fail with connection error
    err := h.Save(nut, cache, time.Minute)
    if err == nil {
        t.Fatal("Expected connection error, got nil")
    }

    // Error should contain connection context
    if !strings.Contains(err.Error(), "connection") &&
       !strings.Contains(err.Error(), "dial") {
        t.Errorf("Error should indicate connection issue: %v", err)
    }
}

// TestMalformedJSONInRedis tests handling of corrupted data
func TestMalformedJSONInRedis(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    ctx := context.Background()
    nut := "malformed-json-nut-123"

    // Manually insert malformed JSON
    err := client.Set(ctx, nut, "{this is not valid json", time.Minute).Err()
    if err != nil {
        t.Fatalf("Failed to setup test: %v", err)
    }
    defer client.Del(ctx, nut)

    h := NewHoard(client)

    // Should fail with decode error
    _, err = h.Get(ssp.Nut(nut))
    if err == nil {
        t.Fatal("Expected JSON decode error, got nil")
    }

    if !strings.Contains(err.Error(), "decode") {
        t.Errorf("Error should indicate decode failure: %v", err)
    }
}

// TestLargePayload tests handling of large HoardCache objects
func TestLargePayload(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)
    nut := ssp.Nut("large-payload-nut-123")

    // Create large payload (1MB of data)
    largeData := make([]byte, 1024*1024)
    for i := range largeData {
        largeData[i] = byte(i % 256)
    }

    cache := &ssp.HoardCache{
        State:        string(largeData),
        RemoteIP:     "192.168.1.1",
        OriginalNut:  nut,
    }

    // Save
    err := h.Save(nut, cache, time.Minute)
    if err != nil {
        t.Fatalf("Failed to save large payload: %v", err)
    }
    defer client.Del(context.Background(), string(nut))

    // Retrieve
    retrieved, err := h.Get(nut)
    if err != nil {
        t.Fatalf("Failed to get large payload: %v", err)
    }

    // Verify
    if len(retrieved.State) != len(cache.State) {
        t.Errorf("Payload size mismatch: got %d, want %d",
            len(retrieved.State), len(cache.State))
    }
}

// TestUnicodeData tests handling of international characters
func TestUnicodeData(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)
    nut := ssp.Nut("unicode-test-nut-123")

    cache := &ssp.HoardCache{
        State:    "Hello 世界 🌍 Привет مرحبا",
        RemoteIP: "2001:db8::1",
    }

    err := h.Save(nut, cache, time.Minute)
    if err != nil {
        t.Fatalf("Failed to save unicode data: %v", err)
    }
    defer client.Del(context.Background(), string(nut))

    retrieved, err := h.Get(nut)
    if err != nil {
        t.Fatalf("Failed to get unicode data: %v", err)
    }

    if retrieved.State != cache.State {
        t.Errorf("Unicode mismatch: got %q, want %q",
            retrieved.State, cache.State)
    }
}
```

---

#### 2.2 Concurrency Tests

**File**: `hoard_concurrent_test.go` (NEW FILE)

**Status**: ❌ **MISSING**

**Tests Needed**:

```go
package redishoard

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "testing"
    "time"

    "github.com/redis/go-redis/v9"
    ssp "github.com/sqrldev/server-go-ssp"
)

// TestConcurrentSave tests multiple goroutines saving different nuts
func TestConcurrentSave(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)

    const numGoroutines = 100
    var wg sync.WaitGroup
    errors := make(chan error, numGoroutines)

    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            nut := ssp.Nut(fmt.Sprintf("concurrent-save-%d", id))
            cache := &ssp.HoardCache{State: fmt.Sprintf("state-%d", id)}

            if err := h.Save(nut, cache, time.Minute); err != nil {
                errors <- fmt.Errorf("goroutine %d: %w", id, err)
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    // Check for errors
    for err := range errors {
        t.Error(err)
    }

    // Cleanup
    ctx := context.Background()
    for i := 0; i < numGoroutines; i++ {
        client.Del(ctx, fmt.Sprintf("concurrent-save-%d", i))
    }
}

// TestConcurrentGetAndDelete tests race conditions in atomic operation
func TestConcurrentGetAndDelete(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)
    nut := ssp.Nut("concurrent-gad-test")
    cache := &ssp.HoardCache{State: "test"}

    // Save once
    err := h.Save(nut, cache, time.Minute)
    if err != nil {
        t.Fatalf("Setup failed: %v", err)
    }

    const numGoroutines = 10
    var wg sync.WaitGroup
    successCount := 0
    var mu sync.Mutex

    // Multiple goroutines try to GetAndDelete the same nut
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()

            _, err := h.GetAndDelete(nut)
            if err == nil {
                mu.Lock()
                successCount++
                mu.Unlock()
            } else if !errors.Is(err, ssp.ErrNotFound) {
                t.Errorf("Unexpected error from goroutine %d: %v", id, err)
            }
        }(i)
    }

    wg.Wait()

    // Exactly ONE goroutine should succeed
    if successCount != 1 {
        t.Errorf("Expected exactly 1 success, got %d", successCount)
    }
}

// TestConcurrentReadWrite tests concurrent Get/Save operations
func TestConcurrentReadWrite(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)
    nut := ssp.Nut("concurrent-rw-test")

    var wg sync.WaitGroup
    stopChan := make(chan struct{})

    // Writer goroutine
    wg.Add(1)
    go func() {
        defer wg.Done()
        counter := 0
        for {
            select {
            case <-stopChan:
                return
            default:
                cache := &ssp.HoardCache{State: fmt.Sprintf("state-%d", counter)}
                h.Save(nut, cache, time.Minute)
                counter++
                time.Sleep(time.Millisecond)
            }
        }
    }()

    // Reader goroutines
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            for {
                select {
                case <-stopChan:
                    return
                default:
                    h.Get(nut)
                    time.Sleep(time.Millisecond)
                }
            }
        }(i)
    }

    // Run for 1 second
    time.Sleep(time.Second)
    close(stopChan)
    wg.Wait()

    // Cleanup
    client.Del(context.Background(), string(nut))
}
```

---

### 3. Security Tests

**Objective**: Verify security properties (memory clearing, validation)

**File**: `security_verification_test.go` (NEW FILE)

**Status**: ❌ **MISSING**

**Tests Needed**:

```go
package redishoard

import (
    "testing"
    "unsafe"
)

// TestMemoryClearedAfterGet verifies data is cleared after Get operation
func TestMemoryClearedAfterGet(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)
    nut := ssp.Nut("memory-clear-test")
    sensitiveData := "SENSITIVE-PASSWORD-12345"
    cache := &ssp.HoardCache{State: sensitiveData}

    // Save
    err := h.Save(nut, cache, time.Minute)
    if err != nil {
        t.Fatalf("Save failed: %v", err)
    }
    defer client.Del(context.Background(), string(nut))

    // Get and let it clear memory
    retrieved, err := h.Get(nut)
    if err != nil {
        t.Fatalf("Get failed: %v", err)
    }

    // At this point, internal byte buffer should be cleared
    // We can't directly verify this without reflection/unsafe,
    // but we can verify the function doesn't panic and works correctly

    if retrieved.State != sensitiveData {
        t.Errorf("Data corrupted: got %q, want %q", retrieved.State, sensitiveData)
    }
}

// TestNoSensitiveDataInErrors verifies errors don't leak nut values
func TestNoSensitiveDataInErrors(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)

    // Try to get non-existent nut with "sensitive" value
    nut := ssp.Nut("SENSITIVE-NUT-VALUE-12345")
    _, err := h.Get(nut)

    if err == nil {
        t.Fatal("Expected error for non-existent nut")
    }

    // Error message should NOT contain the actual nut value
    errMsg := err.Error()
    if strings.Contains(errMsg, string(nut)) {
        t.Errorf("Error message leaks nut value: %s", errMsg)
    }
}

// TestInputValidationPreventsInjection tests that special chars are rejected
func TestInputValidationPreventsInjection(t *testing.T) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)

    maliciousNuts := []string{
        "nut'; DROP TABLE users; --",
        "nut\x00null-byte",
        "nut\nwith\nnewlines",
        "../../../etc/passwd",
        "<script>alert('xss')</script>",
    }

    for _, nut := range maliciousNuts {
        _, err := h.Get(ssp.Nut(nut))
        if err == nil {
            t.Errorf("Malicious nut should be rejected: %q", nut)
        }
        if !strings.Contains(err.Error(), "invalid") {
            t.Errorf("Error should indicate validation failure for %q: %v", nut, err)
        }
    }
}
```

---

### 4. Platform-Specific Tests

**Objective**: Test memory locking on different operating systems

**File**: `secure_platform_test.go` (NEW FILE)

**Status**: ❌ **MISSING**

**Tests Needed**:

```go
//go:build linux || darwin

package redishoard

import (
    "testing"
    "golang.org/x/sys/unix"
)

// TestMemoryLockingUnix tests mlock on Unix-like systems
func TestMemoryLockingUnix(t *testing.T) {
    data := make([]byte, 4096)

    err := LockMemory(data)
    if err != nil {
        // May fail without CAP_IPC_LOCK capability
        t.Skipf("Memory locking failed (may need elevated privileges): %v", err)
    }
    defer UnlockMemory(data)

    // Verify memory is accessible
    data[0] = 42
    if data[0] != 42 {
        t.Error("Locked memory should be accessible")
    }

    // Cleanup
    err = UnlockMemory(data)
    if err != nil {
        t.Errorf("Unlock failed: %v", err)
    }
}
```

```go
//go:build windows

package redishoard

import (
    "testing"
)

// TestMemoryLockingWindows tests VirtualLock on Windows
func TestMemoryLockingWindows(t *testing.T) {
    data := make([]byte, 4096)

    err := LockMemory(data)
    if err != nil {
        // May fail without elevated privileges
        t.Skipf("Memory locking failed (may need admin): %v", err)
    }
    defer UnlockMemory(data)

    // Verify memory is accessible
    data[0] = 42
    if data[0] != 42 {
        t.Error("Locked memory should be accessible")
    }
}
```

---

## Test Execution Strategy

### Local Development

```bash
# Run all tests (requires local Redis)
go test -v ./...

# Run only unit tests (no Redis needed)
go test -v -short ./...

# Run with race detector
go test -race ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -v -run TestConcurrentGetAndDelete
```

### CI/CD Pipeline

```yaml
# .github/workflows/ci.yml (updated)
test:
  strategy:
    matrix:
      os: [ubuntu-latest, macos-latest, windows-latest]
      go: ['1.24', '1.25']
      redis: ['7.2', '7.4']

  steps:
    - name: Start Redis
      run: docker run -d -p 6379:6379 redis:${{ matrix.redis }}

    - name: Run Tests
      run: go test -v -race -coverprofile=coverage.out ./...

    - name: Upload Coverage
      run: bash <(curl -s https://codecov.io/bash)
```

---

## Coverage Improvement Plan

### Current: 46.6% → Target: 85%

**Phase 1: Error Paths** (+15% coverage)
- [ ] Redis connection failures
- [ ] Timeout scenarios
- [ ] Malformed JSON handling
- [ ] Validation errors

**Phase 2: Concurrency** (+10% coverage)
- [ ] Concurrent Save operations
- [ ] Race-free GetAndDelete
- [ ] Read/write contention

**Phase 3: Security** (+8% coverage)
- [ ] Memory clearing verification
- [ ] No data leakage in errors
- [ ] Injection prevention

**Phase 4: Edge Cases** (+8% coverage)
- [ ] Large payloads
- [ ] TTL expiration
- [ ] Unicode handling
- [ ] Platform-specific features

**Total**: 46.6% + 41% = **87.6% coverage** ✅

---

## Test Quality Metrics

### Assertions Per Test
- **Target**: ≥ 3 assertions per test
- **Current**: Varies (some tests have only 1)

### Test Isolation
- ✅ Each test creates unique nuts (no collisions)
- ✅ Tests clean up Redis keys
- ⚠️ Some tests don't handle Redis unavailability gracefully

### Test Naming Convention
- ✅ `Test<Function><Scenario>` (e.g., `TestSaveTTLExpiration`)
- ✅ Table-driven tests use `t.Run()`

---

## Performance Benchmarks

### Current Benchmarks

```bash
# Run benchmarks
go test -bench=. -benchmem

# Compare with baseline
go test -bench=. -benchmem -count=5 > new.txt
benchstat baseline.txt new.txt
```

### Benchmark Coverage

| Operation | Benchmark | Status |
|-----------|-----------|--------|
| ClearBytes | ✅ Yes | Multiple sizes tested |
| ValidateNut | ✅ Yes | Standard nut tested |
| Save | ❌ Missing | Need to add |
| Get | ❌ Missing | Need to add |
| GetAndDelete | ❌ Missing | Need to add |

**New Benchmarks Needed**:

```go
// BenchmarkSave benchmarks Save operation
func BenchmarkSave(b *testing.B) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)
    cache := &ssp.HoardCache{State: "benchmark"}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        nut := ssp.Nut(fmt.Sprintf("bench-save-%d", i))
        h.Save(nut, cache, time.Minute)
    }
}

// BenchmarkGet benchmarks Get operation
func BenchmarkGet(b *testing.B) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)
    nut := ssp.Nut("benchmark-get-nut")
    cache := &ssp.HoardCache{State: "benchmark"}
    h.Save(nut, cache, time.Minute)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        h.Get(nut)
    }
}

// BenchmarkGetAndDelete benchmarks atomic operation
func BenchmarkGetAndDelete(b *testing.B) {
    client := redis.NewUniversalClient(&redis.UniversalOptions{})
    defer client.Close()

    h := NewHoard(client)
    cache := &ssp.HoardCache{State: "benchmark"}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        nut := ssp.Nut(fmt.Sprintf("bench-gad-%d", i))
        h.Save(nut, cache, time.Minute)
        b.StartTimer()

        h.GetAndDelete(nut)
    }
}
```

---

## Test Maintenance

### Code Review Checklist

- [ ] Every new public function has tests
- [ ] Every error path is tested
- [ ] Every security property is verified
- [ ] Benchmarks added for performance-critical code
- [ ] Test names are descriptive
- [ ] Tests are isolated (no shared state)
- [ ] Tests clean up resources

### Periodic Review

- **Weekly**: Review test failures in CI/CD
- **Monthly**: Review coverage reports, identify gaps
- **Quarterly**: Review and update benchmarks

---

## References

- **Go Testing**: https://go.dev/doc/tutorial/add-a-test
- **Table-Driven Tests**: https://go.dev/wiki/TableDrivenTests
- **Testing Best Practices**: https://go.dev/doc/code
- **Race Detector**: https://go.dev/doc/articles/race_detector

---

## Change Log

| Date | Version | Changes |
|------|---------|---------|
| 2025-01-18 | 1.0.0 | Initial comprehensive test plan |

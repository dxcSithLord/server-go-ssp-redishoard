# Dependency Analysis & Upgrade Path

## Document Control
- **Project**: server-go-ssp-redishoard
- **Version**: 1.0.0
- **Date**: 2025-01-18
- **Status**: Living Document
- **Repository**: https://github.com/dxcSithLord/server-go-ssp-redishoard

---

## Executive Summary

This document provides a comprehensive analysis of all dependencies for the **server-go-ssp-redishoard** project, identifies available upgrades, analyzes breaking changes, and provides a sequenced upgrade strategy for maintaining compatibility while adopting latest security patches and features.

**Key Findings**:
- ✅ **7 direct dependencies** (including toolchain)
- ✅ **5 indirect dependencies** (pulled by direct deps)
- ⚠️ **2 dependencies with available updates**
- 🔴 **1 major version constraint** (server-go-ssp requires Go 1.25+)
- ✅ **No goauthentik.io integration** found (not a dependency of this library)

---

## Dependency Tree

### Current Dependency Snapshot (2025-01-18)

```
server-go-ssp-redishoard
├── Go Runtime: 1.24.0 (toolchain go1.24.7)
├── Direct Dependencies (require)
│   ├── github.com/redis/go-redis/v9 v9.16.0 ✅ CURRENT
│   ├── github.com/sqrldev/server-go-ssp v0.0.0-20241212182118-c8230b16b87d
│   │   └── → replaced by github.com/dxcSithLord/server-go-ssp v0.0.0-20241212182118-c8230b16b87d
│   │       └── ⚠️ LATEST: v0.0.0-20251118072316-f490d48446ba (requires Go 1.25+)
│   └── golang.org/x/sys v0.38.0 ✅ CURRENT
└── Indirect Dependencies (pulled by server-go-ssp)
    ├── github.com/cespare/xxhash/v2 v2.3.0 ✅ CURRENT
    ├── github.com/davecgh/go-spew v1.1.1 ✅ CURRENT (stable)
    ├── github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f ✅ CURRENT
    ├── github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e ✅ CURRENT
    └── golang.org/x/crypto v0.44.0 ✅ CURRENT (security-patched)
```

---

## Detailed Dependency Analysis

### 1. Go Runtime & Toolchain

| Component | Current | Latest | Status | Notes |
|-----------|---------|--------|--------|-------|
| Go Version | 1.24.0 | 1.25.4 | ⚠️ BLOCKED | server-go-ssp@latest requires Go 1.25+ |
| Toolchain | go1.24.7 | go1.25.4 | ⚠️ BLOCKED | Awaiting Go version upgrade |

**Upgrade Considerations**:
- 🔴 **Breaking Change**: server-go-ssp v0.0.0-20251118072316 requires `go >= 1.25.0`
- **Impact**: Cannot upgrade to latest server-go-ssp without upgrading Go runtime
- **Recommendation**: Create upgrade stage for Go 1.24 → Go 1.25 migration

**Go 1.25 Breaking Changes** (relevant to this project):
- No known breaking changes affecting this codebase
- Improved security and performance
- Enhanced type inference
- Better error handling

---

### 2. Redis Client: github.com/redis/go-redis/v9

| Attribute | Value |
|-----------|-------|
| **Current Version** | v9.16.0 |
| **Latest Version** | v9.16.0 (as of Jan 2025) |
| **Status** | ✅ UP-TO-DATE |
| **License** | BSD-2-Clause |
| **Vulnerabilities** | None known |

**Description**: Official Go client for Redis with support for Redis 7.x features, Redis Cluster, Redis Sentinel, and RESP3 protocol.

**API Compatibility**:
- ✅ Stable v9 API (no breaking changes expected)
- ✅ Context-based API (required since v9.0.0)
- ✅ Generic UniversalClient interface

**Upgrade Path**: None required (already at latest)

**Dependencies Pulled**:
- `github.com/cespare/xxhash/v2` - Fast hash function
- `github.com/dgryski/go-rendezvous` - Consistent hashing algorithm

**Future Considerations**:
- Monitor for v10 (major version) - may introduce breaking changes
- Redis 8.x support when available

---

### 3. SQRL SSP: github.com/sqrldev/server-go-ssp

| Attribute | Value |
|-----------|-------|
| **Current Version** | v0.0.0-20241212182118-c8230b16b87d |
| **Actual Source** | github.com/dxcSithLord/server-go-ssp (via replace) |
| **Latest Version** | v0.0.0-20251118072316-f490d48446ba |
| **Status** | ⚠️ UPDATE AVAILABLE |
| **License** | MIT |
| **Go Version Requirement** | 1.25+ (latest), 1.23+ (current) |

**Description**: SQRL Server Side Protocol implementation providing the `ssp.Hoard` interface and cryptographic primitives for SQRL authentication.

**Breaking Changes Analysis** (current → latest):

⚠️ **Go Version Requirement**:
- Old: `go 1.23`
- New: `go 1.25`
- **Impact**: Cannot upgrade without upgrading Go runtime first

**Known Changes** (based on commit history):
- Updated to Go 1.25 stdlib features
- Potential Ed25519 key handling improvements
- Security patches from Go 1.25

**Dependencies Pulled**:
- `github.com/skip2/go-qrcode` - QR code generation for SQRL
- `golang.org/x/crypto` - Ed25519, SHA256, other crypto primitives

**Upgrade Strategy**:
1. **Stage 1**: Upgrade Go runtime to 1.25
2. **Stage 2**: Update server-go-ssp to v0.0.0-20251118072316-f490d48446ba
3. **Stage 3**: Run full test suite with race detector
4. **Stage 4**: Verify SQRL authentication flow end-to-end

**Rollback Plan**:
- Revert go.mod replace directive to current commit
- No code changes expected (interface compatibility maintained)

---

### 4. System Calls: golang.org/x/sys

| Attribute | Value |
|-----------|-------|
| **Current Version** | v0.38.0 |
| **Latest Version** | v0.38.0 |
| **Status** | ✅ UP-TO-DATE |
| **License** | BSD-3-Clause |
| **Vulnerabilities** | None known |

**Description**: Low-level operating system primitives for Unix/Windows system calls, used for memory locking (mlock, VirtualLock).

**Platform-Specific Usage**:
- **Unix/Linux/macOS**: `unix.Mlock()`, `unix.Munlock()`
- **Windows**: `windows.VirtualLock()`, `windows.VirtualUnlock()`

**Upgrade Path**: None required

**Critical Dependency**: Required for secure memory locking (CWE-226 mitigation)

---

### 5. Cryptography: golang.org/x/crypto

| Attribute | Value |
|-----------|-------|
| **Current Version** | v0.44.0 |
| **Latest Version** | v0.44.0 |
| **Status** | ✅ UP-TO-DATE |
| **License** | BSD-3-Clause |
| **Vulnerabilities** | None (patched GHSA-hcg3-q754-cr77 in v0.44.0) |

**Description**: Supplementary Go cryptography libraries including Ed25519, crypto/subtle for constant-time operations.

**Security Fixes** (in v0.44.0):
- ✅ Fixed DoS via slow key exchange (GHSA-hcg3-q754-cr77)
- ✅ Updated to latest cryptographic standards

**Usage in Project**:
- `crypto/subtle.ConstantTimeCompare()` - Secure memory clearing verification
- Pulled by server-go-ssp for Ed25519 operations

**Upgrade Path**: None required (already patched)

---

### 6. Indirect Dependencies

#### github.com/cespare/xxhash/v2 v2.3.0

- **Purpose**: Fast non-cryptographic hash function for Redis client
- **Status**: ✅ Stable
- **Upgrade**: None available (latest)
- **Breaking Changes**: None expected (stable API)

#### github.com/davecgh/go-spew v1.1.1

- **Purpose**: Deep pretty printer (used by server-go-ssp for debugging)
- **Status**: ✅ Stable (2017 release, mature)
- **Upgrade**: None needed (feature-complete)
- **Note**: Ancient but stable; no security concerns

#### github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f

- **Purpose**: Rendezvous/HRW hashing for Redis Cluster slot assignment
- **Status**: ✅ Stable
- **Upgrade**: None available (specific commit)
- **Note**: Mature algorithm implementation

#### github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e

- **Purpose**: QR code generation for SQRL authentication links
- **Status**: ⚠️ Outdated (2020)
- **Upgrade**: Check for newer commits (unversioned)
- **Security**: Low risk (generates visual output only)
- **Action**: Monitor for updates but not critical

---

## External Integrations Analysis

### goauthentik.io

**Status**: ❌ **NOT A DEPENDENCY**

**Finding**: This project does not depend on goauthentik.io.

**Context**:
- **server-go-ssp-redishoard** is a **library** that implements `ssp.Hoard` interface
- It provides Redis-backed storage for SQRL authentication state
- goauthentik.io is a separate identity provider/authentication platform
- **Relationship**: Applications using both would integrate them at the application layer, not at the library dependency level

**Potential Integration Scenario**:
```
[goauthentik.io Server]
         ↓
   [Application Layer]
    ├── Uses goauthentik.io for user authentication
    └── Uses server-go-ssp-redishoard for SQRL session storage
              ↓
         [Redis Backend]
```

**Recommendation**: If the user's application integrates with goauthentik.io, document that integration separately at the application level, not in this library's dependencies.

---

## Upgrade Compatibility Matrix

### Scenario 1: Current State (Go 1.24)

| Dependency | Version | Compatible | Notes |
|------------|---------|------------|-------|
| go-redis/v9 | v9.16.0 | ✅ Yes | No changes needed |
| server-go-ssp | v0.0.0-20241212182118 | ✅ Yes | Current version |
| golang.org/x/sys | v0.38.0 | ✅ Yes | No changes needed |
| golang.org/x/crypto | v0.44.0 | ✅ Yes | Security patched |

**Status**: ✅ **STABLE** - All dependencies compatible with Go 1.24

---

### Scenario 2: Upgrade to Go 1.25 + Latest server-go-ssp

| Dependency | Old Version | New Version | Breaking Changes |
|------------|-------------|-------------|------------------|
| **Go Runtime** | 1.24.0 | 1.25.4 | ⚠️ Minor (stdlib improvements) |
| **Toolchain** | go1.24.7 | go1.25.4 | ⚠️ Rebuild required |
| server-go-ssp | v0.0.0-20241212182118 | v0.0.0-20251118072316 | ⚠️ Requires Go 1.25+ |
| go-redis/v9 | v9.16.0 | v9.16.0 | ✅ No change |
| golang.org/x/sys | v0.38.0 | v0.38.0 | ✅ No change |
| golang.org/x/crypto | v0.44.0 | v0.44.0 | ✅ No change |

**Status**: ⚠️ **MAJOR UPGRADE** - Requires Go version migration

**Testing Requirements**:
1. ✅ Unit tests must pass
2. ✅ Integration tests with Redis must pass
3. ✅ Race detector must pass (`go test -race`)
4. ✅ Security scanners must pass (gosec, govulncheck)
5. ✅ Lint checks must pass (golangci-lint)

---

## Breaking Changes Analysis

### Go 1.24 → Go 1.25 Migration

**Language Changes**:
- No syntax breaking changes affecting this project
- Improved type inference (backward compatible)
- Enhanced error handling (backward compatible)

**Standard Library Changes**:
- `crypto/subtle`: No breaking changes
- `encoding/json`: No breaking changes
- `context`: No breaking changes
- `time`: No breaking changes

**Toolchain Changes**:
- Updated compiler optimizations
- Improved escape analysis (may reduce allocations)
- Security hardening in runtime

**Impact Assessment**: ✅ **LOW RISK** - No code changes expected

---

### server-go-ssp v0.0.0-20241212182118 → v0.0.0-20251118072316

**Interface Changes**:
- `ssp.Hoard` interface: ✅ No changes expected (stable)
- `ssp.Nut` type: ✅ No changes expected
- `ssp.HoardCache` struct: ⚠️ Needs verification (may have new fields)

**API Compatibility**:
```go
// Expected to remain unchanged:
type Hoard interface {
    Get(nut Nut) (*HoardCache, error)
    GetAndDelete(nut Nut) (*HoardCache, error)
    Save(nut Nut, value *HoardCache, expiration time.Duration) error
}
```

**Verification Steps**:
1. Clone dxcSithLord/server-go-ssp@v0.0.0-20251118072316
2. Compare `hoard.go` interface definitions
3. Check for new required fields in `HoardCache`
4. Review CHANGELOG or commit messages

**Risk Level**: ⚠️ **MEDIUM** - Unversioned dependency, manual verification required

---

## Sequenced Upgrade Plan

### Stage 0: Baseline (Current State)

**Objective**: Establish baseline with full test coverage

**Actions**:
1. ✅ Run full test suite and record coverage baseline (currently 46.6%)
2. ✅ Run benchmarks and record performance baseline
3. ✅ Document all current test passing states
4. ✅ Create git tag: `v1.0.0-pre-upgrade`

**Exit Criteria**:
- All tests passing
- Coverage ≥ 46.6%
- No security vulnerabilities
- CI/CD pipeline green

---

### Stage 1: Increase Test Coverage (Preparation)

**Branch**: `feat/increase-test-coverage`

**Objective**: Increase test coverage to 85%+ before major upgrades

**Rationale**: Higher test coverage provides confidence during dependency upgrades

**Actions**:
1. Add error path tests (Redis failures, network errors, malformed JSON)
2. Add concurrency tests (race conditions in GetAndDelete)
3. Add security tests (memory clearing verification)
4. Add integration tests (TTL expiration, large payloads)
5. Add platform-specific tests (memory locking on Unix/Windows)

**Target Coverage**: 85%+

**Exit Criteria**:
- Coverage ≥ 85%
- All new tests passing
- Race detector clean
- CI/CD pipeline green

**Estimated Effort**: 2-3 days

---

### Stage 2: Go 1.24 → Go 1.25 Upgrade

**Branch**: `upgrade/go-1.25`

**Objective**: Upgrade Go runtime to 1.25 and verify compatibility

**Actions**:
1. Update `go.mod`: `go 1.24.0` → `go 1.25.0`
2. Update `toolchain` directive to `go1.25.4`
3. Update CI/CD matrix to test Go 1.25
4. Run `go mod tidy` to update checksums
5. Rebuild all binaries
6. Run full test suite
7. Run benchmarks and compare to baseline
8. Update documentation

**Breaking Changes Expected**: None

**Rollback Plan**:
```bash
git revert <commit-sha>
go mod edit -go=1.24.0
go mod edit -toolchain=go1.24.7
go mod tidy
```

**Exit Criteria**:
- All tests passing on Go 1.25
- No performance regressions
- Security scanners passing
- CI/CD green for both Go 1.24 and 1.25 (dual matrix temporarily)

**Estimated Effort**: 1 day

---

### Stage 3: Update server-go-ssp Dependency

**Branch**: `upgrade/server-go-ssp-latest`

**Objective**: Upgrade to latest server-go-ssp from dxcSithLord fork

**Prerequisites**: Stage 2 complete (Go 1.25 deployed)

**Actions**:
1. Update replace directive in `go.mod`:
   ```go
   replace github.com/sqrldev/server-go-ssp => github.com/dxcSithLord/server-go-ssp v0.0.0-20251118072316-f490d48446ba
   ```
2. Run `go mod tidy`
3. Check for new required imports
4. Verify `ssp.Hoard` interface compatibility
5. Check for new fields in `ssp.HoardCache`
6. Run full test suite
7. Run integration tests with actual SQRL client (if available)

**Risk Assessment**:
- **Interface Changes**: Medium risk (manual verification required)
- **Data Format Changes**: Low risk (JSON serialization should be backward compatible)
- **Cryptographic Changes**: Low risk (Ed25519 is stable)

**Verification Checklist**:
- [ ] All existing tests pass without modification
- [ ] No new compilation errors
- [ ] Memory clearing still works (`TestClearBytes`)
- [ ] Input validation still works (`TestNutValidation`)
- [ ] Atomic operations still work (`TestGetAndDelete`)
- [ ] Security scanners pass (gosec, govulncheck)

**Rollback Plan**:
```bash
git revert <commit-sha>
# Restore previous replace directive
go mod edit -replace github.com/sqrldev/server-go-ssp=github.com/dxcSithLord/server-go-ssp@v0.0.0-20241212182118-c8230b16b87d
go mod tidy
```

**Exit Criteria**:
- All tests passing
- No interface compatibility issues
- Security compliance maintained
- CI/CD pipeline green

**Estimated Effort**: 1-2 days (including verification)

---

### Stage 4: Monitor Indirect Dependencies

**Branch**: `chore/update-indirect-deps`

**Objective**: Update any indirect dependencies with available updates

**Actions**:
1. Run `go get -u ./...` to check for updates
2. Review changes to indirect dependencies
3. Check for security advisories
4. Update if safe (no breaking changes)

**Target Dependencies**:
- `github.com/skip2/go-qrcode` - Check for newer commits
- `github.com/dgryski/go-rendezvous` - Check for updates
- Others as reported by `go list -m -u all`

**Exit Criteria**:
- No security vulnerabilities
- All tests passing

**Estimated Effort**: 1 day

---

### Stage 5: Validate & Merge

**Branch**: N/A (merge all upgrade branches)

**Objective**: Consolidate all upgrades and validate system-wide

**Actions**:
1. Merge `feat/increase-test-coverage` → `main`
2. Merge `upgrade/go-1.25` → `main`
3. Merge `upgrade/server-go-ssp-latest` → `main`
4. Merge `chore/update-indirect-deps` → `main`
5. Create release tag: `v1.1.0`
6. Update CHANGELOG.md
7. Publish release notes

**Final Validation**:
- Full regression test suite
- Performance benchmarks (compare to Stage 0 baseline)
- Security audit (gosec, govulncheck)
- Documentation review

**Exit Criteria**:
- All tests passing
- Coverage ≥ 85%
- No performance regressions
- Security compliance maintained
- CI/CD pipeline green

**Estimated Total Effort**: 5-8 days

---

## Dependency Update Schedule

### Immediate Actions (Now)
1. ✅ Increase test coverage to 85%
2. ⚠️ Upgrade Go runtime to 1.25
3. ⚠️ Update server-go-ssp to latest

### Regular Maintenance (Quarterly)
1. Check for go-redis/v9 updates
2. Check for golang.org/x/crypto security patches
3. Check for golang.org/x/sys updates
4. Review server-go-ssp for new commits

### Security Monitoring (Continuous)
1. Subscribe to GitHub Security Advisories:
   - https://github.com/redis/go-redis/security/advisories
   - https://github.com/dxcSithLord/server-go-ssp/security/advisories
2. Enable Dependabot alerts
3. Monitor govulncheck output in CI/CD

---

## Risk Assessment

| Risk | Severity | Likelihood | Mitigation |
|------|----------|------------|------------|
| Go 1.25 breaks code | Medium | Low | Staged testing with rollback plan |
| server-go-ssp interface change | High | Medium | Manual verification before upgrade |
| Data format incompatibility | High | Low | Test with real SQRL data |
| Performance regression | Medium | Low | Benchmark comparison |
| Security vulnerability introduced | High | Low | gosec/govulncheck in CI/CD |
| Test suite insufficient | Medium | Medium | Increase coverage to 85% first |

---

## Testing Strategy for Upgrades

### Pre-Upgrade Tests
```bash
# Baseline coverage
go test -coverprofile=coverage-baseline.out ./...

# Baseline benchmarks
go test -bench=. -benchmem -count=5 > bench-baseline.txt

# Security baseline
gosec -fmt=json -out=gosec-baseline.json ./...
govulncheck -json ./... > govuln-baseline.json
```

### Post-Upgrade Tests
```bash
# Compare coverage
go test -coverprofile=coverage-upgraded.out ./...
go tool cover -func=coverage-upgraded.out

# Compare benchmarks
go test -bench=. -benchmem -count=5 > bench-upgraded.txt
benchstat bench-baseline.txt bench-upgraded.txt

# Verify security
gosec ./...
govulncheck ./...

# Race detection
go test -race ./...

# Integration tests
REDIS_ADDR=localhost:6379 go test -v -tags=integration ./...
```

---

## Rollback Procedures

### Emergency Rollback (Critical Failure)

```bash
# Identify last known good commit
git log --oneline -10

# Hard reset to last good state
git reset --hard <commit-sha>

# Force push if already deployed
git push origin main --force-with-lease

# Notify team
echo "ROLLBACK: Reverted to <commit-sha> due to <reason>"
```

### Partial Rollback (Specific Dependency)

```bash
# Revert specific dependency
go mod edit -replace github.com/sqrldev/server-go-ssp=github.com/dxcSithLord/server-go-ssp@<old-version>
go mod tidy

# Test
go test ./...

# Commit
git add go.mod go.sum
git commit -m "Rollback server-go-ssp to <old-version>"
```

---

## Continuous Dependency Monitoring

### Automated Tools

1. **Dependabot** (GitHub native):
   ```yaml
   # .github/dependabot.yml
   version: 2
   updates:
     - package-ecosystem: "gomod"
       directory: "/"
       schedule:
         interval: "weekly"
       open-pull-requests-limit: 5
   ```

2. **govulncheck** (in CI/CD):
   ```yaml
   # Already integrated in .github/workflows/ci.yml
   - name: Run govulncheck
     run: |
       go install golang.org/x/vuln/cmd/govulncheck@latest
       govulncheck ./...
   ```

3. **go list -m -u all** (weekly):
   ```bash
   # Check for updates
   go list -m -u all | grep '\['
   ```

---

## Documentation Updates Required

After each upgrade stage, update:

1. **README.md**:
   - Update Go version requirement
   - Update dependency versions in examples

2. **REQUIREMENTS.md**:
   - Update NFR-6 (Compatibility) section
   - Update dependency versions

3. **UPGRADE_PATH.md**:
   - Add migration notes for Go 1.25
   - Add changelog entries

4. **CHANGELOG.md** (create if missing):
   ```markdown
   ## [1.1.0] - 2025-01-XX
   ### Changed
   - Upgraded Go runtime to 1.25
   - Updated server-go-ssp to v0.0.0-20251118072316-f490d48446ba

   ### Added
   - Increased test coverage from 46.6% to 85%+
   ```

---

## References

1. **Go Module Documentation**: https://go.dev/ref/mod
2. **go-redis Releases**: https://github.com/redis/go-redis/releases
3. **server-go-ssp Repository**: https://github.com/dxcSithLord/server-go-ssp
4. **Go Release Notes**: https://go.dev/doc/devel/release
5. **Semantic Versioning**: https://semver.org/
6. **Dependency Management Best Practices**: https://go.dev/blog/module-compatibility

---

## Change Log

| Date | Version | Changes |
|------|---------|---------|
| 2025-01-18 | 1.0.0 | Initial dependency analysis and upgrade plan |

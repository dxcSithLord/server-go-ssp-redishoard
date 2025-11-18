# Master TODO List

## Document Control
- **Project**: server-go-ssp-redishoard
- **Version**: 1.0.0
- **Date**: 2025-01-18
- **Status**: Living Document
- **Repository**: https://github.com/dxcSithLord/server-go-ssp-redishoard

---

## Overview

This document consolidates all TODO items, action items, gaps, and recommendations from across the project documentation into a single prioritized task list.

**Sources Consolidated**:
- ✅ REVERSE_ENGINEERED_REQUIREMENTS.md (Gaps & Recommendations)
- ✅ DEPENDENCY_ANALYSIS.md (Sequenced Upgrade Plan)
- ✅ SECURITY_PLAN.md (Implementation Priority Matrix)
- ✅ API_REFERENCE.md (Testing Strategy)
- ✅ ARCHITECTURE_TOGAF.md (Roadmap)
- ✅ Code comments (inline TODOs)

**Priority Levels**:
- 🔴 **CRITICAL** - Security vulnerabilities, breaking issues
- 🟠 **HIGH** - Required for production readiness
- 🟡 **MEDIUM** - Important but not blocking
- 🟢 **LOW** - Nice-to-have improvements

---

## Current Sprint: Production Hardening

### 🔴 CRITICAL Priority

#### CRIT-1: Increase Test Coverage from 46.6% to 85%
**Status**: ⏳ Not Started
**Effort**: 2-3 days
**Owner**: TBD
**Due Date**: Before any dependency upgrades

**Rationale**: Higher coverage provides confidence during Go 1.25 and server-go-ssp upgrades.

**Tasks**:
- [ ] Add error path tests
  - [ ] Redis connection failures
  - [ ] Network timeout errors
  - [ ] Malformed JSON deserialization
  - [ ] Invalid nut formats
  - [ ] Nil HoardCache handling
- [ ] Add concurrency tests
  - [ ] Race conditions in GetAndDelete
  - [ ] Concurrent Save/Get operations
  - [ ] Connection pool exhaustion
- [ ] Add security tests
  - [ ] Verify ClearBytes zeros memory
  - [ ] Verify LockMemory works on supported platforms
  - [ ] Test that sensitive data is not logged
- [ ] Add integration tests
  - [ ] TTL expiration verification
  - [ ] Large payload handling (>1MB)
  - [ ] Redis Cluster connectivity
  - [ ] Redis Sentinel failover
- [ ] Add platform-specific tests
  - [ ] Memory locking on Linux (requires CAP_IPC_LOCK)
  - [ ] Memory locking on Windows (requires elevation)
  - [ ] Graceful degradation on unsupported platforms

**Exit Criteria**:
- Coverage ≥ 85%
- All tests passing
- Race detector clean (`go test -race`)
- CI/CD pipeline green

**References**:
- REVERSE_ENGINEERED_REQUIREMENTS.md § Gap 1
- DEPENDENCY_ANALYSIS.md § Stage 1

---

#### CRIT-2: Verify No Security Regressions
**Status**: ⏳ Not Started
**Effort**: 1 day
**Owner**: TBD
**Due Date**: Before each release

**Tasks**:
- [ ] Run gosec security scanner
- [ ] Run govulncheck for vulnerabilities
- [ ] Verify CWE-200 mitigation (no sensitive logging)
- [ ] Verify CWE-226 mitigation (memory clearing)
- [ ] Review Redis TLS configuration in documentation
- [ ] Test input validation against injection attacks

**Exit Criteria**:
- No HIGH or CRITICAL gosec findings
- No unfixed vulnerabilities in govulncheck
- Security compliance verified

**References**:
- SECURITY_PLAN.md § Verification Checklist

---

### 🟠 HIGH Priority

#### HIGH-1: Upgrade to Go 1.25
**Status**: ⏳ Not Started
**Effort**: 1 day
**Owner**: TBD
**Due Date**: After CRIT-1 complete
**Branch**: `upgrade/go-1.25`

**Prerequisites**:
- ✅ Test coverage ≥ 85% (CRIT-1)
- ✅ Baseline benchmarks recorded

**Tasks**:
- [ ] Update `go.mod`: `go 1.24.0` → `go 1.25.0`
- [ ] Update `toolchain` directive to `go1.25.4`
- [ ] Update CI/CD matrix to include Go 1.25
- [ ] Run `go mod tidy`
- [ ] Rebuild all binaries
- [ ] Run full test suite
- [ ] Run benchmarks and compare to baseline
- [ ] Update README.md with Go 1.25 requirement
- [ ] Update REQUIREMENTS.md NFR-6 section

**Rollback Plan**:
```bash
git revert <commit-sha>
go mod edit -go=1.24.0
go mod edit -toolchain=go1.24.7
go mod tidy
```

**Exit Criteria**:
- All tests passing on Go 1.25
- No performance regressions (>5%)
- CI/CD green for Go 1.25
- Documentation updated

**References**:
- DEPENDENCY_ANALYSIS.md § Stage 2

---

#### HIGH-2: Update server-go-ssp to Latest
**Status**: ⏳ Blocked by HIGH-1
**Effort**: 1-2 days
**Owner**: TBD
**Due Date**: After HIGH-1 complete
**Branch**: `upgrade/server-go-ssp-latest`

**Prerequisites**:
- ✅ Go 1.25 deployed (HIGH-1)

**Tasks**:
- [ ] Update replace directive in `go.mod`:
  ```go
  replace github.com/sqrldev/server-go-ssp => github.com/dxcSithLord/server-go-ssp v0.0.0-20251118072316-f490d48446ba
  ```
- [ ] Run `go mod tidy`
- [ ] Verify `ssp.Hoard` interface compatibility
- [ ] Check for new fields in `ssp.HoardCache`
- [ ] Run full test suite
- [ ] Run integration tests with SQRL client (if available)
- [ ] Update documentation with new version

**Verification Checklist**:
- [ ] All existing tests pass without modification
- [ ] No new compilation errors
- [ ] Memory clearing still works
- [ ] Input validation still works
- [ ] Atomic operations still work
- [ ] Security scanners pass

**Rollback Plan**:
```bash
git revert <commit-sha>
go mod edit -replace github.com/sqrldev/server-go-ssp=github.com/dxcSithLord/server-go-ssp@v0.0.0-20241212182118-c8230b16b87d
go mod tidy
```

**Exit Criteria**:
- All tests passing
- No interface compatibility issues
- Security compliance maintained

**References**:
- DEPENDENCY_ANALYSIS.md § Stage 3

---

#### HIGH-3: Create Production Deployment Guide
**Status**: ⏳ Not Started
**Effort**: 1 day
**Owner**: TBD
**Due Date**: Before v1.0.0 release

**Tasks**:
- [ ] Document Redis production configuration
  - [ ] TLS setup instructions
  - [ ] Authentication configuration
  - [ ] Memory limits (maxmemory, eviction policies)
  - [ ] Connection pool sizing
  - [ ] Monitoring setup
- [ ] Document application integration
  - [ ] Example Dockerfile
  - [ ] Kubernetes deployment YAML
  - [ ] Environment variable configuration
  - [ ] Health check endpoints (at application level)
- [ ] Document operational procedures
  - [ ] Backup and restore (if needed)
  - [ ] Scaling guidelines
  - [ ] Troubleshooting guide
  - [ ] Performance tuning

**Deliverables**:
- [ ] DEPLOYMENT_GUIDE.md
- [ ] Example configurations in `examples/` directory

**References**:
- REVERSE_ENGINEERED_REQUIREMENTS.md § Gap 3

---

### 🟡 MEDIUM Priority

#### MED-1: Add Observability Features
**Status**: ⏳ Not Started
**Effort**: 2-3 days
**Owner**: TBD
**Due Date**: v1.1.0 milestone

**Rationale**: Production monitoring requires metrics for operation counters, latencies, and error rates.

**Tasks**:
- [ ] Design metrics interface (Prometheus-compatible)
- [ ] Add operation counters
  - [ ] `redishoard_save_total` (counter)
  - [ ] `redishoard_get_total` (counter)
  - [ ] `redishoard_getanddelete_total` (counter)
- [ ] Add latency histograms
  - [ ] `redishoard_save_duration_seconds` (histogram)
  - [ ] `redishoard_get_duration_seconds` (histogram)
  - [ ] `redishoard_getanddelete_duration_seconds` (histogram)
- [ ] Add error rate tracking
  - [ ] `redishoard_errors_total{operation,type}` (counter)
- [ ] Create example Prometheus integration
- [ ] Document Grafana dashboard setup

**Design Considerations**:
- Make metrics optional (no required dependency on Prometheus client)
- Use interface injection pattern for metrics recorder
- Provide no-op default implementation
- Document how to integrate with application's metrics

**Exit Criteria**:
- Metrics interface defined
- Example implementation provided
- Documentation complete

**References**:
- REVERSE_ENGINEERED_REQUIREMENTS.md § Gap 2

---

#### MED-2: Improve Error Messages
**Status**: ⏳ Not Started
**Effort**: 1 day
**Owner**: TBD
**Due Date**: v1.1.0 milestone

**Tasks**:
- [ ] Add structured error types (instead of string wrapping)
- [ ] Include error codes for programmatic handling
- [ ] Distinguish transient vs permanent errors
- [ ] Add context to Redis errors (which operation failed)
- [ ] Document all error types in API_REFERENCE.md

**Example**:
```go
type NutValidationError struct {
    Nut    string
    Reason string
    Code   string // "NUT_TOO_SHORT", "NUT_INVALID_CHARS", etc.
}
```

**Exit Criteria**:
- All public functions document possible errors
- Error types are type-assertable
- Application can distinguish error types

**References**:
- API_REFERENCE.md § Error Handling

---

#### MED-3: Add Benchmarking Suite
**Status**: ⏳ Not Started
**Effort**: 1 day
**Owner**: TBD
**Due Date**: v1.1.0 milestone

**Tasks**:
- [ ] Expand existing benchmarks
- [ ] Add benchmarks for different payload sizes
  - [ ] Small (256 bytes)
  - [ ] Medium (4 KB)
  - [ ] Large (64 KB)
- [ ] Add benchmarks for concurrent operations
- [ ] Add benchmarks with Redis latency simulation
- [ ] Document expected performance baselines
- [ ] Add benchmark comparison to CI/CD

**Exit Criteria**:
- Comprehensive benchmark coverage
- Performance baselines documented
- CI tracks performance regressions

**References**:
- REVERSE_ENGINEERED_REQUIREMENTS.md § NFR-1

---

#### MED-4: Enhance Documentation with More Examples
**Status**: ⏳ Not Started
**Effort**: 2 days
**Owner**: TBD
**Due Date**: v1.1.0 milestone

**Tasks**:
- [ ] Create `examples/` directory
- [ ] Add example: Basic SQRL authentication flow
- [ ] Add example: Redis Cluster setup
- [ ] Add example: Redis Sentinel setup
- [ ] Add example: TLS configuration
- [ ] Add example: Docker Compose setup (app + Redis)
- [ ] Add example: Kubernetes deployment
- [ ] Add example: Metrics integration
- [ ] Add example: Custom error handling

**Deliverables**:
- [ ] `examples/basic/` - Simple usage
- [ ] `examples/cluster/` - Redis Cluster
- [ ] `examples/tls/` - TLS configuration
- [ ] `examples/docker/` - Docker Compose
- [ ] `examples/kubernetes/` - K8s manifests

**References**:
- API_REFERENCE.md § Complete Usage Example

---

### 🟢 LOW Priority

#### LOW-1: Update Indirect Dependencies
**Status**: ⏳ Not Started
**Effort**: 1 day
**Owner**: TBD
**Due Date**: v1.2.0 milestone
**Branch**: `chore/update-indirect-deps`

**Tasks**:
- [ ] Run `go get -u ./...` to check for updates
- [ ] Review changes to indirect dependencies
- [ ] Check `github.com/skip2/go-qrcode` for newer commits
- [ ] Update if safe (no breaking changes)
- [ ] Run full test suite
- [ ] Update DEPENDENCY_ANALYSIS.md

**Exit Criteria**:
- No security vulnerabilities in dependencies
- All tests passing

**References**:
- DEPENDENCY_ANALYSIS.md § Stage 4

---

#### LOW-2: Create CHANGELOG.md
**Status**: ⏳ Not Started
**Effort**: 1 hour
**Owner**: TBD
**Due Date**: v1.0.0 release

**Tasks**:
- [ ] Create CHANGELOG.md following Keep a Changelog format
- [ ] Document all changes since project inception
- [ ] Add sections: Added, Changed, Fixed, Security
- [ ] Link to commit SHAs and PRs

**Template**:
```markdown
# Changelog

## [Unreleased]

## [1.0.0] - 2025-01-XX
### Added
- Initial release
- Secure memory handling (CWE-226 mitigation)
- Input validation (injection prevention)
- Atomic GetAndDelete operation

### Security
- Fixed CWE-200: Removed sensitive data logging
- Fixed CWE-226: Added memory clearing
- Added platform-specific memory locking
```

**References**:
- DEPENDENCY_ANALYSIS.md § Documentation Updates Required

---

#### LOW-3: Add Contributing Guidelines
**Status**: ⏳ Not Started
**Effort**: 2 hours
**Owner**: TBD
**Due Date**: v1.0.0 release

**Tasks**:
- [ ] Create CONTRIBUTING.md
- [ ] Document code style guidelines
- [ ] Document commit message format
- [ ] Document PR process
- [ ] Document testing requirements
- [ ] Document security disclosure process

**References**:
- README.md § Contributing

---

#### LOW-4: Setup GitHub Issue Templates
**Status**: ⏳ Not Started
**Effort**: 1 hour
**Owner**: TBD
**Due Date**: v1.0.0 release

**Tasks**:
- [ ] Create `.github/ISSUE_TEMPLATE/bug_report.md`
- [ ] Create `.github/ISSUE_TEMPLATE/feature_request.md`
- [ ] Create `.github/ISSUE_TEMPLATE/security.md`
- [ ] Create `.github/PULL_REQUEST_TEMPLATE.md`

**References**:
- GitHub best practices

---

#### LOW-5: Add Code Coverage Badge
**Status**: ⏳ Not Started
**Effort**: 30 minutes
**Owner**: TBD
**Due Date**: After CRIT-1 complete

**Tasks**:
- [ ] Setup Codecov or Coveralls integration
- [ ] Add coverage badge to README.md
- [ ] Configure coverage reporting in CI/CD

**References**:
- README.md § Badges

---

## Backlog: Future Enhancements

### FUTURE-1: Support for Alternative Backends
**Status**: 💡 Idea
**Effort**: 1-2 weeks
**Motivation**: Allow users to choose backend based on infrastructure

**Potential Backends**:
- Memcached
- etcd
- PostgreSQL with UNLOGGED tables
- DynamoDB

**Design**:
- Create backend interface
- Implement Redis backend (refactor existing)
- Add pluggable backend system

---

### FUTURE-2: Add Distributed Tracing Support
**Status**: 💡 Idea
**Effort**: 1 week
**Motivation**: Observability in microservices environments

**Tasks**:
- Add OpenTelemetry integration
- Trace Redis operations
- Trace serialization/deserialization
- Document Jaeger/Zipkin setup

---

### FUTURE-3: Add Compression Support
**Status**: 💡 Idea
**Effort**: 3 days
**Motivation**: Reduce Redis memory usage for large HoardCache objects

**Tasks**:
- Add optional compression (gzip, snappy, zstd)
- Benchmark compression overhead vs memory savings
- Make compression opt-in via configuration

---

## Dependency Upgrade Roadmap

### Immediate (Before v1.0.0 Release)
- [ ] CRIT-1: Increase test coverage to 85%
- [ ] HIGH-1: Upgrade to Go 1.25
- [ ] HIGH-2: Update server-go-ssp to latest

### Short-term (v1.1.0 - Q1 2025)
- [ ] MED-1: Add observability features
- [ ] MED-2: Improve error messages
- [ ] MED-3: Add benchmarking suite
- [ ] MED-4: Enhance documentation

### Medium-term (v1.2.0 - Q2 2025)
- [ ] LOW-1: Update indirect dependencies
- [ ] Monitor for go-redis/v10 release

### Long-term (v2.0.0 - Q3 2025+)
- [ ] FUTURE-1: Alternative backends
- [ ] FUTURE-2: Distributed tracing
- [ ] FUTURE-3: Compression support

---

## Testing Requirements Matrix

| Task | Unit Tests | Integration Tests | Security Tests | Benchmarks | Documentation |
|------|-----------|-------------------|----------------|------------|---------------|
| CRIT-1 | Required | Required | Required | Optional | Optional |
| HIGH-1 | Required | Required | Required | Required | Required |
| HIGH-2 | Required | Required | Required | Required | Required |
| MED-1 | Required | Optional | N/A | Optional | Required |
| MED-2 | Required | Optional | N/A | N/A | Required |
| MED-3 | N/A | N/A | N/A | Required | Required |
| MED-4 | Optional | Optional | N/A | N/A | Required |

---

## Risk Register

| Task | Risk | Likelihood | Impact | Mitigation |
|------|------|------------|--------|------------|
| HIGH-1 | Go 1.25 breaks code | Low | High | Staged testing with comprehensive test suite |
| HIGH-2 | server-go-ssp interface change | Medium | High | Manual verification before upgrade |
| CRIT-1 | Insufficient test coverage | Low | Medium | Start with error paths, then expand |
| MED-1 | Performance overhead from metrics | Low | Low | Make metrics optional, benchmark before/after |
| FUTURE-1 | Backend abstraction too complex | Medium | Medium | Start with interface design, iterate |

---

## Conflict Resolution

**No conflicting objectives identified** across all documentation sources.

**Alignment Verified**:
- ✅ All security requirements aligned (SECURITY_PLAN.md ↔ REQUIREMENTS.md)
- ✅ All upgrade paths sequenced properly (DEPENDENCY_ANALYSIS.md)
- ✅ All architecture decisions support business objectives (ARCHITECTURE_TOGAF.md)
- ✅ All API documentation matches implementation (API_REFERENCE.md ↔ code)

**Potential Questions for User**:

1. **Test Coverage Target**: Confirm 85% is acceptable, or adjust based on risk tolerance?
2. **Metrics Library**: Prefer Prometheus client, OpenTelemetry, or custom interface?
3. **Release Cadence**: Monthly releases, quarterly, or ad-hoc?
4. **Breaking Changes**: Acceptable for v2.0.0, or maintain strict backward compatibility?
5. **Alternative Backends**: Priority for FUTURE-1, or focus on Redis optimizations?

---

## Progress Tracking

### Completed ✅
- ✅ Reverse-engineered requirements from code
- ✅ Created comprehensive security plan
- ✅ Fixed CWE-200 (sensitive logging)
- ✅ Fixed CWE-226 (memory clearing)
- ✅ Added platform-specific memory locking
- ✅ Created dependency analysis
- ✅ Created staged upgrade plan
- ✅ Created TOGAF architecture documentation
- ✅ Created API reference documentation
- ✅ Updated all sqrldev references to dxcSithLord
- ✅ Configured golangci-lint
- ✅ Setup CI/CD pipeline with security scanning
- ✅ Fixed all lint errors
- ✅ Updated vulnerable dependencies

### In Progress 🔄
- 🔄 Master TODO list consolidation (this document)

### Not Started ⏳
- ⏳ Increase test coverage to 85% (CRIT-1)
- ⏳ Upgrade to Go 1.25 (HIGH-1)
- ⏳ Update server-go-ssp (HIGH-2)
- ⏳ Production deployment guide (HIGH-3)
- ⏳ All MED, LOW, and FUTURE tasks

---

## Next Actions

**Immediate Next Steps**:

1. **Review this TODO list with stakeholders**
   - Confirm priorities
   - Adjust timelines
   - Assign owners

2. **Start CRIT-1: Increase Test Coverage**
   - Create branch: `feat/increase-test-coverage`
   - Begin with error path tests
   - Target: 85% coverage

3. **Prepare for HIGH-1: Go 1.25 Upgrade**
   - Record baseline metrics
   - Review Go 1.25 release notes
   - Plan testing strategy

4. **Create v1.0.0 Release Checklist**
   - [ ] CRIT-1 complete
   - [ ] CRIT-2 complete
   - [ ] HIGH-1 complete
   - [ ] HIGH-2 complete
   - [ ] HIGH-3 complete
   - [ ] LOW-2 complete (CHANGELOG)
   - [ ] LOW-3 complete (CONTRIBUTING)
   - [ ] Documentation review
   - [ ] Security audit
   - [ ] Performance validation

---

## References

### Documentation Sources
- [REVERSE_ENGINEERED_REQUIREMENTS.md](REVERSE_ENGINEERED_REQUIREMENTS.md)
- [DEPENDENCY_ANALYSIS.md](DEPENDENCY_ANALYSIS.md)
- [SECURITY_PLAN.md](SECURITY_PLAN.md)
- [API_REFERENCE.md](API_REFERENCE.md)
- [ARCHITECTURE_TOGAF.md](ARCHITECTURE_TOGAF.md)
- [README.md](README.md)

### External References
- **Keep a Changelog**: https://keepachangelog.com/
- **Semantic Versioning**: https://semver.org/
- **GitHub Issue Templates**: https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests

---

## Change Log

| Date | Version | Changes |
|------|---------|---------|
| 2025-01-18 | 1.0.0 | Initial consolidated TODO list from all documentation sources |

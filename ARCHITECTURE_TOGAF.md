# TOGAF Architecture Documentation

## Document Control
- **Project**: server-go-ssp-redishoard
- **Version**: 1.0.0
- **Date**: 2025-01-18
- **Framework**: TOGAF 9.2
- **Repository**: https://github.com/dxcSithLord/server-go-ssp-redishoard

---

## Executive Summary

This document provides a TOGAF-compliant architecture view of the **server-go-ssp-redishoard** project, mapping business objectives to functional requirements, technical implementations, and deployment architecture.

**TOGAF Architecture Domains Covered**:
1. **Business Architecture** - Business objectives and requirements
2. **Data Architecture** - Data models and flows
3. **Application Architecture** - Component interactions
4. **Technology Architecture** - Platform and infrastructure

---

## Business Architecture

### Business Objectives → Functional Requirements Mapping

```mermaid
graph TD
    %% Business Objectives
    BO1[BO-1: Enable Passwordless Authentication]
    BO2[BO-2: Protect Sensitive Authentication Data]
    BO3[BO-3: Prevent Replay Attacks]
    BO4[BO-4: Production Readiness]

    %% Functional Requirements
    FR1.1[FR-1.1: Save Operation]
    FR1.2[FR-1.2: Get Operation]
    FR1.3[FR-1.3: GetAndDelete Operation]
    FR2.1[FR-2.1: Nut Validation]
    FR3.1[FR-3.1: Secure Memory Clearing]
    FR3.2[FR-3.2: Platform-Specific Memory Locking]

    %% Non-Functional Requirements
    NFR1[NFR-1: Performance &lt;10ms]
    NFR2[NFR-2: Horizontal Scalability]
    NFR3[NFR-3: Security Compliance]
    NFR4[NFR-4: Reliability]
    NFR5[NFR-5: Maintainability]

    %% Relationships
    BO1 --> FR1.1
    BO1 --> FR1.2
    BO1 --> FR1.3
    BO1 --> NFR1
    BO1 --> NFR2

    BO2 --> FR3.1
    BO2 --> FR3.2
    BO2 --> NFR3

    BO3 --> FR1.3
    BO3 --> FR2.1

    BO4 --> NFR4
    BO4 --> NFR5

    %% Styling
    classDef businessObj fill:#e1f5ff,stroke:#01579b,stroke-width:2px
    classDef functionalReq fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    classDef nonfunctionalReq fill:#fff3e0,stroke:#e65100,stroke-width:2px

    class BO1,BO2,BO3,BO4 businessObj
    class FR1.1,FR1.2,FR1.3,FR2.1,FR3.1,FR3.2 functionalReq
    class NFR1,NFR2,NFR3,NFR4,NFR5 nonfunctionalReq
```

---

## Application Architecture

### Component Interaction Diagram

```mermaid
graph TB
    %% External Actors
    APP[Application Layer<br/>SQRL Server]
    CLIENT[SQRL Client<br/>User Device]

    %% Library Components
    subgraph redishoard [redishoard Library]
        HOARD[Hoard<br/>Main Interface]
        SECURE[Security Module<br/>ClearBytes, LockMemory]
        VALIDATE[Validation Module<br/>ValidateNut]
    end

    %% Dependencies
    SSP[server-go-ssp<br/>Interface Definition]
    REDIS_CLIENT[go-redis/v9<br/>Redis Client]

    %% Data Store
    REDIS[(Redis<br/>Data Store)]

    %% Interactions
    CLIENT -->|1. Authentication Request| APP
    APP -->|2. Generate Nut| SSP
    APP -->|3. Save HoardCache| HOARD
    HOARD -->|4. Validate| VALIDATE
    HOARD -->|5. Serialize JSON| HOARD
    HOARD -->|6. Clear Memory| SECURE
    HOARD -->|7. SET key TTL| REDIS_CLIENT
    REDIS_CLIENT -->|8. Store| REDIS

    APP -->|9. Get Session| HOARD
    HOARD -->|10. GET key| REDIS_CLIENT
    REDIS_CLIENT -->|11. Retrieve| REDIS
    HOARD -->|12. Deserialize| HOARD
    HOARD -->|13. Clear Memory| SECURE

    APP -->|14. Authenticate| HOARD
    HOARD -->|15. MULTI/GET/DEL/EXEC| REDIS_CLIENT
    REDIS_CLIENT -->|16. Atomic Operation| REDIS
    HOARD -->|17. Clear Memory| SECURE
    APP -->|18. Authentication Success| CLIENT

    %% Styling
    classDef actor fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    classDef component fill:#e3f2fd,stroke:#1565c0,stroke-width:2px
    classDef dependency fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    classDef datastore fill:#fce4ec,stroke:#c2185b,stroke-width:3px

    class CLIENT,APP actor
    class HOARD,SECURE,VALIDATE component
    class SSP,REDIS_CLIENT dependency
    class REDIS datastore
```

---

## Data Architecture

### Data Flow Diagram: Save Operation

```mermaid
sequenceDiagram
    participant App as Application
    participant Hoard as Hoard.Save()
    participant Validate as ValidateNut()
    participant JSON as json.Marshal()
    participant Secure as ClearBytes()
    participant Redis as Redis Client
    participant Store as Redis Store

    App->>Hoard: Save(nut, cache, TTL)

    activate Hoard
    Hoard->>Validate: ValidateNut(nut)
    activate Validate
    Validate-->>Hoard: nil (valid)
    deactivate Validate

    Hoard->>Hoard: Check cache != nil

    Hoard->>JSON: Marshal(cache)
    activate JSON
    JSON-->>Hoard: jsonBytes
    deactivate JSON

    Note over Hoard: jsonBytes contains<br/>sensitive data

    Hoard->>Redis: SET nut jsonBytes EX TTL
    activate Redis
    Redis->>Store: Store with expiration
    activate Store
    Store-->>Redis: OK
    deactivate Store
    Redis-->>Hoard: nil (success)
    deactivate Redis

    Hoard->>Secure: ClearBytes(jsonBytes)
    activate Secure
    Note over Secure: Zero memory<br/>Prevent reuse
    Secure-->>Hoard: (cleared)
    deactivate Secure

    Hoard-->>App: nil (success)
    deactivate Hoard
```

### Data Flow Diagram: GetAndDelete Operation (Atomic)

```mermaid
sequenceDiagram
    participant App as Application
    participant Hoard as Hoard.GetAndDelete()
    participant Validate as ValidateNut()
    participant Redis as Redis TxPipeline
    participant Store as Redis Store
    participant Secure as ClearBytes()

    App->>Hoard: GetAndDelete(nut)

    activate Hoard
    Hoard->>Validate: ValidateNut(nut)
    activate Validate
    Validate-->>Hoard: nil (valid)
    deactivate Validate

    Note over Hoard,Redis: Begin Atomic Transaction

    Hoard->>Redis: TxPipelined()
    activate Redis

    Redis->>Redis: MULTI
    Redis->>Store: GET nut
    Redis->>Store: DEL nut
    Redis->>Redis: EXEC

    activate Store
    Store-->>Redis: [data, 1]
    deactivate Store

    Redis-->>Hoard: [StringCmd, IntCmd]
    deactivate Redis

    Note over Hoard: Extract data from<br/>StringCmd result

    Hoard->>Hoard: fromBytes(data)

    Note over Hoard: data contains<br/>sensitive info

    Hoard->>Secure: ClearBytes(data)
    activate Secure
    Secure-->>Hoard: (cleared)
    deactivate Secure

    Hoard-->>App: HoardCache, nil
    deactivate Hoard

    Note over Store: Nut now deleted<br/>Cannot be reused
```

### Data Model: HoardCache Structure

```mermaid
classDiagram
    class HoardCache {
        +String State
        +String RemoteIP
        +Nut OriginalNut
        +Nut PagNut
        +*CliRequest LastRequest
        +*SqrlIdentity Identity
        +[]byte LastResponse
    }

    class CliRequest {
        +String Ver
        +String Cmd
        +String Opt
        +[]byte Idk
        +[]byte Pidk
        +[]byte Suk
        +[]byte Vuk
        +[]byte Ids
        +[]byte Pids
        +[]byte Urs
    }

    class SqrlIdentity {
        +[]byte Idk
        +[]byte Suk
        +[]byte Vuk
        +time.Time Created
    }

    class Nut {
        <<type>>
        +String value
    }

    HoardCache --> CliRequest : contains
    HoardCache --> SqrlIdentity : contains
    HoardCache --> Nut : references

    note for HoardCache "Sensitive Fields:\n- LastRequest (signatures)\n- Identity (Ed25519 keys)\n- LastResponse (tokens)\n\n⚠️ Must be cleared after use"
```

### Redis Key Schema

```mermaid
graph LR
    subgraph Redis Storage
        KEY1["{nut-value-1}"]
        KEY2["{nut-value-2}"]
        KEY3["{nut-value-3}"]
    end

    subgraph Value Format
        JSON["JSON: HoardCache<br/>{<br/>  State: '...',<br/>  RemoteIP: '...',<br/>  ...<br/>}"]
    end

    subgraph Expiration
        TTL["TTL: 300s<br/>Auto-delete"]
    end

    KEY1 --> JSON
    KEY1 --> TTL
    KEY2 --> JSON
    KEY2 --> TTL
    KEY3 --> JSON
    KEY3 --> TTL

    style KEY1 fill:#ffebee,stroke:#c62828
    style KEY2 fill:#ffebee,stroke:#c62828
    style KEY3 fill:#ffebee,stroke:#c62828
    style JSON fill:#e8f5e9,stroke:#2e7d32
    style TTL fill:#fff3e0,stroke:#ef6c00
```

---

## Technology Architecture

### Deployment Architecture: Single Redis Instance

```mermaid
graph TB
    subgraph Internet
        CLIENT1[SQRL Client 1]
        CLIENT2[SQRL Client 2]
        CLIENT3[SQRL Client 3]
    end

    subgraph DMZ [Application Layer]
        APP1[App Server 1<br/>Go 1.24+<br/>redishoard]
        APP2[App Server 2<br/>Go 1.24+<br/>redishoard]
        APP3[App Server 3<br/>Go 1.24+<br/>redishoard]
    end

    subgraph Backend [Data Layer]
        REDIS[Redis 7.2+<br/>Port 6379<br/>TLS Enabled]
    end

    CLIENT1 --> APP1
    CLIENT2 --> APP2
    CLIENT3 --> APP3

    APP1 -->|Connection Pool| REDIS
    APP2 -->|Connection Pool| REDIS
    APP3 -->|Connection Pool| REDIS

    style REDIS fill:#fce4ec,stroke:#c2185b,stroke-width:3px
    style APP1 fill:#e3f2fd,stroke:#1565c0,stroke-width:2px
    style APP2 fill:#e3f2fd,stroke:#1565c0,stroke-width:2px
    style APP3 fill:#e3f2fd,stroke:#1565c0,stroke-width:2px
```

### Deployment Architecture: Redis Cluster (High Availability)

```mermaid
graph TB
    subgraph "Application Tier"
        APP1[App Server 1]
        APP2[App Server 2]
        APP3[App Server 3]
    end

    subgraph "Redis Cluster"
        subgraph "Master Nodes"
            M1[Master 1<br/>Slots 0-5461]
            M2[Master 2<br/>Slots 5462-10922]
            M3[Master 3<br/>Slots 10923-16383]
        end

        subgraph "Replica Nodes"
            R1[Replica 1<br/>for Master 1]
            R2[Replica 2<br/>for Master 2]
            R3[Replica 3<br/>for Master 3]
        end
    end

    APP1 --> M1
    APP1 --> M2
    APP1 --> M3
    APP2 --> M1
    APP2 --> M2
    APP2 --> M3
    APP3 --> M1
    APP3 --> M2
    APP3 --> M3

    M1 -.->|Replication| R1
    M2 -.->|Replication| R2
    M3 -.->|Replication| R3

    style M1 fill:#c8e6c9,stroke:#2e7d32,stroke-width:2px
    style M2 fill:#c8e6c9,stroke:#2e7d32,stroke-width:2px
    style M3 fill:#c8e6c9,stroke:#2e7d32,stroke-width:2px
    style R1 fill:#fff9c4,stroke:#f57f17,stroke-width:2px
    style R2 fill:#fff9c4,stroke:#f57f17,stroke-width:2px
    style R3 fill:#fff9c4,stroke:#f57f17,stroke-width:2px
```

### Platform Architecture: Build Targets

```mermaid
graph TB
    subgraph "Source Code"
        GO[Go Source<br/>hoard.go<br/>secure.go]
    end

    subgraph "Build Tags"
        UNIX[secure_unix.go<br/>//+build unix]
        WIN[secure_windows.go<br/>//+build windows]
        OTHER[secure_other.go<br/>//+build !unix,!windows]
    end

    subgraph "Target Platforms"
        LINUX[Linux<br/>amd64, arm64<br/>mlock()]
        MACOS[macOS<br/>amd64, arm64<br/>mlock()]
        WINDOWS[Windows<br/>amd64<br/>VirtualLock()]
        FREEBSD[FreeBSD<br/>amd64<br/>No locking]
    end

    GO --> UNIX
    GO --> WIN
    GO --> OTHER

    UNIX --> LINUX
    UNIX --> MACOS
    WIN --> WINDOWS
    OTHER --> FREEBSD

    style GO fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style UNIX fill:#e3f2fd,stroke:#1565c0,stroke-width:2px
    style WIN fill:#e3f2fd,stroke:#1565c0,stroke-width:2px
    style OTHER fill:#e3f2fd,stroke:#1565c0,stroke-width:2px
```

---

## Security Architecture

### Threat Model and Mitigations

```mermaid
graph TB
    subgraph Threats
        T1[CWE-200<br/>Information Exposure]
        T2[CWE-226<br/>Sensitive Data in Memory]
        T3[CWE-312<br/>Cleartext Storage]
        T4[Replay Attacks]
        T5[Injection Attacks]
        T6[MITM Attacks]
    end

    subgraph Mitigations
        M1[No Sensitive Logging]
        M2[ClearBytes Function]
        M3[Memory Locking]
        M4[Redis TLS]
        M5[GetAndDelete Atomic]
        M6[ValidateNut]
    end

    subgraph Implementation
        I1[hoard.go L44-45]
        I2[secure.go L8-27]
        I3[secure_*.go]
        I4[Redis Client TLS Config]
        I5[hoard.go L50-108]
        I6[secure.go L87-108]
    end

    T1 --> M1
    M1 --> I1

    T2 --> M2
    T2 --> M3
    M2 --> I2
    M3 --> I3

    T3 --> M4
    M4 --> I4

    T4 --> M5
    M5 --> I5

    T5 --> M6
    M6 --> I6

    T6 --> M4

    style T1 fill:#ffebee,stroke:#c62828,stroke-width:2px
    style T2 fill:#ffebee,stroke:#c62828,stroke-width:2px
    style T3 fill:#ffebee,stroke:#c62828,stroke-width:2px
    style T4 fill:#ffebee,stroke:#c62828,stroke-width:2px
    style T5 fill:#ffebee,stroke:#c62828,stroke-width:2px
    style T6 fill:#ffebee,stroke:#c62828,stroke-width:2px

    style M1 fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    style M2 fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    style M3 fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    style M4 fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    style M5 fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    style M6 fill:#fff3e0,stroke:#ef6c00,stroke-width:2px

    style I1 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style I2 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style I3 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style I4 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style I5 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style I6 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
```

---

## Capability Architecture

### Functional Capabilities Matrix

```mermaid
graph LR
    subgraph Business Capabilities
        BC1[Session Management]
        BC2[Authentication State Storage]
        BC3[Replay Prevention]
        BC4[Data Protection]
    end

    subgraph Application Capabilities
        AC1[Save State]
        AC2[Retrieve State]
        AC3[Consume State Atomically]
        AC4[Validate Input]
        AC5[Clear Sensitive Memory]
    end

    subgraph Technical Capabilities
        TC1[Redis SET/GET]
        TC2[Redis MULTI/EXEC]
        TC3[JSON Serialization]
        TC4[Memory Locking]
        TC5[TLS Connection]
    end

    BC1 --> AC1
    BC1 --> AC2
    BC2 --> AC1
    BC2 --> AC2
    BC3 --> AC3
    BC4 --> AC4
    BC4 --> AC5

    AC1 --> TC1
    AC1 --> TC3
    AC2 --> TC1
    AC2 --> TC3
    AC3 --> TC2
    AC4 --> TC1
    AC5 --> TC4

    style BC1 fill:#e1f5ff,stroke:#01579b,stroke-width:2px
    style BC2 fill:#e1f5ff,stroke:#01579b,stroke-width:2px
    style BC3 fill:#e1f5ff,stroke:#01579b,stroke-width:2px
    style BC4 fill:#e1f5ff,stroke:#01579b,stroke-width:2px

    style AC1 fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    style AC2 fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    style AC3 fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    style AC4 fill:#f3e5f5,stroke:#4a148c,stroke-width:2px
    style AC5 fill:#f3e5f5,stroke:#4a148c,stroke-width:2px

    style TC1 fill:#fff3e0,stroke:#e65100,stroke-width:2px
    style TC2 fill:#fff3e0,stroke:#e65100,stroke-width:2px
    style TC3 fill:#fff3e0,stroke:#e65100,stroke-width:2px
    style TC4 fill:#fff3e0,stroke:#e65100,stroke-width:2px
    style TC5 fill:#fff3e0,stroke:#e65100,stroke-width:2px
```

---

## Migration and Transformation Architecture

### Roadmap: Current State → Future State

```mermaid
gantt
    title Upgrade Roadmap
    dateFormat YYYY-MM-DD
    section Stage 0
    Baseline Testing             :done, s0, 2025-01-18, 1d
    section Stage 1
    Increase Test Coverage       :active, s1, 2025-01-19, 3d
    section Stage 2
    Go 1.25 Upgrade             :s2, after s1, 1d
    section Stage 3
    Update server-go-ssp        :s3, after s2, 2d
    section Stage 4
    Update Indirect Deps        :s4, after s3, 1d
    section Stage 5
    Validate & Release          :s5, after s4, 1d
```

### Transformation Architecture: Go 1.24 → Go 1.25

```mermaid
graph LR
    subgraph "Current State"
        C1[Go 1.24.0]
        C2[server-go-ssp<br/>v0.0.0-20241212182118]
        C3[Test Coverage 46.6%]
    end

    subgraph "Transition State"
        T1[Increase Coverage to 85%]
        T2[Update Go to 1.25]
        T3[Verify Compatibility]
    end

    subgraph "Future State"
        F1[Go 1.25.4]
        F2[server-go-ssp<br/>v0.0.0-20251118072316]
        F3[Test Coverage 85%+]
    end

    C1 --> T2
    C2 --> T3
    C3 --> T1

    T1 --> F3
    T2 --> F1
    T3 --> F2

    style C1 fill:#ffebee,stroke:#c62828,stroke-width:2px
    style C2 fill:#ffebee,stroke:#c62828,stroke-width:2px
    style C3 fill:#ffebee,stroke:#c62828,stroke-width:2px

    style T1 fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    style T2 fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    style T3 fill:#fff3e0,stroke:#ef6c00,stroke-width:2px

    style F1 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style F2 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
    style F3 fill:#e8f5e9,stroke:#2e7d32,stroke-width:2px
```

---

## Governance and Compliance

### Architecture Principles Compliance

| Principle | Requirement | Implementation | Status |
|-----------|-------------|----------------|--------|
| **Security by Design** | All sensitive data cleared | ClearBytes(), LockMemory() | ✅ Compliant |
| **Fail-Safe Defaults** | Graceful degradation | Memory locking optional | ✅ Compliant |
| **Least Privilege** | No root required | Best-effort memory locking | ✅ Compliant |
| **Defense in Depth** | Multiple security layers | Validation + Clearing + Locking | ✅ Compliant |
| **Open Standards** | Use SQRL protocol | Implements ssp.Hoard | ✅ Compliant |
| **Separation of Concerns** | Modular design | Hoard / Secure / Validate modules | ✅ Compliant |

### Regulatory Compliance Matrix

| Regulation | Control | Implementation | Evidence |
|------------|---------|----------------|----------|
| **GDPR** (PII Protection) | Art. 32 Security | ClearBytes() | secure.go L8-27 |
| **PCI-DSS** (if handling payment) | Req 3.4 Cryptographic Storage | TLS + Memory Clearing | Redis TLS Config |
| **HIPAA** (if healthcare) | § 164.312(a)(2)(iv) Encryption | TLS Transport | Documentation |
| **SOC 2** (Trust Services) | CC6.1 Logical Access | Input Validation | ValidateNut() |

---

## References

### TOGAF Documentation

- **TOGAF 9.2 Standard**: https://www.opengroup.org/togaf
- **Architecture Development Method (ADM)**: Applied in planning
- **Architecture Building Blocks (ABBs)**: Hoard interface, Security module
- **Solution Building Blocks (SBBs)**: Redis client, Go stdlib

### Project Documentation

- [Requirements](REVERSE_ENGINEERED_REQUIREMENTS.md)
- [Security Plan](SECURITY_PLAN.md)
- [Dependency Analysis](DEPENDENCY_ANALYSIS.md)
- [API Reference](API_REFERENCE.md)

---

## Change Log

| Date | Version | Changes |
|------|---------|---------|
| 2025-01-18 | 1.0.0 | Initial TOGAF architecture documentation with mermaid diagrams |

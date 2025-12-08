# Qubic MicroAuth - Architecture

## Overview

Qubic MicroAuth is a high-performance authentication microservice built on the Qubic blockchain using hexagonal (ports and adapters) architecture. It provides dual protocol support (HTTP REST + gRPC) for maximum flexibility and performance.

## Architecture Principles

### Hexagonal Architecture (Ports & Adapters)

The system follows hexagonal architecture to achieve:
- **Independence**: Business logic independent of frameworks and infrastructure
- **Testability**: Easy to test with mock implementations
- **Flexibility**: Easy to swap implementations (e.g., Redis → DynamoDB)
- **Maintainability**: Clear separation of concerns

### Layers

#### 1. Domain Layer (Core)
- **Location**: `internal/domain/`
- **Purpose**: Pure business logic, no external dependencies
- **Components**:
  - `auth/models.go`: Domain models (WalletAuth, AuthStatus)
  - `auth/service.go`: Business logic implementation

#### 2. Ports (Interfaces)
- **Location**: `internal/ports/`
- **Purpose**: Define contracts for external interactions
- **Components**:
  - `qubic_port.go`: Blockchain interaction interface
  - `wallet_port.go`: Signature verification interface
  - `trust_store_port.go`: Caching interface

#### 3. Primary Adapters (Inbound)
- **Location**: `internal/adapters/primary/`
- **Purpose**: Handle incoming requests
- **Components**:
  - `http/`: Fiber-based REST API
  - `grpc/`: gRPC server implementation

#### 4. Secondary Adapters (Outbound)
- **Location**: `internal/adapters/secondary/`
- **Purpose**: Implement external integrations
- **Components**:
  - `qubic/`: Qubic blockchain client
  - `wallet/`: Signature verification
  - `truststore/`: Redis and in-memory caching

## Data Flow

```
HTTP Request → Handler → Domain Service → Port → Adapter → External System
     ↓                                                            ↓
  Response ← Handler ← Domain Service ← Port ← Adapter ← External System
```

## Caching Strategy

### Three-Layer Cache (L1/L2/L3)

1. **L1: In-Memory Cache**
   - Latency: ~1ms
   - Implementation: `sync.Map` with TTL
   - Use case: Hot data, frequently accessed wallets

2. **L2: Redis Cache**
   - Latency: ~5-10ms
   - Implementation: Redis with pipelining
   - Use case: Distributed cache, shared across instances

3. **L3: Blockchain**
   - Latency: ~100-500ms
   - Implementation: Qubic RPC client
   - Use case: Source of truth, cache miss fallback

### Cache Flow

```
Request → L1 (Memory) → L2 (Redis) → L3 (Blockchain)
            ↓ Hit          ↓ Hit        ↓ Always succeeds
         Response       Response      Response + Cache Update
```

## Performance Optimizations

### 1. Fast JSON Serialization
- **Library**: Sonic (bytedance/sonic)
- **Benefit**: JIT-compiled JSON, 2-3x faster than standard library
- **Usage**: Configured in Fiber app

### 2. Connection Pooling
- **Redis**: 100 connections pre-warmed
- **gRPC**: HTTP/2 multiplexing
- **Qubic**: Keep-alive connections

### 3. Batch Operations
- **Batch Get**: Retrieve multiple wallet statuses in one call
- **Pipeline**: Redis pipelining for batch operations
- **Parallel Queries**: Concurrent blockchain queries

### 4. Zero-Allocation Logging
- **Library**: zerolog
- **Benefit**: No heap allocations during logging
- **Configuration**: JSON format for production, pretty for dev

## Smart Contract Design

### Upgradable Pattern

The smart contract supports upgrades through a pointer pattern:

```
Contract v1 → set_next_contract(v2_address)
              ↓
           Contract v2 (new features)
```

### State Management

```cpp
map<wallet_address, WalletAuthData> {
  status: ACTIVE | BLOCKED | REVIEW
  trust_score: 0-100
  updated_at: timestamp
}
```

### Admin Controls

- Only admin can update statuses
- Admin can transfer ownership
- Admin can set next contract address

## Monitoring & Observability

### Metrics (Prometheus)

- **HTTP Metrics**: Request count, duration, status codes
- **gRPC Metrics**: RPC count, duration, errors
- **Cache Metrics**: Hit/miss rates for L1/L2
- **Blockchain Metrics**: Request count, duration, errors

### Logging

- **Structured Logging**: JSON format with context
- **Log Levels**: Debug, Info, Warn, Error
- **Correlation IDs**: Track requests across services

### Health Checks

- `/health`: Overall system health
- Checks: Qubic node connectivity, Redis connectivity

## Security Considerations

### 1. Wallet Verification
- Signature verification before status updates
- Qubic address validation (60 uppercase A-Z)

### 2. Admin Authorization
- Admin signature required for status changes
- Admin address stored in smart contract

### 3. Rate Limiting
- TODO: Implement rate limiting per wallet
- TODO: DDoS protection

### 4. Input Validation
- All inputs validated before processing
- Trust score range: 0-100
- Address format validation

## Scalability

### Horizontal Scaling

- **Stateless Design**: No session state, easy to scale
- **Load Balancing**: Multiple instances behind load balancer
- **Shared Cache**: Redis for distributed caching

### Vertical Scaling

- **Connection Pools**: Maximize resource utilization
- **Goroutine Pools**: Limit concurrent operations
- **Memory Management**: Efficient data structures

## Future Enhancements

1. **L1 Cache**: Implement in-memory cache layer
2. **Circuit Breaker**: Protect against blockchain failures
3. **Rate Limiting**: Per-wallet and global limits
4. **Metrics Dashboard**: Grafana dashboards
5. **Tracing**: Distributed tracing with Jaeger
6. **Multi-Region**: Deploy across multiple regions
7. **Smart Contract Events**: Listen to blockchain events
8. **Webhook Support**: Notify on status changes

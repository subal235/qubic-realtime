# Realtime

> **Real-Time Infrastructure Layer for Qubic Blockchain**

Realtime provides the foundational primitives for building true real-time apps on Qubic: **identity in milliseconds** (TurboAuth) and **value movement in milliseconds** (TurboRoute).

---

## 🎯 Vision

Modern applications demand instant responses. Realtime brings millisecond-level authentication and payment routing to the Qubic ecosystem, enabling developers to build truly responsive decentralized applications.

---

## 🏗️ Architecture

This is a **monorepo** containing multiple high-performance microservices:

```
Realtime/
├── services/
│   ├── turboauth/          # 🔐 Real-Time Identity Engine
│   └── turboroute/         # 💸 Real-Time Payment Routing
│
├── shared/
│   ├── proto/              # Shared gRPC definitions
│   ├── docker/             # Shared Docker configurations
│   ├── scripts/            # Build & deployment scripts
│   └── docs/               # Shared documentation
│
└── infrastructure/
    ├── docker-compose.yml  # Production orchestration
    └── prometheus.yml      # Monitoring configuration
```

---

## 🚀 Services

### 🔐 TurboAuth – Real-Time Identity Engine

**One-liner**: Millisecond authentication and wallet trust engine for Qubic-powered apps, APIs, and devices.

#### What TurboAuth Does

- ✅ Verifies wallet + signature in **single-digit to low-double-digit ms**
- ✅ Returns clear auth decision: `status` (ACTIVE/BLOCKED/UNKNOWN) + `trustScore` (0–100)
- ✅ **Triple-layer caching** for blazing speed:
  - **L1**: In-memory cache (~1ms)
  - **L2**: Redis shared cache (~5-10ms)
  - **L3**: Qubic smart contract (~100-500ms)
- ✅ Simple API (HTTP REST + gRPC):
  - `POST /api/v1/verify` (REST)
  - `AuthService.VerifyWallet` (gRPC)

#### Who Uses TurboAuth?

**Qubic dApp Developers** who want "login with wallet" with:
- Instant response
- Shared trust/reputation layer
- No custom implementation needed

**Backend/API Developers** who don't want to:
- Implement signature verification themselves
- Build trust logic from scratch
- Manage wallet reputation

**IoT / Access Control / Device Developers** who need:
- "Is this wallet allowed to open this door/use this device?"
- Response in **< 30ms** over HTTP
- Offline-capable with cached decisions

#### Key Features

**Architecture**:
- Hexagonal (Ports & Adapters) in Go
- Domain: `auth` (models, service, ports)
- Adapters: HTTP (Fiber), gRPC, Qubic client, Redis, Memory cache

**Smart Contract** (C++):
- `TurboAuth_v1` on Qubic blockchain
- Wallet → {status, trustScore} mapping
- Admin-only updates
- Upgradable via "next contract" pattern

**Infrastructure**:
- Docker containerized
- Prometheus metrics + Grafana dashboards
- Centralized config and logging

**Status**: ✅ **Production Ready** (needs Qubic SDK integration)

[📖 TurboAuth Documentation](./services/turboauth/README.md)

---

### 💸 TurboRoute – Real-Time Payment Routing Engine

**One-liner**: Low-latency payment routing and execution engine that finds and triggers the best way to move value in real time.

#### What TurboRoute Will Do

**Receives a payment intent**:
```json
{
  "from": "WALLET_A...",
  "to": "WALLET_B...",
  "amount": 1000,
  "preferences": {
    "priority": "speed",      // or "cost" or "privacy"
    "max_fee": 10,
    "timeout_ms": 5000
  }
}
```

**Returns the best route**:
```json
{
  "route_id": "direct_transfer",
  "estimated_fee": 2,
  "estimated_time_ms": 150,
  "hops": ["WALLET_A", "WALLET_B"],
  "execution_tx": "0x..."
}
```

**Key Capabilities**:
- 🚀 **Route Discovery**: Finds all possible payment paths
- ⚡ **Smart Selection**: Chooses optimal route based on preferences
- 🔄 **Auto-Execution**: Triggers the transaction on-chain
- 📊 **Real-Time Metrics**: Tracks route health and performance
- 💰 **Cost Optimization**: Minimizes fees while meeting requirements

#### Who Uses TurboRoute?

**dApp Developers** who want:
- Single, simple payment API
- No manual route management
- Automatic optimization

**Games & Real-Time Apps** that need:
- Micro-payments (pay-per-use)
- Instant tipping
- Low-latency transactions
- No routing complexity

**Future DEX / DeFi Systems** that want to:
- Plug in as "routes" into the engine
- Provide liquidity paths
- Compete on speed/cost

#### Planned Architecture

**Domain**: `payment` / `route`
- Models: `PaymentIntent`, `RouteOption`, `RouteDecision`
- Service: Route discovery, selection, execution logic

**Ports**:
- `QubicPaymentPort` – Execute transactions on-chain
- `RouteCatalogPort` – List/score available routes
- `RouteHealthPort` – Monitor route performance

**Adapters**:
- HTTP: `POST /api/v1/pay`, `POST /api/v1/route`
- gRPC: `RouteService.FindRoute`, `RouteService.ExecutePayment`
- Redis: Cache route metrics and health data

**Smart Contracts** (Future):
- `RouteRegistry` – Where all routes are registered
- `RouteExecutor` – On-chain executor that moves funds

**Status**: 🚧 **Planned for Phase 3**

[📖 TurboRoute Documentation](./services/turboroute/README.md)

---

## 🎯 Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- Make (optional, for convenience)

### Run All Services

```bash
# Development mode (with hot reload)
make dev

# Production mode
make deploy

# Or manually
cd infrastructure
docker-compose up -d
```

### Run Individual Services

```bash
# TurboAuth only
cd services/turboauth/backend
go run ./cmd/api

# Or with Docker
docker-compose up turboauth
```

---

## � Use Cases

### TurboAuth Use Cases

**1. Wallet-Based Login**
```javascript
// Frontend
const response = await fetch('/api/v1/verify', {
  method: 'POST',
  body: JSON.stringify({
    wallet_address: userWallet,
    signature: signedMessage,
    message: challengeMessage
  })
});
// Response in < 10ms
```

**2. API Access Control**
```go
// Backend middleware
func AuthMiddleware(c *fiber.Ctx) error {
    wallet := c.Get("X-Wallet-Address")
    status := turboauth.GetStatus(wallet)
    if status.Status != "ACTIVE" || status.TrustScore < 50 {
        return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
    }
    return c.Next()
}
```

**3. IoT Device Access**
```python
# Smart lock controller
def can_unlock(wallet_address):
    response = requests.get(f'http://turboauth:8080/api/v1/status/{wallet_address}')
    auth = response.json()
    return auth['status'] == 'ACTIVE' and auth['trust_score'] >= 80
```

### TurboRoute Use Cases (Planned)

**1. In-Game Payments**
```javascript
// Pay for in-game item
const payment = await turboroute.pay({
  to: merchantWallet,
  amount: 100,
  preferences: { priority: 'speed', max_fee: 5 }
});
// Executed in < 200ms
```

**2. Micro-Tipping**
```javascript
// Tip content creator
await turboroute.pay({
  to: creatorWallet,
  amount: 1,
  preferences: { priority: 'cost' }  // Minimize fees
});
```

**3. Automated Payments**
```go
// Subscription payment
route, _ := turboroute.FindRoute(PaymentIntent{
    From: userWallet,
    To: serviceWallet,
    Amount: monthlyFee,
    Preferences: RoutePreferences{Priority: "reliability"}
})
turboroute.Execute(route)
```

---

## 📦 Monorepo Benefits

✅ **Unified Vision**: Services designed to work together  
✅ **Shared Infrastructure**: One Docker setup, one CI/CD pipeline  
✅ **Code Reuse**: Shared utilities, SDKs, and proto definitions  
✅ **Atomic Changes**: Update multiple services in one commit  
✅ **Easier Development**: Single clone, single setup  
✅ **Consistent Tooling**: Same build tools, linters, formatters  

---

## 🛠️ Development

### Common Commands

```bash
# Build all services
make build

# Run tests
make test

# Generate proto files
make proto

# Clean build artifacts
make clean

# View logs
make logs

# Health check
make health
```

### Project Structure

```
Realtime/
├── services/turboauth/
│   ├── backend/           # Go service
│   ├── contracts/         # C++ smart contracts
│   └── api/              # Proto definitions
│
├── services/turboroute/   # (Coming soon)
│
├── shared/
│   ├── proto/            # Shared gRPC schemas
│   ├── scripts/          # Build & deployment scripts
│   └── docs/             # Architecture documentation
│
└── infrastructure/
    ├── docker-compose.yml
    └── monitoring/       # Prometheus, Grafana configs
```

---

## 📊 Performance Targets

### TurboAuth

| Metric | Target | Notes |
|--------|--------|-------|
| Verification (cached) | < 10ms | L2 Redis hit |
| Verification (cold) | < 100ms | L3 Blockchain query |
| Throughput | 50k+ req/s | Single instance |
| Availability | 99.9% | With Redis HA |

### TurboRoute (Planned)

| Metric | Target | Notes |
|--------|--------|-------|
| Route Discovery | < 50ms | All available routes |
| Route Selection | < 20ms | Optimal route chosen |
| Payment Execution | < 200ms | End-to-end |
| Throughput | 10k+ payments/s | Batch processing |

---

## 📚 Documentation

- **Architecture**: [shared/docs/architecture.md](./shared/docs/architecture.md)
- **TurboAuth**: [services/turboauth/README.md](./services/turboauth/README.md)
- **TurboRoute**: [services/turboroute/README.md](./services/turboroute/README.md)
- **Contributing**: [CONTRIBUTING.md](./CONTRIBUTING.md)
- **Monorepo Guide**: [MONOREPO_GUIDE.md](./MONOREPO_GUIDE.md)

---

## 🗺️ Roadmap

### Phase 1: Foundation ✅
- [x] TurboAuth service architecture
- [x] Hexagonal design pattern
- [x] Docker infrastructure
- [x] Monitoring setup (Prometheus + Grafana)
- [x] Triple-layer caching (L1/L2/L3)

### Phase 2: Integration 🚧
- [ ] Real Qubic SDK integration
- [ ] Wallet signature verification (crypto library)
- [ ] Smart contract deployment to Qubic
- [ ] Comprehensive test suite
- [ ] Production deployment guide

### Phase 3: TurboRoute 📋
- [ ] TurboRoute service architecture
- [ ] Route discovery algorithms
- [ ] Payment execution engine
- [ ] Route health monitoring
- [ ] Cross-service communication (TurboAuth + TurboRoute)

### Phase 4: Production 🎯
- [ ] Kubernetes deployment
- [ ] Multi-region support
- [ ] Advanced monitoring & alerting
- [ ] Load testing & optimization
- [ ] Public API documentation

---

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines.

---

## 📄 License

MIT

---

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/your-org/realtime/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/realtime/discussions)
- **Documentation**: [Wiki](https://github.com/your-org/realtime/wiki)

---

**Built with ❤️ for the Qubic Ecosystem**

*Realtime: Identity and value movement in milliseconds.*

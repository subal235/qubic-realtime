# Realtime Platform

> Real-time authentication and payment routing infrastructure for Qubic blockchain

**Status:** ✅ Production Ready (TurboAuth) | 🚧 In Development (TurboRoute)

---

## 🎯 What is Realtime?

Realtime is a **monorepo** containing microservices for the Qubic blockchain ecosystem:

1. **TurboAuth** - Decentralized wallet authentication and trust scoring
2. **TurboRoute** - Intelligent payment routing and optimization (planned)

All services follow **hexagonal (ports & adapters) architecture** for maximum flexibility and testability.

---

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make

### Run Locally

```bash
# 1. Clone and setup
git clone <your-repo>
cd Realtime
make init

# 2. Start infrastructure (Redis, Prometheus, Grafana)
make dev

# 3. Run services (in separate terminals)
make dev-turboauth    # Terminal 1
make dev-turboroute   # Terminal 2
```

### Deploy with Docker

```bash
# Build and deploy everything
make deploy

# Check health
make health

# View logs
make logs
```

---

## 📦 Services

### TurboAuth (✅ Production Ready)

**Purpose:** Decentralized wallet authentication and trust scoring

**Features:**
- ✅ HTTP REST API (port 8080)
- ✅ gRPC API (port 9090)
- ✅ Multi-layer caching (L1/L2/L3)
- ✅ Smart contract integration
- ✅ Trust score calculation
- ✅ Rate limiting & session management

**Documentation:** [services/turboauth/README.md](./services/turboauth/README.md)

---

### TurboRoute (🚧 In Development)

**Purpose:** Intelligent payment routing and optimization

**Features:**
- ✅ HTTP REST API (port 8081)
- ✅ Route discovery algorithms
- ✅ Mock payment execution
- 🚧 Smart contracts (planned Phase 3.4)
- 🚧 Multi-hop routing
- 🚧 Real-time optimization

**Documentation:** [services/turboroute/README.md](./services/turboroute/README.md)

---

## 🏗️ Architecture

### Monorepo Structure

```
Realtime/
├── services/              # Microservices
│   ├── turboauth/        # Authentication service
│   │   ├── backend/      # Go service (hexagonal)
│   │   └── contracts/    # Smart contracts
│   └── turboroute/       # Routing service
│       └── backend/      # Go service (hexagonal)
├── infrastructure/       # Docker Compose, monitoring
├── documents/            # Documentation
├── shared/               # Shared utilities
├── Makefile              # Build automation
└── .env                  # Configuration
```

### Hexagonal Architecture

Both services follow **hexagonal (ports & adapters) architecture**:

```
┌─────────────────────────────────────┐
│   PRIMARY ADAPTERS (Inbound)        │
│     HTTP, gRPC, CLI                 │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│      DOMAIN (Hexagon)               │
│   Business Logic & Rules            │
│   • Models                          │
│   • Services                        │
│   • Port Interfaces                 │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│  SECONDARY ADAPTERS (Outbound)      │
│  Qubic, Redis, Database, APIs       │
└─────────────────────────────────────┘
```

**Benefits:**
- ✅ Business logic isolated from infrastructure
- ✅ Easy to test (mock adapters)
- ✅ Easy to swap implementations
- ✅ Framework independent

**Learn more:** [documents/HEXAGONAL_ARCHITECTURE_GUIDE.md](./documents/HEXAGONAL_ARCHITECTURE_GUIDE.md)

---

## 🛠️ Development

### Available Commands

```bash
# Development
make dev              # Start infrastructure only
make dev-turboauth    # Run TurboAuth locally
make dev-turboroute   # Run TurboRoute locally

# Building
make build            # Build all services
make build-turboauth  # Build TurboAuth only
make build-turboroute # Build TurboRoute only

# Docker
make deploy           # Deploy all services
make stop             # Stop all services
make logs             # View logs

# Testing
make test             # Run all tests
make health           # Check service health

# Utilities
make check-ports      # Check which ports are in use
make ps               # Show running processes
make kill             # Kill all local processes
make clean            # Clean build artifacts
```

### Scripts

```bash
./start.sh            # Start all services in background
./stop.sh             # Stop all services
./check-ports.sh      # Check port usage
```

---

## 📊 Service Ports

| Service | HTTP | gRPC | Metrics |
|---------|------|------|---------|
| TurboAuth | 8080 | 9090 | 2112 |
| TurboRoute | 8081 | 9091 | 2113 |
| Redis | 6379 | - | - |
| Prometheus | 9093 | - | - |
| Grafana | 3000 | - | - |

---

## 📚 Documentation

### Getting Started
- [Development Guide](./documents/DEVELOPMENT_GUIDE.md)
- [Port Management](./documents/PORT_MANAGEMENT.md)
- [Stopping Services](./documents/STOPPING_SERVICES.md)

### Architecture
- [Architecture Analysis](./documents/ARCHITECTURE_ANALYSIS.md)
- [Hexagonal Architecture Guide](./documents/HEXAGONAL_ARCHITECTURE_GUIDE.md)
- [Dockerfile Structure](./documents/DOCKERFILE_STRUCTURE.md)

### Reference
- [Documentation Index](./documents/README.md)
- [Implementation Checklist](./documents/CHECKLIST.md)
- [Fix Summary](./documents/FIX_SUMMARY.md)
- [All Fixed Summary](./ALL_FIXED.md)

---

## 🔧 Configuration

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# Qubic Network
QUBIC_NODE_URL=http://qubic-node:21841
QUBIC_CONTRACT_ADDRESS=

# TurboAuth
TURBOAUTH_HTTP_PORT=8080
TURBOAUTH_GRPC_PORT=9090
TURBOAUTH_LOG_LEVEL=info

# TurboRoute
TURBOROUTE_HTTP_PORT=8081
TURBOROUTE_GRPC_PORT=9091
TURBOROUTE_LOG_LEVEL=info

# Infrastructure
REDIS_PASSWORD=
GRAFANA_PASSWORD=admin
```

---

## 🧪 Testing

### Unit Tests

```bash
# Test all services
make test

# Test specific service
make test-turboauth
make test-turboroute
```

### Health Checks

```bash
# Check all services
make health

# Or manually
curl http://localhost:8080/health  # TurboAuth
curl http://localhost:8081/health  # TurboRoute
```

---

## 📈 Monitoring

### Prometheus

Access metrics at: `http://localhost:9093`

**Metrics exposed:**
- HTTP request counts
- Request duration histograms
- Error rates
- Custom business metrics

### Grafana

Access dashboards at: `http://localhost:3000`

**Default credentials:** admin/admin

---

## 📝 Project Status

### Completed ✅

- ✅ Hexagonal architecture implementation
- ✅ TurboAuth service (production ready)
- ✅ TurboAuth smart contract
- ✅ TurboRoute core service
- ✅ Docker deployment
- ✅ Monitoring setup
- ✅ Comprehensive documentation

### In Progress 🚧

- 🚧 TurboAuth extended features (sessions, webhooks)
- 🚧 TurboRoute advanced routing
- 🚧 Unit tests
- 🚧 Integration tests

### Planned 📋

- 📋 TurboRoute smart contracts
- 📋 API documentation (OpenAPI/Swagger)
- 📋 Admin dashboard
- 📋 SDKs (Go, JS, Python)

**See:** [documents/CHECKLIST.md](./documents/CHECKLIST.md)

---

## 🐛 Troubleshooting

### Port Conflicts

```bash
# Check what's using ports
make check-ports

# Kill processes
make kill
```

### Services Won't Start

```bash
# Check logs
make logs

# Restart services
make stop
make deploy
```

### Build Issues

```bash
# Clean and rebuild
make clean
make build
```

**See:** [documents/DEVELOPMENT_GUIDE.md](./documents/DEVELOPMENT_GUIDE.md)

---

## 🎯 Quick Reference

```bash
# Start everything
make deploy

# Development mode
make dev              # Infrastructure
make dev-turboauth    # Service 1
make dev-turboroute   # Service 2

# Check status
make health
make check-ports

# Stop everything
make stop
make kill
```

---

## 📄 License

MIT

---

**Built with ❤️ for the Qubic ecosystem**

*Realtime Platform - Authentication and routing in milliseconds.*

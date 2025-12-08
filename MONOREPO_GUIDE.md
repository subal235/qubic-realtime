# Realtime Monorepo Guide

## 🎯 Quick Reference

### Project Structure

```
Realtime/                          # Root monorepo
├── services/                      # All microservices
│   ├── turboauth/                # Authentication service ✅
│   └── turboroute/                # Routing service 🚧
│
├── shared/                        # Shared code & utilities
│   ├── proto/                    # Shared gRPC definitions
│   ├── docker/                   # Shared Docker configs
│   ├── scripts/                  # Build & deployment scripts
│   └── docs/                     # Architecture docs
│
└── infrastructure/               # Deployment configs
    ├── docker-compose.yml        # Production
    ├── docker-compose.dev.yml    # Development
    └── prometheus.yml            # Monitoring
```

### Common Commands

```bash
# Development
make dev                 # Start all services (dev mode)
make dev-turboauth       # Start only TurboAuth

# Building
make build              # Build all services
make build-turboauth    # Build TurboAuth only

# Testing
make test               # Run all tests
make test-coverage      # With coverage report

# Docker
make deploy             # Deploy all services
make stop               # Stop all services
make logs               # View all logs

# Utilities
make health             # Check service health
make metrics            # View metrics
make clean              # Clean build artifacts
```

### Service Ports

| Service | HTTP | gRPC | Metrics |
|---------|------|------|---------|
| TurboAuth | 8080 | 9090 | 2112 |
| TurboRoute | 8081 | 9091 | 2113 |

| Infrastructure | Port |
|----------------|------|
| Redis | 6379 |
| Prometheus | 9091 |
| Grafana | 3000 |

## 🏗️ Adding a New Service

1. Create service directory:
   ```bash
   mkdir -p services/myservice/{backend,api,contracts}
   ```

2. Follow TurboAuth structure:
   - Hexagonal architecture
   - HTTP REST + gRPC
   - Docker support
   - README.md

3. Update root files:
   - `Makefile` - Add build/test commands
   - `infrastructure/docker-compose.yml` - Add service
   - `README.md` - List new service

4. Add to shared proto if needed:
   ```bash
   cp api/proto/*.proto ../../shared/proto/
   ```

## 📝 Development Workflow

### Working on TurboAuth

```bash
# Navigate to service
cd services/turboauth/backend

# Install dependencies
go mod download

# Run locally
go run ./cmd/api

# Run tests
go test ./...

# Build
go build -o turboauth ./cmd/api
```

### Working on Shared Code

```bash
# Navigate to shared
cd shared/proto

# Update proto files
vim myservice.proto

# Regenerate (from root)
make proto
```

### Working on Infrastructure

```bash
# Navigate to infrastructure
cd infrastructure

# Test locally
docker-compose up

# Deploy
docker-compose up -d
```

## 🔄 Monorepo Benefits

✅ **Single Source of Truth**: One repo, one clone  
✅ **Shared Infrastructure**: Docker, CI/CD, configs  
✅ **Code Reuse**: Shared utilities and proto files  
✅ **Atomic Changes**: Update multiple services together  
✅ **Unified Versioning**: Tag releases across all services  
✅ **Easier Onboarding**: New devs clone once  

## 📦 Dependency Management

### Service Dependencies

Each service has its own `go.mod`:
```
services/turboauth/backend/go.mod
services/turboroute/backend/go.mod
```

### Shared Dependencies

Shared code can have its own modules:
```
shared/utils/go.mod
shared/sdk/go.mod
```

## 🚀 Deployment Strategies

### Development

```bash
make dev
# Hot reload, debug logging, local volumes
```

### Staging

```bash
cd infrastructure
docker-compose -f docker-compose.yml up -d
# Production-like, with monitoring
```

### Production

```bash
# Kubernetes (future)
kubectl apply -f infrastructure/k8s/

# Or Docker Swarm
docker stack deploy -c docker-compose.yml realtime
```

## 🧪 Testing Strategy

### Unit Tests

```bash
# Per service
cd services/turboauth/backend
go test ./internal/domain/...

# All services
make test
```

### Integration Tests

```bash
# With real dependencies
make test-integration
```

### E2E Tests

```bash
# Across services
cd tests/e2e
go test -v
```

## 📊 Monitoring

### Prometheus Metrics

```bash
# View metrics
curl http://localhost:2112/metrics

# Prometheus UI
open http://localhost:9091
```

### Grafana Dashboards

```bash
# Open Grafana
open http://localhost:3000

# Login: admin/admin
```

### Logs

```bash
# All services
make logs

# Specific service
make logs-turboauth

# Follow logs
docker-compose logs -f turboauth
```

## 🔐 Security

### Environment Variables

- Never commit `.env` files
- Use `.env.example` as template
- Store secrets in secret manager (production)

### API Keys

```bash
# Development
export QUBIC_API_KEY=dev_key

# Production
# Use Kubernetes secrets or similar
```

## 📚 Documentation

- **Root README**: Overview, quick start
- **Service READMEs**: Service-specific docs
- **CONTRIBUTING.md**: Development guidelines
- **shared/docs/**: Architecture, API reference

## 🎓 Learning Path

1. **Start Here**: Root README
2. **Understand Structure**: This guide
3. **Pick a Service**: Read service README
4. **Study Code**: Start with `cmd/api/main.go`
5. **Make Changes**: Follow CONTRIBUTING.md
6. **Deploy**: Use infrastructure/

## ❓ FAQ

**Q: Why a monorepo?**  
A: Shared infrastructure, easier development, unified vision.

**Q: Can services be deployed independently?**  
A: Yes! Each service has its own Dockerfile.

**Q: How do services communicate?**  
A: gRPC for internal, HTTP REST for external.

**Q: Where do I add shared utilities?**  
A: `shared/` directory with its own Go module.

**Q: How do I add a new service?**  
A: Follow the "Adding a New Service" section above.

---

**Happy coding! 🚀**

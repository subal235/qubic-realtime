# Contributing to Realtime

Thank you for your interest in contributing to the Realtime project! This document provides guidelines for contributing to this monorepo.

---

## 📋 Table of Contents

- [Getting Started](#getting-started)
- [Monorepo Structure](#monorepo-structure)
- [Development Workflow](#development-workflow)
- [Code Standards](#code-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)

---

## Getting Started

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make
- protoc (Protocol Buffers compiler)
- Git

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/your-org/realtime.git
cd Realtime

# Initialize dependencies
make init

# Start development environment
make dev
```

---

## Monorepo Structure

This is a **monorepo** containing multiple services. Understanding the structure is crucial:

```
Realtime/
├── services/           # Individual microservices
│   ├── turboauth/     # Authentication service
│   └── turboroute/     # Routing service (planned)
│
├── shared/            # Shared code and utilities
│   ├── proto/        # Shared gRPC definitions
│   ├── docker/       # Shared Docker configs
│   ├── scripts/      # Build & deployment scripts
│   └── docs/         # Shared documentation
│
└── infrastructure/   # Deployment & orchestration
    ├── docker-compose.yml
    └── k8s/         # Kubernetes configs (future)
```

### Where to Make Changes

| Type of Change | Location |
|----------------|----------|
| Service-specific feature | `services/<service-name>/` |
| Shared utilities | `shared/` |
| Infrastructure | `infrastructure/` |
| Documentation | `shared/docs/` or service README |
| Build scripts | `shared/scripts/` |

---

## Development Workflow

### 1. Create a Branch

```bash
# Feature branch
git checkout -b feature/turboauth-rate-limiting

# Bug fix
git checkout -b fix/turboauth-cache-issue

# Documentation
git checkout -b docs/update-api-reference
```

### 2. Make Changes

Follow the [Code Standards](#code-standards) below.

### 3. Test Locally

```bash
# Run tests
make test

# Run specific service tests
cd services/turboauth/backend
go test ./...

# Integration tests
make test-integration
```

### 4. Commit

Use conventional commits:

```bash
git commit -m "feat(turboauth): add rate limiting middleware"
git commit -m "fix(turboroute): resolve cache invalidation bug"
git commit -m "docs: update API reference"
```

**Commit Format**:
```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types**:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance

**Scopes**:
- `turboauth`
- `turboroute`
- `shared`
- `infra`

### 5. Push & Create PR

```bash
git push origin feature/turboauth-rate-limiting
```

Then create a Pull Request on GitHub.

---

## Code Standards

### Go Code

#### Style Guide

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `golangci-lint` for linting

```bash
# Format code
make fmt

# Run linter
make lint
```

#### Hexagonal Architecture

All services follow hexagonal (ports & adapters) architecture:

```
internal/
├── domain/          # Business logic (NO external dependencies)
│   └── auth/
│       ├── models.go
│       ├── service.go
│       └── ports.go  # Interfaces
│
└── adapters/
    ├── primary/     # Inbound (HTTP, gRPC)
    └── secondary/   # Outbound (DB, Cache, External APIs)
```

**Rules**:
1. Domain layer has NO imports from adapters
2. Ports (interfaces) are defined in domain
3. Adapters implement ports
4. Dependency injection in `main.go`

#### Naming Conventions

```go
// Interfaces: end with Port or Service
type QubicPort interface { ... }
type AuthService struct { ... }

// Implementations: descriptive names
type RedisStore struct { ... }
type FiberHandler struct { ... }

// Private functions: camelCase
func validateAddress(addr string) bool { ... }

// Public functions: PascalCase
func NewService(...) *Service { ... }
```

### Testing

#### Unit Tests

```go
// service_test.go
func TestGetStatus_CacheHit(t *testing.T) {
    // Arrange
    mockCache := &MockTrustStore{}
    svc := NewService(nil, nil, mockCache, time.Minute)
    
    // Act
    result, err := svc.GetStatus(ctx, "WALLET...")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, StatusActive, result.Status)
}
```

#### Integration Tests

```go
// +build integration

func TestGetStatus_Integration(t *testing.T) {
    // Test with real Redis, etc.
}
```

### Documentation

- **Code Comments**: Explain WHY, not WHAT
- **Package Comments**: Document package purpose
- **Public APIs**: Full godoc comments
- **READMEs**: Keep up-to-date

```go
// Good: Explains reasoning
// We use a 5-minute TTL to balance freshness with blockchain load

// Bad: States the obvious
// Set TTL to 5 minutes
```

---

## Testing

### Running Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Integration tests
make test-integration

# Benchmarks
make benchmark
```

### Test Requirements

- **Unit tests**: Required for all business logic
- **Integration tests**: Required for adapters
- **Coverage**: Aim for >80% on domain layer
- **Benchmarks**: For performance-critical paths

---

## Pull Request Process

### Before Submitting

- [ ] Tests pass (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Linter passes (`make lint`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated (if applicable)

### PR Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Checklist
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

### Review Process

1. **Automated Checks**: CI/CD runs tests, linters
2. **Code Review**: At least 1 approval required
3. **Testing**: Reviewer tests locally if needed
4. **Merge**: Squash and merge to main

---

## Monorepo Best Practices

### Adding a New Service

1. Create directory: `services/<service-name>/`
2. Follow TurboAuth structure
3. Update root `Makefile`
4. Update `infrastructure/docker-compose.yml`
5. Add README
6. Update root README

### Shared Code

- Put in `shared/` directory
- Create Go module if needed
- Document in `shared/README.md`

### Dependencies

- Service-specific: `services/<service>/go.mod`
- Shared: `shared/<module>/go.mod`
- Keep dependencies minimal

---

## Questions?

- **Issues**: [GitHub Issues](https://github.com/your-org/realtime/issues)
- **Discussions**: [GitHub Discussions](https://github.com/your-org/realtime/discussions)
- **Chat**: [Discord/Slack]

---

**Thank you for contributing to Realtime! 🚀**

# TurboRoute

> Real-time payment routing and optimization service for Qubic blockchain

Part of the [Realtime](../../README.md) monorepo.

---

## 🎯 Vision

**One-liner**: Low-latency payment routing and execution engine that finds and triggers the best way to move value in real time.

TurboRoute lets developers say **"pay this wallet under these preferences"** and handles the smart, low-latency routing for them.

---

## Status

🚧 **Planned for Phase 3**

This service will be developed after TurboAuth reaches production maturity.

---

## What TurboRoute Will Do

### Input: Payment Intent

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

### Output: Optimal Route + Execution

```json
{
  "route_id": "direct_transfer",
  "estimated_fee": 2,
  "estimated_time_ms": 150,
  "hops": ["WALLET_A", "WALLET_B"],
  "execution_tx": "0x...",
  "status": "executed"
}
```

---

## Key Capabilities

### 🚀 Route Discovery
- Finds all possible payment paths
- Considers direct transfers, multi-hop routes, liquidity pools
- Real-time route availability checking

### ⚡ Smart Selection
- Chooses optimal route based on user preferences:
  - **Speed**: Fastest execution time
  - **Cost**: Lowest fees
  - **Privacy**: Most anonymous path
  - **Reliability**: Highest success rate

### 🔄 Auto-Execution
- Triggers the transaction on-chain automatically
- Handles retries and fallbacks
- Provides real-time status updates

### 📊 Real-Time Metrics
- Tracks route health and performance
- Monitors success rates
- Measures actual vs estimated times/fees

### 💰 Cost Optimization
- Minimizes fees while meeting requirements
- Batch processing for micro-payments
- Dynamic fee estimation

---

## Who Uses TurboRoute?

### dApp Developers
Want a single, simple payment API instead of:
- Managing multiple contracts/routes manually
- Implementing routing logic themselves
- Handling edge cases and failures

**Example**:
```javascript
// Instead of this complexity:
const route = await findBestRoute(from, to, amount);
const contract = await getContract(route.contractAddress);
const tx = await contract.transfer(to, amount);
await waitForConfirmation(tx);

// Just do this:
await turboroute.pay({ to, amount, preferences: { priority: 'speed' } });
```

### Games & Real-Time Apps
Need instant payments for:
- **Micro-payments**: Pay-per-use features
- **Tipping**: Instant creator rewards
- **In-game purchases**: Items, upgrades, etc.

Without managing routing complexity.

**Example**:
```javascript
// Tip a streamer
await turboroute.pay({
  to: streamerWallet,
  amount: 5,
  preferences: { priority: 'speed', max_fee: 0.1 }
});
// Executed in < 200ms
```

### Future DEX / DeFi Systems
Want to plug in as "routes" into the engine:
- Provide liquidity paths
- Compete on speed/cost
- Earn routing fees

---

## Planned Architecture

Following the same hexagonal pattern as TurboAuth:

```
┌─────────────────────────────────────────┐
│   Primary Adapters (Inbound)           │
│   ┌─────────┐      ┌─────────┐        │
│   │  HTTP   │      │  gRPC   │        │
│   │ :8081   │      │ :9091   │        │
│   └────┬────┘      └────┬────┘        │
│        │                │              │
│        └────────┬───────┘              │
│                 ▼                      │
│        ┌────────────────┐              │
│        │  Domain Layer  │              │
│        │ (Routing Logic)│              │
│        └────────┬───────┘              │
│                 │                      │
│        ┌────────┴────────┐             │
│        ▼        ▼         ▼            │
│   ┌─────┐  ┌──────┐  ┌──────┐        │
│   │Qubic│  │Route │  │Redis │        │
│   │ Pay │  │Catalog│  │Cache │        │
│   └─────┘  └──────┘  └──────┘        │
│   Secondary Adapters (Outbound)       │
└─────────────────────────────────────────┘
```

### Domain Layer

**Models**:
```go
type PaymentIntent struct {
    From        string
    To          string
    Amount      int64
    Preferences RoutePreferences
}

type RouteOption struct {
    RouteID       string
    Hops          []string
    EstimatedFee  int64
    EstimatedTime time.Duration
    SuccessRate   float64
}

type RouteDecision struct {
    SelectedRoute RouteOption
    Reason        string
    Alternatives  []RouteOption
}
```

**Service**:
```go
type Service struct {
    paymentPort  QubicPaymentPort
    catalogPort  RouteCatalogPort
    cachePort    RouteHealthPort
}

func (s *Service) FindRoute(intent PaymentIntent) (*RouteDecision, error)
func (s *Service) ExecutePayment(route RouteOption) (string, error)
func (s *Service) GetRouteHealth(routeID string) (*RouteHealth, error)
```

### Ports (Interfaces)

**QubicPaymentPort**:
```go
type QubicPaymentPort interface {
    ExecuteTransfer(from, to string, amount int64) (txHash string, err error)
    GetTransactionStatus(txHash string) (TransactionStatus, error)
    EstimateFee(from, to string, amount int64) (int64, error)
}
```

**RouteCatalogPort**:
```go
type RouteCatalogPort interface {
    ListRoutes(from, to string) ([]RouteOption, error)
    ScoreRoute(route RouteOption, preferences RoutePreferences) (float64, error)
    RegisterRoute(route RouteOption) error
}
```

**RouteHealthPort**:
```go
type RouteHealthPort interface {
    GetHealth(routeID string) (*RouteHealth, error)
    UpdateMetrics(routeID string, metrics RouteMetrics) error
    GetCachedRoute(from, to string) (*RouteOption, error)
}
```

### Adapters

**HTTP** (`adapters/primary/http/`):
- `POST /api/v1/pay` - Execute payment with auto-routing
- `POST /api/v1/route` - Find best route (no execution)
- `GET /api/v1/routes/:from/:to` - List all available routes
- `GET /api/v1/health/:routeID` - Get route health metrics

**gRPC** (`adapters/primary/grpc/`):
- `RouteService.FindRoute` - Discover optimal route
- `RouteService.ExecutePayment` - Execute payment
- `RouteService.GetRouteHealth` - Health metrics

**Qubic Client** (`adapters/secondary/qubic/`):
- Execute on-chain transfers
- Query transaction status
- Estimate fees

**Redis Cache** (`adapters/secondary/cache/`):
- Cache route health metrics
- Store recent route decisions
- Track success rates

---

## Smart Contracts (Future)

### RouteRegistry Contract (C++)

Stores all available payment routes:

```cpp
struct Route {
    string routeID;
    string[] hops;
    int64 baseFee;
    bool isActive;
};

map<string, Route> routes;

// Admin can register new routes
void registerRoute(Route route);

// Anyone can query routes
Route[] getRoutes(string from, string to);
```

### RouteExecutor Contract (C++)

Executes multi-hop payments:

```cpp
struct PaymentExecution {
    string routeID;
    string from;
    string to;
    int64 amount;
    string[] hops;
};

// Execute payment along specified route
string executePayment(PaymentExecution payment);

// Get execution status
ExecutionStatus getStatus(string txHash);
```

---

## API Examples

### REST API

**Find and Execute Payment**:
```bash
curl -X POST http://localhost:8081/api/v1/pay \
  -H "Content-Type: application/json" \
  -d '{
    "from": "WALLET_A...",
    "to": "WALLET_B...",
    "amount": 1000,
    "preferences": {
      "priority": "speed",
      "max_fee": 10
    }
  }'
```

**Response**:
```json
{
  "route_id": "direct_transfer",
  "estimated_fee": 2,
  "estimated_time_ms": 150,
  "execution_tx": "0x123...",
  "status": "executed"
}
```

**Find Route (No Execution)**:
```bash
curl -X POST http://localhost:8081/api/v1/route \
  -H "Content-Type: application/json" \
  -d '{
    "from": "WALLET_A...",
    "to": "WALLET_B...",
    "amount": 1000,
    "preferences": { "priority": "cost" }
  }'
```

### gRPC API

```go
client := pb.NewRouteServiceClient(conn)

route, err := client.FindRoute(ctx, &pb.FindRouteRequest{
    From: "WALLET_A...",
    To: "WALLET_B...",
    Amount: 1000,
    Preferences: &pb.RoutePreferences{
        Priority: "speed",
        MaxFee: 10,
    },
})
```

---

## Performance Targets

| Metric | Target | Notes |
|--------|--------|-------|
| Route Discovery | < 50ms | All available routes |
| Route Selection | < 20ms | Optimal route chosen |
| Payment Execution | < 200ms | End-to-end |
| Throughput | 10k+ payments/s | With batch processing |
| Success Rate | > 99% | With retries |

---

## Development Timeline

### Phase 3.1: Foundation
- [ ] Project structure (mirror TurboAuth)
- [ ] Domain models and service
- [ ] Port interfaces
- [ ] Mock adapters for testing

### Phase 3.2: Core Routing
- [ ] Route discovery algorithm
- [ ] Route scoring/selection logic
- [ ] Basic HTTP API
- [ ] Redis caching

### Phase 3.3: Execution
- [ ] Qubic payment integration
- [ ] Transaction monitoring
- [ ] Retry logic
- [ ] gRPC API

### Phase 3.4: Advanced Features
- [ ] Multi-hop routing
- [ ] Route health tracking
- [ ] Batch payments
- [ ] Smart contract deployment

---

## Integration with TurboAuth

TurboRoute will integrate with TurboAuth for:

**Payment Authorization**:
```go
// Before executing payment, verify wallet
authStatus := turboauth.GetStatus(paymentIntent.From)
if authStatus.Status != "ACTIVE" || authStatus.TrustScore < 50 {
    return errors.New("wallet not authorized for payments")
}
```

**Trust-Based Routing**:
```go
// Prefer routes through high-trust wallets
for _, route := range routes {
    for _, hop := range route.Hops {
        trustScore := turboauth.GetStatus(hop).TrustScore
        route.Score += trustScore * 0.1
    }
}
```

---

## Contributing

TurboRoute is planned for Phase 3. If you're interested in contributing:

1. Review the TurboAuth implementation
2. Understand the hexagonal architecture
3. Check the roadmap above
4. Join discussions in [GitHub Discussions](https://github.com/your-org/realtime/discussions)

---

## License

MIT

---

**Part of the [Realtime](../../README.md) ecosystem**

*TurboRoute: Value movement in milliseconds.*

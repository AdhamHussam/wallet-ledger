# Wallet Ledger Service

A high-performance, ACID-compliant financial ledger and wallet service written in Go and backed by PostgreSQL, paired with a modern React + TypeScript dashboard. Designed for high-concurrency environments with strict money conservation guarantees and deadlock-free transactions.

---

## Architecture Overview

```
                      +----------------------------------+
                      |         Web Dashboard            |
                      |    (React 19 + TypeScript +      |
                      |   Tailwind CSS + Lucide Icons)   |
                      +-----------------+----------------+
                                        |
                      +-----------------v----------------+
                      |         HTTP REST API            |
                      |   (Go standard net/http router   |
                      |    + CORS + Recovery Middleware) |
                      +-----------------+----------------+
                                        |
                      +-----------------v----------------+
                      |          Service Layer           |
                      |   - Input & Currency Validation  |
                      |   - Account & Transfer Rules     |
                      +-----------------+----------------+
                                        |
                      +-----------------v----------------+
                      |           Store Layer            |
                      |   - Deadlock-Free TransferTx     |
                      |   - Atomic Multi-Record Commits  |
                      +-----------------+----------------+
                                        |
                      +-----------------v----------------+
                      |        PostgreSQL Database       |
                      |   - Accounts, Entries, Transfers |
                      |   - Non-Negative Balance Checks  |
                      +----------------------------------+
```

### Core Design Principles

1. **Double-Entry Ledger Invariant**: Every transfer creates an immutable record linking two accounts, plus two balancing entries (negative for debit, positive for credit).
2. **Conservation of Money**: Under any number of concurrent transactions across multiple accounts, $\sum \text{Balance}_{\text{final}} = \sum \text{Balance}_{\text{initial}}$.
3. **Deadlock Prevention**: In concurrent bidirectional transfers between accounts ($A \rightarrow B$ and $B \rightarrow A$), deadlocks are mathematically eliminated by enforcing consistent global acquisition order (always locking and updating the lower account ID first).
4. **Database Check Constraints**: Enforces `balance >= 0` at the database level (`balance_non_negative`), protecting against any overdraft race condition.
5. **Zero External Router Dependencies**: Handlers and routing leverage Go's built-in `net/http` router with middleware for structured logging (`log/slog`), CORS support, and panic recovery.
6. **Embedded Database Migrations**: Uses `//go:embed` with `golang-migrate` to automatically apply migrations on service startup.

---

## Database Schema

```
 +--------------------+       +--------------------+       +--------------------+
 |      accounts      |       |      entries       |       |     transfers      |
 +--------------------+       +--------------------+       +--------------------+
 | id (BIGSERIAL, PK) |<--+   | id (BIGSERIAL, PK) |       | id (BIGSERIAL, PK) |
 | owner (VARCHAR)    |   +---| account_id (FK)    |   +---| from_account_id(FK)|
 | balance (BIGINT)   |   |   | amount (BIGINT)    |   +---| to_account_id (FK) |
 | currency (VARCHAR) |   |   | created_at (TZ)    |   |   | amount (BIGINT > 0)|
 | created_at (TZ)    |   |   +--------------------+   |   | created_at (TZ)    |
 +--------------------+   +----------------------------+   +--------------------+
```

---

## API Documentation

### 1. Health Check
```http
GET /health
```
**Response (200 OK):**
```json
{
  "status": "up",
  "service": "wallet-ledger"
}
```

### 2. Create Account
```http
POST /accounts
Content-Type: application/json

{
  "owner": "alice",
  "currency": "USD",
  "initial_balance": 1000
}
```
**Response (201 Created):**
```json
{
  "id": 1,
  "owner": "alice",
  "balance": 1000,
  "currency": "USD",
  "created_at": "2026-10-04T15:16:11.510845+03:00"
}
```

### 3. Get Account by ID
```http
GET /accounts/1
```
**Response (200 OK):**
```json
{
  "id": 1,
  "owner": "alice",
  "balance": 1000,
  "currency": "USD",
  "created_at": "2026-10-04T15:16:11.510845+03:00"
}
```

### 4. List Accounts (Paginated)
```http
GET /accounts?page_id=1&page_size=10
```
**Response (200 OK):**
```json
[
  {
    "id": 1,
    "owner": "alice",
    "balance": 1000,
    "currency": "USD",
    "created_at": "2026-10-04T15:16:11.510845+03:00"
  }
]
```

### 5. Transfer Money
```http
POST /transfers
Content-Type: application/json

{
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 250,
  "currency": "USD"
}
```
**Response (201 Created):**
```json
{
  "transfer": {
    "id": 1,
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 250,
    "created_at": "2026-10-04T15:16:11.527005+03:00"
  },
  "from_account": {
    "id": 1,
    "owner": "alice",
    "balance": 750,
    "currency": "USD",
    "created_at": "2026-10-04T15:16:11.510845+03:00"
  },
  "to_account": {
    "id": 2,
    "owner": "bob",
    "balance": 250,
    "currency": "USD",
    "created_at": "2026-10-04T15:16:11.518103+03:00"
  },
  "from_entry": {
    "id": 1,
    "account_id": 1,
    "amount": -250,
    "created_at": "2026-10-04T15:16:11.527005+03:00"
  },
  "to_entry": {
    "id": 2,
    "account_id": 2,
    "amount": 250,
    "created_at": "2026-10-04T15:16:11.527005+03:00"
  }
}
```

### 6. Get Transfer by ID
```http
GET /transfers/1
```
**Response (200 OK):**
```json
{
  "id": 1,
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 250,
  "created_at": "2026-10-04T15:16:11.527005+03:00"
}
```

---

## Configuration

Settings are loaded from environment variables (see [`.env.example`](.env.example)):

| Variable | Default | Description |
|---|---|---|
| `ENVIRONMENT` | `development` | Deployment environment |
| `DB_SOURCE` | `postgres://postgres:secretpassword@localhost:5432/ledger_db?sslmode=disable` | PostgreSQL connection URI |
| `HTTP_SERVER_ADDRESS` | `:8080` | Port / interface for HTTP server |
| `HTTP_READ_TIMEOUT` | `5s` | Maximum duration reading request |
| `HTTP_WRITE_TIMEOUT` | `10s` | Maximum duration writing response |
| `HTTP_IDLE_TIMEOUT` | `60s` | Keep-alive timeout |
| `HTTP_SHUTDOWN_TIMEOUT`| `10s` | Maximum duration for graceful shutdown |
| `RUN_MIGRATION_ON_BOOT`| `true` | Apply pending DB migrations at startup |

---

## Getting Started

### Prerequisites
- Go 1.22+
- Node.js 18+ & npm
- Docker & Docker Compose

### 1. Start PostgreSQL
```bash
docker compose up -d
```

### 2. Run the Backend API Service
```bash
go run ./cmd/api
```
The service will automatically connect to PostgreSQL, apply all embedded database migrations, and start listening on `:8080`.

### 3. Run the Frontend Dashboard
```bash
cd frontend
npm install
npm run dev
```
Open **[http://localhost:3000](http://localhost:3000)** in your browser. The Vite development server automatically proxies API requests to the Go backend at `:8080`.

---

## Testing & Verification

### Run Backend Unit & Integration Tests with Race Detection
```bash
go test -v -race ./...
```

### Run Backend Benchmarks
```bash
go test -bench=. -benchmem -run=^$ ./internal/repository/db
```

### Build Frontend for Production
```bash
cd frontend
npm run build
```
Outputs static assets into `frontend/dist/`.

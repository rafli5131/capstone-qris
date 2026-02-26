# Plan: QRIS Payment API (Go + Fiber + PostgreSQL + Redis)

**Overview:** Build a QRIS payment microservice in Go following the `khannedy/golang-clean-architecture` clean architecture pattern. The app exposes 3 REST endpoints (inquiry, payment, status), secured by HMAC-SHA256 header signature. Merchant data is cached in Redis. Transactions use optimistic locking. The boilerplate's MySQL driver is swapped for PostgreSQL, Kafka is removed entirely, and Swagger + k6 load tests are added. Everything runs via Docker Compose.

---

## API Endpoints

### Headers (all requests)

| Header | Description | Example |
|---|---|---|
| `X-Timestamp` | Request time in ISO8601 | `2026-02-25T20:30:00Z` |
| `X-Client-Key` | Unique system/merchant identifier | `MK-9921-X` |
| `X-Signature` | `HMAC-SHA256(SecretKey, Payload)` | `a5f8e...` (hex string) |

### Sample QRIS Payload
```
00020101021126690021ID.CO.BANKMANDIRI.WWW01189360000801299399930211712993999340303UKE51440014ID.CO.QRIS.WWW0215ID10232756067300303UKE5204274153033605802ID5912M Ivan Store6015Jakarta Timur (61051364062070703A0163045F26
```

### `GET /api/qris/inquiry/{qris_payload}` (Redis-cached)
```json
{
  "status": "success",
  "data": {
    "merchant_id": "MICH-001",
    "merchant_name": "Toko Berkah Mandiri",
    "terminal_id": "T001",
    "city": "Malang",
    "fixed_amount": 0,
    "inquiry_id": "inq_789abc"
  },
  "metadata": {
    "latency_ms": 45,
    "source": "cache"
  }
}
```

### `POST /api/qris/payment`
**Request:**
```json
{
  "inquiry_id": "inq_789abc",
  "user_id": "user_123",
  "amount": 50000,
  "payment_method": "balance",
  "pincode": "******"
}
```
**Response:**
```json
{
  "status": "processing",
  "transaction_id": "UUID",
  "message": "Transaksi sedang diproses",
  "estimated_completion": "200ms"
}
```

### `GET /api/transaction/status/{transaction_id}`
```json
{
  "transaction_id": "UUID",
  "status": "SUCCESS",
  "final_balance": 450000,
  "timestamp": "2026-02-25T20:15:30Z"
}
```

---

## Steps

### 1. Scaffold Project Structure

Create the project at `/home/rafli/Project/capstone-qris` mirroring the boilerplate layout:

```
capstone-qris/
├── cmd/web/main.go
├── db/migrations/
├── internal/
│   ├── config/          # app.go, fiber.go, gorm.go, redis.go, viper.go, logrus.go, validator.go
│   ├── delivery/http/   # controllers, middleware/, route/
│   ├── entity/
│   ├── model/           # DTOs + converter/
│   ├── repository/
│   └── usecase/
├── api/swagger.yaml
├── test/k6/load_test.js
├── Dockerfile
├── docker-compose.yml
└── config.json
```

---

### 2. Database Migrations (`db/migrations/`)

Four sequential migration pairs (`.up.sql` / `.down.sql`) via `golang-migrate`:

| # | Table | Key fields |
|---|---|---|
| 001 | `api_clients` | `client_id VARCHAR PK`, `client_secret`, `status ENUM('ACTIVE','INACTIVE')`, `created_at` |
| 002 | `merchants` | `merchant_id VARCHAR PK`, `merchant_name`, `mcc`, `city`, `is_active BOOLEAN` |
| 003 | `accounts` | `account_id VARCHAR PK`, `balance DECIMAL(15,2)`, `currency`, `version INT` (optimistic lock) |
| 004 | `transactions` | `transaction_id UUID PK`, `trace_id VARCHAR (indexed)`, `account_id FK`, `merchant_id FK`, `amount DECIMAL`, `status ENUM('PENDING','SUCCESS','FAILED','REVERSED')`, `created_at` |

---

### 3. Entities (`internal/entity/`)

Four entity files with GORM struct tags, matching migration schemas:
- `entity/api_client_entity.go`
- `entity/merchant_entity.go`
- `entity/account_entity.go`
- `entity/transaction_entity.go`

---

### 4. Models / DTOs (`internal/model/`)

- `model.go` — generic `WebResponse[T]` + `Metadata` struct (latency_ms, source)
- `qris_model.go` — `InquiryRequest`, `InquiryResponse` matching the spec
- `payment_model.go` — `PaymentRequest` (inquiry_id, user_id, amount, payment_method, pincode), `PaymentResponse` (status, transaction_id, estimated_completion)
- `transaction_model.go` — `TransactionStatusResponse`
- `converter/` — entity ↔ DTO converters

---

### 5. Configuration (`internal/config/`)

| File | Purpose |
|---|---|
| `viper.go` | Load `config.json` (db, redis, web ports) |
| `gorm.go` | GORM + `gorm.io/driver/postgres`, connection pool |
| `redis.go` | `go-redis/v9` client setup |
| `fiber.go` | Fiber app, global error handler |
| `logrus.go` | Structured logger |
| `validator.go` | go-playground/validator |
| `app.go` | `Bootstrap()` — wires all layers in order |

---

### 6. Repositories (`internal/repository/`)

- `repository.go` — generic `Repository[T]` base (CRUD via generics)
- `api_client_repository.go` — `FindByClientID`
- `merchant_repository.go` — `FindByMerchantID`
- `account_repository.go` — `FindByAccountID`, `UpdateWithOptimisticLock` (uses `version` field + `WHERE version = ?`)
- `transaction_repository.go` — `Create`, `FindByTransactionID`

---

### 7. QRIS Parser (`internal/util/qris_parser.go`)

A custom EMV TLV parser that processes the QRIS payload string (e.g., `000201...`) to extract:
- `merchant_id`
- `merchant_name`
- `city`
- `mcc`
- `fixed_amount`

Used in the inquiry use case on cache miss.

---

### 8. Use Cases (`internal/usecase/`)

**`qris_usecase.go` — Inquiry flow:**
1. Hash the raw QRIS payload → Redis key `qris:inquiry:{hash}`
2. Cache HIT → return with `"source": "cache"`, record latency
3. Cache MISS → parse QRIS payload → query `MerchantRepository` → store in Redis (TTL 5 min) → return with `"source": "database"`
4. Generate and store `inquiry_id` (UUID) in Redis pointing to merchant data

**`payment_usecase.go` — Payment flow:**
1. Validate `inquiry_id` exists in Redis
2. Load `Account` by `user_id`
3. Check balance ≥ amount
4. Start DB transaction; `UPDATE accounts SET balance = balance - ?, version = version + 1 WHERE account_id = ? AND version = ?` (optimistic lock; retry on conflict)
5. Create `Transaction` record with status `PENDING`
6. Commit; launch goroutine to update status to `SUCCESS` asynchronously
7. Return `"status": "processing"` immediately

**`transaction_usecase.go`:**
- `GetStatus(transaction_id)` → query DB, return current status + final_balance

---

### 9. Auth Middleware (`internal/delivery/http/middleware/`)

`signature_middleware.go`:
1. Read `X-Client-Key`, `X-Timestamp`, `X-Signature` headers
2. Look up `ApiClient` by `X-Client-Key`; check `status = ACTIVE`
3. Validate `X-Timestamp` within ±5 minutes (replay attack prevention)
4. Compute `HMAC-SHA256(client_secret, request_body)` → compare hex with `X-Signature`
5. On failure → `401 Unauthorized`

---

### 10. Controllers & Routes (`internal/delivery/http/`)

**`qris_controller.go`:**
- `GET /api/qris/inquiry/:qris_payload` → calls `QrisUseCase.Inquiry()`

**`payment_controller.go`:**
- `POST /api/qris/payment` → calls `PaymentUseCase.Pay()`
- `GET /api/transaction/status/:transaction_id` → calls `TransactionUseCase.GetStatus()`

**`route/route.go`:**
- All 3 routes protected by `SignatureMiddleware`
- Register Swagger UI route via `swagger` fiber middleware

---

### 11. Swagger (`api/swagger.yaml` + `cmd/web/main.go`)

Full OpenAPI 3.0 spec covering:
- All 3 endpoints with request/response schemas
- Request header parameters (X-Timestamp, X-Client-Key, X-Signature)
- Error response schemas (401, 400, 404, 500)

Served live via `github.com/gofiber/swagger` at `GET /swagger/*`.

---

### 12. Docker

**`Dockerfile`** — multi-stage build:
- Stage 1: `golang:1.23-alpine` → build binary
- Stage 2: `alpine:latest` → copy binary + `config.json`

**`docker-compose.yml`** — 3 services:

| Service | Image | Ports |
|---|---|---|
| `app` | built from Dockerfile | `3000:3000` |
| `postgres` | `postgres:16-alpine` | `5432:5432` |
| `redis` | `redis:7-alpine` | `6379:6379` |

Includes health checks, volume mounts for postgres data persistence, and a `migrate` one-shot service to run DB migrations on startup.

---

### 13. k6 Load Test (`test/k6/load_test.js`)

Scenario: ramp to **1000 RPS** over 30s, hold 2 minutes, ramp down.

Tests all 3 endpoints with correct HMAC-SHA256 signature generation (using k6's `crypto` module):
1. `GET /api/qris/inquiry/{sample_payload}` — measures cache hit vs miss
2. `POST /api/qris/payment` — measures P95 latency
3. `GET /api/transaction/status/{id}` — verifies final state

Thresholds: `http_req_duration p(95)<500ms`, `http_req_failed rate<0.01`.

---

### 14. `go.mod` Dependencies

```
github.com/gofiber/fiber/v2
github.com/gofiber/swagger
gorm.io/gorm
gorm.io/driver/postgres
github.com/redis/go-redis/v9
github.com/spf13/viper
github.com/go-playground/validator/v10
github.com/sirupsen/logrus
github.com/golang-migrate/migrate/v4
github.com/google/uuid
github.com/swaggo/swag
```

---

## Verification

1. `docker compose up --build` — all 3 services healthy, migrations applied
2. Hit `GET /swagger/` to validate all endpoints documented
3. Test auth rejection: call without headers → `401`
4. Test inquiry cache: first call → `"source": "database"`, second → `"source": "cache"`
5. Test payment + polling status endpoint for `SUCCESS`
6. `k6 run test/k6/load_test.js` → confirm 1000 RPS, p95 < 500ms

---

## Decisions

- **Framework:** Fiber (fasthttp, chosen for performance target of 1000 RPS)
- **PostgreSQL driver:** `gorm.io/driver/postgres` replacing boilerplate's MySQL
- **Kafka:** removed entirely; payment async handled via goroutine + DB status polling
- **QRIS parsing:** custom EMV TLV parser (no external lib needed — format is deterministic)
- **Optimistic locking:** `version` field on `accounts` table to prevent double-spend without pessimistic DB locks
- **Swagger:** `gofiber/swagger` + `swaggo/swag` annotations, spec at `GET /swagger/*`
- **Migrations:** run automatically via `docker-compose` one-shot `migrate` service at startup

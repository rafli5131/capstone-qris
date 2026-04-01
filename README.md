# QRIS Payment API

High-performance QRIS (Quick Response Code Indonesian Standard) payment microservice built with Go, featuring Redis caching, optimistic locking, and signature-based authentication.

## 🚀 Features

- **QRIS Processing**: Parse and decode QRIS codes from image files or raw strings
- **Payment Processing**: Complete payment flow with inquiry and transaction endpoints
- **Redis Caching**: Fast inquiry lookups with configurable TTL
- **Optimistic Locking**: Race-condition prevention for concurrent transactions
- **Signature Authentication**: HMAC-based request signing for API security
- **Swagger Documentation**: Auto-generated API docs at `/swagger/`
- **Health Checks**: Liveness endpoint for monitoring
- **Docker Support**: Full containerization with Docker Compose
- **Database Migrations**: Versioned schema management with golang-migrate
- **Load Testing**: K6 scripts for performance testing

## 📋 Tech Stack

- **Language**: Go 1.23
- **Web Framework**: [Fiber v2](https://gofiber.io/)
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **ORM**: [GORM](https://gorm.io/)
- **Migration**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Validation**: go-playground/validator
- **QR Decoding**: gozxing
- **Logging**: Logrus
- **API Docs**: Swagger/OpenAPI

## 📦 Prerequisites

- **Go** 1.23+ (for local development)
- **Docker** & **Docker Compose** (recommended)
- **Make** (optional, for convenience commands)

## 🛠️ Installation

### Clone the Repository

```bash
git clone https://github.com/rafli5131/capstone-qris.git
cd capstone-qris
```

### Setup with Docker (Recommended)

```bash
# Start all services (PostgreSQL, Redis, Migrations, App)
make docker-up

# Or without make:
docker compose up --build -d
```

The API will be available at `http://localhost:3000`

### Setup for Local Development

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Start PostgreSQL and Redis:**
   ```bash
   # Using Docker Compose for just the databases:
   docker compose up postgres redis -d
   ```

3. **Run migrations:**
   ```bash
   make migrate-up
   ```

4. **Update configuration:**
   Edit `config.json` if needed (default values work with Docker services).

5. **Run the application:**
   ```bash
   make run
   # Or: go run ./cmd/web/main.go
   ```

## ⚙️ Configuration

Configuration is loaded from `config.json` and can be overridden by environment variables.

**Example `config.json`:**

```json
{
  "web": {
    "port": 3000,
    "prefork": false
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "qris_user",
    "password": "qris_pass",
    "name": "qris_db",
    "sslmode": "disable"
  },
  "redis": {
    "addr": "localhost:6379",
    "inquiry_ttl_seconds": 300
  },
  "app": {
    "log_level": "info",
    "jwt_secret": "supersecretkey"
  }
}
```

**Environment Variable Mapping:**

- `DATABASE_HOST` → `database.host`
- `DATABASE_PORT` → `database.port`
- `DATABASE_USERNAME` → `database.username`
- `DATABASE_PASSWORD` → `database.password`
- `DATABASE_NAME` → `database.name`
- `REDIS_ADDR` → `redis.addr`
- `WEB_PORT` → `web.port`
- `APP_JWT_SECRET` → `app.jwt_secret`

## 🔌 API Endpoints

### Health Check
- `GET /livez` - Liveness probe

### QRIS Operations
- `GET /api/qris/inquiry/{qris_payload}` - Query merchant info by QRIS payload
- `POST /api/qris/inquiry/image` - Upload QR image and decode to inquiry
- `POST /api/qris/merchant/image` - Register or reactivate a merchant from QRIS image
- `POST /api/qris/payment` - Submit payment
- `GET /api/qris/status/{transaction_id}` - Check transaction status

### Merchant Operations
- `GET /api/merchant/{merchant_id}/income` - View merchant income and transaction summary
- `GET /api/merchant/{merchant_id}/transactions` - List all transactions for a merchant

### Admin Operations
- `GET /api/admin/transactions` - List all transactions
- `PUT /api/admin/transactions/{transaction_id}` - Edit transaction status or amount
- `GET /api/admin/api-clients` - List API clients
- `POST /api/admin/api-clients` - Create API client
- `PUT /api/admin/api-clients/{client_id}` - Update client secret and status

### Authentication

Public authentication endpoints:
- `POST /api/auth/register` - Create a user account with `username`, `password`, and `initial_balance`
- `POST /api/auth/login` - Authenticate with `username` and `password`

All other `/api/*` endpoints require a valid JWT in the `Authorization` header.

**Headers:**
- `Authorization`: `Bearer {token}`

**Seeded accounts** (migrated automatically):
- `admin` / `adminpass` → role `ADMIN`
- `user` / `userpass` → role `USER`
- `merchant` / `merchantpass` → role `MERCHANT`

## 📖 API Documentation

Interactive Swagger documentation is available at:

```
http://localhost:3000/swagger/
```

To regenerate Swagger docs:

```bash
# Install swag CLI:
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs:
make swagger
```

## 📁 Project Structure

```
capstone-qris/
├── cmd/
│   └── web/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration & DI
│   │   ├── app.go
│   │   ├── fiber.go
│   │   ├── gorm.go
│   │   ├── redis.go
│   │   └── validator.go
│   ├── delivery/
│   │   └── http/                # HTTP handlers (controllers)
│   │       ├── payment_controller.go
│   │       ├── qris_controller.go
│   │       ├── middleware/      # Request middlewares
│   │       └── route/           # Route definitions
│   ├── usecase/                 # Business logic
│   │   ├── payment_usecase.go
│   │   ├── qris_usecase.go
│   │   └── transaction_usecase.go
│   ├── repository/              # Data access layer
│   │   ├── account_repository.go
│   │   ├── api_client_repository.go
│   │   ├── merchant_repository.go
│   │   └── transaction_repository.go
│   ├── entity/                  # Domain entities
│   ├── model/                   # Request/Response DTOs
│   └── util/                    # Utility functions
│       ├── qr_image.go          # QR image decoder
│       └── qris_parser.go       # QRIS payload parser
├── db/
│   └── migrations/              # Database migrations
├── docs/                        # Swagger generated docs
├── test/
│   └── k6/
│       └── load_test.js         # Load testing script
├── config.json                  # Configuration file
├── docker-compose.yml           # Docker orchestration
├── Dockerfile                   # Application container
├── Makefile                     # Build automation
└── go.mod                       # Go dependencies
```

## 🗄️ Database Migrations

### Run Migrations

```bash
# Up migrations
make migrate-up

# Down migrations (rollback all)
make migrate-down
```

Migrations are automatically run during `docker compose up` via the `migrate` service.

## 🧪 Testing

### Load Testing with K6

```bash
# Install k6: https://k6.io/docs/getting-started/installation/

# Run load test
make test-k6

# Run with JSON output
make test-k6-json
```

## 🚦 Development

### Available Make Commands

```bash
make run          # Run locally
make build        # Build binary
make tidy         # Tidy dependencies
make swagger      # Generate Swagger docs
make docker-up    # Start all services
make docker-down  # Stop all services
make migrate-up   # Run migrations
make migrate-down # Rollback migrations
make test-k6      # Run load tests
```

### Running Without Docker

1. Ensure PostgreSQL and Redis are running locally
2. Update `config.json` with correct connection details
3. Run migrations: `make migrate-up`
4. Start the app: `make run`

## 🐳 Docker Services

The `docker-compose.yml` defines:

- **postgres**: PostgreSQL 16 database
- **redis**: Redis 7 cache
- **migrate**: One-shot migration runner
- **app**: Main QRIS API service

All services are on the `qris_net` bridge network.

## 📝 Environment Variables for Docker

Set in `docker-compose.yml` under `app.environment`:

```yaml
DATABASE_HOST: postgres
DATABASE_PORT: "5432"
DATABASE_USERNAME: qris_user
DATABASE_PASSWORD: qris_pass
DATABASE_NAME: qris_db
REDIS_ADDR: redis:6379
WEB_PORT: "3000"
```

## 🔒 Security

- **Signature Verification**: All API requests must be signed with HMAC-SHA256
- **Timestamp Validation**: Prevents replay attacks with configurable tolerance
- **Environment Variables**: Sensitive configs can be injected via env vars

## 📊 Monitoring

- **Health Endpoint**: `GET /livez` returns `OK` if service is healthy
- **Docker Health Checks**: Configured for all services in `docker-compose.yml`
- **Logging**: Structured JSON logging with Logrus

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👥 Authors

- **Rafli** - [GitHub](https://github.com/rafli)

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Fast and flexible Go web framework
- [GORM](https://gorm.io/) - Fantastic ORM library for Go
- [gozxing](https://github.com/makiuchi-d/gozxing) - QR code decoding library

---

**Built with ❤️ using Go**

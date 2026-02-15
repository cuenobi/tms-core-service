# TMS Core Service

Transportation Management System Core Service - A Golang microservice built with Clean Architecture.

## Tech Stack

- **Framework**: Gin (HTTP router)
- **CLI**: Cobra
- **Database**: PostgreSQL with GORM
- **Cache**: Redis
- **Migration**: golang-migrate
- **Authentication**: JWT
- **Documentation**: Swagger (swaggo)

## Project Structure

```
tms-core-service/
├── cmd/                    # CLI commands (Cobra)
├── db/migrations/          # Database migrations
├── docs/                   # Swagger documentation
├── internal/               # Internal application code
│   ├── api/               # Delivery layer (HTTP handlers, routes)
│   ├── config/            # Configuration
│   ├── domain/            # Domain layer (entities, interfaces)
│   ├── infra/             # Infrastructure layer (DB, Redis)
│   ├── server/            # Server setup
│   ├── usecase/           # Application business logic
│   └── util/              # Utility functions
├── pkg/                   # Shared packages
├── docker-compose.yaml
├── Dockerfile
├── env.yaml               # Configuration file
├── go.mod
├── main.go
└── Makefile
```

## Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Make (optional, for convenience)

## Quick Start

### 1. Install Dependencies

```bash
go mod download
```

### 2. Start Infrastructure (PostgreSQL & Redis)

```bash
docker-compose up -d postgres redis
```

### 3. Run Migrations

```bash
make migrate-up
# or
go run main.go migrate up
```

### 4. Generate Swagger Docs

```bash
make swagger
# or
swag init -g main.go -o ./docs
```

### 5. Start Server

```bash
make run
# or
go run main.go serve
```

The server will start at `http://localhost:8080`

## Available Commands

### Makefile

```bash
make build              # Build the application
make run                # Run the application
make test               # Run tests
make swagger            # Generate Swagger docs
make migrate-up         # Run migrations up
make migrate-down       # Roll back migrations
make migrate-create     # Create new migration (usage: make migrate-create name=create_users)
make docker-build       # Build Docker image
make docker-up          # Start all Docker services
make docker-down        # Stop all Docker services
make clean              # Clean build artifacts
make install-tools      # Install dev tools (swag, migrate)
```

### CLI Commands

```bash
# Show help
go run main.go --help

# Start server
go run main.go serve

# Run migrations
go run main.go migrate up
go run main.go migrate down

# Create new migration
go run main.go new-migration create_users_table

# Print current configuration
go run main.go print-config
```

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login

### Protected Endpoints (Require JWT)

- Protected routes will be added as the application grows

### Swagger Documentation

Visit `http://localhost:8080/swagger/index.html` to view the API documentation.

## Bruno API Collection

This project includes a [Bruno](https://www.usebruno.com) API collection for easy API testing located in the `bruno/` directory.

**Quick Start:**
1. Install Bruno from https://www.usebruno.com/downloads
2. Open Bruno → "Open Collection" → Select the `bruno` folder
3. Run requests: Health Check → Register → Login (tokens auto-saved)

The collection includes pre-configured requests with automatic JWT token management. See [bruno/README.md](bruno/README.md) for details.

## Configuration

The application uses `env.yaml` for configuration. Key settings:

- **Server**: Port, mode, timeouts
- **Database**: PostgreSQL connection settings
- **Redis**: Cache configuration
- **JWT**: Secret and token expiry

For production, consider using environment variables or secrets management.

## Development

### Project Layout Philosophy

This project follows Clean Architecture principles:

- **Domain Layer** (`internal/domain/`): Business entities and repository interfaces
- **Use Case Layer** (`internal/usecase/`): Application business logic
- **Infrastructure Layer** (`internal/infra/`): External implementations (DB, Redis)
- **Delivery Layer** (`internal/api/`): HTTP handlers and routes

### Adding a New Feature

1. Define entity in `internal/domain/entity/`
2. Define repository interface in `internal/domain/repository/`
3. Implement repository in `internal/infra/db/repository/`
4. Create use case in `internal/usecase/`
5. Create HTTP handler in `internal/api/http/handler/`
6. Register routes in `internal/api/http/route/`
7. Wire dependencies in `internal/server/dependency.go`

## Testing

```bash
make test
# or
go test -v ./...
```

## Docker Deployment

### Build and run with Docker Compose

```bash
docker-compose up --build
```

This will start PostgreSQL, Redis, and the application.

## License

MIT

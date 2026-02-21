# TMS Core Service — AI Rules

> These rules guide all AI-assisted code generation, modification, and review for the `tms-core-service` project. Always follow them when implementing any feature, fix, or refactor.

---

## 1. Project Overview

| Key              | Value                                              |
| ---------------- | -------------------------------------------------- |
| **Module**       | `tms-core-service`                                 |
| **Go version**   | 1.24+                                              |
| **Architecture** | Clean Architecture (4 layers)                      |
| **HTTP**         | GoFiber v2 (`github.com/gofiber/fiber/v2`)         |
| **ORM**          | GORM v1 (`gorm.io/gorm`) + PostgreSQL              |
| **Cache**        | Redis (`github.com/redis/go-redis/v9`)             |
| **CLI**          | Cobra (`github.com/spf13/cobra`)                   |
| **Config**       | Viper (`github.com/spf13/viper`) — YAML + env vars |
| **Auth**         | JWT (`github.com/golang-jwt/jwt/v5`), bcrypt, Google OAuth, LINE OAuth |
| **Docs**         | Swagger via `swag` annotations                     |
| **Migrations**   | `golang-migrate/migrate/v4`                        |
| **Validation**   | `go-playground/validator/v10`                      |

---

## 2. Directory Structure & Layer Responsibilities

```
├── cmd/                          # CLI commands (Cobra) — entry points only
├── db/migrations/                # SQL migration files
├── docs/                         # Generated Swagger documentation
├── internal/                     # All application-internal code
│   ├── api/http/                 # Delivery Layer (HTTP)
│   │   ├── dto/                  # Request/Response DTOs (JSON tags + validation)
│   │   ├── handler/<domain>/     # HTTP handlers grouped by domain
│   │   ├── middleware/           # HTTP middleware (JWT, CORS, Trace, Recover)
│   │   └── route/                # Route registration & Dependencies struct
│   ├── config/                   # AppConfig struct + Viper loader
│   ├── domain/                   # Domain Layer (pure, zero dependencies)
│   │   ├── cache/                # Cache interfaces
│   │   ├── entity/               # Domain entities (plain Go structs)
│   │   ├── errs/                 # Domain sentinel errors + ValidationErrors
│   │   ├── repository/           # Repository interfaces
│   │   └── service/              # Service interfaces (HashService, TokenService)
│   ├── infra/                    # Infrastructure Layer (implementations)
│   │   ├── db/
│   │   │   ├── connection.go     # GORM connection factory
│   │   │   ├── transactor.go     # Context-based transaction management
│   │   │   ├── model/            # GORM models with ToEntity()/FromEntity()
│   │   │   └── repository/       # Repository implementations per domain
│   │   ├── redis/                # Redis connection + CacheRepository impl
│   │   └── service/              # Service implementations (hash, token)
│   ├── server/                   # Server bootstrap
│   │   ├── server.go             # Fiber app creation + startup + shutdown
│   │   ├── middleware.go         # Global middleware application
│   │   └── dependency.go         # Manual dependency injection (wire)
│   ├── usecase/                  # Use Case Layer (business logic)
│   │   └── <domain>/             # Use case structs per domain
│   └── util/                     # Internal utilities
│       ├── apierror/             # APIError struct + error codes
│       ├── httpresponse/         # Standardized HTTP response helpers
│       └── validator/            # Struct validation wrapper
├── pkg/                          # Shared, externally-importable packages
│   ├── hash/                     # bcrypt helpers
│   └── jwt/                      # JWT service + claims
├── main.go                       # Entry point with Swagger annotations
├── env.yaml / env.example.yaml   # Configuration files
├── Makefile                      # Common dev & build tasks
├── Dockerfile                    # Multi-stage Docker build
└── docker-compose.yaml           # Local dev environment
```

---

## 3. Clean Architecture Rules

### 3.1 Dependency Rule (CRITICAL)

```
Domain ← UseCase ← Infrastructure
                  ← Delivery (api/http)
```

- **Domain** (`internal/domain/`) MUST have **zero** imports from any other `internal/` package. It defines only pure Go structs, interfaces, and sentinel errors.
- **UseCase** (`internal/usecase/`) imports **only** from `domain/`. Never import infra, handler, or any framework package.
- **Infrastructure** (`internal/infra/`) implements domain interfaces. It imports `domain/` and external libraries (GORM, Redis, etc.).
- **Delivery** (`internal/api/http/`) imports `usecase/` and `util/`. Handlers MUST NOT import `domain/repository/` or `infra/` directly.
- **Server** (`internal/server/`) is the composition root — it is the ONLY package that wires infra → usecase → handler.

### 3.2 Never Violate

- Do NOT import `gorm`, `fiber`, `redis`, or any framework in `domain/` or `usecase/`.
- Do NOT call repository methods directly from handlers — always go through a use case.
- Do NOT put business logic in handlers — handlers should only parse requests, call use cases, and format responses.

---

## 4. Implementation Patterns

### 4.1 Adding a New Feature (End-to-End Checklist)

When adding a new domain feature (e.g., `vehicle`), create files in this order:

1. **Domain Entity** — `internal/domain/entity/vehicle.go`
   - Pure Go struct, no GORM tags, no JSON tags.
   - Use `uuid.UUID` for IDs, `time.Time` for timestamps, pointer types for nullable fields.

2. **Domain Repository Interface** — `internal/domain/repository/vehicle_repository.go`
   - Define the interface with `context.Context` as first arg in every method.
   - Return `*entity.Vehicle` (pointer) for single entities, `[]*entity.Vehicle` for lists.
   - Return `errs.ErrNotFound` when entity does not exist.

3. **Domain Service Interface** (if needed) — `internal/domain/service/interfaces.go`
   - Add new interface methods here for cross-cutting abstractions.

4. **Use Case** — `internal/usecase/vehicle/`
   - `vehicle_usecase.go` — Struct with constructor `NewVehicleUseCase(...)`, depends only on domain interfaces.
   - `io.go` — Input/Output structs for the use case (never share DTOs with handlers).
   - Wrap errors with `fmt.Errorf("context: %w", err)`.
   - Use `errors.Is(err, errs.ErrNotFound)` for domain error checking.

5. **DB Model** — `internal/infra/db/model/vehicle_model.go`
   - GORM struct with tags. Implement `TableName()`, `ToEntity()`, and `FromEntity()`.
   - Map `gorm.DeletedAt` ↔ `*time.Time` for soft deletes.

6. **Repository Implementation** — `internal/infra/db/repository/vehicle/vehicle_repo.go`
   - Struct implementing domain interface. Use `db.FromContext(ctx, r.db)` for transaction support.
   - Map `gorm.ErrRecordNotFound` → `errs.ErrNotFound`, `gorm.ErrDuplicatedKey` → `errs.ErrConflict`.
   - Constructor returns domain interface type: `func NewVehicleRepository(db *gorm.DB) repository.VehicleRepository`.

7. **DTO** — `internal/api/http/dto/vehicle_dto.go`
   - Request structs with `json:"..."` and `validate:"..."` tags.
   - Response structs with `json:"..."` tags only.

8. **Handler** — `internal/api/http/handler/vehicle/handler.go`
   - Struct with constructor accepting use case pointers.
   - Each method: parse request → validate → call use case → map to DTO response.
   - Use `httpresponse.Success()`, `httpresponse.Created()`, `httpresponse.Error()`, `httpresponse.Paginated()`.
   - Add Swagger godoc annotations above every handler method.
   - Use `c.Context()` for context, `validator.Validate(req)` for validation.

9. **Routes** — `internal/api/http/route/route.go`
   - Add handler to `Dependencies` struct.
   - Register routes in `SetupRoutes()`. Use `protected` group for JWT-authenticated routes.

10. **Dependency Wiring** — `internal/server/dependency.go`
    - Instantiate repo → use case → handler. Pass to `route.Dependencies`.

11. **Migration** — `db/migrations/` via `make migrate-create name=create_vehicles`

### 4.2 Naming Conventions

| Item                     | Convention                                      | Example                         |
| ------------------------ | ----------------------------------------------- | ------------------------------- |
| Package names            | `lowercase`, singular                           | `vehicle`, `auth`, `healthcheck`|
| File names               | `snake_case.go`                                 | `vehicle_usecase.go`            |
| Struct names             | `PascalCase`                                    | `VehicleUseCase`                |
| Interface names          | `PascalCase`, suffix with concept                | `VehicleRepository`             |
| Constructor functions    | `New<StructName>`                               | `NewVehicleUseCase`             |
| Handler packages         | One package per domain under `handler/`          | `handler/vehicle/`              |
| Repository packages      | One package per domain under `repository/`       | `repository/vehicle/`           |
| Use case packages        | One package per domain under `usecase/`           | `usecase/vehicle/`              |
| DB model packages        | All models in `infra/db/model/`                  | `model/vehicle_model.go`        |
| Domain error vars        | `Err<Name>` sentinel errors                      | `errs.ErrNotFound`              |
| API error codes          | `SCREAMING_SNAKE_CASE` ErrorCode constants        | `CodeNotFound`                  |

### 4.3 Error Handling

- **Domain errors** live in `internal/domain/errs/errors.go` as sentinel `var` declarations.
- **Use cases** wrap errors with `fmt.Errorf("<context>: %w", err)` for chain tracing.
- **Handlers** call `httpresponse.Error(c, err)` which automatically maps domain errors → `APIError` → HTTP status codes via `mapToAPIError()`.
- **API errors** (`internal/util/apierror/`) have `Code`, `Message`, `StatusCode`, `TraceID`, and optional `Errors` map.
- Never expose internal/infrastructure errors to the client — they are mapped to `INTERNAL_SERVER_ERROR`.

### 4.4 HTTP Response Format

All API responses use the standard `httpresponse.Response` struct:

```json
{
  "success": true|false,
  "message": "Human-readable message",
  "data": { ... },
  "error": "...",
  "errors": { "field": ["tag1", "tag2"] }
}
```

- Success: `httpresponse.Success(c, data, "message")`
- Created: `httpresponse.Created(c, data, "message")`
- Error: `httpresponse.Error(c, err)`
- Paginated: `httpresponse.Paginated(c, data, total, limit, offset)`

### 4.5 Validation

- Use `go-playground/validator/v10` struct tags on DTO request structs.
- Call `validator.Validate(req)` in handlers before calling use cases.
- Returns `errs.ValidationErrors` (a `map[string][]string`) which maps to `VALIDATION_ERROR` API error.

### 4.6 Middleware

- **Trace**: Adds `X-Trace-ID` header/context to every request.
- **Recover**: Catches panics and returns standardized 500 error.
- **CORS**: Configured in `internal/api/http/middleware/cors.go`.
- **JWT Auth**: Applied via `middleware.JWTAuth(deps.JWTService)` on protected route groups.
- JWT context values: `GetUserID(c)`, `GetUserEmail(c)`, `GetUserStatus(c)`.

### 4.7 Database & Transactions

- Use `db.FromContext(ctx, r.db)` in ALL repository methods to support context-carried transactions.
- Transaction support via `Transactor.WithTransaction(ctx, func(ctx) error)`.
- GORM connection pool is configured from `env.yaml` pool settings.
- All times stored in UTC (`time.Now().UTC()` via GORM NowFunc).

### 4.8 Model ↔ Entity Mapping

- DB models live in `internal/infra/db/model/` with GORM struct tags.
- Each model has:
  - `TableName() string` — explicit table name.
  - `ToEntity() *entity.X` — convert model → domain entity.
  - `FromEntity(e *entity.X) *model.X` — convert domain entity → model (package-level func).
- Soft deletes: `gorm.DeletedAt` in models ↔ `*time.Time` in entities.

### 4.9 Configuration

- Primary config: `env.yaml` (Viper, YAML format).
- Environment variable overrides use `TMS_` prefix or direct bindings for K8s/Docker (e.g., `DATABASE_HOST`, `JWT_SECRET`).
- Sensitive values (secrets, OAuth credentials) must be overridden via environment variables, NEVER hardcoded.

### 4.10 Swagger Documentation

- Main annotations in `main.go` (`@title`, `@version`, `@host`, `@BasePath`, `@securityDefinitions.apikey`).
- Handler method annotations with godoc-style `@Summary`, `@Description`, `@Tags`, `@Param`, `@Success`, `@Failure`, `@Router`.
- Regenerate after changes: `make swagger`.

---

## 5. Code Style & Quality

- **Imports**: Group in this order: (1) standard library, (2) project imports, (3) third-party. Separate with blank lines.
- **Comments**: Every exported type, function, and interface MUST have a godoc comment.
- **Error wrapping**: Always use `fmt.Errorf("<layer>: <operation>: %w", err)`.
- **Context propagation**: All repository/service/use case methods take `context.Context` as first argument.
- **Constructor injection**: Dependencies are injected through constructor functions (`New*`), never via globals or service locators.
- **No circular imports**: Enforce by keeping layers isolated as described in Section 3.

---

## 6. Infrastructure Service Pattern

When wrapping `pkg/` implementations to satisfy domain interfaces:

1. Create wrapper in `internal/infra/service/<name>/`.
2. Struct holds reference to the `pkg/` concrete type.
3. Constructor returns the domain interface type (e.g., `service.HashService`).
4. Methods delegate to `pkg/` implementation.

Example: `internal/infra/service/hash/bcrypt.go` wraps `pkg/hash/` to implement `service.HashService`.

---

## 7. Do & Don't Quick Reference

### ✅ DO

- Follow the 4-layer Clean Architecture dependency rule strictly.
- Create separate Input/Output structs in `usecase/<domain>/io.go`.
- Use domain sentinel errors from `errs` package.
- Use `db.FromContext()` in every repository method.
- Add Swagger annotations to every handler method.
- Use `httpresponse.*` helpers for all HTTP responses.
- Return domain interface types from infra constructors.
- Soft delete via GORM's `DeletedAt` column.
- Use `context.Context` as first parameter everywhere.

### ❌ DON'T

- Import `gorm`, `fiber`, or `redis` in `domain/` or `usecase/`.
- Share DTOs between handler and use case layers (use `io.go` structs in use case).
- Expose raw infrastructure errors to API clients.
- Put business logic in handlers (keep them thin).
- Hardcode secrets in config files or code.
- Skip validation in handlers before calling use cases.
- Use magic strings for error codes — use `apierror.Code*` constants.
- Create circular package dependencies.

---

## 8. Common Commands

```bash
make run                    # Start the server (go run main.go serve)
make build                  # Build binary
make test                   # Run tests
make swagger                # Regenerate Swagger docs
make migrate-up             # Run pending migrations
make migrate-down           # Rollback last migration
make migrate-create name=X  # Create new migration
make docker-up              # Start Docker Compose (PostgreSQL, Redis)
make docker-down            # Stop Docker Compose
make tidy                   # go mod tidy
```

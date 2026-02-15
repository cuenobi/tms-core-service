├── cmd/                          # Application entry points and CLI commands
│   ├── root.go                   # Root command for the CLI
│   ├── serve_cmd.go              # Command to start the server
│   ├── migrate_cmd.go            # Command to run database migrations
│   ├── new_migration_cmd.go      # Command to create new migration files
│   ├── print_config_cmd.go       # Command to print the current configuration
│   └── generate_sql_builder.go   # Command to generate SQL builder files
│   └── ...                       # Other commands
├── db/                           # Database-related files
│   └── migrations/               # Database migration files
├── docs/                         # Documentation files
│   ├── swagger.yaml              # API documentation in Swagger format
│   └── ...                       # Other design and architecture documents
├── internal/                     # Internal application code
│   ├── api/                      # Delivery layer
│   │   └── http/                 # HTTP API components
│   │       ├── handler/          # HTTP handlers by domain
│   │       │   ├── healthcheck/  # Health check handlers
│   │       │   └── ...           # Other domain handlers
│   │       ├── middleware/       # HTTP middleware
│   │       └── route/            # HTTP route definitions
│   ├── config/                   # Application configuration
│   ├── domain/                   # Domain layer
│   │   ├── cache/                # Cache-related interfaces
│   │   ├── entity/               # Domain models (e.g., Concert, Reservation)
│   │   ├── errs/                 # Domain-specific errors
│   │   └── repository/           # Repository interfaces
│   ├── infra/                    # Infrastructure layer
│   │   ├── db/                   # Database implementations
│   │   │   ├── connection.go     # Database connection setup
│   │   │   ├── sql_execer.go     # SQL execution interface
│   │   │   ├── transactor.go     # Transaction management
│   │   │   ├── mocks/            # Mock implementations
│   │   │   ├── model_gen/        # Generated models from DB schema
│   │   │   └── repository/       # Repository implementations
│   │   │       ├── healthcheck/  # Health check repository
│   │   │       └── ...           # Other repositories
│   │   └── redis/                # Redis implementations
│   │       ├── client.go         # Redis client interface
│   │       ├── connection.go     # Redis connection setup
│   │       ├── mocks/            # Mock implementations
│   │       └── repository/       # Repository implementations
│   │           ├── seat/         # Seat locking and cache implementations
│   │           └── ...           # Other repositories
│   ├── server/                   # Server setup and initialization
│   │   ├── dependency.go         # Dependency injection
│   │   ├── middleware.go         # Server middleware
│   │   └── server.go             # HTTP server setup
│   ├── usecase/                  # Application business logic
│   │   ├── healthcheck/          # Health check use case
│   │   └── ...                   # Other use cases
│   └── util/                     # Utility functions
│       ├── httpresponse/         # HTTP response helpers
│       └── ...                   # Other utility functions
├── pkg/                          # Shared helper packages
├── docker-compose.yaml           # Docker Compose for local development
├── Dockerfile                    # Docker build definition
├── env.yaml                      # Environment variables configuration
├── go.mod                        # Go module definition
├── go.sum                        # Go module dependencies
├── main.go                       # Main application entry point
├── Makefile                      # Makefile for common tasks
└── README.md                     # Project documentation
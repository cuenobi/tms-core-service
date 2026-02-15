.PHONY: build run test swagger migrate-up migrate-down migrate-create docker-build docker-up docker-down clean

APP_NAME=tms-core-service
BUILD_DIR=./bin
MIGRATION_DIR=./db/migrations

build:
	@echo "Building application..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) main.go

run:
	@echo "Running application..."
	@go run main.go serve

test:
	@echo "Running tests..."
	@go test -v ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g main.go -o ./docs

migrate-up:
	@echo "Running migrations up..."
	@go run main.go migrate up

migrate-down:
	@echo "Running migrations down..."
	@go run main.go migrate down

migrate-create:
	@echo "Creating new migration..."
	@go run main.go new-migration $(name)

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-up:
	@echo "Starting Docker Compose services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker Compose services..."
	@docker-compose down

docker-logs:
	@docker-compose logs -f

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -rf ./docs

install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

tidy:
	@echo "Tidying go modules..."
	@go mod tidy

deps:
	@echo "Downloading dependencies..."
	@go mod download

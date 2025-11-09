# Essential Commands for Fusion Development

## Setup & Installation

### Install Development Tools
```bash
make install-tools
```
Installs: Air (hot reload), Swag (Swagger), golangci-lint (linter)

## Development Server

### Run Development Server (with hot reload)
```bash
air
```

### Run Server Directly
```bash
go run cmd/server/main.go serve
```

### Run with Custom Config
```bash
go run cmd/server/main.go serve --config=./configs/config.yaml
```

### Run with Custom Port
```bash
go run cmd/server/main.go serve --port=3000
```

## Building

### Build Production Binary
```bash
make build
# Output: bin/fusion
```

### Run Built Binary
```bash
./bin/fusion serve
```

## Docker

### Build Docker Image (Single Architecture)
```bash
make docker-build
```

### Build Multi-Architecture Docker Image
```bash
make docker-build-multiarch
```

### Build and Push to Registry
```bash
make docker-push
```

### Run with Docker Compose
```bash
make docker-up          # Start all services
make docker-logs         # View logs
make docker-down         # Stop services
```

### Manual Docker Commands
```bash
# Start services
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down
```

## Database

### Run Migrations
```bash
go run cmd/server/main.go migrate up
```

### Generate New Ent Schema
```bash
make ent-new name=ModelName
```

### Generate Ent Code After Schema Changes
```bash
make generate-ent
```

## Code Quality

### Run Linter
```bash
make lint
```

### Generate Swagger Documentation
```bash
make generate-swagger
```

### Format Swagger Annotations
```bash
make format-swagger
```

## Testing

### Run All Tests
```bash
make test
```

### Run Unit Tests Only
```bash
make test-unit
```

### Run Integration Tests
```bash
make test-integration          # Uses SQLite
make test-integration-postgres # Uses PostgreSQL
```

### Run E2E Tests
```bash
make test-e2e                  # Uses SQLite
make test-e2e-postgres         # Uses PostgreSQL
```

### Run All Tests with Coverage
```bash
make test-coverage
```

### Run Short Tests
```bash
make test-short
```
Skips long-running tests

## Configuration

### Using Environment Variables
All config can be overridden with environment variables:
```bash
FUSION_SERVER_PORT=8080
FUSION_DATABASE_DSN="postgres://..."
FUSION_JWT_SECRET="your-secret"
```

### Config File Location
- Main config: `configs/config.yaml`
- Environment-specific: `configs/config.{env}.yaml`
- Override with: `--config=./configs/config.yaml`

## Swagger Documentation
- Swagger UI: `/swagger/*` route
- Generate docs: `make generate-swagger`
- Uses Swag annotations in handler files

# Customer Service Microservice

A Go-based microservice for customer management with MongoDB, extracted from the NestJS TC2-G38 project.

## Features

- Clean Architecture with separation of concerns
- Domain-driven design
- CPF and Email validation
- MongoDB as NoSQL database
- RESTful API with Gin framework
- Comprehensive unit and integration tests
- Docker and docker-compose support
- Same endpoint contracts as the original NestJS service

## Architecture

```
customer-service-go/
├── api/              # Application entry point
├── internal/
│   ├── domain/          # Business entities and rules
│   ├── usecase/         # Business logic
│   ├── repository/      # Data persistence layer
│   └── handler/         # HTTP handlers
├── pkg/
│   ├── validator/       # Validation utilities (CPF, Email)
│   └── errors/          # Custom error types
└── test/                # Integration tests
```

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized deployment)
- MongoDB 7.0+ (if running locally without Docker)

## Installation

### Using Docker Compose (Recommended)

1. Clone the repository and navigate to the service directory:
```bash
cd customer-service-go
```

2. Create a `.env` file or use environment variables (see Configuration section)

3. Start the services:
```bash
docker-compose up -d
```

The service will be available at `http://localhost:8080`

### Local Development

1. Install dependencies:
```bash
go mod download
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Start MongoDB (if not using Docker):
```bash
docker run -d -p 27017:27017 --name mongodb mongo:7.0
```

4. Run the application:
```bash
go run api/main.go
```

Or using Make:
```bash
make run
```

## Configuration

The service uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `MONGODB_URI` | MongoDB connection URI | `mongodb://mongodb:27017` |
| `MONGODB_DATABASE` | Database name | `customer_db` |
| `MONGODB_PORT` | MongoDB port (for docker-compose) | `27017` |
| `PORT` | Server port | `8080` |

### Local Development

Create a `.env` file:
```bash
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=customer_db
MONGODB_PORT=27017
PORT=8080
```

### Production/CI/CD

Set these as environment variables or GitHub secrets:
```bash
export MONGODB_URI="mongodb://your-production-host:27017"
export MONGODB_DATABASE="customer_db"
export PORT="8080"
```

### GitHub Actions/CI

Set these as repository secrets:
- `MONGODB_URI` - Production MongoDB URI
- `MONGODB_DATABASE` - Database name
- `MONGODB_PORT` - MongoDB port
- `PORT` - Application port
- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub password/token

## API Endpoints

All endpoints maintain the same contracts as the original NestJS service.

### Create Customer
```http
POST /customer
Content-Type: application/json

{
  "name": "John Doe",
  "cpf": "111.444.777-35",
  "email": "john@example.com"
}
```

**Response (201 Created):**
```json
{
  "id": "uuid",
  "name": "John Doe",
  "cpf": "11144477735",
  "email": "john@example.com",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Get Customer by CPF
```http
GET /customer/:cpf
```

**Response (200 OK):**
```json
{
  "id": "uuid",
  "name": "John Doe",
  "cpf": "11144477735",
  "email": "john@example.com",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z"
}
```

### Update Customer
```http
PATCH /customer/:id
Content-Type: application/json

{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

**Response (200 OK):**
```json
{
  "id": "uuid",
  "name": "Jane Doe",
  "cpf": "11144477735",
  "email": "jane@example.com",
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T12:00:00Z"
}
```

### Delete Customer
```http
DELETE /customer/:id
```

**Response (204 No Content)**

### Health Check
```http
GET /health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "service": "customer-service"
}
```

## Error Responses

All errors follow the same format as the NestJS service:

```json
{
  "message": "Error description",
  "statusCode": 400,
  "error": "ERROR_CODE"
}
```

### Error Codes

- `NAME_EMPTY` (400): Name cannot be empty
- `INVALID_CPF` (400): Invalid CPF format
- `INVALID_EMAIL` (400): Invalid email format
- `CUSTOMER_ALREADY_EXISTS` (409): Customer with same CPF or email already exists
- `CUSTOMER_NOT_FOUND` (404): Customer not found
- `INTERNAL_ERROR` (500): Internal server error

## Testing

### Run All Tests
```bash
make test
```

Or:
```bash
go test -v ./...
```

### Run Unit Tests Only
```bash
make test-unit
```

Or:
```bash
go test -v -short ./...
```

### Run Integration Tests

Integration tests require MongoDB. You can run them with the provided script:

```bash
./test-integration.sh
```

Or manually:
```bash
# Start MongoDB for testing
docker run -d --name mongodb-test -p 27018:27017 mongo:7.0

# Run integration tests
TEST_MONGODB_URI="mongodb://localhost:27018" go test -v -tags=integration ./internal/repository/

# Cleanup
docker stop mongodb-test && docker rm mongodb-test
```

### Test Coverage
```bash
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -html=coverage.txt
```

## Build

### Build Binary
```bash
make build
```

Or:
```bash
go build -o bin/customer-service ./api
```

### Build Docker Image
```bash
make docker-build
```

Or:
```bash
docker build -t customer-service .
```

## Development

### Project Structure

- **api**: Application entry point and main function
- **internal/domain**: Business entities (Customer, CPF, Email value objects)
- **internal/usecase**: Business logic (Create, Update, Delete, GetByCPF)
- **internal/repository**: Data access layer with MongoDB implementation
- **internal/handler**: HTTP handlers and routing
- **pkg/validator**: Reusable validation functions
- **pkg/errors**: Custom error types

### Adding New Features

1. Define domain entities in `internal/domain`
2. Create use cases in `internal/usecase`
3. Implement repository methods in `internal/repository`
4. Add HTTP handlers in `internal/handler`
5. Update routes in `internal/handler/routes.go`
6. Write tests for all layers

## Migration from NestJS

This service maintains 100% API compatibility with the original NestJS customer module:

- Same endpoints: `POST /customer`, `GET /customer/:cpf`, `PATCH /customer/:id`, `DELETE /customer/:id`
- Same request/response formats
- Same validation rules (CPF, Email)
- Same error codes and messages
- Same business logic

### Differences

- **Database**: PostgreSQL → MongoDB
- **Language**: TypeScript → Go
- **Framework**: NestJS → Gin
- **Architecture**: Clean Architecture maintained in both

## Docker Compose Services

The `docker-compose.yml` includes:

- **mongodb**: MongoDB 7.0 database
- **customer-service**: The Go microservice

Both services are connected via the `customer-network` bridge network.

## Makefile Commands

```bash
make help              # Show available commands
make build             # Build the application
make run               # Run the application locally
make test              # Run all tests
make test-unit         # Run unit tests only
make test-integration  # Run integration tests
make clean             # Clean build artifacts
make docker-build      # Build Docker image
make docker-up         # Start docker-compose services
make docker-down       # Stop docker-compose services
make docker-logs       # View docker-compose logs
```

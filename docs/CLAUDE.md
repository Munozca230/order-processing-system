# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a microservices-based order processing system with the following components:

- **Order Worker** (Java 21 + Spring Boot WebFlux): Kafka consumer that processes orders by enriching them with product and customer data, then persisting to MongoDB
- **Product API** (Go): RESTful service for product information
- **Customer API** (Go): RESTful service for customer information
- **Infrastructure**: Kafka for messaging, MongoDB for persistence, Redis for distributed locking and retry management

The system implements an event-driven architecture with reactive programming patterns, focusing on resilience, scalability, and observability.

## Development Workflows

### 🧪 Testing (Development)
```bash
# Run all tests (uses Testcontainers - independent of Docker Compose)
cd order-worker && mvn test

# Run specific test class
cd order-worker && mvn test -Dtest=OrderKafkaConsumerTest

# Run integration tests (requires Docker for Testcontainers)
cd order-worker && mvn test -Dtest=OrderIntegrationTest

# Clean and run all tests
cd order-worker && mvn clean test
```

### 🚀 Production-like Environment (Docker Compose)
```bash
# Start complete system (recommended for full validation)
cd infra && docker-compose up -d

# View application logs
cd infra && docker-compose logs -f order-worker

# View all services status
cd infra && docker-compose ps

# Stop all services
cd infra && docker-compose down

# Clean restart (removes volumes/data)
cd infra && docker-compose down -v && docker-compose up -d
```

### ⚡ Local Development (Optional)
```bash
# Start only infrastructure services
cd infra && docker-compose up -d zookeeper kafka mongo redis

# Run application locally for debugging
cd order-worker && mvn spring-boot:run

# Build the project
cd order-worker && mvn clean compile

# Package the application
cd order-worker && mvn package
```

### Go APIs
```bash
# Build Product API
cd product-api && go build -o product-api .

# Build Customer API
cd customer-api && go build -o customer-api .

# Run tests
cd product-api && go test ./...
cd customer-api && go test ./...
```

### Testing vs Production Commands Summary

**For Testing & Development:**
```bash
cd order-worker && mvn test  # Fast, isolated, uses Testcontainers
```

**For Production-like Validation:**
```bash
cd infra && docker-compose up -d  # Complete system, real networking
```

## Project Structure Insights

### Order Worker Architecture
- **Consumer**: `OrderKafkaConsumer` handles Kafka message consumption
- **Services**: 
  - `EnrichmentService`: Calls external APIs to enrich orders
  - `ValidationService`: Validates enriched order data
  - `OrderLockService`: Manages distributed locking with Redis
- **Repository**: `OrderRepository` for MongoDB operations
- **Models**: Domain objects for orders, customers, and products
- **Events**: Application events for processing notifications

### Testing Strategy
- **Unit Tests**: Mock-based testing for individual components
- **Integration Tests**: Full end-to-end testing with Testcontainers
- **Test Infrastructure**: Uses WireMock for external API mocking, Testcontainers for Kafka/MongoDB/Redis

### Key Configuration
- Application config in `order-worker/src/main/resources/application.yml`
- External API URLs configurable via properties
- Kafka topics and MongoDB collections are auto-created

## Development Notes

### Working with the Order Worker
- Uses reactive programming with Project Reactor
- Implements distributed locking to prevent duplicate processing
- Handles retries and error scenarios gracefully
- All tests are designed to work with `@Testcontainers(disabledWithoutDocker = true)`

### Integration Testing
- Integration tests spin up real Kafka, MongoDB, and Redis instances
- External APIs are mocked with WireMock
- Tests verify complete end-to-end message processing

### Common Development Patterns
- Services use reactive `Mono` and `Flux` types
- Repository operations return reactive types
- Event-driven architecture with Spring Application Events
- Distributed locking pattern for preventing duplicate processing

## Environment Setup

### Prerequisites
- **Java 21**: Required for the Order Worker (Spring Boot application)
- **Go 1.22**: Required for Product API and Customer API
- **Docker**: Required for integration tests and infrastructure services
- **Maven**: For Java project management (usually bundled with Java installations)

### Java Installation
The project requires Java 21. If you encounter `JAVA_HOME` not found errors:

```bash
# On Ubuntu/Debian
sudo apt update
sudo apt install -y openjdk-21-jdk

# Set JAVA_HOME (add to ~/.bashrc or ~/.zshrc)
export JAVA_HOME=/usr/lib/jvm/java-21-openjdk-amd64
export PATH=$JAVA_HOME/bin:$PATH

# Verify installation
java -version
mvn -version
```

### Troubleshooting
- If `mvn test` fails with "JAVA_HOME not defined", ensure Java 21 is installed and JAVA_HOME is set
- Integration tests will be skipped if Docker is not available (`@Testcontainers(disabledWithoutDocker = true)`)
- External API calls in tests are mocked with WireMock to avoid network dependencies
- **Health Checks**: All services use optimized health checks with 5s intervals. The order-worker uses process-based health checking (`pgrep java`) for reliability
- If services show "Started" but not "Healthy", wait ~30 seconds for health checks to complete. All services should reach "Healthy" status for proper dependency management
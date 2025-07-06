# Cross-platform Order Processing System
# Works on Linux, macOS, Windows (with make installed)
# Eliminates PowerShell dependency - Native technology testing

.PHONY: help frontend backend build test test-unit test-integration test-e2e test-go clean status logs docker-clean

# Default target
help:
	@echo "🚀 Order Processing System - Cross-Platform Commands"
	@echo ""
	@echo "📋 Deployment:"
	@echo "  frontend     - Deploy full stack with web interface"
	@echo "  backend      - Deploy backend-only (APIs + Worker)"
	@echo "  build        - Build all services without starting"
	@echo ""
	@echo "🧪 Testing (Native Technologies):"
	@echo "  test         - Run all tests (unit + integration + e2e)"
	@echo "  test-unit    - Run unit tests (Go + Java)"
	@echo "  test-integration - Run integration tests with Testcontainers"
	@echo "  test-e2e     - Run end-to-end system tests"
	@echo "  test-go      - Run Go API tests (containerized)"
	@echo "  test-go-native - Run Go tests natively (requires Go installed)"
	@echo ""
	@echo "🔧 Operations:"
	@echo "  status       - Check service health"
	@echo "  logs         - View order-worker logs"
	@echo "  clean        - Clean restart with fresh data"
	@echo "  reset        - Quick reset (preserve some data)"
	@echo "  docker-clean - Fix Docker network issues"
	@echo ""
	@echo "🛠️ Prerequisites:"
	@echo "  - Docker & Docker Compose installed"
	@echo "  - Java 21 (for unit tests)"
	@echo "  - Go 1.22 (for Go tests)"
	@echo "  - make command available"
	@echo ""
	@echo "📖 Examples:"
	@echo "  make frontend    # Full deployment"
	@echo "  make test        # All tests"
	@echo "  make test-unit   # Fast unit tests only"

# Build all services
build:
	@echo "🔨 Building all services..."
	@cd infra && docker-compose build
	@echo "✅ Build complete!"

# Deploy frontend (full stack)
frontend:
	@echo "🌐 Deploying Frontend Profile (Full Stack)..."
	@echo "📋 Services: Backend + Order API + Nginx Frontend"
	@echo "🧹 Cleaning Docker networks..."
	@docker network prune -f 2>/dev/null || true
	@cd infra && docker-compose down -v --remove-orphans 2>/dev/null || true
	@cd infra && docker-compose --profile frontend up -d
	@echo "⏳ Waiting for services to be healthy..."
	@$(MAKE) wait-healthy
	@echo "✅ Frontend deployment complete!"
	@echo "🌐 Access: http://localhost:8080"

# Deploy backend only
backend:
	@echo "⚙️ Deploying Backend Profile (APIs + Worker)..."
	@echo "📋 Services: Kafka + MongoDB + Redis + APIs + Worker"
	@echo "🧹 Cleaning Docker networks..."
	@docker network prune -f 2>/dev/null || true
	@cd infra && docker-compose down -v --remove-orphans 2>/dev/null || true
	@cd infra && docker-compose up -d
	@echo "⏳ Waiting for services to be healthy..."
	@$(MAKE) wait-healthy
	@echo "✅ Backend deployment complete!"
	@echo "📊 APIs: Product(8081) Customer(8082)"

# Wait for services to be healthy (internal helper)
wait-healthy:
	@for i in $$(seq 1 12); do \
		HEALTHY=$$(cd infra && docker-compose ps --format "table {{.Status}}" | grep -c "healthy"); \
		TOTAL=$$(cd infra && docker-compose ps --format "table {{.Name}}" | grep -v kafka-setup | tail -n +2 | wc -l | tr -d '\n'); \
		echo "🔍 Health Check: $$HEALTHY/$$TOTAL services healthy (attempt $$i/12)"; \
		if [ "$$TOTAL" -gt "0" ] && [ "$$HEALTHY" -eq "$$TOTAL" ]; then \
			echo "🎉 All services are healthy!"; \
			break; \
		fi; \
		sleep 5; \
	done

# ==========================================
# TESTING TARGETS - NATIVE TECHNOLOGIES
# ==========================================

# Run all tests
test: test-unit test-integration test-e2e
	@echo "🎉 All tests completed!"

# Unit tests - Native technology runners
test-unit:
	@echo "🧪 Running Unit Tests (Native Technologies)..."
	@echo ""
	@echo "📊 Java Tests (Maven + JUnit):"
	@cd services/order-worker && mvn test -q
	@echo "✅ Java unit tests passed!"
	@echo ""
	@echo "📊 Go Tests (Go native test runner):"
	@$(MAKE) test-go
	@echo "✅ All unit tests completed!"

# Go tests (inside containers for cross-platform compatibility)
test-go:
	@echo "🐹 Running Go API Tests (inside containers)..."
	@echo "📊 Testing Product API..."
	@cd services/product-api && docker build -t product-api-test --target builder .
	@docker run --rm product-api-test go test ./... -v
	@echo "📊 Testing Customer API..."
	@cd services/customer-api && docker build -t customer-api-test --target builder .
	@docker run --rm customer-api-test go test ./... -v
	@echo "✅ Go tests passed!"

# Go tests (native - requires Go installed on host)
test-go-native:
	@echo "🐹 Running Go API Tests (native)..."
	@echo "⚠️  Requires Go 1.22+ installed on host"
	@cd services/product-api && go test ./... -v
	@cd services/customer-api && go test ./... -v
	@echo "✅ Go tests passed!"

# Integration tests with Testcontainers
test-integration:
	@echo "🔗 Running Integration Tests (Testcontainers)..."
	@echo "📊 These tests spin up real infrastructure..."
	@cd services/order-worker && mvn test -Dtest=*IntegrationTest* -q
	@echo "✅ Integration tests passed!"

# End-to-end system tests
test-e2e:
	@echo "🌐 Running End-to-End Tests..."
	@echo "📋 Starting fresh system..."
	@cd infra && docker-compose down -v 2>/dev/null || true
	@cd infra && docker-compose up -d
	@$(MAKE) wait-healthy
	@echo ""
	@echo "📊 Testing Data Initialization..."
	@cd infra && docker-compose exec mongo mongosh catalog --eval "print('Products:', db.products.countDocuments()); print('Customers:', db.customers.countDocuments());" || echo "❌ MongoDB data test failed"
	@echo ""
	@echo "🔍 Testing API Health..."
	@curl -sf http://localhost:8081/health > /dev/null && echo "✅ Product API healthy" || echo "❌ Product API failed"
	@curl -sf http://localhost:8082/health > /dev/null && echo "✅ Customer API healthy" || echo "❌ Customer API failed"
	@echo ""
	@echo "📦 Testing Product API..."
	@curl -sf http://localhost:8081/products > /dev/null && echo "✅ Product API responding" || echo "❌ Product API failed"
	@echo ""
	@echo "👥 Testing Customer API..."
	@curl -sf http://localhost:8082/customers > /dev/null && echo "✅ Customer API responding" || echo "❌ Customer API failed"
	@echo ""
	@echo "📨 Testing Order Processing Pipeline..."
	@cd infra && echo '{"orderId":"e2e-test-001","customerId":"customer-premium","products":[{"productId":"product-1"}]}' | \
		docker-compose exec -T kafka kafka-console-producer.sh --bootstrap-server kafka:9092 --topic orders
	@echo "⏳ Waiting for order processing..."
	@sleep 5
	@cd infra && docker-compose exec mongo mongosh orders --eval "print('Processed Orders:', db.orders.countDocuments());" || echo "❌ Order processing test failed"
	@echo "✅ End-to-end tests completed!"

# ==========================================
# OPERATIONS
# ==========================================

# Check service status
status:
	@echo "📊 Service Status:"
	@cd infra && docker-compose ps

# View logs
logs:
	@echo "📝 Order Worker Logs (last 50 lines):"
	@cd infra && docker-compose logs --tail=50 order-worker

# Clean restart
clean:
	@echo "🧹 Clean restart with fresh data..."
	@echo "🧹 Cleaning Docker resources..."
	@docker network prune -f 2>/dev/null || true
	@docker volume prune -f 2>/dev/null || true
	@cd infra && docker-compose down -v --remove-orphans 2>/dev/null || true
	@cd infra && docker-compose up -d
	@$(MAKE) wait-healthy
	@echo "✅ Clean restart complete!"

# Quick reset
reset:
	@echo "🔄 Quick reset..."
	@cd infra && docker-compose restart
	@$(MAKE) wait-healthy
	@echo "✅ Reset complete!"

# ==========================================
# ADVANCED TESTING TARGETS
# ==========================================

# Run specific test class (Java)
test-java-class:
	@if [ -z "$(CLASS)" ]; then \
		echo "❌ Usage: make test-java-class CLASS=OrderKafkaConsumerTest"; \
		exit 1; \
	fi
	@echo "🧪 Running Java test class: $(CLASS)"
	@cd services/order-worker && mvn test -Dtest=$(CLASS) -q

# Run tests with coverage
test-coverage:
	@echo "📊 Running tests with coverage..."
	@cd services/order-worker && mvn test jacoco:report -q
	@echo "✅ Java coverage report: services/order-worker/target/site/jacoco/index.html"
	@cd services/product-api && go test -coverprofile=coverage.out ./...
	@cd services/product-api && go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Go coverage report: services/product-api/coverage.html"

# Performance test
test-performance:
	@echo "⚡ Running basic performance tests..."
	@$(MAKE) backend
	@echo "📊 Testing API response times..."
	@for i in $$(seq 1 100); do \
		curl -s -w "%{time_total}\n" -o /dev/null http://localhost:8081/products; \
	done | awk '{sum+=$$1; count++} END {print "Average response time:", sum/count "s"}'

# Memory test
test-memory:
	@echo "💾 Running memory tests..."
	@$(MAKE) backend
	@echo "📊 Checking container memory usage..."
	@cd infra && docker stats --no-stream --format "table {{.Container}}\t{{.MemUsage}}\t{{.MemPerc}}"

# Fix Docker network issues
docker-clean:
	@echo "🧹 Fixing Docker network issues..."
	@echo "📋 Stopping all containers..."
	@cd infra && docker-compose down -v --remove-orphans 2>/dev/null || true
	@echo "📋 Cleaning orphaned networks..."
	@docker network prune -f 2>/dev/null || true
	@echo "📋 Cleaning orphaned volumes..."
	@docker volume prune -f 2>/dev/null || true
	@echo "📋 Cleaning system resources..."
	@docker system prune -f 2>/dev/null || true
	@echo "✅ Docker cleanup complete! Try 'make frontend' again."
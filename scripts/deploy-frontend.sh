#!/bin/bash
# Cross-platform frontend deployment script
# Works on Linux, macOS, Windows (Git Bash/WSL)

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${CYAN}[STEP]${NC} $1"
}

# Help function
show_help() {
    echo "Frontend Profile Deployment - Cross-Platform"
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --preserve-data    Preserve existing data volumes"
    echo "  --help            Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                 # Clean deployment (default)"
    echo "  $0 --preserve-data # Preserve existing data"
    exit 0
}

# Parse arguments
PRESERVE_DATA=false
while [[ $# -gt 0 ]]; do
    case $1 in
        --preserve-data)
            PRESERVE_DATA=true
            shift
            ;;
        --help)
            show_help
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            ;;
    esac
done

print_status "Frontend Profile Deployment (Full Stack + Web Interface)"
echo "============================================================="

if [ "$PRESERVE_DATA" = true ]; then
    print_warning "MODE: Preserving existing data volumes"
else
    print_status "MODE: Clean deployment (volumes will be reset)"
fi

# Check if Docker is running
print_step "Checking Docker status..."
if ! docker version > /dev/null 2>&1; then
    print_error "Docker is not running or not accessible"
    echo ""
    echo "Please ensure Docker is running:"
    echo "  1. Start Docker Desktop application"
    echo "  2. Wait for the green 'Engine running' status"
    echo "  3. Verify with: docker --version"
    echo ""
    echo "If Docker is not installed:"
    echo "  Download from: https://www.docker.com/products/docker-desktop"
    exit 1
fi

DOCKER_VERSION=$(docker version --format '{{.Client.Version}}' 2>/dev/null)
print_status "Docker is running (Version: $DOCKER_VERSION)"

# Navigate to infra directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
INFRA_PATH="$PROJECT_ROOT/infra"

if [ ! -d "$INFRA_PATH" ]; then
    print_error "Infrastructure directory not found: $INFRA_PATH"
    exit 1
fi

cd "$INFRA_PATH"
print_status "Working directory: $(pwd)"

# Stop existing services
if [ "$PRESERVE_DATA" = true ]; then
    print_step "Stopping existing services (preserving data)..."
    docker-compose --profile frontend --profile backend down || print_warning "No existing services to stop"
    print_status "Existing services stopped (data preserved)"
else
    print_step "Stopping existing services and cleaning volumes..."
    docker-compose --profile frontend --profile backend down -v || print_warning "No existing services to stop"
    print_status "Existing services stopped and volumes cleaned"
    
    print_step "Cleaning orphaned Docker volumes..."
    docker volume prune -f || print_warning "No orphaned volumes to clean"
    print_status "Orphaned volumes cleaned"
fi

# Start frontend profile
print_step "Starting Frontend Profile services..."
print_status "Backend Services:"
echo "  - Zookeeper (Kafka coordination)"
echo "  - Kafka (Message broker)"
echo "  - MongoDB (Database + sample data)"
echo "  - Redis (Cache & distributed locks)"
echo "  - Order Worker (Java - Core processing)"
echo "  - Product API (Go - Product catalog)"
echo "  - Customer API (Go - Customer management)"
echo ""
print_status "Frontend Services:"
echo "  - Order API (Node.js - Frontend bridge to Kafka)"
echo "  - Nginx Frontend (Web server + proxy)"

if ! docker-compose --profile frontend up -d; then
    print_error "Failed to start frontend services"
    exit 1
fi

print_status "Frontend profile services started successfully!"

# Wait for services to be ready
print_step "Waiting for services to be healthy..."
MAX_WAIT=60
WAITED=0
CHECK_INTERVAL=5

while [ $WAITED -lt $MAX_WAIT ]; do
    sleep $CHECK_INTERVAL
    WAITED=$((WAITED + CHECK_INTERVAL))
    
    # Count healthy services
    HEALTHY_COUNT=$(docker-compose ps --format "table {{.Status}}" | grep -c "healthy" || echo "0")
    TOTAL_SERVICES=$(docker-compose ps --format "table {{.Name}}" | grep -v kafka-setup | tail -n +2 | wc -l)
    
    if [ "$HEALTHY_COUNT" -eq "$TOTAL_SERVICES" ]; then
        print_status "All services are healthy!"
        break
    else
        print_warning "Healthy: $HEALTHY_COUNT/$TOTAL_SERVICES (waited ${WAITED}s)"
    fi
done

if [ $WAITED -ge $MAX_WAIT ]; then
    print_warning "Timeout reached. Some services may still be starting."
fi

# Check service status
print_step "Checking service status..."
docker-compose ps

# Health checks
print_step "Health checks..."
services=(
    "Frontend Web:http://localhost:8080/health"
    "Order API:http://localhost:3000/health" 
    "Product API:http://localhost:8081/health"
    "Customer API:http://localhost:8082/health"
)

for service in "${services[@]}"; do
    name="${service%%:*}"
    url="${service##*:}"
    
    if curl -s -f "$url" > /dev/null 2>&1; then
        print_status "$name is healthy"
    else
        print_warning "$name is not responding"
    fi
done

# Display useful information
print_status "Frontend Profile Ready!"
echo "========================"
echo ""
print_status "Primary Access Point:"
echo "  Frontend Web Interface: http://localhost:8080"
echo "     Complete visual interface for order testing"
echo ""
print_status "API Endpoints:"
echo "  Order API: http://localhost:3000"
echo "  Product API: http://localhost:8081"
echo "  Customer API: http://localhost:8082"
echo ""
print_status "Database Access:"
echo "  MongoDB: mongodb://localhost:27017"
echo "  Redis: redis://localhost:6379"
echo "  Kafka: kafka://localhost:9092"
echo ""
print_status "Testing Options:"
echo "  Web Interface: Open http://localhost:8080 in browser"
echo "  Postman Collection: Import postman/*.json files"
echo "  Direct API Test:"
echo '    curl -X POST http://localhost:3000/api/orders -H "Content-Type: application/json" -d'"'"'{"orderId":"test","customerId":"customer-1","products":[{"productId":"product-1"}]}'"'"
echo ""
print_status "Monitoring:"
echo "  Service Status: docker-compose ps"
echo "  Order Worker Logs: docker-compose logs -f order-worker"
echo "  Frontend Logs: docker-compose logs -f nginx-frontend"
echo ""
print_status "To switch to Backend only:"
echo "  docker-compose down && docker-compose up -d"
echo ""
print_status "To stop all services:"
echo "  docker-compose --profile frontend down"

print_status "Frontend Profile deployment completed!"
print_status "Access your order processing system at: http://localhost:8080"
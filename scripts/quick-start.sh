#!/bin/bash

# MCP Core Inventory - Quick Start Script
# Version: 1.2.0
# Purpose: Fast deployment for testing/development

set -euo pipefail

# Configuration
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly PROJECT_ROOT="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Colors
readonly RED='\033[0;31m'
readonly GREEN='\033[0;32m'
readonly YELLOW='\033[1;33m'
readonly BLUE='\033[0;34m'
readonly NC='\033[0m'

log() {
    echo -e "${BLUE}[$(date '+%H:%M:%S')] $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Quick deployment using local binary
quick_start() {
    log_header "Quick Start - MCP Core Inventory v1.2.0"
    
    # Create environment
    if [ ! -f "$PROJECT_ROOT/.env" ]; then
        log_warning "Creating default .env file"
        cp "$PROJECT_ROOT/env.example" "$PROJECT_ROOT/.env"
    fi
    
    # Start required services (PostgreSQL, Redis, NATS)
    log "Starting core services..."
    
    # PostgreSQL (if not running)
    if ! docker ps --filter "name=postgres" --format "{{.Names}}" | grep -q postgres; then
        log "Starting PostgreSQL..."
        docker run -d \
            --name mcp-core-postgres \
            -e POSTGRES_DB=mcp_core_inventory \
            -e POSTGRES_USER=postgres \
            -e POSTGRES_PASSWORD=postgres \
            -p 5432:5432 \
            -v postgres_data:/var/lib/postgresql/data \
            postgres:15-alpine
        log_success "PostgreSQL started"
    else
        log_success "PostgreSQL already running"
    fi
    
    # Redis (if not running)
    if ! docker ps --filter "name=redis" --format "{{.Names}}" | grep -q redis; then
        log "Starting Redis..."
        docker run -d \
            --name mcp-core-redis \
            -p 6379:6379 \
            -v redis_data:/data \
            redis:7-alpine
        log_success "Redis started"
    else
        log_success "Redis already running"
    fi
    
    # NATS (if not running)
    if ! docker ps --filter "name=nats" --format "{{.Names}}" | grep -q nats; then
        log "Starting NATS..."
        docker run -d \
            --name mcp-core-nats \
            -p 4222:4222 \
            -p 8222:8222 \
            -v nats_data:/data \
            nats:2.10-alpine \
            -js
        log_success "NATS started"
    else
        log_success "NATS already running"
    fi
    
    # Wait for services to be ready
    log "Waiting for services to be ready..."
    sleep 10
    
    # Start application locally
    log "Starting MCP Core Inventory..."
    
    # Set environment variables
    export DATABASE_URL="postgres://postgres:postgres@localhost:5432/mcp_core_inventory?sslmode=disable"
    export REDIS_URL="redis://localhost:6379"
    export NATS_URL="nats://localhost:4222"
    export SERVER_HOST="0.0.0.0"
    export SERVER_PORT="8080"
    export LOGGING_LEVEL="info"
    export JWT_SECRET="quick-start-secret-key"
    
    # Check if binary exists and is executable
    if [ -f "$PROJECT_ROOT/bin/mcp-core-inventory" ]; then
        log "Using compiled binary..."
        cd "$PROJECT_ROOT"
        
        # Try to run binary
        if ./bin/mcp-core-inventory --version &> /dev/null; then
            log_success "Starting application..."
            ./bin/mcp-core-inventory &
            APP_PID=$!
            log_success "Application started (PID: $APP_PID)"
            
            # Wait for application to be ready
            sleep 5
            
            # Test health endpoint
            if curl -s http://localhost:8080/health &> /dev/null; then
                log_success "Application is healthy!"
                log "API available at: http://localhost:8080"
                log "Health check: http://localhost:8080/health"
                log "OpenAPI docs: http://localhost:8080/v1/docs"
                
                # Keep the script running
                wait $APP_PID
            else
                log_warning "Application may not be ready yet"
                log "Check logs for details"
                wait $APP_PID
            fi
        else
            log_error "Binary is not compatible with this OS"
            fallback_to_go_run
        fi
    else
        log_warning "Binary not found, building..."
        fallback_to_go_run
    fi
}

# Fallback to go run
fallback_to_go_run() {
    log "Fallback: Building with Go..."
    cd "$PROJECT_ROOT"
    
    # Check if Go is available
    if command -v go &> /dev/null; then
        log "Building and running with go run..."
        go run ./cmd/main.go &
        APP_PID=$!
        log_success "Application started (PID: $APP_PID)"
        
        # Wait for application
        wait $APP_PID
    else
        log_error "Go is not installed and binary is not available"
        log "Please install Go or compile the application first"
        exit 1
    fi
}

# Health check function
health_check() {
    log "Running health checks..."
    
    local services=(
        "PostgreSQL:5432:psql postgres -h localhost -c 'SELECT 1'"
        "Redis:6379:redis-cli ping"
        "NATS:4222:curl -s http://localhost:8222/varz"
        "Application:8080:curl -s http://localhost:8080/health"
    )
    
    for service in "${services[@]}"; do
        local name=$(echo "$service" | cut -d: -f1)
        local port=$(echo "$service" | cut -d: -f2)
        local command=$(echo "$service" | cut -d: -f3)
        
        log "Checking $name (port $port)..."
        if eval "$command" &> /dev/null; then
            log_success "$name: OK"
        else
            log_warning "$name: Failed"
        fi
    done
}

# Cleanup function
cleanup() {
    log "Cleaning up..."
    
    # Stop containers
    docker stop mcp-core-postgres mcp-core-redis mcp-core-nats 2>/dev/null || true
    docker rm mcp-core-postgres mcp-core-redis mcp-core-nats 2>/dev/null || true
    
    # Remove volumes
    docker volume rm postgres_data redis_data nats_data 2>/dev/null || true
    
    log_success "Cleanup completed"
}

# Help function
help() {
    cat << EOF
MCP Core Inventory - Quick Start Script v1.2.0

Usage:
  $0 [COMMAND]

Commands:
  start       Start all services and the application
  health      Run health checks on all services
  cleanup     Stop and remove all containers and volumes
  help        Show this help message

Examples:
  $0 start      # Start everything
  $0 health     # Check service health
  $0 cleanup    # Clean up all resources

Environment:
  - .env file will be created if not exists
  - Default database: postgres://postgres:postgres@localhost:5432/mcp_core_inventory
  - Default Redis: redis://localhost:6379
  - Default NATS: nats://localhost:4222
  - Default API: http://localhost:8080

EOF
}

log_header() {
    echo
    echo -e "${BLUE}═════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE} $1${NC}"
    echo -e "${BLUE}═════════════════════════════════════════════════════════════════${NC}"
    echo
}

# Main execution
case "${1:-start}" in
    start)
        quick_start
        ;;
    health)
        health_check
        ;;
    cleanup)
        cleanup
        ;;
    help|--help|-h)
        help
        ;;
    *)
        log_error "Unknown command: $1"
        help
        exit 1
        ;;
esac
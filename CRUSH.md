# CRUSH - mcp-core-inventory Development Guide

This document contains essential information for agents working with the mcp-core-inventory codebase.

## Project Overview

mcp-core-inventory is a comprehensive Model Context Protocol (MCP) template and generation engine. It's a Go-based system that follows Clean Architecture principles. The project has been renamed from MCP-HULK and provides a core inventory management system with MCP protocol support, multiple service entry points, and extensive tooling for generation, validation, and deployment.

## Development Commands

### Build & Run
```bash
# Build the application
make build

# Run with default configuration (main HTTP server)
make run

# Install dependencies
make deps

# Build for production (multi-platform)
make build-prod

# Development setup (installs golangci-lint, mockgen, air)
make dev-setup

# Hot reload for development
make dev
```

### Testing
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Test-specific commands
make test-unit         # Domain and app layer tests
make test-integration  # Adapter and integration tests
make test-e2e          # End-to-end tests
make test-concurrency  # Race condition tests

# Run all checks (lint, vet, test)
make check
```

### Security & Quality
```bash
# Security scanning
make security-scan

# Code quality
make quality

# Benchmarking
make benchmark

# Profiling
make profile-cpu
make profile-mem
```

### CI/CD Pipeline
```bash
# Full CI pipeline
make ci

# Quality gate
make quality-gate

# Configuration validation
make validate-config

# CI/CD tools setup
make ci-setup

# PERFORMANCE TESTING (NEW)
./scripts/load-testing/black-friday-certification.sh

# Performance automation via GitHub Actions
# Triggers: .github/workflows/performance-testing.yml
```

### Docker
```bash
# Build Docker image
make docker

# Run Docker container
make docker-run

# Deploy with Docker
make deploy
make deploy-push
```

### Code Generation
```bash
# Generate code (go generate)
make generate

# Install locally
make install

# Generate documentation
make docs
```

### Health & Operations
```bash
# Health checks
make health           # Full health check
make health-quick     # Quick service check

# SLO monitoring
make slos

# Load testing
make load-test        # Black Friday scenario
make stress-test      # High-load stress test

# Chaos engineering
make chaos

# Backup and recovery
make backup
make rollback
```

## Project Structure

The project follows a strict Clean Architecture pattern:

```
├── cmd/                    # Application entry points
│   ├── main.go             # Main HTTP server
│   ├── mcp-cli/            # MCP CLI tool
│   ├── mcp-server/         # MCP server
│   ├── core-inventory/     # Core inventory service
│   ├── tools-validator/    # Validation tools
│   ├── tools-generator/    # Generation tools
│   ├── tools-deployer/     # Deployment tools
│   ├── thor/               # Development utility
│   └── mcp-init/           # MCP project initializer
├── internal/               # Private application code
│   ├── core/               # Performance engine, config, scheduler
│   ├── ai/                 # AI & knowledge management
│   ├── state/              # State management with event sourcing
│   ├── observability/      # Observability and health checks
│   ├── versioning/         # Version control for data/models/knowledge
│   ├── mcp/                # MCP protocol implementation & generators
│   ├── services/           # Application services layer
│   ├── domain/             # Domain entities, repositories, value objects
│   ├── application/        # Use cases and DTOs
│   ├── infrastructure/     # External integrations (databases, APIs)
│   ├── interfaces/         # Input/output adapters (HTTP, gRPC, CLI)
│   ├── adapters/           # Additional adapters
│   ├── app/                # Application-level components
│   └── security/           # Authentication, authorization, encryption
├── pkg/                    # Public libraries
│   ├── glm/               # GLM AI client
│   ├── httpserver/        # HTTP server utilities
│   ├── validator/         # Validation utilities
│   ├── profiler/          # Profiling tools
│   ├── logger/            # Logging utilities
│   ├── knowledge/         # Knowledge management
│   ├── optimizer/         # Optimization tools
│   └── mcp/               # MCP utilities
├── templates/              # Code generation templates
│   ├── base/              # Base template
│   ├── go/                # Go templates
│   ├── wasm/              # WebAssembly templates
│   ├── tinygo/            # TinyGo templates
│   ├── mcp-go-premium/    # Premium MCP Go templates
│   ├── k8s/               # Kubernetes templates
│   ├── docker-compose/    # Docker Compose templates
│   └── web/               # Web templates
├── tools/                  # Development tools
│   ├── validators/        # Validation tools
│   ├── generators/        # Generation tools
│   ├── deployers/         # Deployment tools
│   ├── converters/        # Conversion tools
│   └── analyzers/         # Analysis tools
├── config/                 # Configuration files
├── scripts/                # Automation scripts
├── docs/                   # Documentation
├── contracts/              # API contracts and schemas
└── ops/                    # Operations files (k6, prometheus, etc.)
```

## Key Architectural Patterns

### Clean Architecture
- **Domain**: Core business logic (entities, value objects)
- **Application**: Use cases and application-specific rules
- **Infrastructure**: External concerns (databases, APIs)
- **Interfaces**: Adapters for input/output (HTTP, gRPC, CLI)

### Multi-Service Architecture
The system provides multiple entry points:
- **main.go**: Primary HTTP server with full feature set
- **mcp-server**: Dedicated MCP protocol server
- **core-inventory**: Specialized inventory management service
- **tools-***: Specialized tooling for validation, generation, deployment
- **thor**: Development utility and helper tools

### MCP Protocol Implementation
The system implements the Model Context Protocol with:
- JSON-RPC 2.0 messaging via NATS JetStream
- Tool registration and execution engine
- Resource management and caching
- Prompt templating and generation

### Template Generation System
Support for multiple technology stacks:
- **Go**: Standard Go backend services with Clean Architecture
- **TinyGo**: Go compiled to WebAssembly for edge computing
- **Rust WASM**: Rust compiled to WebAssembly for performance
- **Web**: React/TypeScript frontend applications
- **Kubernetes**: Container orchestration templates
- **Docker Compose**: Local development environments

## Configuration

Configuration is managed through YAML files in `config/` directory:

### Primary Config
- `config/config.yaml` - Main application configuration
- Uses Viper for loading with environment variable overrides
- Prefix: `HULK_` for environment variables (not `MCP_HULK_`)

### Key Configuration Sections
- **Server**: HTTP server settings (port, timeouts)
- **Engine**: Worker pool and execution settings
- **NATS**: Message broker configuration
- **Logging**: Structured logging with Zap
- **Telemetry**: OpenTelemetry tracing and metrics
- **MCP**: Protocol-specific settings

### Default Values
The system provides sensible defaults in `internal/core/config/config.go`:
- Server port: 8080
- Workers: "auto" (NumCPU * 2)
- Queue size: 2000
- Cache L1 size: 5000 entries
- Database: PostgreSQL on localhost:5432
- AI Provider: GLM with glm-4.6-z.ai model
- NATS: localhost:4222

## Core Dependencies

### Go 1.23 Runtime
- **Echo v4**: HTTP server framework
- **NATS**: Message broker and JetStream streaming
- **Viper**: Configuration management with env var support
- **Cobra**: CLI framework for command-line tools
- **Zap**: Structured logging with JSON output

### Observability Stack
- **OpenTelemetry**: Distributed tracing with Jaeger
- **Prometheus**: Metrics collection and alerting
- **Custom metrics**: Performance and resource monitoring

### Data & Storage
- **Badger v4**: Embedded key-value store for caching
- **PostgreSQL**: Primary relational database
- **Multi-level caching**: L1 (memory), L2 (Redis), L3 (Badger)

### AI/ML Integration
- **Multiple LLM providers**: OpenAI, Gemini, GLM
- **Vector database support**: Qdrant, Pinecone, Weaviate
- **Knowledge graph capabilities**: Semantic search and retrieval

### Container & Cloud Native
- **Kubernetes**: Full k8s API support for deployment
- **Docker**: Multi-stage builds and optimization
- **Docker Compose**: Local development environments

## Naming Conventions

### Files
- **Go files**: `snake_case.go` describing purpose
  - Examples: `mcp_http_handler.go`, `model_registry.go`, `state_snapshot.go`
- **Handlers**: Suffix with type (`*_http_handler.go`, `*_grpc_server.go`)
- **Repositories**: `*_repository.go` for interfaces, `postgres_*_repository.go` for implementations
- **Scripts**: Category prefixed (`setup_*.sh`, `deploy_*.sh`, `generate_*.sh`)

### Directories
- Lowercase English names with underscores if needed
- Fixed top-level directories: `internal/`, `pkg/`, `templates/`, `scripts/`, `config/`, `docs/`

### Code Style
- Follow Go standard conventions
- Use Zap for structured logging
- Error wrapping with context
- Interface-based design in domain layer

## Critical Rules (from .cursor policy)

### Forbidden
- Creating files/directories not in the official tree structure
- Renaming files/directories without updating documentation
- Generic names like `utils.go`, `helpers.go` in critical layers
- Creating unapproved top-level directories
- Editing template files directly in `templates/` for specific cases

### Required
- All new artifacts must be documented in the official tree
- Standardized function comments for Go code
- Follow the Clean Architecture layer separation

## Testing

### Structure
- Unit tests: `*_test.go` alongside source files
- Integration tests: Organized by feature area
- Coverage reporting available via `make test-coverage`

### Key Test Patterns
- Table-driven tests for multiple scenarios
- Mock interfaces for infrastructure concerns
- Test fixtures in `testdata/` directories where needed

## Development Workflow

### Adding Features
1. Identify the appropriate layer (domain/application/infrastructure/interfaces)
2. Create domain entities and value objects first
3. Implement use cases in application layer
4. Add infrastructure implementations as needed
5. Create interface adapters (HTTP/gRPC/CLI)
6. Add tests at each layer
7. Run `make quality-gate` before committing

### Template Development
1. Create template in `templates/` with `manifest.yaml`
2. Implement generators in `internal/mcp/generators/`
3. Add validation logic in `internal/mcp/validators/`
4. Update template registry in MCP configuration
5. Test with `make test-integration`

### Service Development
1. Choose appropriate entry point in `cmd/`
2. Follow existing patterns in similar services
3. Implement core logic in `internal/` layers
4. Add configuration to `config/` if needed
5. Add tests and documentation
6. Use `make dev` for hot reload during development

### Configuration Changes
1. Update config structs in `internal/core/config/config.go`
2. Add defaults in `setDefaults()` method
3. Update `config/config.yaml` with new sections
4. Add validation rules where needed
5. Test with `make validate-config`

### Deployment Process
1. Build with `make build-prod`
2. Create Docker image with `make docker`
3. Test deployment with `make load-test`
4. Deploy with `make deploy`
5. Monitor with `make health` and `make slos`

## Performance Considerations

### Execution Engine
- Worker pools with configurable sizes
- Circuit breakers for resilience
- Multi-level caching (L1/L2/L3)
- NATS JetStream for reliable messaging

### Memory Management
- Custom memory optimizer in AI layer
- Episodic, semantic, and working memory types
- Memory consolidation processes

### Caching Strategy
- L1: In-memory cache for hot data
- L2: Distributed cache (Redis)  
- L3: Persistent cache (Badger)

### Security & Compliance

### Authentication & Authorization
- JWT token management with secure signing
- OAuth 2.0 provider integrations
- Role-based access control (RBAC)
- Session management with secure cookies

### Encryption & Security
- AES-256 encryption at rest and in transit
- Certificate management and rotation
- Secure key storage with hardware security modules
- Security scanning with `make security-scan`

### Compliance Features
- GDPR compliance features with data subject rights
- Comprehensive audit logging
- Configurable data retention policies
- Privacy-by-design architecture

## Observability Stack

### Metrics & Monitoring
- **Prometheus-compatible metrics**: Custom performance metrics and resource utilization tracking
- **OpenTelemetry distributed tracing**: Jaeger integration with request correlation across services
- **Structured JSON logging**: Multiple log levels with context-aware logging via Zap

### Health & Reliability
- **Health checks**: Built-in health endpoints and comprehensive health monitoring
- **SLO monitoring**: Service level objectives tracking and alerting
- **Circuit breakers**: Resilience patterns for distributed system failures

## Key Tools & Utilities

### Development Tools
- **Thor**: Development utility and helper tools (`cmd/thor/`)
- **Air**: Hot reload for development (`make dev`)
- **Mockgen**: Interface mocking for testing
- **Validate-tree**: Project structure validation tool

### Operations Tools
- **Load testing**: Black Friday scenario simulations
- **Stress testing**: High-load performance testing
- **Chaos engineering**: Resilience testing capabilities
- **Profiling**: CPU and memory profiling tools

## CI/CD Pipeline

### GitLab CI Integration
The project uses GitLab CI with comprehensive pipeline:
- **validate_tree**: Validates project structure against documented tree
- **build**: Builds the application
- **test**: Runs all tests with coverage reporting

### Pipeline Commands
```bash
# Full CI pipeline (equivalent to GitLab CI)
make ci

# Quality gate (pre-commit checks)
make quality-gate

# Local validation
make validate-config
```

### Validation Requirements
- Tree structure compliance: 95% minimum
- Zero missing files from documented structure
- All tests must pass
- Coverage thresholds enforced

### Security & Compliance
```bash
# Security vulnerability scanning
make security-scan

# Dependency checking
govulncheck ./...
nancy sleuth
```

## Critical Gotchas & Important Notes

### Project Structure Validation
- **Mandatory**: GitLab CI enforces 95% tree structure compliance
- **Zero tolerance**: No missing files from documented structure
- **Documentation**: All new artifacts must be documented in official tree

### Environment Variable Prefix
- **Correct**: `HULK_` prefix for environment variables
- **Incorrect**: `MCP_HULK_` (documented elsewhere but incorrect)

### Module Name
- **Current**: `github.com/vertikon/mcp-core-inventory`
- **Note**: Project renamed from MCP-HULK, some references may still use old name

### Multi-Service Architecture
- **Not a monolith**: Multiple entry points in `cmd/` directory
- **Service-specific**: Each tool has its own binary and purpose
- **No unified CLI**: Access tools via individual binaries, not single `hulk` command

### Testing Requirements
- **Concurrency testing**: Critical for ledger-like components
- **Quality gate**: Must pass before any commit
- **Coverage enforcement**: Minimum thresholds enforced in CI

### Performance & Scalability
- **Multi-level caching**: L1 (memory), L2 (Redis), L3 (Badger)
- **Worker pools**: Configurable, defaults to NumCPU * 2
- **Circuit breakers**: Essential for resilience
- **Load testing**: Black Friday scenarios built-in

## Migration from MCP-HULK

### Key Changes
- Module renamed from `mcp-hulk` to `mcp-core-inventory`
- Environment variable prefix changed to `HULK_`
- Multi-service architecture instead of monolithic CLI
- Enhanced CI/CD with tree structure validation
- Expanded template system and tooling

### Migration Checklist
- [ ] Update import paths in external dependencies
- [ ] Change environment variable prefixes
- [ ] Update CI/CD pipelines to use new structure
- [ ] Review service entry points for appropriate binary
- [ ] Validate tree structure compliance
- [ ] Update documentation and references

## Environment Configuration

### Environment Variables
- **Prefix**: `HULK_` (not `MCP_HULK_` as mentioned elsewhere)
- **Database**: `HULK_DATABASE_URL`, `HULK_DATABASE_PASSWORD`
- **AI**: `HULK_AI_API_KEY`
- **Module**: Currently uses `github.com/vertikon/mcp-core-inventory`

### Key File Locations
- **Main Config**: `config/config.yaml`
- **Makefile**: Comprehensive build/deploy commands
- **CLI Entry Points**: `cmd/` directory with multiple services

## Updated CLI Commands

The actual CLI is accessed via individual binaries, not a unified `hulk` command:

```bash
# Main HTTP server
make run                    # or: go run ./cmd/main.go

# MCP CLI tool
go run ./cmd/mcp-cli/main.go [commands]

# MCP Server
go run ./cmd/mcp-server/main.go

# Core Inventory Service
go run ./cmd/core-inventory/main.go

# Tools
go run ./cmd/tools-validator/main.go
go run ./cmd/tools-generator/main.go
go run ./cmd/tools-deployer/main.go

# Development tool (thor)
go run ./cmd/thor/main.go
```

## Testing

### Structure
- Unit tests: `*_test.go` alongside source files
- Integration tests: `./internal/adapters/...`, `./tests/integration/...`
- E2E tests: `./tests/e2e/...`
- Coverage reporting available via `make test-coverage` or `make coverage`

### Key Test Commands
```bash
# All tests
make test

# Specific test types
make test-unit         # Domain and app layer
make test-integration  # Adapters and integration
test-e2e              # End-to-end
make test-concurrency  # Race condition tests (important for ledger)

# Quality gates
make quality-gate      # Full test suite + security
```

## Troubleshooting

### Common Issues
- **NATS connection failed**: Check NATS server is running on localhost:4222
- **Configuration not found**: Ensure `config/config.yaml` exists or set HULK_* environment variables
- **Build failures**: Run `make deps` to ensure dependencies are current
- **Lint failures**: Install golangci-lint via `make dev-setup`
- **Tree validation failures**: Check `.cursor/` directory for documented structure requirements
- **Template generation errors**: Verify `templates/` directory structure and manifest.yaml files

### Debug Mode
Set `HULK_LOGGING_LEVEL=debug` or modify `config/config.yaml` for verbose logging output.

### Performance Issues
- Use `make profile-cpu` and `make profile-mem` to identify bottlenecks
- Check worker pool configuration in `config/config.yaml`
- Monitor NATS JetStream performance and queue sizes
- Verify cache hit rates and memory usage patterns

### Health & Monitoring
```bash
# Quick health checks
make health-quick     # Service health
make health           # Full system health
make slos             # SLO monitoring

# Logs and status
make logs             # Service logs
make status           # Service status
```

## Critical Artifacts Generated

### 1. OpenAPI Contract ✅ COMPLETE
- **Location**: `contracts/openapi/core-inventory-v1.yaml`
- **Purpose**: Complete API specification with business rules, error scenarios, examples
- **Features**: ACID operations, idempotency protection, race condition documentation

### 2. Load Testing Script ✅ COMPLETE
- **Location**: `ops/k6/black-friday-reservation-load-test.js`
- **Purpose**: Black Friday scenario testing for race conditions and performance
- **Capabilities**: 1000+ concurrent VUs, race detection, consistency verification

### 3. CI/CD Pipeline Integration ✅ COMPLETE
- **Location**: `.github/workflows/performance-testing.yml`
- **Features**: Automated load testing, multi-environment, performance gate checking
- **Triggers**: Push, PRs, scheduled daily, manual dispatch

### 4. Grafana Dashboard ✅ COMPLETE
- **Location**: `ops/grafana-dashboards/performance-dashboard.json`
- **Features**: Real-time race condition monitoring, stock consistency alerts, performance metrics
- **Alerts**: Critical (PagerDuty) for races, warnings for performance degradation

### 5. Load Testing Certification ✅ COMPLETE
- **Location**: `scripts/load-testing/black-friday-certification.sh`
- **Features**: Automated certification scenarios, performance gates, executive summary
- **Outcomes**: CERTIFIED/CONDITIONAL/REJECTED with pass/fail criteria

**Usage Examples:**
```bash
# Execute full Black Friday certification
./scripts/load-testing/black-friday-certification.sh

# Run automated performance test via CI/CD
gh workflow run performance-testing.yml -f test_type=black-friday

# Monitor in real-time with Grafana
# Import: ops/grafana-dashboards/performance-dashboard.json
```

---

*This CRUSH.md is actively maintained and reflects the current state of the mcp-core-inventory project, including the newly generated critical artifacts for production readiness.*
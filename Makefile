.PHONY: build test clean lint run deps docker deploy health load-test chaos rollback slos

# BLOCO-1 - Core Inventory Makefile
# Alinhado com o BLOCO-1 Blueprint - Core Inventory (Ledger ACID)

# Variáveis
SERVICE_NAME = mcp-core-inventory
VERSION ?= latest
REGISTRY = vertikon
DOCKER_COMPOSE_FILE = docker-compose.prod.yaml

# Build da aplicação
build:
	@echo "🔨 Build do BLOCO-1 Core Inventory..."
	go build -ldflags='-w -s -extldflags "-static"' -o bin/$(SERVICE_NAME) ./cmd
	@echo "✅ Build concluído: bin/$(SERVICE_NAME)"

# Build para produção (multi-plataforma)
build-prod:
	@echo "🏭 Build para produção (multi-plataforma)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -o bin/$(SERVICE_NAME)-linux ./cmd
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -o bin/$(SERVICE_NAME)-darwin ./cmd
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags='-w -s -extldflags "-static"' -o bin/$(SERVICE_NAME).exe ./cmd
	@echo "✅ Build produção concluído"

# Build Docker
build-docker:
	@echo "🐳 Build da imagem Docker..."
	docker build -t $(REGISTRY)/$(SERVICE_NAME):$(VERSION) .
	@echo "✅ Imagem Docker construída: $(REGISTRY)/$(SERVICE_NAME):$(VERSION)"

# Push da imagem Docker
push-docker:
	@echo "📤 Enviando imagem para registry..."
	docker push $(REGISTRY)/$(SERVICE_NAME):$(VERSION)
	@echo "✅ Imagem enviada para registry"

# Rodar testes
test:
	@echo "🧪 Executando testes do BLOCO-1..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Testes concluídos"

# Testes de unidade
test-unit:
	@echo "🧪 Testes de unidade..."
	go test -v -race -short ./internal/domain/ledger/...
	go test -v -race -short ./internal/app/...

# Testes de integração
test-integration:
	@echo "🧪 Testes de integração..."
	go test -v -race -tags=integration ./internal/adapters/... ./tests/integration/...

# Testes E2E
test-e2e:
	@echo "🧪 Testes E2E..."
	go test -v -race -tags=e2e ./tests/e2e/...

# Testes de concorrência (BLOCO-1 específico)
test-concurrency:
	@echo "🧪 Testes de concorrência (Race Conditions)..."
	go test -v -race -count=100 -parallel=10 ./internal/domain/ledger/...

# Verificar cobertura
coverage:
	@echo "📊 Verificando cobertura de testes..."
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Relatório de cobertura gerado: coverage.html"

# Limpar artefatos
clean:
	@echo "🧹 Limpando artefatos..."
	rm -rf bin/ coverage.out coverage.html
	docker system prune -f
	docker volume prune -f
	@echo "✅ Limpeza concluída"

# Lint
lint:
	@echo "🔍 Executando linters..."
	golangci-lint run
	go vet ./...
	go fmt ./...
	@echo "✅ Lint concluído"

# Instalar dependências
deps:
	@echo "📦 Instalando dependências..."
	go mod tidy
	go mod download
	go mod verify
	@echo "✅ Dependências instaladas"

# Gerar código
generate:
	@echo "⚙️ Gerando código..."
	go generate ./...
	@echo "✅ Código gerado"

# Rodar aplicação localmente
run:
	@echo "🚀 Rodando BLOCO-1 Core Inventory localmente..."
	go run ./cmd/main.go

# Deploy completo
deploy:
	@echo "🚀 Deploy do BLOCO-1..."
	./deploy.sh deploy $(VERSION)

# Deploy com push
deploy-push:
	@echo "🚀 Deploy com push do BLOCO-1..."
	./deploy.sh deploy $(VERSION) --push

# Health check
health:
	@echo "🏥 Verificando saúde do BLOCO-1..."
	./health-check.sh all

# Health check rápido
health-quick:
	@echo "🏥 Health check rápido..."
	./health-check.sh service

# Verificar SLOs
slos:
	@echo "📊 Verificando SLOs do BLOCO-1..."
	./health-check.sh slos

# Teste de carga (Black Friday)
load-test:
	@echo "🧪 Executando teste de carga (Black Friday)..."
	./deploy.sh load-test

# Teste de estresse
stress-test:
	@echo "🔥 Executando teste de estresse..."
	k6 run --vus 500 --duration 5m ops/k6/bloco-1-core/reserve-load-test.js

# Chaos Engineering
chaos:
	@echo "🌀 Executando testes de Chaos..."
	./deploy.sh chaos

# Rollback
rollback:
	@echo "🔄 Executando rollback..."
	./deploy.sh rollback

# Backup
backup:
	@echo "💾 Criando backup..."
	./deploy.sh backup

# Logs
logs:
	@echo "📋 Exibindo logs..."
	./deploy.sh logs

# Status dos serviços
status:
	@echo "📊 Status dos serviços..."
	./deploy.sh status

# Limpeza de recursos
cleanup:
	@echo "🧹 Limpando recursos..."
	./deploy.sh cleanup

# Setup de desenvolvimento
dev-setup:
	@echo "⚙️ Setup de desenvolvimento..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang/mock/mockgen@latest
	go install github.com/air-verse/air@latest
	@echo "✅ Setup de desenvolvimento concluído"

# Hot reload (desenvolvimento)
dev:
	@echo "🔥 Hot reload..."
	air -c .air.toml

# Verificar dependências vulneráveis
security-scan:
	@echo "🔒 Escaneando vulnerabilidades..."
	govulncheck ./...
	nancy sleuth

# Benchmark
benchmark:
	@echo "⚡ Executando benchmarks..."
	go test -bench=. -benchmem ./...

# Profile de CPU
profile-cpu:
	@echo "📊 Profile de CPU..."
	go test -cpuprofile=cpu.prof -bench=. ./...
	go tool pprof cpu.prof

# Profile de memória
profile-mem:
	@echo "📊 Profile de memória..."
	go test -memprofile=mem.prof -bench=. ./...
	go tool pprof mem.prof

# Gerar documentação
docs:
	@echo "📚 Gerando documentação..."
	godoc -http=:6060
	@echo "📚 Documentação disponível em: http://localhost:6060"

# Verificar qualidade do código
quality:
	@echo "📊 Verificando qualidade do código..."
	golangci-lint run
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	@echo "📊 Relatório de qualidade gerado"

# Validar configurações
validate-config:
	@echo "✅ Validando configurações..."
	./deploy.sh validate-config

# Instalar ferramentas de CI/CD
ci-setup:
	@echo "🔧 Setup de CI/CD..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/govulncheck@latest
	go install github.com/sonatypecommunity/nancy@latest
	@echo "✅ Setup de CI/CD concluído"

# Pipeline completo de CI
ci:
	@echo "🔄 Pipeline CI completo..."
	make deps
	make lint
	make test
	make security-scan
	make build
	@echo "✅ Pipeline CI concluído"

# Pipeline completo de qualidade
quality-gate:
	@echo "🚪 Quality Gate..."
	make test
	make test-concurrency
	make coverage
	make security-scan
	@echo "✅ Quality Gate concluído"

# Ajuda
help:
	@echo "BLOCO-1 - Core Inventory Makefile"
	@echo ""
	@echo "Targets principais:"
	@echo "  build          - Build da aplicação"
	@echo "  build-prod     - Build para produção (multi-plataforma)"
	@echo "  build-docker   - Build da imagem Docker"
	@echo "  deploy         - Deploy completo"
	@echo "  deploy-push    - Deploy com push para registry"
	@echo "  test           - Rodar todos os testes"
	@echo "  test-unit      - Testes de unidade"
	@echo "  test-integration - Testes de integração"
	@echo "  test-e2e       - Testes E2E"
	@echo "  test-concurrency - Testes de concorrência"
	@echo "  coverage       - Verificar cobertura"
	@echo "  health         - Health check completo"
	@echo "  slos           - Verificar SLOs"
	@echo "  load-test      - Teste de carga (Black Friday)"
	@echo "  stress-test    - Teste de estresse"
	@echo "  chaos          - Testes de Chaos Engineering"
	@echo "  rollback       - Rollback para versão anterior"
	@echo "  backup         - Criar backup"
	@echo "  logs           - Exibir logs"
	@echo "  status         - Status dos serviços"
	@echo "  cleanup        - Limpar recursos"
	@echo ""
	@echo "Targets de desenvolvimento:"
	@echo "  dev-setup      - Setup de desenvolvimento"
	@echo "  dev            - Hot reload"
	@echo "  run            - Rodar localmente"
	@echo ""
	@echo "Targets de qualidade:"
	@echo "  lint           - Executar linters"
	@echo "  security-scan  - Escanear vulnerabilidades"
	@echo "  benchmark      - Executar benchmarks"
	@echo "  quality        - Verificar qualidade do código"
	@echo ""
	@echo "Targets de CI/CD:"
	@echo "  ci             - Pipeline CI completo"
	@echo "  quality-gate   - Quality Gate"
	@echo "  validate-config - Validar configurações"
	@echo ""
	@echo "Exemplos:"
	@echo "  make build"
	@echo "  make deploy"
	@echo "  make load-test"
	@echo "  make slos"
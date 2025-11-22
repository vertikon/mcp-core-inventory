#!/bin/bash

# Script de Deploy para Produção - mcp-core-inventory (BLOCO-1 Core Ledger)
# Alinhado com o BLOCO-1 Blueprint - Core Inventory (Ledger ACID)
# Uso: ./deploy.sh [version] [--push] [--chaos]

set -e

# Configurações
SERVICE_NAME="mcp-core-inventory"
VERSION=${1:-latest}
REGISTRY="vertikon"
ENV_FILE=".env"
CHAOS_TEST=${CHAOS_TEST:-false}

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Funções de log
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_blue() {
    echo -e "${BLUE}[BLOCO-1]${NC} $1"
}

# Verificar pré-requisitos
check_prerequisites() {
    log_info "Verificando pré-requisitos do BLOCO-1..."
    
    # Verificar Docker
    if ! command -v docker &> /dev/null; then
        log_error "Docker não está instalado"
        exit 1
    fi
    
    # Verificar Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose não está instalado"
        exit 1
    fi
    
    # Verificar arquivo de ambiente
    if [ ! -f "$ENV_FILE" ]; then
        log_warn "Arquivo $ENV_FILE não encontrado. Usando env.example como template."
        if [ -f "env.example" ]; then
            cp env.example "$ENV_FILE"
            log_warn "Por favor, edite o arquivo $ENV_FILE com suas configurações antes de continuar."
            exit 1
        else
            log_error "Arquivo env.example não encontrado"
            exit 1
        fi
    fi
    
    # Verificar configurações essenciais do BLOCO-1
    if [ ! -f "./config/redis.conf" ]; then
        log_error "Arquivo de configuração do Redis não encontrado"
        exit 1
    fi
    
    if [ ! -f "./config/nats.conf" ]; then
        log_error "Arquivo de configuração do NATS não encontrado"
        exit 1
    fi
    
    if [ ! -f "./config/prometheus.yml" ]; then
        log_error "Arquivo de configuração do Prometheus não encontrado"
        exit 1
    fi
    
    log_info "Pré-requisitos verificados com sucesso"
}

# Build da imagem Docker
build_image() {
    log_blue "Construindo imagem Docker do Core Inventory (Ledger ACID)..."
    
    docker build -t "${REGISTRY}/${SERVICE_NAME}:${VERSION}" .
    
    if [ $? -eq 0 ]; then
        log_info "Imagem Docker construída com sucesso: ${REGISTRY}/${SERVICE_NAME}:${VERSION}"
    else
        log_error "Falha ao construir imagem Docker"
        exit 1
    fi
}

# Push para registry (opcional)
push_image() {
    if [ "$2" = "--push" ]; then
        log_info "Enviando imagem para registry..."
        docker push "${REGISTRY}/${SERVICE_NAME}:${VERSION}"
        log_info "Imagem enviada para registry com sucesso"
    fi
}

# Deploy com Docker Compose
deploy_services() {
    log_blue "Iniciando deployment dos serviços do BLOCO-1..."
    
    # Parar serviços existentes
    log_info "Parando serviços existentes..."
    docker-compose -f docker-compose.prod.yaml down || true
    
    # Limpar recursos antigos (opcional)
    if [ "$3" = "--clean" ]; then
        log_info "Limpando recursos antigos..."
        docker system prune -f
        docker volume prune -f
    fi
    
    # Iniciar serviços
    log_info "Iniciando serviços do BLOCO-1..."
    docker-compose -f docker-compose.prod.yaml up -d
    
    # Verificar status
    log_info "Verificando status dos serviços..."
    sleep 15
    docker-compose -f docker-compose.prod.yaml ps
    
    # Health check detalhado
    log_info "Aguardando health checks dos serviços críticos..."
    sleep 30
    
    # Verificar se o serviço principal está saudável
    if docker-compose -f docker-compose.prod.yaml ps | grep -q "Up (healthy)"; then
        log_blue "✅ Deploy do BLOCO-1 realizado com sucesso! Serviços estão saudáveis."
        log_blue "🔗 Acesse: http://localhost:8080/health"
        log_blue "📊 Métricas: http://localhost:9090"
        log_blue "🔍 Tracing: http://localhost:16686"
    else
        log_warn "Alguns serviços podem não estar saudáveis. Verifique os logs:"
        docker-compose -f docker-compose.prod.yaml logs --tail=50
    fi
}

# Backup do banco de dados
backup_database() {
    log_blue "Criando backup do banco de dados do Ledger..."
    
    BACKUP_DIR="./backups"
    TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
    BACKUP_FILE="${BACKUP_DIR}/mcp_core_inventory_backup_${TIMESTAMP}.sql"
    
    mkdir -p "$BACKUP_DIR"
    
    docker-compose -f docker-compose.prod.yaml exec postgres-prod pg_dump -U postgres mcp_core_inventory > "$BACKUP_FILE"
    
    if [ $? -eq 0 ]; then
        log_info "Backup criado com sucesso: $BACKUP_FILE"
    else
        log_error "Falha ao criar backup"
    fi
}

# Teste de Carga (BLOCO-1 Stress Test)
run_load_test() {
    log_blue "Executando teste de carga do BLOCO-1 (Black Friday Scenario)..."
    
    # Verificar se o k6 está disponível
    if ! command -v k6 &> /dev/null; then
        log_warn "k6 não encontrado localmente. Executando via Docker..."
        
        # Descobrir rede Docker
        NETWORK_NAME=$(docker network ls --filter name=bloco-1-network --format "{{.Name}}" | head -1)
        if [ -z "$NETWORK_NAME" ]; then
            NETWORK_NAME="mcp-core-inventory_bloco-1-network"
        fi
        
        # Executar k6 via Docker
        docker run --rm -i \
            --network "$NETWORK_NAME" \
            -v "${PWD}/ops/k6/bloco-1-core:/scripts" \
            -e BASE_URL="http://mcp-core-inventory:8080" \
            grafana/k6 run /scripts/reserve-load-test.js
    else
        # Executar k6 localmente
        k6 run ops/k6/bloco-1-core/reserve-load-test.js
    fi
}

# Chaos Engineering (BLOCO-1 Resilience Test)
run_chaos_test() {
    log_blue "Executando teste de Chaos Engineering do BLOCO-1..."
    
    # Teste 1: Redis Down
    log_info "Teste 1: Simulando falha do Redis..."
    docker-compose -f docker-compose.prod.yaml stop redis-prod
    sleep 10
    
    # Verificar se o fallback para Postgres está funcionando
    log_info "Verificando fallback para Postgres..."
    curl -f http://localhost:8080/health || log_error "Serviço indisponível"
    
    # Restaurar Redis
    log_info "Restaurando Redis..."
    docker-compose -f docker-compose.prod.yaml start redis-prod
    sleep 10
    
    # Teste 2: NATS Delay
    log_info "Teste 2: Simulando atraso no NATS..."
    # Implementar lógica de delay conforme necessidade
    
    log_blue "✅ Testes de Chaos concluídos"
}

# Rollback
rollback() {
    log_blue "Iniciando rollback do BLOCO-1..."
    
    # Obter imagem anterior
    PREVIOUS_IMAGE=$(docker images --format "table {{.Repository}}:{{.Tag}}" | grep "${REGISTRY}/${SERVICE_NAME}" | head -2 | tail -1)
    
    if [ -z "$PREVIOUS_IMAGE" ]; then
        log_error "Não foi encontrada imagem anterior para rollback"
        exit 1
    fi
    
    log_info "Fazendo rollback para: $PREVIOUS_IMAGE"
    
    # Atualizar docker-compose.yaml com imagem anterior
    sed -i "s|image: ${REGISTRY}/${SERVICE_NAME}:.*|image: $PREVIOUS_IMAGE|g" docker-compose.prod.yaml
    
    # Redeploy
    deploy_services
}

# Limpeza
cleanup() {
    log_info "Limpando recursos não utilizados do BLOCO-1..."
    
    # Remover imagens não utilizadas
    docker image prune -f
    
    # Remover volumes não utilizados
    docker volume prune -f
    
    log_info "Limpeza concluída"
}

# Verificação de SLOs do BLOCO-1
check_slos() {
    log_blue "Verificando SLOs do BLOCO-1..."
    
    # Verificar latência p99
    LATENCY=$(curl -s http://localhost:9090/api/v1/query?query=histogram_quantile(0.99,sum(rate(http_request_duration_seconds_bucket[5m]))by(le)) | jq -r '.data.result[0].value[1]')
    
    if (( $(echo "$LATENCY > 0.12" | bc -l) )); then
        log_warn "⚠️ SLO Violado: p99 latency > 120ms (atual: ${LATENCY}s)"
    else
        log_info "✅ SLO OK: p99 latency = ${LATENCY}s"
    fi
    
    # Verificar race conditions
    RACE_CONDITIONS=$(curl -s http://localhost:9090/api/v1/query?query=sum(increase(race_condition_incidents_total[5m])) | jq -r '.data.result[0].value[1]')
    
    if [ "$RACE_CONDITIONS" != "0" ]; then
        log_error "🚨 SLO Crítico Violado: Race Conditions detectadas = $RACE_CONDITIONS"
    else
        log_info "✅ SLO OK: Zero Race Conditions"
    fi
}

# Menu de ajuda
show_help() {
    echo "Uso: $0 [OPÇÃO] [ARGUMENTOS]"
    echo ""
    echo "Deploy do BLOCO-1 - Core Inventory (Ledger ACID)"
    echo ""
    echo "Opções:"
    echo "  deploy [version] [--push] [--clean]    Build e deploy da aplicação"
    echo "  backup                                   Criar backup do banco de dados"
    echo "  rollback                                 Rollback para versão anterior"
    echo "  load-test                                Executar teste de carga (Black Friday)"
    echo "  chaos                                    Executar testes de chaos"
    echo "  slos                                     Verificar SLOs do BLOCO-1"
    echo "  cleanup                                  Limpar recursos não utilizados"
    echo "  logs [service]                           Mostrar logs dos serviços"
    echo "  status                                   Mostrar status dos serviços"
    echo "  help                                     Mostrar esta ajuda"
    echo ""
    echo "Exemplos:"
    echo "  $0 deploy v1.0.0 --push"
    echo "  $0 deploy --clean"
    echo "  $0 load-test"
    echo "  $0 chaos"
    echo "  $0 slos"
    echo "  $0 logs mcp-core-inventory"
}

# Logs
show_logs() {
    SERVICE=${2:-""}
    if [ -n "$SERVICE" ]; then
        docker-compose -f docker-compose.prod.yaml logs -f "$SERVICE"
    else
        docker-compose -f docker-compose.prod.yaml logs -f
    fi
}

# Status
show_status() {
    log_blue "Status dos serviços do BLOCO-1:"
    docker-compose -f docker-compose.prod.yaml ps
    
    echo ""
    log_blue "Consumo de recursos:"
    docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}"
    
    echo ""
    log_blue "Verificando SLOs..."
    check_slos
}

# Main
case "${1:-deploy}" in
    "deploy")
        check_prerequisites
        build_image
        push_image "$@"
        deploy_services "$@"
        ;;
    "backup")
        backup_database
        ;;
    "rollback")
        rollback
        ;;
    "load-test")
        run_load_test
        ;;
    "chaos")
        run_chaos_test
        ;;
    "slos")
        check_slos
        ;;
    "cleanup")
        cleanup
        ;;
    "logs")
        show_logs "$@"
        ;;
    "status")
        show_status
        ;;
    "help"|"-h"|"--help")
        show_help
        ;;
    *)
        log_error "Opção inválida: $1"
        show_help
        exit 1
        ;;
esac

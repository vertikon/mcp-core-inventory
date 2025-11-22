#!/bin/bash

# Health Check Script - BLOCO-1 Core Inventory
# Verifica a saúde de todos os componentes críticos do BLOCO-1

set -e

# Configurações
SERVICE_NAME="mcp-core-inventory"
BASE_URL="http://localhost:8080"
PROMETHEUS_URL="http://localhost:9090"
JAEGER_URL="http://localhost:16686"

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

# Verificar saúde do serviço principal
check_core_service() {
    log_blue "Verificando saúde do Core Inventory..."
    
    # Health check básico
    if curl -f -s "${BASE_URL}/health" > /dev/null; then
        log_info "✅ Core Inventory: Saudável"
        return 0
    else
        log_error "❌ Core Inventory: Indisponível"
        return 1
    fi
}

# Verificar endpoints críticos
check_critical_endpoints() {
    log_blue "Verificando endpoints críticos..."
    
    local endpoints=("/health" "/metrics" "/v1/available")
    local all_healthy=true
    
    for endpoint in "${endpoints[@]}"; do
        if curl -f -s "${BASE_URL}${endpoint}" > /dev/null; then
            log_info "✅ ${endpoint}: OK"
        else
            log_error "❌ ${endpoint}: FALHA"
            all_healthy=false
        fi
    done
    
    if [ "$all_healthy" = true ]; then
        log_info "✅ Todos os endpoints críticos estão saudáveis"
        return 0
    else
        log_error "❌ Alguns endpoints estão com falha"
        return 1
    fi
}

# Verificar dependências
check_dependencies() {
    log_blue "Verificando dependências..."
    
    # PostgreSQL
    if docker-compose -f docker-compose.prod.yaml exec -T postgres-prod pg_isready -U postgres > /dev/null 2>&1; then
        log_info "✅ PostgreSQL: Saudável"
    else
        log_error "❌ PostgreSQL: Falha"
        return 1
    fi
    
    # Redis
    if docker-compose -f docker-compose.prod.yaml exec -T redis-prod redis-cli ping > /dev/null 2>&1; then
        log_info "✅ Redis: Saudável"
    else
        log_error "❌ Redis: Falha"
        return 1
    fi
    
    # NATS
    if curl -f -s http://localhost:8222/varz > /dev/null; then
        log_info "✅ NATS: Saudável"
    else
        log_error "❌ NATS: Falha"
        return 1
    fi
    
    return 0
}

# Verificar SLOs
check_slos() {
    log_blue "Verificando SLOs do BLOCO-1..."
    
    # Verificar latência p99
    local latency=$(curl -s "${PROMETHEUS_URL}/api/v1/query?query=histogram_quantile(0.99,sum(rate(http_request_duration_seconds_bucket[5m]))by(le))" | jq -r '.data.result[0].value[1]' 2>/dev/null || echo "null")
    
    if [ "$latency" != "null" ] && (( $(echo "$latency > 0.12" | bc -l) )); then
        log_warn "⚠️ SLO Violado: p99 latency > 120ms (atual: ${latency}s)"
    else
        log_info "✅ SLO OK: p99 latency = ${latency}s"
    fi
    
    # Verificar race conditions
    local race_conditions=$(curl -s "${PROMETHEUS_URL}/api/v1/query?query=sum(increase(race_condition_incidents_total[5m]))" | jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    
    if [ "$race_conditions" != "0" ]; then
        log_error "🚨 SLO Crítico Violado: Race Conditions detectadas = $race_conditions"
        return 1
    else
        log_info "✅ SLO OK: Zero Race Conditions"
    fi
    
    # Verificar error rate
    local error_rate=$(curl -s "${PROMETHEUS_URL}/api/v1/query?query=rate(http_requests_total{status=~\"5..\"}[5m])" | jq -r '.data.result[0].value[1]' 2>/dev/null || echo "0")
    
    if [ "$error_rate" != "null" ] && (( $(echo "$error_rate > 0.01" | bc -l) )); then
        log_warn "⚠️ SLO Violado: Error rate > 1% (atual: $(echo "$error_rate * 100" | bc -l)%)"
    else
        log_info "✅ SLO OK: Error rate = $(echo "$error_rate * 100" | bc -l)%"
    fi
    
    return 0
}

# Verificar recursos
check_resources() {
    log_blue "Verificando utilização de recursos..."
    
    # Verificar uso de CPU e memória
    local stats=$(docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}" mcp-core-inventory)
    
    echo "$stats" | while IFS=$'\t' read -r container cpu mem; do
        # Extrair valor numérico da porcentagem de CPU
        local cpu_value=$(echo "$cpu" | sed 's/%//')
        
        # Verificar CPU
        if (( $(echo "$cpu_value > 80" | bc -l) )); then
            log_warn "⚠️ Alto uso de CPU: ${cpu}"
        else
            log_info "✅ CPU OK: ${cpu}"
        fi
        
        # Verificar memória (simplificado)
        if [[ "$mem" == *"GiB"* ]]; then
            local mem_value=$(echo "$mem" | sed 's/GiB//')
            if (( $(echo "$mem_value > 1.5" | bc -l) )); then
                log_warn "⚠️ Alto uso de memória: ${mem}"
            else
                log_info "✅ Memória OK: ${mem}"
            fi
        fi
    done
}

# Verificar logs de erros
check_error_logs() {
    log_blue "Verificando logs de erros recentes..."
    
    # Verificar logs de erros nos últimos 5 minutos
    local error_count=$(docker-compose -f docker-compose.prod.yaml logs --since=5m mcp-core-inventory 2>&1 | grep -i "error\|exception\|fatal" | wc -l)
    
    if [ "$error_count" -gt 0 ]; then
        log_warn "⚠️ Encontrados ${error_count} logs de erro nos últimos 5 minutos"
        return 1
    else
        log_info "✅ Nenhum log de erro nos últimos 5 minutos"
        return 0
    fi
}

# Verificar conectividade externa
check_connectivity() {
    log_blue "Verificando conectividade externa..."
    
    # Verificar se consegue se conectar ao mcp-fulfillment-ops (quando existir)
    if docker network ls | grep -q "bloco-1-network"; then
        # Testar resolução de nome na rede
        if docker run --rm --network bloco-1-network alpine ping -c 1 mcp-fulfillment-ops > /dev/null 2>&1; then
            log_info "✅ Conectividade com mcp-fulfillment-ops: OK"
        else
            log_warn "⚠️ mcp-fulfillment-ops não está acessível (pode ser normal se não implantado)"
        fi
    fi
    
    return 0
}

# Relatório final
generate_report() {
    local exit_code=$1
    
    echo ""
    log_blue "=== RELATÓRIO DE SAÚDE DO BLOCO-1 ==="
    echo "Serviço: $SERVICE_NAME"
    echo "Timestamp: $(date)"
    echo "Status: $([ $exit_code -eq 0 ] && echo 'SAUDÁVEL' || echo 'PROBLEMAS DETECTADOS')"
    echo ""
    
    if [ $exit_code -eq 0 ]; then
        log_info "🎉 BLOCO-1 está operacional e saudável!"
        echo ""
        echo "Acessos:"
        echo "  • API: ${BASE_URL}"
        echo "  • Métricas: ${PROMETHEUS_URL}"
        echo "  • Tracing: ${JAEGER_URL}"
    else
        log_error "🚨 BLOCO-1 apresenta problemas que requerem atenção!"
        echo ""
        echo "Ações recomendadas:"
        echo "  • Verificar logs: docker-compose -f docker-compose.prod.yaml logs mcp-core-inventory"
        echo "  • Verificar status: docker-compose -f docker-compose.prod.yaml ps"
        echo "  • Consultar runbooks: docs/runbooks/"
    fi
    
    echo ""
    echo "=== FIM DO RELATÓRIO ==="
}

# Menu de ajuda
show_help() {
    echo "Uso: $0 [OPÇÃO]"
    echo ""
    echo "Health Check do BLOCO-1 - Core Inventory"
    echo ""
    echo "Opções:"
    echo "  all           Verifica todos os componentes (padrão)"
    echo "  service       Verifica apenas o serviço principal"
    echo "  endpoints     Verifica endpoints críticos"
    echo "  dependencies  Verifica dependências (PostgreSQL, Redis, NATS)"
    echo "  slos          Verifica SLOs (latência, race conditions, error rate)"
    echo "  resources     Verifica utilização de recursos"
    echo "  logs          Verifica logs de erros"
    echo "  connectivity  Verifica conectividade externa"
    echo "  help          Mostra esta ajuda"
    echo ""
    echo "Exemplos:"
    echo "  $0"
    echo "  $0 all"
    echo "  $0 slos"
    echo "  $0 service"
}

# Main
case "${1:-all}" in
    "all")
        check_core_service
        check_critical_endpoints
        check_dependencies
        check_slos
        check_resources
        check_error_logs
        check_connectivity
        generate_report $?
        ;;
    "service")
        check_core_service
        generate_report $?
        ;;
    "endpoints")
        check_critical_endpoints
        generate_report $?
        ;;
    "dependencies")
        check_dependencies
        generate_report $?
        ;;
    "slos")
        check_slos
        generate_report $?
        ;;
    "resources")
        check_resources
        generate_report $?
        ;;
    "logs")
        check_error_logs
        generate_report $?
        ;;
    "connectivity")
        check_connectivity
        generate_report $?
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

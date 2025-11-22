# BLOCO-1 - Core Inventory Deployment Guide

## 📋 Visão Geral

O **mcp-core-inventory** é o componente **Ledger ACID** do BLOCO-1, responsável por manter a **fonte única da verdade** sobre estoque (saldo, reserva, alocação, lotes, validade).

Este documento guia o deployment em produção seguindo as especificações do [BLOCO-1 Blueprint](../../../cursor/BLOCOS/BLOCO-1-BLUEPRINT-mcp-core-inventory.md).

---

## 🏗️ Arquitetura do Deployment

```
┌─────────────────────────────────────────────────────────────────┐
│                    BLOCO-1 - Core Inventory               │
├─────────────────────────────────────────────────────────────────┤
│  mcp-core-inventory (Ledger ACID)                        │
│  ├── API REST (Porta 8080)                               │
│  ├── NATS Events (Eventos inventory.*)                     │
│  ├── Redis Cache (Locks e Cache de Saldo)                  │
│  └── PostgreSQL (Persistência ACID)                          │
├─────────────────────────────────────────────────────────────────┤
│  Infraestrutura de Observabilidade                             │
│  ├── Prometheus (Métricas)                                   │
│  ├── Jaeger (Tracing Distribuído)                             │
│  ├── Grafana (Dashboards)                                     │
│  └── Nginx (Gateway/API Gateway)                              │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Deploy Rápido

### 1. Pré-requisitos

- Docker 20.10+
- Docker Compose 2.0+
- 4GB+ RAM
- 20GB+ Disco

### 2. Configuração

```bash
# Copiar template de ambiente
cp env.example .env

# Editar com suas configurações
nano .env
```

Variáveis essenciais:
```bash
DB_PASSWORD=senha_forte_postgres
REDIS_PASSWORD=senha_forte_redis
JWT_SECRET=seu_jwt_secreto_minimo_32_caracteres
AI_API_KEY=sua_chave_api_ai
```

### 3. Deploy

```bash
# Build e deploy
./deploy.sh deploy v1.0.0

# Deploy com push para registry
./deploy.sh deploy v1.0.0 --push

# Deploy com limpeza de recursos
./deploy.sh deploy v1.0.0 --clean
```

---

## 🔧 Componentes do Deployment

### Serviços Principais

| Serviço | Descrição | Porta | Saúde |
|----------|-------------|--------|--------|
| mcp-core-inventory | Ledger ACID do estoque | 8080 | `/health` |
| postgres-prod | Banco de dados PostgreSQL | 5432 | `pg_isready` |
| redis-prod | Cache e locks distribuídos | 6379 | `redis-cli ping` |
| nats-prod | Message broker (eventos) | 4222/8222 | `/varz` |

### Observabilidade

| Serviço | Descrição | Porta | Acesso |
|----------|-------------|--------|---------|
| prometheus-prod | Coletor de métricas | 9090 | http://localhost:9090 |
| jaeger-prod | Tracing distribuído | 16686 | http://localhost:16686 |
| nginx | Reverse proxy | 80/443 | http://localhost |

---

## 📊 SLOs e Métricas

### SLOs Críticos (P0)

| Métrica | Objetivo | Alerta |
|----------|-----------|---------|
| p99 `/reserve` | < 120ms | > 120ms por 5min |
| Error rate | < 1% | > 1% por 1min |
| Race conditions | = 0 | > 0 imediato |
| Redis latency | < 50ms | > 50ms por 1min |
| Postgres lock wait | < 2s | > 2s por 1min |

### Dashboards Principais

1. **Core Inventory Health**: Métricas de negócio e performance
2. **Race Conditions Monitor**: Detecção de violações de concorrência
3. **Resource Utilization**: CPU, memória, I/O dos serviços

---

## 🧪 Testes de Carga

### Black Friday Scenario

```bash
# Executar teste de carga
./deploy.sh load-test

# Ou manualmente via Docker
docker run --rm -i \
  --network bloco-1-network \
  -v "${PWD}/ops/k6/bloco-1-core:/scripts" \
  -e BASE_URL="http://mcp-core-inventory:8080" \
  grafana/k6 run /scripts/reserve-load-test.js
```

### Cenários Testados

1. **Alta Concorrência**: 1000 usuários tentando reservar mesmo SKU
2. **Volume Sustentado**: 1500 RPS por 10 minutos
3. **Chaos Engineering**: Falha de Redis, latência de Postgres

---

## 🛡️ Segurança

### Configurações de Produção

- **RBAC**: Controle de acesso por papel
- **TLS/SSL**: Comunicação criptografada
- **Rate Limiting**: Limitação por IP/usuário
- **Input Validation**: Validação rigorosa de entradas
- **Audit Logging**: Logs estruturados com trace_id

### Variáveis de Segurança

```bash
# JWT
JWT_SECRET=minimo_32_caracteres_aqui

# Database
DB_PASSWORD=senha_super_forte_aqui
POSTGRES_SSL_MODE=require

# Redis
REDIS_PASSWORD=senha_forte_redis_aqui

# API
AI_API_KEY=sua_chave_api_secreta
```

---

## 🔄 Operações Dia-a-Dia

### Backup

```bash
# Backup manual
./deploy.sh backup

# Backup automático (via cron)
0 2 * * * /path/to/mcp-core-inventory/deploy.sh backup
```

### Monitoramento

```bash
# Verificar status
./deploy.sh status

# Verificar SLOs
./deploy.sh slos

# Ver logs
./deploy.sh logs mcp-core-inventory
```

### Rollback

```bash
# Rollback para versão anterior
./deploy.sh rollback
```

---

## 🚨 Incident Response

### Runbooks

1. **Race Condition Detectada**: [Runbook](docs/runbooks/race-condition.md)
2. **Redis Indisponível**: [Runbook](docs/runbooks/redis-down.md)
3. **Alta Latência**: [Runbook](docs/runbooks/high-latency.md)
4. **Overselling Detectado**: [Runbook](docs/runbooks/overselling.md)

### Comandos de Emergência

```bash
# Escalar serviço
docker-compose -f docker-compose.prod.yaml up -d --scale mcp-core-inventory=3

# Reiniciar serviço específico
docker-compose -f docker-compose.prod.yaml restart redis-prod

# Ver logs em tempo real
docker-compose -f docker-compose.prod.yaml logs -f mcp-core-inventory
```

---

## 📈 Escalabilidade

### Horizontal Scaling

```bash
# Escalar para 3 réplicas
docker-compose -f docker-compose.prod.yaml up -d --scale mcp-core-inventory=3

# Escalar com limites de recursos
docker-compose -f docker-compose.prod.yaml up -d \
  --scale mcp-core-inventory=3 \
  --scale redis-prod=2
```

### Vertical Scaling

Ajustar recursos em `docker-compose.prod.yaml`:

```yaml
deploy:
  resources:
    limits:
      cpus: '2.0'
      memory: 2G
    reservations:
      cpus: '1.0'
      memory: 1G
```

---

## 🔧 Troubleshooting

### Problemas Comuns

#### 1. Serviço não inicia

```bash
# Verificar logs
docker-compose -f docker-compose.prod.yaml logs mcp-core-inventory

# Verificar configurações
docker-compose -f docker-compose.prod.yaml config
```

#### 2. Conexão com banco falha

```bash
# Testar conexão
docker-compose -f docker-compose.prod.yaml exec postgres-prod pg_isready -U postgres

# Verificar credenciais
docker-compose -f docker-compose.prod.yaml exec mcp-core-inventory env | grep DATABASE
```

#### 3. Redis não responde

```bash
# Testar Redis
docker-compose -f docker-compose.prod.yaml exec redis-prod redis-cli ping

# Verificar configuração
docker-compose -f docker-compose.prod.yaml exec redis-prod redis-cli config get "*"
```

### Performance Issues

```bash
# Verificar uso de recursos
docker stats

# Verificar métricas
curl http://localhost:9090/api/v1/query?query=rate(http_requests_total[5m])

# Verificar traces
curl http://localhost:16686/api/traces?service=mcp-core-inventory
```

---

## 📚 Referências

- [BLOCO-1 Blueprint](../../../cursor/BLOCOS/BLOCO-1-BLUEPRINT-mcp-core-inventory.md)
- [API Documentation](./docs/api/openapi.yaml)
- [Architecture](./docs/architecture/clean-architecture.md)
- [Monitoring Guide](./docs/monitoring/setup.md)

---

## 🆘 Suporte

Em caso de incidentes críticos:

1. Verificar [Runbooks](docs/runbooks/)
2. Consultar [Dashboards](http://localhost:3000)
3. Escalar equipe SRE via canal `#bloco-1-alerts`

---

**Status**: ✅ Production Ready  
**Versão**: v1.0.0  
**Última Atualização**: 2025-11-22

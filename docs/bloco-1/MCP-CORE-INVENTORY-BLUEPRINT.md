**Blueprint dedicado** para  
`E:\vertikon\.endurance\internal\services\bloco-1-core\mcp-core-inventory`
(agora versionado em `docs/bloco-1/`).  
Pronto para guiar auditorias, desenvolvimento e validação contínua.

---

## 0. Missão do MCP

> **mcp-core-inventory = livro-razão ACID do estoque**  
> É o **único** lugar que altera o saldo consolidado (saldo, reserva, alocação, lotes, validade) em toda a plataforma.

Ele responde **“quanto tem, onde, e o que está reservado”** – independentemente de canal, checkout ou fulfillment.

---

## 1. Contexto no BLOCO-1

No BLOCO-1 você tem dois motores principais:

* **`mcp-core-inventory`** → **Ledger ACID** (saldo/reserva/ajustes/lotes)
* **`mcp-fulfillment-ops`** → execução física (recebimento, expedição, devolução, transferência), que **nunca grava estoque direto**, apenas **solicita** pro core.

**Princípio fundador:**

* Fulfillment **solicita** alteração via API/evento
* Core Inventory **executa** a mudança e publica o novo estado

---

## 2. Escopo, responsabilidades e invariantes

### 2.1. Escopo Positivo (o que este MCP faz)

1. Manter o **livro-razão de estoque**:

   * `StockLedger`, `StockEntry`, `Reservation`, `StockMovement`, `Batch`.
2. Expor operações síncronas de estoque para outros blocos:

   * Reserva, confirmação, liberação, ajustes administrativos, consulta de saldo disponível.
3. Aplicar **regras de negócio críticas**:

   * Não vender sem saldo.
   * Garantir ordem correta reserva → confirmação/expedição → baixa.
   * Política de lotes **FIFO/FEFO** para consumo.
4. Manter **idempotência** de comandos críticos (principalmente reserva/confirm).
5. Publicar eventos padrão NATS (via schemas):

   * `inventory.reserve.request`, `inventory.reserve.confirmed`, e eventos de fulfillment correlatos.

### 2.2. Fora de Escopo (o que NÃO faz)

* Não calcula **preço** (isso é bloco de vendas).
* Não se preocupa com **pagamento** ou status financeiro.
* Não contém **PII de cliente** (somente IDs técnicos).
* Não orquestra fluxo físico (isso é `mcp-fulfillment-ops`).

### 2.3. Invariantes

* **Saldo nunca negativo**.
* **Race condition = 0** (meta hard); qualquer ocorrência vira incidente P0.
* Toda reserva tem:

  * TTL
  * `idempotency_key`
  * lock via Redis+Lua + fallback Postgres ACID.

---

## 3. Interfaces oficiais (contratos)

### 3.1. HTTP / OpenAPI

Diretório de contrato:

* `contracts/openapi/bloco-1-core/core-inventory-v1.yaml`
* `contracts/openapi/bloco-1-core/core-inventory-v2.yaml`

**Operações mínimas v1:**

* `POST /reserve` → cria reserva de estoque
* `POST /confirm` → confirma reserva (baixa definitiva)
* `POST /release` → libera reserva (volta pro saldo)
* `POST /adjust` → ajuste administrativo (entrada/saída manual)
* `GET /available` → consulta de saldo disponível para um SKU/local

**v2**: mantém as mesmas operações, mas com:

* headers de **Deprecation/Sunset**
* campos adicionais para lote/validade/extensões.

### 3.2. Eventos NATS (Protobuf/Avro)

Pasta:

* `contracts/nats-schemas/bloco-1-core/`

Eventos principais:

* `inventory.reserve.request`
* `inventory.reserve.confirmed`
* (e eventos de `fulfillment.*` consumidos/gerados em coordenação com `mcp-fulfillment-ops`).

### 3.3. Consumer–Provider externo

Pactos com Checkout/OMS:

* `contracts/pact/bloco-1-core-checkout/*.json`

---

## 4. Arquitetura interna de código

### 4.1. Layout de diretório (núcleo do MCP)

```
internal/services/bloco-1-core/
└── mcp-core-inventory/
    ├── cmd/core-inventory/main.go
    ├── internal/
    │   ├── app/            # casos de uso
    │   ├── domain/ledger/  # regras de negócio puras
    │   ├── adapters/       # postgres, redis, nats
    │   └── interfaces/http/# handlers HTTP/gRPC
    └── tests/ledger/       # testes de domínio table-driven
```

### 4.2. Papéis por camada

- **`cmd/core-inventory/main.go`**  
  Bootstrap, DI, carregamento de configs, wiring de HTTP + NATS + Telemetry.

- **`internal/app/`**  
  Casos de uso orquestradores:
  - `reserve_stock.go`
  - `confirm_reservation.go`
  - `release_reservation.go`
  - `adjust_stock.go`
  - `query_available.go`

- **`internal/domain/ledger/`**  
  Regra de negócio pura:
  - `stock_ledger.go` (agregado principal)
  - `reservation.go`
  - `stock_movement.go`
  - `policies.go` (não vender sem saldo, etc.)
  - pasta `batches/` para lotes/validade.

- **`internal/adapters/postgres/`**  
  - `ledger_repository.go` + `migrations/*.sql` (0001 ledger, 0002 batches…).

- **`internal/adapters/redis/`**  
  - `stock_cache.go` (cache de saldo por SKU/local)
  - `reservation_lock.go` (locks distribuídos via Lua).

- **`internal/adapters/nats/`**  
  - publisher/consumer para eventos `inventory.*` e integração com fulfillment.

- **`internal/interfaces/http/`**  
  - Handlers dos endpoints descritos no OpenAPI.

- **`tests/`**  
  - foco em **table-driven tests** no domínio `ledger/` e testes E2E de fluxo de reserva/alocação.

---

## 5. Fluxos funcionais principais

### 5.1. `ReserveStock`

1. Recebe comando `/reserve`.  
2. Gera `idempotency_key` (ou valida uma recebida).  
3. Aplica lock Redis+Lua para o SKU/local.  
4. Verifica saldo disponível no ledger.  
5. Cria `Reservation` com TTL.  
6. Persiste movimento no Postgres (transação ACID).  
7. Eventualmente publica `inventory.reserve.request`.

**Falha:**  
– Se não houver saldo, responde erro de negócio (não técnico).

### 5.2. `ConfirmReservation`

1. Recebe comando `/confirm` referenciando uma reserva.  
2. Garante invariantes de saldo (não pode negativar).  
3. Converte a reserva em baixa definitiva.  
4. Grava no ledger.  
5. Publica `inventory.reserve.confirmed`.

### 5.3. `ReleaseReservation`

- Libera reservas expiradas ou canceladas:
  - devolve saldo ao ledger
  - invalida cache
  - registra movimento de “liberação”.

### 5.4. `AdjustStock`

- Ajustes administrativos (inventário, correção, auditoria):
  - gera `StockMovement` de ajuste
  - impacta ledger com rastreabilidade total.

### 5.5. `QueryAvailable`

- Consulta de saldo disponível:
  - tenta cache (`stock_cache.go`)
  - em cache miss, calcula a partir do ledger.

---

## 6. Infra, observabilidade e SLOs

### 6.1. Métricas e dashboards

- `ops/grafana-dashboards/bloco-1-core/core-inventory-health.json`
- `ops/prometheus-rules/bloco-1-core-rules.yaml`

**Obrigatório acompanhar:**

- p99 de `/reserve` e `/available`
- `stock_reservation_errors_total`
- `race_condition_incidents_total`
- `redis_latency_ms`, `postgres_lock_wait_ms`

Alertas P0: p99 alto, error_rate>1%, race_condition>0, Redis down>20s, lock>2s.

### 6.2. Testes de carga oficiais

- `ops/k6/bloco-1-core/reserve-load-test.js` (1.500 RPS /reserve, p99<120ms)
- `ops/k6/bloco-1-core/available-read-test.js` (3.000 RPS /available, p99<50ms)
- `ops/k6/bloco-1-core/mixed-scenario-test.js` (mix reserva+leitura+devolução).

---

## 7. Maturidade por nível

### P0 — GO LIVE

- Ledger ACID implementado (Postgres + migrations).
- Reserve/Confirm/Release/Adjust/Query operando.
- Redis + fallback Postgres testados (chaos Redis down).
- Métricas e alertas mínimos em produção.
- Testes de carga essenciais (`reserve-load-test`, `available-read-test`).

### P1 — 30 dias pós-produção

- Pact Broker ativo entre Core Inventory e Checkout.
- Dashboards comparando v1/v2.
- Limpeza automática de reservas expiradas (TTLs).

### P2 — Maturidade avançada

- Chaos Engineering recorrente: cenários Redis down, latência Postgres, NATS com delay.
- Métrica de “dívida de consistência” (reservas órfãs vs ledger).
- E2E multi-bloco (Checkout → Core Inventory → Fulfillment → Logística).

---

## 8. Como usar esse blueprint na prática

1. Criar/atualizar `docs/bloco-1/MCP-CORE-INVENTORY-BLUEPRINT.md` (este arquivo).
2. Garantir que a árvore de diretórios do MCP bate com o blueprint (principalmente `app/domain/adapters/interfaces/tests`).
3. Marcar o status P0/P1/P2 de cada item no quadro operacional do time.

Se necessário, derivar checklists operacionais, runbooks e scripts específicos a partir deste blueprint.

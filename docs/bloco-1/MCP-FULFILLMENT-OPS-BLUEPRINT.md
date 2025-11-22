Blueprint oficial do **`MCP-FULFILLMENT-OPS`**, espelhado no modelo do `mcp-core-inventory` e agora versionado em `docs/bloco-1/`.

---

## 0. Missão do MCP

> **mcp-fulfillment-ops = cérebro operacional da movimentação física de estoque**  
> Orquestra **entrada, saída, transferências, devoluções e contagens físicas**, garantindo que **todo movimento real** gere o comando correto para o `mcp-core-inventory` (livro-razão) sem nunca manipular saldo diretamente.

É o MCP que responde: **“o que está acontecendo fisicamente com o estoque (onde, quando, como e por quê)”**.

---

## 1. Contexto no BLOCO-1

No BLOCO-1 temos dois motores principais:

* **`mcp-core-inventory`** → **Ledger ACID** (saldo, reserva, lotes, validade).
* **`mcp-fulfillment-ops`** → **Execução física dos movimentos** (inbound, outbound, devolução, transferência, contagem).

**Princípio fundador:**

* Fulfillment **NUNCA altera saldo**.
* Fulfillment **solicita** via API/evento.
* Core Inventory **executa** a mudança e publica o novo estado.

---

## 2. Escopo, responsabilidades e invariantes

### 2.1 Escopo positivo (o que este MCP faz)

1. **Orquestrar operações físicas de estoque:**
   * Recebimento de mercadorias (CD, loja, cross-dock).
   * Separação (picking) e embalagem (packing).
   * Expedição / saída logística.
   * Devoluções (cliente → loja/CD).
   * Transferências entre locais.
   * Contagens cíclicas e inventários físicos.
2. **Traduzir o mundo físico em comandos para o ledger:**
   * Ao confirmar um recebimento → chama `core-inventory` para **entrada**.
   * Ao expedir pedido já reservado → chama `core-inventory` para **baixa definitiva**.
   * Ao registrar devolução → chama `core-inventory` para **entrada de retorno / ajuste**.
   * Ao concluir transferência → chama `core-inventory` para **saída de origem + entrada no destino**.
3. **Manter o estado operacional de workflows de fulfillment:**
   * `FulfillmentOrder`, `InboundShipment`, `OutboundShipment`, `TransferOrder`, `ReturnOrder`, `CycleCountTask`.
   * Estados: `pending`, `in_progress`, `blocked`, `completed`, `cancelled`.
4. **Integrar com OMS/WMS/TMS e canais de venda/logística.**
5. **Gerar trilha de auditoria física.**

### 2.2 Fora de escopo

* Não calcula preço/promoção.
* Não controla pagamento/financeiro.
* Não guarda PII sensível.
* Não decide política de estoque (min/max etc.).

### 2.3 Invariantes

* Nenhuma operação altera saldo sem passar pelo `mcp-core-inventory`.
* Todo movimento físico gera/consome reserva no ledger.
* Estados de workflow são idempotentes/reprocessáveis.
* Operações críticas exigem `idempotency_key` e correlação (`trace_id`, `correlation_id`).

---

## 3. Interfaces oficiais (contratos)

### 3.1 HTTP / OpenAPI

Pasta sugerida:

* `contracts/openapi/bloco-1-core/fulfillment-ops-v1.yaml`
* `contracts/openapi/bloco-1-core/fulfillment-ops-v2.yaml`

**Operações v1:**

- **Inbound**: `POST /inbound/start`, `POST /inbound/confirm`.
- **Outbound**: `POST /outbound/start_picking`, `POST /outbound/ship`.
- **Transfer**: `POST /transfer/create`, `POST /transfer/complete`.
- **Returns**: `POST /returns/register`, `POST /returns/complete`.
- **Cycle Count**: `POST /cycle_count/open`, `POST /cycle_count/submit`.

v2 adiciona lote/validade/extensões e headers de depreciação.

### 3.2 Eventos NATS (AsyncAPI)

* `fulfillment.inbound.created/received`
* `fulfillment.outbound.picking_started/shipped`
* `fulfillment.transfer.created/completed`
* `fulfillment.return.registered/completed`
* `fulfillment.cycle_count.opened/completed`

Integra-se com eventos do ledger (`inventory.reserve.confirmed`, etc.).

### 3.3 Pactos

* `contracts/pact/bloco-1-core-fulfillment-oms/`
* `contracts/pact/bloco-1-core-fulfillment-logistics/`

---

## 4. Arquitetura interna de código

### 4.1 Layout

```
internal/services/bloco-1-core/
└── mcp-fulfillment-ops/
    ├── cmd/fulfillment-ops/main.go
    ├── internal/
    │   ├── app/                 # casos de uso
    │   ├── domain/fulfillment/  # agregados e políticas
    │   ├── adapters/            # postgres, nats, wms
    │   └── interfaces/http/     # handlers
    └── tests/
        ├── domain/
        └── e2e/
```

### 4.2 Camadas

- `cmd/` – bootstrap, DI, config, wiring.
- `internal/app/` – `receive_goods.go`, `ship_order.go`, `register_return.go`, `complete_transfer.go`, `open_cycle_count.go`, `submit_cycle_count.go`.
- `internal/domain/fulfillment/` – entidades `FulfillmentOrder`, `InboundShipment`, `OutboundShipment`, `TransferOrder`, `ReturnOrder`, `CycleCountTask`, `policies.go`.
- `internal/adapters/postgres/` – persistência de workflows.
- `internal/adapters/nats/` – publishers/subscribers.
- `internal/adapters/wms/` – integrações externas.
- `internal/interfaces/http/` – handlers REST.
- `tests/domain` + `tests/e2e`.

---

## 5. Fluxos funcionais principais

### 5.1 `ReceiveGoods`

1. OMS/WMS emite `inbound.created`.
2. `mcp-fulfillment-ops` cria `InboundShipment`.
3. Operação física confere itens.
4. `/inbound/confirm` chama caso de uso `ReceiveGoods`.
5. Valida estado, registra conferência, chama `mcp-core-inventory` para entrada, atualiza shipment para `completed`, publica `fulfillment.inbound.received`.

### 5.2 `ShipOrder`

1. OMS envia `order.ready_to_pick`.
2. MCP cria `FulfillmentOrder`.
3. Picking/packing executados.
4. `/outbound/ship` → caso `ShipOrder`.
5. Valida linhas, chama Core para confirmar reservas, atualiza order para `shipped`, publica `fulfillment.outbound.shipped`.

### 5.3 `TransferStock`

1. Cria `TransferOrder`.
2. Ao despachar: chama Core para saída.
3. Ao receber destino: chama Core para entrada, atualiza order para `completed`.

### 5.4 `RegisterReturn`

1. `/returns/register` cria `ReturnOrder`.
2. `/returns/complete` aplica lógica (estoque vendável ou não), chama Core para entrada/ajuste, marca `completed`.

### 5.5 `CycleCount`

1. Scheduler cria `CycleCountTask`.
2. Operador conta e envia `/cycle_count/submit`.
3. Caso de uso compara com ledger, gera ajustes via Core e registra divergências.

---

## 6. Infra, observabilidade e SLOs

### 6.1 Métricas

- `fulfillment_inbound_duration_seconds`
- `fulfillment_outbound_duration_seconds`
- `fulfillment_return_duration_seconds`
- `fulfillment_transfer_duration_seconds`
- `fulfillment_sync_errors_core_inventory_total`
- `fulfillment_state_inconsistencies_total`

Alertas P0: p99 acima de SLA, erros de sincronização com Core, filas pendentes.

### 6.2 Testes de carga

- `inbound_burst_test`
- `outbound_peak_test`
- `core_down_degraded_mode`

---

## 7. Maturidade por nível

### P0 — GO LIVE

- Fluxos básicos implementados (inbound/outbound/transfer/returns).
- Integração confiável com `mcp-core-inventory`.
- Métricas e alertas mínimos.
- Testes E2E OMS → Fulfillment → Core.

### P1 — 30 dias pós-produção

- Contagem cíclica integrada ao Core.
- Relatórios de divergência física vs ledger.
- Pact Broker ativo com OMS/Logística.

### P2 — Maturidade avançada

- Chaos engineering recorrente (Core down, NATS delay, WMS offline).
- Otimização de roteamento operacional.
- KPIs de eficiência operacional.

---

## 8. Como usar este blueprint

1. Manter esta cópia sob controle de versão.
2. Garantir que a árvore real (`internal/app`, `internal/domain`, `internal/adapters`, etc.) siga este desenho.
3. Usar os requisitos P0/P1/P2 como checklist em reuniões de go-live e auditorias.

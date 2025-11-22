# BLOCO-1 – Plano de Implementação dos 415 Artefatos Pendentes

Este blueprint descreve **como** entregar os 415 arquivos exigidos pelo _Tree Validator_ mas ausentes no repositório. Ele se apoia no relatório `validation-result.json` e nos blueprints oficiais (Core Inventory, Fulfillment Ops e Executivos) para priorizar a implementação.

---

## 1. Visão Geral do Gap

| Categoria                       | Total | Exemplos                                                                                  |
|---------------------------------|-------|-------------------------------------------------------------------------------------------|
| Contratos / Documentação        | 118   | `openapi.yaml`, `contracts/openapi/bloco-1-core/core-inventory-v1.yaml`, `docs/guides/*`   |
| Código de Domínio / Serviços    | 144   | `internal/domain/ledger/*`, `internal/services/state_app_service.go`, `internal/mcp/*`     |
| Ferramentas / Scripts           | 56    | `scripts/validation/*.sh`, `tools/generators/*`, `ops/k6/*`                                |
| Templates / Generators          | 62    | `templates/go/*`, `templates/web/*`, `templates/mcp-go-premium/*`                          |
| Infra / Config                  | 35    | `config/mcp/*.yaml`, `config/monitoring/*.yaml`, `ops/prometheus/*`, `docker-compose/*`    |

> **Meta:** reduzir `missing_files` de 415 para 0 e `compliance_percent` para ≥ 1.00.

---

## 2. Estratégia de Execução em 10 Fases

| Fase | Objetivo Principal | Entregáveis-Chave |
|------|-------------------|-------------------|
| **F0 – Preparação** | Congelar estado atual e medir baseline. | Branch `feature/bloco1-blueprint`, execução do `tools/validate_tree.go`, lista priorizada de arquivos ausentes. |
| **F1 – Conversão de Blueprints para Docs** | Trazer os blueprints do `.cursor` para `docs/bloco-1/`. | `docs/bloco-1/MCP-CORE-INVENTORY-BLUEPRINT.md`, `docs/bloco-1/MCP-FULFILLMENT-OPS-BLUEPRINT.md`, `docs/bloco-1/README.md`. |
| **F2 – Contratos OpenAPI/AsyncAPI** | Criar contratos HTTP e eventos. | Arquivos em `contracts/openapi/bloco-1-core/` e `contracts/nats-schemas/bloco-1-core/`. |
| **F3 – Domínio (Ledger + Fulfillment)** | Implementar entidades e políticas. | `internal/domain/ledger/*`, `internal/domain/fulfillment/*`, testes table-driven básicos. |
| **F4 – Casos de Uso e Serviços** | Implementar `internal/app` e `internal/services`. | Use cases core + serviços auxiliares com testes. |
| **F5 – Interfaces (HTTP/gRPC/CLI)** | Handlers, routers e clientes. | `internal/interfaces/http/*`, `internal/interfaces/grpc/*`, `cmd/*` wiring. |
| **F6 – MCP Generators e Ferramentas** | Restaurar `internal/mcp`, `tools/generators`, `tools/validators`. | Código gerador/validador funcional. |
| **F7 – Templates e Scripts** | Repor `templates/*` e `scripts/*`. | Templates Go/Web/TinyGo/MCP, scripts de setup/deploy/validation. |
| **F8 – Infra & Observabilidade** | Configs, prometheus/grafana, docker-compose. | `config/**/*`, `ops/prometheus/*`, `ops/grafana/*`, `ops/k6/*`. |
| **F9 – Testes, Automação e Validação Final** | Garantir qualidade e conformidade 100%. | `go test ./...`, execução de k6, atualização de relatórios/README, `validation-result.json` zerado. |

Cada fase depende da anterior liberando insumos (ex.: F2 depende da doc consolidada na F1). Depois da F9, reexecutar a validação completa e atualizar os indicadores de compliance.

### Status Atual das Fases

| Fase | Situação | Observações |
|------|----------|-------------|
| F0 | ✅ Concluída | Baseline capturado via `validation-result.json`. |
| F1 | ✅ Concluída | Blueprints do Core Inventory e Fulfillment Ops versionados em `docs/bloco-1/` e documentados em `docs/bloco-1/README.md`. |
| F2 | ✅ Concluída | Contratos OpenAPI e schemas NATS versionados em `contracts/`. |
| F3 | ✅ Concluída | Domínio `internal/domain/ledger` validado e `internal/domain/fulfillment` criado com testes básicos. |
| F4–F9 | ⏳ Pendente | Serão iniciadas sequencialmente após aprovação das fases anteriores. |

---

## 3. Detalhamento por Trilha

### 3.1 Contratos & Documentação

| Entrega                                 | Responsável | Critério de Aceite                                                       | Arquivo de Saída                                        |
|-----------------------------------------|-------------|--------------------------------------------------------------------------|---------------------------------------------------------|
| `core-inventory-v1.yaml` / `v2.yaml`    | API         | `POST /reserve`, `/confirm`, `/release`, `/adjust`, `/available` descritos com exemplos, códigos de erro e headers de Idempotency. | `contracts/openapi/bloco-1-core/core-inventory-v1.yaml`<br>`contracts/openapi/bloco-1-core/core-inventory-v2.yaml` |
| `fulfillment-ops-v1.yaml`               | API         | Endpoints inbound/outbound/transfer/returns/cycle_count com modelos de item, batching e estados. | `contracts/openapi/bloco-1-core/fulfillment-ops-v1.yaml` |
| `contracts/nats-schemas/*.json`         | Backend     | Schemas `inventory.reserve.*`, `fulfillment.*` e eventos `health.*` com validação via `jsonschema`. | `contracts/nats-schemas/bloco-1-core/*.json`             |
| `docs/bloco-1/*.md`                     | Arq/Prod    | Blueprint convertido para doc versionada com seções P0/P1/P2.             | `docs/bloco-1/MCP-CORE-INVENTORY-BLUEPRINT.md` etc.      |

### 3.2 Domínio & Serviços

| Entrega                        | Blocos                                 | Arquivos/Checkpoints                                                                                 |
|--------------------------------|----------------------------------------|------------------------------------------------------------------------------------------------------|
| `internal/domain/ledger/*`     | Entidades (`StockLedger`, `Reservation`, `Movement`, `Batch`).| `internal/domain/ledger/stock_ledger.go`, `reservation.go`, `stock_movement.go`, `policy.go`, `batches/*.go`. |
| `internal/app/*`               | Use cases (Reserve, Confirm, Release, Adjust, Query).        | `internal/app/reserve_stock.go`, `confirm_reservation.go`, `release_reservation.go`, `adjust_stock.go`, `query_available.go`. |
| `internal/services/*`          | Serviços agregadores (State, AI, Monitoring).               | `internal/services/state_app_service.go`, `internal/services/ai_app_service.go`, `internal/services/monitoring_app_service.go`.|
| `internal/mcp/*`               | MCP Generators, Validators, Handlers HTTP/gRPC/CLI.         | `internal/mcp/generators/*`, `internal/mcp/validators/*`, `internal/mcp/handlers/*`.                                                      |
| Testes de Domínio              | QA                                      | Pastas `internal/domain/ledger/*.test.go`, `internal/app/*_test.go` com cobertura ≥70%.                                                   |

### 3.3 Ferramentas & Templates

| Entrega                        | Observações                                          | Arquivos-Chave                                                                                               |
|--------------------------------|------------------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `tools/generators/*.go`        | Generators de configs, MCPs, templates.               | `tools/generators/mcp_generator.go`, `tools/generators/config_generator.go`, `tools/generators/code_generator.go`. |
| `tools/validators/*.go`        | Validadores de árvore, dependências, estrutura.       | `tools/validators/mcp_validator.go`, `tools/validators/dependency_validator.go`, `tools/validators/structure_validator.go`. |
| `templates/*`                  | Go/TinyGo/Web/MCP Premium com arquivos `.tmpl`.       | `templates/go/*`, `templates/web/*`, `templates/tinygo/*`, `templates/mcp-go-premium/*`.                        |
| `scripts/*`                    | Setup, deploy, validation, monitoring, optimization.  | `scripts/validation/validate_tree.sh`, `scripts/generation/generate_docs.sh`, `scripts/maintenance/backup.sh`.  |

### 3.4 Infraestrutura & Observabilidade

| Entrega                        | Observações                                          | Arquivos Necessários                                                                                |
|--------------------------------|------------------------------------------------------|------------------------------------------------------------------------------------------------------|
| `config/**/*.yaml`             | Núcleos: core, ai, state, security, monitoring.       | `config/mcp/tools.yaml`, `config/ai/learning.yaml`, `config/state/state_cache.yaml`, `config/security/rbac.yaml`. |
| `ops/prometheus/*`             | `prometheus.yml`, regras, dashboards.                 | `ops/prometheus/prometheus.yml`, `ops/prometheus-rules/bloco-1-core-rules.yaml`.                     |
| `ops/grafana/*`                | JSON dashboards (Core Inventory Health, Fulfillment). | `ops/grafana-dashboards/bloco-1-core/*.json`.                                                       |
| `ops/k6/*`                     | Scripts `reserve-load-test`, `available-read-test`.   | `ops/k6/bloco-1-core/reserve-load-test.js`, `ops/k6/bloco-1-core/available-read-test.js`.            |
| `docker-compose` especializados| Compose para Inventário, Fulfillment, Observabilidade.| `docker-compose.inventory.yaml`, `docker-compose.fulfillment.yaml`, `docker-compose.ops.yaml`.       |

---

## 4. Timeline Sugerido

| Semana | Objetos Principais                           | Métrica de Saída                              |
|--------|----------------------------------------------|-----------------------------------------------|
| S0     | Fase 0 + 1                                   | OpenAPI v1/v2 + docs publicados               |
| S1     | Fase 2 (Domínio + Serviços)                  | 100% dos `internal/domain/app/services` criados|
| S2     | Fase 3 (Ferramentas + Templates)             | `tools/*`, `templates/*`, `scripts/*` validados|
| S3     | Fase 4 + 5                                   | Observabilidade ativa + `validation_tree` verde|

---

## 5. Checklist de Conclusão

- [ ] `validation-result.json` → `missing_files == 0`
- [ ] `go test ./...` → sucesso (CI)
- [ ] `ops/k6/bloco-1-core/reserve-load-test.js` executado contra serviço Docker
- [ ] Prometheus + Grafana + Alertas rodando (`docker-compose.ops.yaml`)
- [ ] Blueprints convertidos para `docs/bloco-1/*` e linkados no README
- [ ] Contratos OpenAPI/AsyncAPI versionados sob `contracts/`
- [ ] Pact tests publicados no broker (quando configurado)

> Após todas as etapas, atualizar `validation-report.md`, `README-BLOCO-1.md` e `DEPLOYMENT-SUMMARY.md` com a nova conformidade.

---

## 6. Referências

- `validation-result.json`
- `BLOCO-1-BLUEPRINT-MCP-CORE-INVENTORY.md`
- `BLOCO-1-BLUEPRINT-MCP-FULFILLMENT-OPS.md`
- `README-BLOCO-1.md`

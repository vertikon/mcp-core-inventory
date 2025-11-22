# NATS Subjects Documentation

## Overview
Este documento descreve todos os subjects NATS utilizados no mcp-core-inventory para comunicação entre serviços.

## Subjects por Categoria

### 1. Engine Management
- `engine.task.submit` - Submissão de tarefas para execução
- `engine.task.schedule` - Agendamento de tarefas
- `engine.task.cancel` - Cancelamento de tarefas
- `engine.status` - Status do engine de execução
- `engine.metrics` - Métricas de performance

### 2. MCP Protocol
- `mcp.tool.register` - Registro de ferramentas MCP
- `mcp.tool.unregister` - Remoção de ferramentas MCP
- `mcp.tool.execute` - Execução de ferramentas MCP
- `mcp.server.start` - Inicialização de servidor MCP
- `mcp.server.stop` - Parada de servidor MCP
- `mcp.client.connect` - Conexão de cliente MCP
- `mcp.client.disconnect` - Desconexão de cliente MCP

### 3. Project Management
- `project.create` - Criação de novos projetos
- `project.update` - Atualização de projetos existentes
- `project.delete` - Exclusão de projetos
- `project.list` - Listagem de projetos
- `project.get` - Obtenção de detalhes de projeto
- `project.validate` - Validação de estrutura de projeto

### 4. Template Management
- `template.create` - Criação de templates
- `template.update` - Atualização de templates
- `template.delete` - Exclusão de templates
- `template.list` - Listagem de templates
- `template.get` - Obtenção de detalhes de template

### 5. Registry Management
- `registry.register` - Registro de entidades
- `registry.unregister` - Remoção de registros
- `registry.list` - Listagem de registros
- `registry.get` - Obtenção de detalhes de registro
- `registry.update` - Atualização de registros

### 6. Cache Management
- `cache.set` - Armazenamento de dados no cache
- `cache.get` - Obtenção de dados do cache
- `cache.delete` - Exclusão de dados do cache
- `cache.clear` - Limpeza completa do cache
- `cache.stats` - Estatísticas de uso do cache

### 7. Monitoring & Observability
- `monitoring.health` - Verificação de saúde de serviços
- `monitoring.metrics` - Coleta de métricas
- `monitoring.alerts` - Gerenciamento de alertas
- `monitoring.traces` - Rastreamento distribuído
- `monitoring.logs` - Coleta de logs estruturados

### 8. AI & Knowledge Management
- `ai.knowledge.index` - Indexação de conhecimento
- `ai.knowledge.search` - Busca no conhecimento
- `ai.knowledge.update` - Atualização de conhecimento
- `ai.memory.store` - Armazenamento de memória
- `ai.memory.retrieve` - Recuperação de memória
- `ai.model.train` - Treinamento de modelos
- `ai.model.predict` - Predição com modelos

### 9. Security & Authentication
- `security.auth.login` - Autenticação de usuários
- `security.auth.logout` - Logout de usuários
- `security.auth.validate` - Validação de tokens
- `security.audit.log` - Registro de auditoria
- `security.permission.check` - Verificação de permissões

### 10. Configuration Management
- `config.update` - Atualização de configurações
- `config.reload` - Recarregamento de configurações
- `config.validate` - Validação de configurações
- `config.backup` - Backup de configurações
- `config.restore` - Restauração de configurações

## Formato dos Subjects

### Naming Convention
```
{category}.{entity}.{action}
{category}.{subsystem}.{action}
{service}.{operation}
```

### Exemplos de Uso
```go
// Publicar tarefa no engine
nc.Publish("engine.task.submit", taskData)

// Solicitar execução de ferramenta MCP
nc.Publish("mcp.tool.execute", toolRequest)

// Registrar novo projeto
nc.Publish("project.create", projectData)

// Armazenar no cache
nc.Publish("cache.set", cacheData)
```

## Considerações

1. **Wildcard Subscriptions**: Use `>` para múltiplos níveis (ex: `engine.>`)
2. **Reply Subjects**: Use `.reply` para respostas (ex: `mcp.tool.execute.reply`)
3. **Error Handling**: Use `.error` para falhas (ex: `project.create.error`)
4. **Batching**: Agrupe mensagens relacionadas para eficiência
5. **QoS**: Configure Quality of Service por tipo de mensagem

## Implementação

### Go Example
```go
// Publisher
func PublishTaskSubmit(nc *nats.Conn, task *Task) error {
    data := map[string]interface{}{
        "task_id": task.ID,
        "type": task.Type,
        "payload": task.Data,
    }
    
    return nc.Publish("engine.task.submit", data)
}

// Subscriber
func SubscribeToEngineEvents(nc *nats.Conn) (*nats.Subscription, error) {
    return nc.Subscribe("engine.>", func(msg *nats.Msg) {
        handleEngineEvent(msg)
    })
}
```

---
*Gerado automaticamente pelo mcp-core-inventory*

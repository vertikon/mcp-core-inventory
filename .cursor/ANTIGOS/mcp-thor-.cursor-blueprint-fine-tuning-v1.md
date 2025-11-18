

# Proposta de Fine-tuning para GLM-4.6 com Conhecimento Local de Negócio

## 📋 Visão Geral

Com base na análise da documentação do MCP Thor, proponho uma arquitetura para fine-tuning do GLM-4.6 que integre conhecimento local de negócio e persista entre instâncias, aproveitando a estrutura já existente no projeto.

## 🏗️ Arquitetura Proposta

### 1. Estrutura de Armazenamento de Conhecimento

```
mcp-thor/
├── internal/
│   ├── ai/
│   │   ├── finetuning/              # Lógica de fine-tuning
│   │   │   ├── knowledge_store.go   # Armazenamento de conhecimento
│   │   │   ├── prompt_builder.go    # Construtor de prompts com contexto
│   │   │   └── memory_manager.go    # Gerenciamento de memória persistente
│   │   ├── specialist/              # Especialistas de IA existentes
│   │   └── local_knowledge/         # NOVO: Conhecimento local de negócio
│   │       ├── domain_knowledge.go  # Conhecimento do domínio
│   │       ├── business_rules.go    # Regras de negócio
│   │       └── context_cache.go     # Cache de contexto
│   └── infrastructure/
│       └── persistence/              # NOVO: Persistência de conhecimento
│           ├── vector_store.go       # Armazenamento vetorial
│           └── knowledge_db.go       # Banco de conhecimento
```

### 2. Componente de Conhecimento Local

```go
// internal/ai/local_knowledge/domain_knowledge.go
package local_knowledge

import (
    "context"
    "encoding/json"
    "time"
)

type DomainKnowledge struct {
    ID          string                 `json:"id"`
    Domain      string                 `json:"domain"`
    Concepts    map[string]Concept     `json:"concepts"`
    Rules       []BusinessRule         `json:"rules"`
    Examples    []Example              `json:"examples"`
    LastUpdated time.Time              `json:"last_updated"`
}

type Concept struct {
    Name        string   `json:"name"`
    Definition  string   `json:"definition"`
    Attributes  []string `json:"attributes"`
    Relations   []string `json:"relations"`
}

type BusinessRule struct {
    ID          string `json:"id"`
    Description string `json:"description"`
    Condition   string `json:"condition"`
    Action      string `json:"action"`
    Priority    int    `json:"priority"`
}

type Example struct {
    Input       string `json:"input"`
    Output      string `json:"output"`
    Explanation string `json:"explanation"`
}

type KnowledgeStore interface {
    Store(ctx context.Context, knowledge *DomainKnowledge) error
    Retrieve(ctx context.Context, domain string) (*DomainKnowledge, error)
    Update(ctx context.Context, knowledge *DomainKnowledge) error
    Search(ctx context.Context, query string) ([]*DomainKnowledge, error)
}
```

### 3. Gerenciador de Memória Persistente

```go
// internal/ai/finetuning/memory_manager.go
package finetuning

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/vertikon/mcp-thor/internal/ai/local_knowledge"
    "github.com/vertikon/mcp-thor/internal/infrastructure/persistence"
)

type MemoryManager struct {
    store     local_knowledge.KnowledgeStore
    vectorDB  persistence.VectorStore
    cache     persistence.Cache
    maxSize   int
    ttl       time.Duration
}

func NewMemoryManager(store local_knowledge.KnowledgeStore, vectorDB persistence.VectorStore, cache persistence.Cache) *MemoryManager {
    return &MemoryManager{
        store:    store,
        vectorDB: vectorDB,
        cache:    cache,
        maxSize:  1000, // Máximo de entradas na memória
        ttl:      24 * time.Hour, // TTL de 24 horas
    }
}

func (m *MemoryManager) StoreInteraction(ctx context.Context, domain, input, output string) error {
    // Criar exemplo de interação
    example := local_knowledge.Example{
        Input:       input,
        Output:      output,
        Explanation: "Interação do usuário com assistente de IA",
    }
    
    // Recuperar conhecimento existente
    knowledge, err := m.store.Retrieve(ctx, domain)
    if err != nil {
        // Criar novo conhecimento se não existir
        knowledge = &local_knowledge.DomainKnowledge{
            ID:          domain,
            Domain:      domain,
            Concepts:    make(map[string]local_knowledge.Concept),
            Rules:       []local_knowledge.BusinessRule{},
            Examples:    []local_knowledge.Example{example},
            LastUpdated: time.Now(),
        }
    } else {
        // Adicionar exemplo ao conhecimento existente
        knowledge.Examples = append(knowledge.Examples, example)
        knowledge.LastUpdated = time.Now()
    }
    
    // Armazenar conhecimento atualizado
    return m.store.Update(ctx, knowledge)
}

func (m *MemoryManager) GetRelevantContext(ctx context.Context, domain, query string) (string, error) {
    // Verificar cache primeiro
    cacheKey := fmt.Sprintf("context:%s:%s", domain, query)
    if cached, err := m.cache.Get(ctx, cacheKey); err == nil {
        return cached, nil
    }
    
    // Buscar conhecimento relevante
    knowledge, err := m.store.Retrieve(ctx, domain)
    if err != nil {
        return "", err
    }
    
    // Construir contexto a partir do conhecimento
    contextBuilder := NewPromptBuilder()
    contextBuilder.AddDomainKnowledge(knowledge)
    
    // Buscar exemplos similares usando busca vetorial
    similarExamples, err := m.vectorDB.SearchSimilar(ctx, query, 5)
    if err == nil {
        contextBuilder.AddExamples(similarExamples)
    }
    
    context := contextBuilder.Build()
    
    // Armazenar em cache
    m.cache.Set(ctx, cacheKey, context, m.ttl)
    
    return context, nil
}
```

### 4. Construtor de Prompts com Contexto

```go
// internal/ai/finetuning/prompt_builder.go
package finetuning

import (
    "fmt"
    "strings"
    
    "github.com/vertikon/mcp-thor/internal/ai/local_knowledge"
)

type PromptBuilder struct {
    sections []string
}

func NewPromptBuilder() *PromptBuilder {
    return &PromptBuilder{
        sections: make([]string, 0),
    }
}

func (pb *PromptBuilder) AddSystemPrompt(prompt string) *PromptBuilder {
    pb.sections = append(pb.sections, fmt.Sprintf("### Sistema\n%s", prompt))
    return pb
}

func (pb *PromptBuilder) AddDomainKnowledge(knowledge *local_knowledge.DomainKnowledge) *PromptBuilder {
    concepts := make([]string, 0, len(knowledge.Concepts))
    for name, concept := range knowledge.Concepts {
        concepts = append(concepts, fmt.Sprintf("- %s: %s", name, concept.Definition))
    }
    
    rules := make([]string, 0, len(knowledge.Rules))
    for _, rule := range knowledge.Rules {
        rules = append(rules, fmt.Sprintf("- %s: %s", rule.ID, rule.Description))
    }
    
    section := fmt.Sprintf(
        "### Conhecimento do Domínio: %s\n\n**Conceitos:**\n%s\n\n**Regras de Negócio:**\n%s",
        knowledge.Domain,
        strings.Join(concepts, "\n"),
        strings.Join(rules, "\n"),
    )
    
    pb.sections = append(pb.sections, section)
    return pb
}

func (pb *PromptBuilder) AddExamples(examples []local_knowledge.Example) *PromptBuilder {
    examplesText := make([]string, 0, len(examples))
    for i, example := range examples {
        examplesText = append(examplesText, fmt.Sprintf(
            "**Exemplo %d:**\nEntrada: %s\nSaída: %s\nExplicação: %s",
            i+1, example.Input, example.Output, example.Explanation,
        ))
    }
    
    section := fmt.Sprintf("### Exemplos Relevantes\n%s", strings.Join(examplesText, "\n\n"))
    pb.sections = append(pb.sections, section)
    return pb
}

func (pb *PromptBuilder) AddUserQuery(query string) *PromptBuilder {
    pb.sections = append(pb.sections, fmt.Sprintf("### Consulta do Usuário\n%s", query))
    return pb
}

func (pb *PromptBuilder) Build() string {
    return strings.Join(pb.sections, "\n\n")
}
```

### 5. Integração com GLM-4.6

```go
// internal/ai/specialist/glm_specialist.go
package specialist

import (
    "context"
    "fmt"
    
    "github.com/vertikon/mcp-thor/internal/ai/finetuning"
    "github.com/vertikon/mcp-thor/internal/ai/local_knowledge"
    "github.com/vertikon/mcp-thor/internal/config"
)

type GLMSpecialist struct {
    client         GLMClient
    memoryManager  *finetuning.MemoryManager
    config         *config.AIConfig
    systemPrompt   string
}

func NewGLMSpecialist(client GLMClient, memoryManager *finetuning.MemoryManager, config *config.AIConfig) *GLMSpecialist {
    return &GLMSpecialist{
        client:        client,
        memoryManager: memoryManager,
        config:        config,
        systemPrompt:  "Você é um assistente especializado no ecossistema Vertikon. Use o conhecimento fornecido para responder às perguntas do usuário.",
    }
}

func (s *GLMSpecialist) ProcessRequest(ctx context.Context, domain, query string) (string, error) {
    // Obter contexto relevante do conhecimento local
    context, err := s.memoryManager.GetRelevantContext(ctx, domain, query)
    if err != nil {
        // Continuar sem contexto se houver erro
        context = ""
    }
    
    // Construir prompt com contexto
    promptBuilder := finetuning.NewPromptBuilder()
    promptBuilder.
        AddSystemPrompt(s.systemPrompt).
        AddUserQuery(query)
    
    if context != "" {
        promptBuilder.AddDomainKnowledge(&local_knowledge.DomainKnowledge{
            Domain: domain,
        })
        promptBuilder.sections = append([]string{context}, promptBuilder.sections...)
    }
    
    prompt := promptBuilder.Build()
    
    // Enviar requisição para GLM-4.6
    response, err := s.client.CreateChatCompletion(ctx, GLMRequest{
        Model:       s.config.Model,
        Messages:    []GLMMessage{{Role: "user", Content: prompt}},
        Temperature: s.config.Temperature,
        MaxTokens:   s.config.MaxTokens,
    })
    
    if err != nil {
        return "", fmt.Errorf("erro ao chamar GLM API: %w", err)
    }
    
    // Armazenar interação para aprendizado futuro
    err = s.memoryManager.StoreInteraction(ctx, domain, query, response.Content)
    if err != nil {
        // Log erro, mas não falhar a requisição
        fmt.Printf("Erro ao armazenar interação: %v\n", err)
    }
    
    return response.Content, nil
}
```

### 6. Armazenamento Vetorial para Busca Semântica

```go
// internal/infrastructure/persistence/vector_store.go
package persistence

import (
    "context"
    "encoding/json"
    
    "github.com/vertikon/mcp-thor/internal/ai/local_knowledge"
)

type VectorStore interface {
    Store(ctx context.Context, id string, vector []float32, metadata map[string]interface{}) error
    Search(ctx context.Context, queryVector []float32, limit int) ([]VectorSearchResult, error)
    SearchByText(ctx context.Context, text string, limit int) ([]VectorSearchResult, error)
}

type VectorSearchResult struct {
    ID       string                 `json:"id"`
    Score    float32                `json:"score"`
    Metadata map[string]interface{} `json:"metadata"`
}

type InMemoryVectorStore struct {
    vectors map[string][]float32
    metadata map[string]map[string]interface{}
}

func NewInMemoryVectorStore() *InMemoryVectorStore {
    return &InMemoryVectorStore{
        vectors:  make(map[string][]float32),
        metadata: make(map[string]map[string]interface{}),
    }
}

func (s *InMemoryVectorStore) Store(ctx context.Context, id string, vector []float32, metadata map[string]interface{}) error {
    s.vectors[id] = vector
    s.metadata[id] = metadata
    return nil
}

func (s *InMemoryVectorStore) Search(ctx context.Context, queryVector []float32, limit int) ([]VectorSearchResult, error) {
    // Implementação simplificada de busca vetorial
    results := make([]VectorSearchResult, 0)
    
    for id, vector := range s.vectors {
        // Calcular similaridade de cosseno (simplificado)
        score := cosineSimilarity(queryVector, vector)
        
        results = append(results, VectorSearchResult{
            ID:       id,
            Score:    score,
            Metadata: s.metadata[id],
        })
    }
    
    // Ordenar por score e limitar resultados
    // Implementação simplificada
    if len(results) > limit {
        results = results[:limit]
    }
    
    return results, nil
}

func (s *InMemoryVectorStore) SearchByText(ctx context.Context, text string, limit int) ([]VectorSearchResult, error) {
    // Em uma implementação real, converteríamos o texto em um vetor
    // Para simplificação, vamos buscar por metadados
    results := make([]VectorSearchResult, 0)
    
    for id, metadata := range s.metadata {
        if textMatch, ok := metadata["text"].(string); ok && contains(textMatch, text) {
            results = append(results, VectorSearchResult{
                ID:       id,
                Score:    1.0, // Score perfeito para correspondência exata
                Metadata: metadata,
            })
        }
    }
    
    if len(results) > limit {
        results = results[:limit]
    }
    
    return results, nil
}

// Funções auxiliares (implementações simplificadas)
func cosineSimilarity(a, b []float32) float32 {
    if len(a) != len(b) {
        return 0
    }
    
    var dotProduct, normA, normB float32
    for i := 0; i < len(a); i++ {
        dotProduct += a[i] * b[i]
        normA += a[i] * a[i]
        normB += b[i] * b[i]
    }
    
    if normA == 0 || normB == 0 {
        return 0
    }
    
    return dotProduct / (float32(math.Sqrt(float64(normA))) * float32(math.Sqrt(float64(normB))))
}

func contains(text, substr string) bool {
    return strings.Contains(strings.ToLower(text), strings.ToLower(substr))
}
```

### 7. CLI para Gerenciamento de Conhecimento

```go
// internal/cli/knowledge/knowledge.go
package knowledge

import (
    "context"
    "fmt"
    "os"
    
    "github.com/spf13/cobra"
    "github.com/vertikon/mcp-thor/internal/ai/local_knowledge"
    "github.com/vertikon/mcp-thor/internal/config"
)

var knowledgeCmd = &cobra.Command{
    Use:   "knowledge",
    Short: "Gerenciar conhecimento local de negócio",
    Long:  "Comandos para gerenciar o conhecimento local de negócio usado pelo GLM-4.6",
}

var importCmd = &cobra.Command{
    Use:   "import [file]",
    Short: "Importar conhecimento de um arquivo",
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        filePath := args[0]
        
        // Ler arquivo
        data, err := os.ReadFile(filePath)
        if err != nil {
            return fmt.Errorf("erro ao ler arquivo: %w", err)
        }
        
        // Desserializar conhecimento
        var knowledge local_knowledge.DomainKnowledge
        if err := json.Unmarshal(data, &knowledge); err != nil {
            return fmt.Errorf("erro ao desserializar conhecimento: %w", err)
        }
        
        // Carregar configuração
        cfg, err := config.Load()
        if err != nil {
            return fmt.Errorf("erro ao carregar configuração: %w", err)
        }
        
        // Inicializar store de conhecimento
        store, err := initKnowledgeStore(cfg)
        if err != nil {
            return fmt.Errorf("erro ao inicializar store: %w", err)
        }
        
        // Armazenar conhecimento
        ctx := context.Background()
        if err := store.Store(ctx, &knowledge); err != nil {
            return fmt.Errorf("erro ao armazenar conhecimento: %w", err)
        }
        
        fmt.Printf("Conhecimento do domínio '%s' importado com sucesso\n", knowledge.Domain)
        return nil
    },
}

var exportCmd = &cobra.Command{
    Use:   "export [domain] [file]",
    Short: "Exportar conhecimento de um domínio para um arquivo",
    Args:  cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
        domain := args[0]
        filePath := args[1]
        
        // Carregar configuração
        cfg, err := config.Load()
        if err != nil {
            return fmt.Errorf("erro ao carregar configuração: %w", err)
        }
        
        // Inicializar store de conhecimento
        store, err := initKnowledgeStore(cfg)
        if err != nil {
            return fmt.Errorf("erro ao inicializar store: %w", err)
        }
        
        // Recuperar conhecimento
        ctx := context.Background()
        knowledge, err := store.Retrieve(ctx, domain)
        if err != nil {
            return fmt.Errorf("erro ao recuperar conhecimento: %w", err)
        }
        
        // Serializar conhecimento
        data, err := json.MarshalIndent(knowledge, "", "  ")
        if err != nil {
            return fmt.Errorf("erro ao serializar conhecimento: %w", err)
        }
        
        // Escrever arquivo
        if err := os.WriteFile(filePath, data, 0644); err != nil {
            return fmt.Errorf("erro ao escrever arquivo: %w", err)
        }
        
        fmt.Printf("Conhecimento do domínio '%s' exportado para '%s'\n", domain, filePath)
        return nil
    },
}

func init() {
    knowledgeCmd.AddCommand(importCmd)
    knowledgeCmd.AddCommand(exportCmd)
}
```

## 🔧 Configuração

### 1. Configuração de IA

```yaml
# config/ai.yaml
ai:
  provider: "glm"
  model: "glm-4.6-z.ai"
  api_key: "${GLM_API_KEY}"
  endpoint: "https://api.z.ai/v1"
  temperature: 0.3
  max_tokens: 4000
  
  knowledge:
    store_type: "postgres"  # postgres, mongodb, file
    vector_store_type: "chroma"  # chroma, pinecone, in_memory
    cache_ttl: "24h"
    max_memory_entries: 1000
    
  fine_tuning:
    enabled: true
    auto_update: true
    update_interval: "6h"
    min_examples_for_update: 10
```

### 2. Configuração de Persistência

```yaml
# config/persistence.yaml
persistence:
  postgres:
    host: "${DB_HOST}"
    port: "${DB_PORT}"
    user: "${DB_USER}"
    password: "${DB_PASSWORD}"
    database: "${DB_NAME}"
    sslmode: "require"
    
  vector_store:
    chroma:
      host: "${CHROMA_HOST}"
      port: "${CHROMA_PORT}"
      collection: "mcp-knowledge"
      
  cache:
    redis:
      host: "${REDIS_HOST}"
      port: "${REDIS_PORT}"
      password: "${REDIS_PASSWORD}"
      db: 1
```

## 🚀 Implementação

### 1. Script de Setup

```powershell
# scripts/setup-ai-knowledge.ps1
param(
    [Parameter(Mandatory=$true)]
    [string]$InstallPath,
    
    [Parameter(Mandatory=$false)]
    [string]$ConfigPath = ".\config\ai.yaml"
)

Write-Host "Configurando ambiente de conhecimento local para GLM-4.6..." -ForegroundColor Green

# Criar diretórios necessários
New-Item -ItemType Directory -Force -Path "$InstallPath\internal\ai\local_knowledge"
New-Item -ItemType Directory -Force -Path "$InstallPath\internal\ai\finetuning"
New-Item -ItemType Directory -Force -Path "$InstallPath\internal\infrastructure\persistence"

# Copiar arquivos de configuração
Copy-Item -Path ".\config\ai.yaml" -Destination "$InstallPath\config\ai.yaml" -Force
Copy-Item -Path ".\config\persistence.yaml" -Destination "$InstallPath\config\persistence.yaml" -Force

# Instalar dependências Go
Write-Host "Instalando dependências Go..." -ForegroundColor Yellow
go get github.com/chromedp/chromedp
go get github.com/go-redis/redis/v8
go get github.com/lib/pq
go get github.com/jackc/pgx/v4

Write-Host "Configuração concluída com sucesso!" -ForegroundColor Green
```

### 2. Exemplo de Conhecimento de Domínio

```json
// examples/inventory-knowledge.json
{
  "id": "inventory",
  "domain": "inventory",
  "concepts": {
    "produto": {
      "name": "produto",
      "definition": "Item mantido em estoque para venda ou uso",
      "attributes": ["sku", "nome", "descrição", "preço", "quantidade"],
      "relations": ["categoria", "fornecedor", "localização"]
    },
    "reserva": {
      "name": "reserva",
      "definition": "Alocação temporária de itens de estoque para um pedido",
      "attributes": ["id_pedido", "id_produto", "quantidade", "validade"],
      "relations": ["pedido", "produto"]
    }
  },
  "rules": [
    {
      "id": "inventory-001",
      "description": "Não é possível reservar mais itens do que o disponível em estoque",
      "condition": "quantidade_reservada > quantidade_disponível",
      "action": "rejeitar_reserva",
      "priority": 1
    },
    {
      "id": "inventory-002",
      "description": "Reservas expiram após 24 horas se não confirmadas",
      "condition": "tempo_reserva > 24h AND status = 'pendente'",
      "action": "liberar_reserva",
      "priority": 2
    }
  ],
  "examples": [
    {
      "input": "Como faço para reservar 10 unidades do produto SKU-123?",
      "output": "Para reservar 10 unidades do produto SKU-123, você deve usar o endpoint POST /api/v1/reservas com o corpo da requisição contendo o SKU do produto e a quantidade desejada. O sistema verificará a disponibilidade em estoque antes de confirmar a reserva.",
      "explanation": "O processo de reserva envolve verificar a disponibilidade e alocar temporariamente os itens."
    }
  ],
  "last_updated": "2025-11-15T12:00:00Z"
}
```

## 📋 Benefícios

1. **Conhecimento Persistente**: O sistema mantém o conhecimento adquirido entre sessões
2. **Contexto de Negócio**: Incorpora conhecimento específico do domínio do Vertikon
3. **Aprendizado Contínuo**: Aprende com as interações para melhorar respostas futuras
4. **Busca Semântica**: Utiliza armazenamento vetorial para encontrar exemplos relevantes
5. **Integração Transparente**: Funciona com a estrutura existente do MCP Thor

## 🔄 Fluxo de Trabalho

1. **Inicialização**: Carrega conhecimento existente do armazenamento persistente
2. **Consulta**: Recebe uma consulta do usuário
3. **Contextualização**: Busca conhecimento relevante usando busca semântica
4. **Processamento**: Envia prompt enriquecido com contexto para GLM-4.6
5. **Resposta**: Retorna resposta ao usuário
6. **Aprendizado**: Armazena interação para aprendizado futuro

Essa abordagem permite que o GLM-4.6 tenha conhecimento local do negócio Vertikon e persista esse conhecimento entre instâncias, melhorando continuamente suas respostas com base nas interações anteriores.
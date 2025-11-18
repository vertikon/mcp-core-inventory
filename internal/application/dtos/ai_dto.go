package dtos

// ChatRequest represents a chat request
type ChatRequest struct {
	Message     string                 `json:"message" validate:"required"`
	Context     []string               `json:"context"`
	Model       string                 `json:"model"`
	Temperature float64                `json:"temperature"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// ChatResponse represents a chat response
type ChatResponse struct {
	Response string                 `json:"response"`
	Model    string                 `json:"model"`
	Tokens   int                    `json:"tokens"`
	Metadata map[string]interface{} `json:"metadata"`
}

// GenerateRequest represents a generation request
type GenerateRequest struct {
	Prompt      string                 `json:"prompt" validate:"required"`
	Model       string                 `json:"model"`
	Temperature float64                `json:"temperature"`
	MaxTokens   int                    `json:"max_tokens"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GenerateResponse represents a generation response
type GenerateResponse struct {
	Content  string                 `json:"content"`
	Model    string                 `json:"model"`
	Tokens   int                    `json:"tokens"`
	Metadata map[string]interface{} `json:"metadata"`
}

// EmbedRequest represents an embedding request
type EmbedRequest struct {
	Text     string                 `json:"text" validate:"required"`
	Model    string                 `json:"model"`
	Metadata map[string]interface{} `json:"metadata"`
}

// EmbedResponse represents an embedding response
type EmbedResponse struct {
	Embedding []float64              `json:"embedding"`
	Model     string                 `json:"model"`
	Metadata  map[string]interface{} `json:"metadata"`
}

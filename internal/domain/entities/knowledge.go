// Package entities provides domain entities
package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Knowledge represents a knowledge entity for AI/RAG
type Knowledge struct {
	id          string
	name        string
	description string
	documents   []*Document
	embeddings  map[string]*Embedding
	version     int
	createdAt   time.Time
	updatedAt   time.Time
}

// Document represents a document in the knowledge base
type Document struct {
	id        string
	content   string
	metadata  map[string]interface{}
	createdAt time.Time
}

// Embedding represents an embedding vector for a document
type Embedding struct {
	documentID string
	vector     []float64
	dimension  int
	model      string
	createdAt  time.Time
}

// NewKnowledge creates a new Knowledge entity
func NewKnowledge(name string, description string) (*Knowledge, error) {
	if name == "" {
		return nil, fmt.Errorf("knowledge name cannot be empty")
	}

	now := time.Now()
	return &Knowledge{
		id:          uuid.New().String(),
		name:        name,
		description: description,
		documents:   make([]*Document, 0),
		embeddings:  make(map[string]*Embedding),
		version:     1,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ID returns the knowledge ID
func (k *Knowledge) ID() string {
	return k.id
}

// Name returns the knowledge name
func (k *Knowledge) Name() string {
	return k.name
}

// Description returns the knowledge description
func (k *Knowledge) Description() string {
	return k.description
}

// Version returns the knowledge version
func (k *Knowledge) Version() int {
	return k.version
}

// Documents returns a copy of the documents list
func (k *Knowledge) Documents() []*Document {
	documents := make([]*Document, len(k.documents))
	copy(documents, k.documents)
	return documents
}

// Embeddings returns a copy of the embeddings map
func (k *Knowledge) Embeddings() map[string]*Embedding {
	embeddings := make(map[string]*Embedding)
	for k, v := range k.embeddings {
		embeddings[k] = &Embedding{
			documentID: v.documentID,
			vector:     append([]float64{}, v.vector...),
			dimension:  v.dimension,
			model:      v.model,
			createdAt:  v.createdAt,
		}
	}
	return embeddings
}

// CreatedAt returns the creation timestamp
func (k *Knowledge) CreatedAt() time.Time {
	return k.createdAt
}

// UpdatedAt returns the last update timestamp
func (k *Knowledge) UpdatedAt() time.Time {
	return k.updatedAt
}

// AddDocument adds a document to the knowledge base
func (k *Knowledge) AddDocument(content string, metadata map[string]interface{}) (*Document, error) {
	if content == "" {
		return nil, fmt.Errorf("document content cannot be empty")
	}

	doc := &Document{
		id:        uuid.New().String(),
		content:   content,
		metadata:  copyMetadata(metadata),
		createdAt: time.Now(),
	}

	k.documents = append(k.documents, doc)
	k.touch()
	return doc, nil
}

// RemoveDocument removes a document by ID
func (k *Knowledge) RemoveDocument(documentID string) error {
	for i, doc := range k.documents {
		if doc.id == documentID {
			k.documents = append(k.documents[:i], k.documents[i+1:]...)
			delete(k.embeddings, documentID)
			k.touch()
			return nil
		}
	}
	return fmt.Errorf("document %s not found", documentID)
}

// AddEmbedding adds an embedding for a document
func (k *Knowledge) AddEmbedding(documentID string, vector []float64, model string) error {
	if documentID == "" {
		return fmt.Errorf("document ID cannot be empty")
	}
	if len(vector) == 0 {
		return fmt.Errorf("embedding vector cannot be empty")
	}
	if model == "" {
		return fmt.Errorf("embedding model cannot be empty")
	}

	// Verify document exists
	found := false
	for _, doc := range k.documents {
		if doc.id == documentID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("document %s not found", documentID)
	}

	k.embeddings[documentID] = &Embedding{
		documentID: documentID,
		vector:     append([]float64{}, vector...),
		dimension:  len(vector),
		model:      model,
		createdAt:  time.Now(),
	}
	k.touch()
	return nil
}

// GetEmbedding retrieves an embedding by document ID
func (k *Knowledge) GetEmbedding(documentID string) (*Embedding, error) {
	embedding, exists := k.embeddings[documentID]
	if !exists {
		return nil, fmt.Errorf("embedding for document %s not found", documentID)
	}
	return &Embedding{
		documentID: embedding.documentID,
		vector:     append([]float64{}, embedding.vector...),
		dimension:  embedding.dimension,
		model:      embedding.model,
		createdAt:  embedding.createdAt,
	}, nil
}

// IncrementVersion increments the knowledge version
func (k *Knowledge) IncrementVersion() {
	k.version++
	k.touch()
}

// touch updates the updatedAt timestamp
func (k *Knowledge) touch() {
	k.updatedAt = time.Now()
}

// DocumentID returns the document ID
func (d *Document) ID() string {
	return d.id
}

// Content returns the document content
func (d *Document) Content() string {
	return d.content
}

// Metadata returns a copy of the document metadata
func (d *Document) Metadata() map[string]interface{} {
	return copyMetadata(d.metadata)
}

// CreatedAt returns the document creation timestamp
func (d *Document) CreatedAt() time.Time {
	return d.createdAt
}

// DocumentID returns the embedding document ID
func (e *Embedding) DocumentID() string {
	return e.documentID
}

// Vector returns a copy of the embedding vector
func (e *Embedding) Vector() []float64 {
	return append([]float64{}, e.vector...)
}

// Dimension returns the embedding dimension
func (e *Embedding) Dimension() int {
	return e.dimension
}

// Model returns the embedding model
func (e *Embedding) Model() string {
	return e.model
}

// CreatedAt returns the embedding creation timestamp
func (e *Embedding) CreatedAt() time.Time {
	return e.createdAt
}

// copyMetadata creates a deep copy of metadata map
func copyMetadata(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return nil
	}
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

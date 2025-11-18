package knowledge

import (
	"context"
	"fmt"
)

// KnowledgeGraph manages knowledge graph operations
type KnowledgeGraph struct {
	graphClient GraphClient
}

// NewKnowledgeGraph creates a new knowledge graph manager
func NewKnowledgeGraph(graphClient GraphClient) *KnowledgeGraph {
	return &KnowledgeGraph{
		graphClient: graphClient,
	}
}

// CreateEntity creates a new entity in the knowledge graph
func (kg *KnowledgeGraph) CreateEntity(ctx context.Context, entityID string, entityType string, properties map[string]interface{}) error {
	if kg.graphClient == nil {
		return fmt.Errorf("graph client not available")
	}

	props := copyMetadataKG(properties)
	props["type"] = entityType

	return kg.graphClient.CreateNode(ctx, "entities", entityID, props)
}

// CreateRelation creates a relation between entities
func (kg *KnowledgeGraph) CreateRelation(ctx context.Context, fromID string, toID string, relationType string, properties map[string]interface{}) error {
	if kg.graphClient == nil {
		return fmt.Errorf("graph client not available")
	}

	return kg.graphClient.CreateEdge(ctx, fromID, toID, relationType, properties)
}

// Traverse traverses the graph starting from a node
func (kg *KnowledgeGraph) Traverse(ctx context.Context, startID string, relationTypes []string, depth int, limit int) ([]*GraphNode, error) {
	if kg.graphClient == nil {
		return nil, fmt.Errorf("graph client not available")
	}

	if limit <= 0 {
		limit = 10
	}
	if depth <= 0 {
		depth = 2
	}

	// Build Cypher query for traversal
	cypher := kg.buildTraversalQuery(startID, relationTypes, depth, limit)
	params := map[string]interface{}{
		"start_id": startID,
		"limit":    limit,
	}

	results, err := kg.graphClient.Query(ctx, cypher, params)
	if err != nil {
		return nil, fmt.Errorf("failed to traverse graph: %w", err)
	}

	// Convert results to GraphNode
	nodes := make([]*GraphNode, 0)
	nodeMap := make(map[string]*GraphNode)

	for _, result := range results {
		for _, nodeData := range result.Nodes {
			nodeID, ok := nodeData["id"].(string)
			if !ok {
				continue
			}

			if _, exists := nodeMap[nodeID]; !exists {
				node := &GraphNode{
					ID:         nodeID,
					Type:       getString(nodeData, "type"),
					Properties: nodeData,
				}
				nodeMap[nodeID] = node
				nodes = append(nodes, node)
			}
		}
	}

	return nodes, nil
}

// Query performs a Cypher query on the graph
func (kg *KnowledgeGraph) Query(ctx context.Context, cypher string, params map[string]interface{}) ([]*GraphNode, error) {
	if kg.graphClient == nil {
		return nil, fmt.Errorf("graph client not available")
	}

	results, err := kg.graphClient.Query(ctx, cypher, params)
	if err != nil {
		return nil, fmt.Errorf("failed to query graph: %w", err)
	}

	nodes := make([]*GraphNode, 0)
	nodeMap := make(map[string]*GraphNode)

	for _, result := range results {
		for _, nodeData := range result.Nodes {
			nodeID, ok := nodeData["id"].(string)
			if !ok {
				continue
			}

			if _, exists := nodeMap[nodeID]; !exists {
				node := &GraphNode{
					ID:         nodeID,
					Type:       getString(nodeData, "type"),
					Properties: nodeData,
				}
				nodeMap[nodeID] = node
				nodes = append(nodes, node)
			}
		}
	}

	return nodes, nil
}

// FindRelated finds entities related to a given entity
func (kg *KnowledgeGraph) FindRelated(ctx context.Context, entityID string, relationType string, limit int) ([]*GraphNode, error) {
	if limit <= 0 {
		limit = 10
	}

	cypher := `
		MATCH (start {id: $start_id})-[r:` + relationType + `]->(related)
		RETURN related
		LIMIT $limit
	`

	params := map[string]interface{}{
		"start_id": entityID,
		"limit":    limit,
	}

	return kg.Query(ctx, cypher, params)
}

// GraphNode represents a node in the knowledge graph
type GraphNode struct {
	ID         string
	Type       string
	Properties map[string]interface{}
}

// buildTraversalQuery builds a Cypher query for graph traversal
func (kg *KnowledgeGraph) buildTraversalQuery(startID string, relationTypes []string, depth int, limit int) string {
	relationPattern := ""
	if len(relationTypes) > 0 {
		relationPattern = ":" + relationTypes[0]
	}

	query := fmt.Sprintf(`
		MATCH path = (start {id: $start_id})-[*1..%d%s]-(related)
		RETURN DISTINCT related
		LIMIT $limit
	`, depth, relationPattern)

	return query
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func copyMetadataKG(src map[string]interface{}) map[string]interface{} {
	if src == nil {
		return make(map[string]interface{})
	}
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

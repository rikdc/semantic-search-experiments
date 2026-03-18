package embedder

import (
	"context"

	"github.com/ollama/ollama/api"
)

// OllamaClient wraps the Ollama API for generating embeddings locally
type OllamaClient struct {
	client *api.Client
	model  string
}

// NewOllamaClient creates a new Ollama embedding client
// Reads OLLAMA_HOST from environment (defaults to http://localhost:11434)
func NewOllamaClient(model string) *OllamaClient {
	c, _ := api.ClientFromEnvironment()
	return &OllamaClient{client: c, model: model}
}

// CreateEmbedding generates an embedding vector for the given text
func (o *OllamaClient) CreateEmbedding(text string) (Vector, error) {
	req := &api.EmbeddingRequest{
		Model:  o.model,
		Prompt: text,
	}

	resp, err := o.client.Embeddings(context.Background(), req)
	if err != nil {
		return nil, err
	}

	// Convert from []float64 (Ollama SDK) to []float32 (our internal format)
	vec := make(Vector, len(resp.Embedding))
	for i, v := range resp.Embedding {
		vec[i] = float32(v)
	}
	return vec, nil
}

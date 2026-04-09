// Package embedder provides a common interface for generating text embeddings,
// with implementations for OpenAI and Ollama.
package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// Embedder converts a text string into a vector of float32 values.
type Embedder interface {
	CreateEmbedding(ctx context.Context, text string) ([]float32, error)
}

// openAIClient wraps the OpenAI embeddings API.
type openAIClient struct {
	client *openai.Client
	model  string
}

// NewOpenAIClient returns an Embedder backed by the OpenAI API.
// The API key is read from the OPENAI_API_KEY environment variable.
func NewOpenAIClient(model string) Embedder {
	return &openAIClient{
		client: openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		model:  model,
	}
}

func (c *openAIClient) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	resp, err := c.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.EmbeddingModel(c.model),
	})
	if err != nil {
		return nil, fmt.Errorf("openai embedding: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("openai returned empty embedding data")
	}
	return resp.Data[0].Embedding, nil
}

// ollamaClient wraps the Ollama local embedding API.
type ollamaClient struct {
	host  string
	model string
}

// NewOllamaClient returns an Embedder backed by a locally running Ollama instance.
// The host defaults to http://localhost:11434 if OLLAMA_HOST is not set.
func NewOllamaClient(model string) Embedder {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}
	return &ollamaClient{host: host, model: model}
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (c *ollamaClient) CreateEmbedding(ctx context.Context, text string) ([]float32, error) {
	body, err := json.Marshal(ollamaRequest{Model: c.model, Prompt: text})
	if err != nil {
		return nil, fmt.Errorf("marshalling ollama request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.host+"/api/embeddings", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating ollama request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status %d", resp.StatusCode)
	}

	var result ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding ollama response: %w", err)
	}
	return result.Embedding, nil
}

// NewClient returns an Embedder for the given provider and model.
// provider must be "openai" or "ollama".
func NewClient(provider, model string) (Embedder, error) {
	switch provider {
	case "openai":
		return NewOpenAIClient(model), nil
	case "ollama":
		return NewOllamaClient(model), nil
	default:
		return nil, fmt.Errorf("unknown provider %q: must be openai or ollama", provider)
	}
}

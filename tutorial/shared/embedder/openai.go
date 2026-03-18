package embedder

import (
	"context"

	"github.com/openai/openai-go"
)

type OpenAIClient struct {
	client openai.Client
	model  string
}

func NewOpenAIClient(model string) *OpenAIClient {
	client := openai.NewClient()
	return &OpenAIClient{
		client: client,
		model:  model,
	}
}

func (o *OpenAIClient) CreateEmbedding(text string) (Vector, error) {
	params := openai.EmbeddingNewParams{
		Input: openai.EmbeddingNewParamsInputUnion{
			OfString: openai.String(text),
		},
		Model: openai.EmbeddingModel(o.model),
	}

	res, err := o.client.Embeddings.New(context.Background(), params)
	if err != nil {
		return nil, err
	}

	// Boundary Conversion: SDK ([]float64) -> Internal ([]float32)
	raw := res.Data[0].Embedding
	vec := make(Vector, len(raw))
	for i, v := range raw {
		vec[i] = float32(v)
	}

	return vec, nil
}

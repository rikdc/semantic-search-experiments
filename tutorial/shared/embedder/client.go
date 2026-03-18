package embedder

// Vector is our standard representation for embeddings (float32 for efficiency)
type Vector = []float32

// Client is the interface that all embedding providers must implement
type Client interface {
	CreateEmbedding(text string) (Vector, error)
}

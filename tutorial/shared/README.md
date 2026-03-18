# Shared Tutorial Code

This directory contains **boilerplate code** that's provided for all tutorial parts. You don't need to implement these—they're infrastructure to support the learning objectives.

## What's Included

### `embedder/` - Embedding API Clients

Pre-built clients for generating embeddings from text:

- **`openai.go`** - OpenAI API client (text-embedding-3-small, text-embedding-3-large)
- **`ollama.go`** - Ollama local client (qwen2.5:latest, nomic-embed-text, etc.)
- **`client.go`** - Common interface and Vector type

**Usage in tutorial code:**

```go
import "embed/shared/embedder"

// OpenAI (requires OPENAI_API_KEY environment variable)
client := embedder.NewOpenAIClient("text-embedding-3-small")

// OR Ollama (requires Ollama running locally)
client := embedder.NewOllamaClient("qwen2.5:latest")

// Both implement the same interface
embedding, err := client.CreateEmbedding("Hello world")
// embedding is []float32 (Vector type)
```

## Why Shared Code?

The tutorial focuses on **understanding embeddings**, not API integration details. These clients:

- Handle API authentication and requests
- Convert between SDK types ([]float64) and our internal format ([]float32)
- Provide consistent interface for both providers
- Let you focus on the core concepts: similarity, centroids, search, etc.

## How to Use in Each Part

Each part's `go.mod` should reference the shared embedder:

```go
module embed

go 1.21

require (
    github.com/openai/openai-go v0.1.0-alpha.36
    github.com/ollama/ollama v0.1.0  // Optional if using Ollama
    github.com/spf13/cobra v1.8.1
)

replace embed/shared => ../shared
```

Then import in your code:

```go
import "embed/shared/embedder"
```

## What You'll Actually Implement

The tutorial has you implement the **interesting parts**:

- **Part 1**: Cosine similarity, centroid calculation
- **Part 2**: Z-score anomaly detection
- **Part 3**: AST parsing, vector stores, semantic search
- **Part 4**: Document metadata and enrichment
- **Part 5**: Vector arithmetic and query steering
- **Part 6**: PCA visualization, semantic profiles

The `embedder` package is just infrastructure to support these learning goals.

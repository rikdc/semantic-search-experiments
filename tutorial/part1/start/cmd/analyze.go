package cmd

import (
	_ "embedtutorial/math" // Will be used when implementing similarity calculations
	"embedtutorial/shared/embedder"
	_ "encoding/json" // Will be used for loading inputs.json
	"fmt"
	_ "os" // Will be used for file operations
	"strings"

	"github.com/spf13/cobra"
)

// Transaction represents a financial transaction with label and text
type Transaction struct {
	Label    string `json:"label"`
	Category string `json:"category"`
	Text     string `json:"text"`
}

// TransactionEmbedding pairs a transaction with its embedding vector
type TransactionEmbedding struct {
	Transaction Transaction
	Embedding   []float32
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze transaction similarities using embeddings",
	Long: `Generate embeddings for transactions and analyze their semantic relationships.

Calculates:
- Pairwise cosine similarities (how similar each transaction is to every other)
- Centroid (average embedding representing typical transactions)
- Distance from centroid (how typical each transaction is)`,
	RunE: runAnalyze,
}

func Register(root *cobra.Command) {
	root.AddCommand(analyzeCmd)
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	fmt.Println("\n📊 Analyzing Transaction Embeddings")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Step 1: Load transactions from inputs.json
	// TODO: Read the inputs.json file and parse into []Transaction
	//
	// Hint: Use os.ReadFile() and json.Unmarshal()
	//
	// Example:
	//   data, err := os.ReadFile("inputs.json")
	//   if err != nil {
	//       return fmt.Errorf("failed to read inputs.json: %w", err)
	//   }
	//
	//   var transactions []Transaction
	//   if err := json.Unmarshal(data, &transactions); err != nil {
	//       return fmt.Errorf("failed to parse JSON: %w", err)
	//   }

	var transactions []Transaction // TODO: Load from file

	fmt.Printf("Loaded %d transactions from inputs.json\n\n", len(transactions))

	// Step 2: Get configuration from flags
	provider, _ := cmd.Flags().GetString("provider")
	model, _ := cmd.Flags().GetString("model")

	fmt.Printf("Using %s with model %s\n\n", provider, model)

	// Step 3: Create embedding client
	// The embedder package is provided in ../shared/embedder
	var client embedder.Client
	switch provider {
	case "openai":
		client = embedder.NewOpenAIClient(model)
	case "ollama":
		client = embedder.NewOllamaClient(model)
	default:
		return fmt.Errorf("unsupported provider: %s", provider)
	}
	_ = client // Will be used when generating embeddings

	// Step 4: Generate embeddings for all transactions
	// TODO: Create embeddings for each transaction
	//
	// Hint: Loop through transactions and call client.CreateEmbedding()
	//
	// Example:
	//   embeddings := make([]TransactionEmbedding, len(transactions))
	//   for i, tx := range transactions {
	//       fmt.Printf("  [%d/%d] %s", i+1, len(transactions), tx.Label)
	//
	//       embedding, err := client.CreateEmbedding(tx.Text)
	//       if err != nil {
	//           return fmt.Errorf("failed to embed %s: %w", tx.Label, err)
	//       }
	//
	//       embeddings[i] = TransactionEmbedding{
	//           Transaction: tx,
	//           Embedding:   embedding,
	//       }
	//
	//       fmt.Printf(" → [%.4f, %.4f, ...] (%d dims)\n", 
	//           embedding[0], embedding[1], len(embedding))
	//   }

	fmt.Println("Generating embeddings...")
	var embeddings []TransactionEmbedding // TODO: Generate embeddings
	fmt.Println()

	// Step 5: Calculate pairwise similarities
	// TODO: Calculate similarity matrix (transaction vs transaction)
	//
	// Hint: Use nested loops and math.CosineSimilarity()
	// The diagonal (comparing a vector to itself) should be exactly 1.0000
	//
	// Example:
	//   fmt.Println("=== Pairwise Similarity Matrix ===")
	//   fmt.Println()
	//
	//   // Print header
	//   fmt.Printf("%-15s", "")
	//   for _, emb := range embeddings {
	//       fmt.Printf("%-10s ", emb.Transaction.Label)
	//   }
	//   fmt.Println()
	//
	//   // Print rows (full symmetric matrix)
	//   for _, emb1 := range embeddings {
	//       fmt.Printf("%-15s", emb1.Transaction.Label)
	//
	//       for _, emb2 := range embeddings {
	//           sim, _ := math.CosineSimilarity(emb1.Embedding, emb2.Embedding)
	//           fmt.Printf("%-10.4f ", sim)
	//       }
	//       fmt.Println()
	//   }

	// TODO: Implement pairwise similarity calculation
	fmt.Println()

	// Step 6: Calculate centroid
	// TODO: Calculate the centroid of all embeddings
	//
	// Hint: Extract just the embedding vectors and call math.CalculateCentroid()
	//
	// Example:
	//   fmt.Println("=== Calculating Centroid ===")
	//   fmt.Println()
	//
	//   vectors := make([][]float32, len(embeddings))
	//   for i, emb := range embeddings {
	//       vectors[i] = emb.Embedding
	//   }
	//
	//   centroid, err := math.CalculateCentroid(vectors)
	//   if err != nil {
	//       return fmt.Errorf("failed to calculate centroid: %w", err)
	//   }
	//
	//   fmt.Printf("Centroid calculated: [%.4f, %.4f, ...] (%d dims)\n", 
	//       centroid[0], centroid[1], len(centroid))

	var centroid []float32 // TODO: Calculate centroid

	// Step 7: Calculate similarity to centroid
	// TODO: For each transaction, calculate similarity to centroid
	//
	// Hint: Loop through embeddings and call math.CosineSimilarity()
	//
	// Example:
	//   fmt.Println("=== Similarity to Centroid ===")
	//   fmt.Println()
	//
	//   for _, emb := range embeddings {
	//       sim, err := math.CosineSimilarity(emb.Embedding, centroid)
	//       if err != nil {
	//           return fmt.Errorf("failed to calculate similarity: %w", err)
	//       }
	//
	//       category := ""
	//       if sim > 0.85 {
	//           category = " (very typical)"
	//       } else if sim > 0.75 {
	//           category = " (typical)"
	//       } else {
	//           category = " (unusual)"
	//       }
	//
	//       fmt.Printf("%-15s: %.4f %s\n", emb.Transaction.Label, sim, category)
	//   }

	// TODO: Implement centroid similarity calculation

	// Done!
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("✅ Analysis complete!")
	fmt.Println()

	return nil
}

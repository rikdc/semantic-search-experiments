package cmd

import (
	"embedtutorial/shared/embedder"
	"embedtutorial/math"
	"encoding/json"
	"fmt"
	"os"
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
	data, err := os.ReadFile("inputs.json")
	if err != nil {
		return fmt.Errorf("failed to read inputs.json: %w\nMake sure inputs.json exists in the current directory", err)
	}

	var transactions []Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("Loaded %d transactions from inputs.json\n\n", len(transactions))

	// Step 2: Get provider and model from flags
	provider, _ := cmd.Flags().GetString("provider")
	model, _ := cmd.Flags().GetString("model")

	fmt.Printf("Using %s with model %s\n\n", provider, model)

	// Step 3: Create embedding client
	var client embedder.Client
	switch provider {
	case "openai":
		client = embedder.NewOpenAIClient(model)
	case "ollama":
		client = embedder.NewOllamaClient(model)
	default:
		return fmt.Errorf("unsupported provider: %s (use 'openai' or 'ollama')", provider)
	}

	// Step 4: Generate embeddings for all transactions
	fmt.Println("Generating embeddings...")
	embeddings := make([]TransactionEmbedding, len(transactions))

	for i, tx := range transactions {
		fmt.Printf("  [%d/%d] %s", i+1, len(transactions), tx.Label)

		embedding, err := client.CreateEmbedding(tx.Text)
		if err != nil {
			return fmt.Errorf("failed to embed transaction %s: %w", tx.Label, err)
		}

		embeddings[i] = TransactionEmbedding{
			Transaction: tx,
			Embedding:   embedding,
		}

		// Show first few dimensions as preview
		fmt.Printf(" → [%.4f, %.4f, ...] (%d dims)\n", embedding[0], embedding[1], len(embedding))
	}

	fmt.Println()

	// Step 5: Calculate pairwise similarities
	fmt.Println("=== Pairwise Similarity Matrix ===")
	fmt.Println()

	// Print header
	fmt.Printf("%-15s", "")
	for _, emb := range embeddings {
		fmt.Printf("%-10s ", emb.Transaction.Label)
	}
	fmt.Println()

	// Print rows (full symmetric matrix)
	for _, emb1 := range embeddings {
		fmt.Printf("%-15s", emb1.Transaction.Label)

		for _, emb2 := range embeddings {
			sim, err := math.CosineSimilarity(emb1.Embedding, emb2.Embedding)
			if err != nil {
				return fmt.Errorf("failed to calculate similarity: %w", err)
			}
			fmt.Printf("%-10.4f ", sim)
		}
		fmt.Println()
	}

	fmt.Println()

	// Step 6: Calculate centroid
	fmt.Println("=== Calculating Centroid ===")
	fmt.Println()

	vectors := make([][]float32, len(embeddings))
	for i, emb := range embeddings {
		vectors[i] = emb.Embedding
	}

	centroid, err := math.CalculateCentroid(vectors)
	if err != nil {
		return fmt.Errorf("failed to calculate centroid: %w", err)
	}

	fmt.Printf("Centroid calculated: [%.4f, %.4f, ...] (%d dims)\n", centroid[0], centroid[1], len(centroid))
	fmt.Println("This represents the 'average' or 'typical' transaction")
	fmt.Println()

	// Step 7: Calculate similarity to centroid
	fmt.Println("=== Similarity to Centroid ===")
	fmt.Println()

	var sumSim float32
	similarities := make([]float32, len(embeddings))

	for i, emb := range embeddings {
		sim, err := math.CosineSimilarity(emb.Embedding, centroid)
		if err != nil {
			return fmt.Errorf("failed to calculate centroid similarity: %w", err)
		}

		similarities[i] = sim
		sumSim += sim

		// Categorize based on similarity
		category := ""
		if sim > 0.85 {
			category = " (very typical)"
		} else if sim > 0.75 {
			category = " (typical)"
		} else {
			category = " (unusual)"
		}

		fmt.Printf("%-15s: %.4f %s\n", emb.Transaction.Label, sim, category)
	}

	avgSim := sumSim / float32(len(embeddings))
	fmt.Printf("\nAverage similarity to centroid: %.4f\n", avgSim)

	// Done!
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("✅ Analysis complete!")
	fmt.Println()
	fmt.Println("💡 Insights:")
	fmt.Println("  - High pairwise similarity (>0.9) = very related transactions")
	fmt.Println("  - Low centroid similarity (<0.75) = unusual/outlier transactions")
	fmt.Println("  - Similar categories (e.g., Coffee_A, Coffee_B) should cluster together")
	fmt.Println()

	return nil
}

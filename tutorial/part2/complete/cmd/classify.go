package cmd

import (
	"embedtutorial/math"
	"embedtutorial/shared/embedder"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

// Transaction represents a financial transaction with label, category, and text
type Transaction struct {
	Label    string `json:"label"`
	Category string `json:"category"`
	Text     string `json:"text"`
}

var classifyCmd = &cobra.Command{
	Use:   "classify [query]",
	Short: "Classify a transaction by nearest category centroid",
	Long: `Classify a new transaction description by comparing its embedding
to per-category centroids computed from training data in inputs.json.

Loads training data, groups entries by category, computes a centroid
for each category, then finds which centroid the query is closest to.`,
	Args: cobra.ExactArgs(1),
	RunE: runClassify,
}

func Register(root *cobra.Command) {
	root.AddCommand(classifyCmd)
}

func runClassify(cmd *cobra.Command, args []string) error {
	query := args[0]

	// Step A: Load and unmarshal inputs.json
	data, err := os.ReadFile("inputs.json")
	if err != nil {
		return fmt.Errorf("failed to read inputs.json: %w\nMake sure inputs.json exists in the current directory", err)
	}

	var transactions []Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return fmt.Errorf("failed to parse inputs.json: %w", err)
	}

	// Step B: Initialise the embedder client
	provider, _ := cmd.Flags().GetString("provider")
	model, _ := cmd.Flags().GetString("model")

	fmt.Printf("Generating embeddings for training data via %s (%s)...\n", provider, model)

	var client embedder.Client
	switch provider {
	case "openai":
		client = embedder.NewOpenAIClient(model)
	case "ollama":
		client = embedder.NewOllamaClient(model)
	default:
		return fmt.Errorf("unsupported provider: %s (use 'openai' or 'ollama')", provider)
	}

	// Step C: Generate embeddings for all inputs
	type labeledEmbedding struct {
		category  string
		embedding []float32
	}

	labeled := make([]labeledEmbedding, len(transactions))
	for i, tx := range transactions {
		embedding, err := client.CreateEmbedding(tx.Text)
		if err != nil {
			return fmt.Errorf("failed to embed '%s': %w", tx.Label, err)
		}
		labeled[i] = labeledEmbedding{
			category:  tx.Category,
			embedding: embedding,
		}
	}

	// Step D: Group vectors by category
	groups := make(map[string][][]float32)
	for _, le := range labeled {
		groups[le.category] = append(groups[le.category], le.embedding)
	}

	// Step E: Calculate centroid for each category
	centroids := make(map[string][]float32)
	for category, vectors := range groups {
		centroid, err := math.CalculateCentroid(vectors)
		if err != nil {
			return fmt.Errorf("failed to calculate centroid for category %s: %w", category, err)
		}
		centroids[category] = centroid
	}

	// Step F: Embed the query text
	fmt.Printf("\nClassifying: '%s'\n", query)

	queryEmbedding, err := client.CreateEmbedding(query)
	if err != nil {
		return fmt.Errorf("failed to embed query: %w", err)
	}

	// Step G: Calculate similarity between query and each centroid
	type categoryScore struct {
		category   string
		similarity float32
	}

	scores := make([]categoryScore, 0, len(centroids))
	for category, centroid := range centroids {
		sim, err := math.CosineSimilarity(queryEmbedding, centroid)
		if err != nil {
			return fmt.Errorf("failed to calculate similarity for category %s: %w", category, err)
		}
		scores = append(scores, categoryScore{category: category, similarity: sim})
	}

	// Sort by similarity descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].similarity > scores[j].similarity
	})

	// Step H: Print results and declare the best category
	fmt.Println("\nCategory Similarities:")
	fmt.Println("----------------------")
	for _, s := range scores {
		fmt.Printf("Category: %-12s | Similarity: %.4f\n", s.category, s.similarity)
	}

	best := scores[0]
	fmt.Printf("\nResult: '%s' is likely in the [%s] category.\n", query, best.category)

	return nil
}

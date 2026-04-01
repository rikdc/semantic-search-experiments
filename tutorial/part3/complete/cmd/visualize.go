package cmd

import (
	"embedtutorial/math"
	"embedtutorial/shared/embedder"
	"embedtutorial/viz"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Transaction represents a financial transaction with label, category, and text.
type Transaction struct {
	Label    string `json:"label"`
	Category string `json:"category"`
	Text     string `json:"text"`
}

var visualizeCmd = &cobra.Command{
	Use:   "visualize",
	Short: "Project embeddings to 2D and write a scatter plot",
	Long: `Project transaction embeddings to 2D using PCA and write output.svg.

Loads training data from inputs.json, generates embeddings, and uses
Principal Component Analysis to reduce the vectors to two dimensions.
Pass --query to include a specific transaction in the plot.`,
	RunE: runVisualize,
}

func Register(root *cobra.Command) {
	root.AddCommand(visualizeCmd)
}

func init() {
	visualizeCmd.Flags().String("query", "", "Optional transaction text to include in the plot as a query point")
	visualizeCmd.Flags().String("output", "output.svg", "Output SVG filename")
}

func runVisualize(cmd *cobra.Command, args []string) error {
	query, _ := cmd.Flags().GetString("query")
	output, _ := cmd.Flags().GetString("output")

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

	fmt.Printf("Generating embeddings via %s (%s)...\n", provider, model)

	var client embedder.Client
	switch provider {
	case "openai":
		client = embedder.NewOpenAIClient(model)
	case "ollama":
		client = embedder.NewOllamaClient(model)
	default:
		return fmt.Errorf("unsupported provider: %s (use 'openai' or 'ollama')", provider)
	}

	// Step C: Generate embeddings for all training inputs
	vectors := make([][]float32, len(transactions))
	labels := make([]string, len(transactions))
	categories := make([]string, len(transactions))

	for i, tx := range transactions {
		vec, err := client.CreateEmbedding(tx.Text)
		if err != nil {
			return fmt.Errorf("failed to embed '%s': %w", tx.Label, err)
		}
		vectors[i] = vec
		labels[i] = tx.Label
		categories[i] = tx.Category
		fmt.Printf("  [%d/%d] %s\n", i+1, len(transactions), tx.Label)
	}

	// Step D: Append query vector if provided.
	// The query must be included BEFORE PCA runs — it must share the same
	// PCA axes as the training points to be comparable on the plot.
	if query != "" {
		fmt.Printf("\nEmbedding query: '%s'\n", query)
		queryVec, err := client.CreateEmbedding(query)
		if err != nil {
			return fmt.Errorf("failed to embed query: %w", err)
		}
		vectors = append(vectors, queryVec)
		labels = append(labels, "[query]")
		categories = append(categories, "[query]")
	}

	// Step E: Project all vectors to 2D via PCA
	fmt.Println("\nRunning PCA...")
	points, err := math.ProjectTo2D(vectors)
	if err != nil {
		return fmt.Errorf("PCA failed: %w", err)
	}

	// Step F: Write the SVG scatter plot
	if err := viz.SaveEmbeddingPlot(output, labels, categories, points); err != nil {
		return fmt.Errorf("failed to save plot: %w", err)
	}

	fmt.Printf("Wrote %s (%d points)\n", output, len(points))
	return nil
}

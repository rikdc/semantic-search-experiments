package cmd

import (
	"embedtutorial/math"
	"embedtutorial/shared/embedder"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect anomalous transactions by distance from the Normal centroid",
	Long: `Load transactions from inputs.json, embed each one, and flag any that
fall below the similarity threshold relative to the centroid of Normal transactions.

Use --raw to embed only the dollar amount instead of a synthetic sentence.
This demonstrates how context-free embedding degrades detection quality.`,
	RunE: runDetect,
}

func Register(root *cobra.Command) {
	root.AddCommand(detectCmd)
}

func init() {
	detectCmd.Flags().Bool("raw", false, "Embed raw amount only, without merchant/time context")
	detectCmd.Flags().Float32("threshold", 0.80, "Similarity threshold below which a transaction is flagged")
}

func runDetect(cmd *cobra.Command, args []string) error {
	raw, _ := cmd.Flags().GetBool("raw")
	threshold, _ := cmd.Flags().GetFloat32("threshold")

	// Step A: Load inputs.json
	data, err := os.ReadFile("inputs.json")
	if err != nil {
		return fmt.Errorf("failed to read inputs.json: %w\nMake sure inputs.json exists in the current directory", err)
	}
	var transactions []math.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return fmt.Errorf("failed to parse inputs.json: %w", err)
	}

	// Step B: Initialise the embedder client
	provider, _ := cmd.Flags().GetString("provider")
	model, _ := cmd.Flags().GetString("model")

	mode := "synthetic string"
	if raw {
		mode = "raw amount"
	}
	fmt.Printf("Generating embeddings via %s (%s) [%s mode]...\n", provider, model, mode)

	var client embedder.Client
	switch provider {
	case "openai":
		client = embedder.NewOpenAIClient(model)
	case "ollama":
		client = embedder.NewOllamaClient(model)
	default:
		return fmt.Errorf("unsupported provider: %s (use 'openai' or 'ollama')", provider)
	}

	// Step C: Format and embed each transaction.
	// FormatTransaction builds a descriptive sentence; FormatRaw strips context down
	// to just the dollar amount. The comparison shows how much context matters.
	embeddings := make([][]float32, len(transactions))
	for i, tx := range transactions {
		var text string
		if raw {
			text = math.FormatRaw(tx)
		} else {
			text = math.FormatTransaction(tx)
		}

		vec, err := client.CreateEmbedding(text)
		if err != nil {
			return fmt.Errorf("failed to embed '%s': %w", tx.Label, err)
		}
		embeddings[i] = vec
		fmt.Printf("  [%d/%d] %s → %q\n", i+1, len(transactions), tx.Label, text)
	}

	// Step D: Collect Normal embeddings and compute centroid.
	// Only transactions labeled "Normal" contribute to the reference point.
	var normalVecs [][]float32
	for i, tx := range transactions {
		if tx.Category == "Normal" {
			normalVecs = append(normalVecs, embeddings[i])
		}
	}
	if len(normalVecs) == 0 {
		return fmt.Errorf("no transactions labeled 'Normal' found in inputs.json")
	}

	centroid, err := math.CalculateCentroid(normalVecs)
	if err != nil {
		return fmt.Errorf("centroid: %w", err)
	}

	// Step E: Compare every transaction to the Normal centroid and flag outliers.
	fmt.Printf("\nNormal centroid built from %d transactions. Threshold: %.2f\n\n", len(normalVecs), threshold)
	fmt.Printf("%-14s %-18s %s\n", "Label", "Category", "Similarity")
	fmt.Println("----------------------------------------------")

	for i, tx := range transactions {
		sim, err := math.CosineSimilarity(embeddings[i], centroid)
		if err != nil {
			return fmt.Errorf("similarity for %s: %w", tx.Label, err)
		}

		flag := "✓ normal"
		if sim < threshold {
			flag = "✗ FLAGGED"
		}

		fmt.Printf("%-14s (%-16s) sim: %.4f  %s\n", tx.Label, tx.Category, sim, flag)
	}

	return nil
}

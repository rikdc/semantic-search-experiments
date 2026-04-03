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
are anomalous relative to the centroid of Normal transactions.

By default, flags transactions whose similarity falls below a fixed threshold.
Use --zscores to flag based on standard deviations from the Normal cluster mean
instead — this adapts to your data's spread rather than a hand-tuned cutoff.

Use --raw to embed only the dollar amount instead of a synthetic sentence.`,
	RunE: runDetect,
}

func Register(root *cobra.Command) {
	root.AddCommand(detectCmd)
}

func init() {
	detectCmd.Flags().Bool("raw", false, "Embed raw amount only, without merchant/time context")
	detectCmd.Flags().Float32("threshold", 0.80, "Similarity threshold below which a transaction is flagged (fixed-threshold mode)")
	detectCmd.Flags().Bool("zscores", false, "Use z-score thresholding instead of a fixed similarity threshold")
	detectCmd.Flags().Float64("zscore-threshold", -2.0, "Z-score below which a transaction is flagged (z-score mode)")
}

func runDetect(cmd *cobra.Command, args []string) error {
	raw, _ := cmd.Flags().GetBool("raw")
	threshold, _ := cmd.Flags().GetFloat32("threshold")
	useZScores, _ := cmd.Flags().GetBool("zscores")
	zscoreThreshold, _ := cmd.Flags().GetFloat64("zscore-threshold")

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

	// Step C: Format and embed each transaction
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

	// Step D: Collect Normal embeddings and compute centroid
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

	// Step E: Compute similarity of every transaction to the Normal centroid
	sims := make([]float32, len(transactions))
	for i := range transactions {
		sim, err := math.CosineSimilarity(embeddings[i], centroid)
		if err != nil {
			return fmt.Errorf("similarity for %s: %w", transactions[i].Label, err)
		}
		sims[i] = sim
	}

	// Step F: Collect Normal similarities and compute z-score parameters if needed
	var mean, stddev float64
	if useZScores {
		normalSims := make([]float32, 0, len(normalVecs))
		for i, tx := range transactions {
			if tx.Category == "Normal" {
				normalSims = append(normalSims, sims[i])
			}
		}
		mean, stddev = math.CalculateMeanStdDev(normalSims)
	}

	// Step G: Print results
	if useZScores {
		fmt.Printf("\nNormal centroid built from %d transactions. Z-score threshold: %.2f\n\n", len(normalVecs), zscoreThreshold)
		fmt.Printf("%-14s %-18s %-12s %s\n", "Label", "Category", "Similarity", "Z-Score")
		fmt.Println("------------------------------------------------------")

		for i, tx := range transactions {
			z := (float64(sims[i]) - mean) / stddev
			flag := "✓ normal"
			if z < zscoreThreshold {
				flag = "✗ FLAGGED"
			}
			fmt.Printf("%-14s (%-16s) sim: %.4f  z: %6.2f  %s\n",
				tx.Label, tx.Category, sims[i], z, flag)
		}
	} else {
		fmt.Printf("\nNormal centroid built from %d transactions. Threshold: %.2f\n\n", len(normalVecs), threshold)
		fmt.Printf("%-14s %-18s %s\n", "Label", "Category", "Similarity")
		fmt.Println("----------------------------------------------")

		for i, tx := range transactions {
			flag := "✓ normal"
			if sims[i] < threshold {
				flag = "✗ FLAGGED"
			}
			fmt.Printf("%-14s (%-16s) sim: %.4f  %s\n",
				tx.Label, tx.Category, sims[i], flag)
		}
	}

	return nil
}

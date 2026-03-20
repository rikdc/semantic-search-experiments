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

	// TODO A: Load and unmarshal inputs.json
	// Read inputs.json and parse it into a []Transaction slice.
	// Hint: use os.ReadFile("inputs.json") then json.Unmarshal into transactions
	var transactions []Transaction
	_ = transactions

	// TODO B: Initialise the embedder client (switch on provider flag)
	// Create an embedder.Client based on the --provider flag value.
	// Hint: get provider/model with cmd.Flags().GetString(), then switch on provider
	provider, _ := cmd.Flags().GetString("provider")
	model, _ := cmd.Flags().GetString("model")

	fmt.Printf("Generating embeddings for training data via %s (%s)...\n", provider, model)

	var client embedder.Client
	_ = client

	// TODO C: Generate embeddings for all inputs
	// Loop through transactions and call client.CreateEmbedding(tx.Text) for each.
	// Hint: store results as a slice of structs pairing category with embedding

	// TODO D: Group vectors by category into map[string][][]float32
	// Build a map where keys are category names and values are slices of embedding vectors.
	// Hint: range over your labeled results and append each embedding to groups[category]
	groups := make(map[string][][]float32)
	_ = groups

	// TODO E: Call CalculateCentroid for each category group
	// Compute the centroid for each category's group of vectors.
	// Hint: range over groups, call math.CalculateCentroid(vectors) for each
	centroids := make(map[string][]float32)
	_ = centroids

	// TODO F: Embed the query text
	// Generate an embedding for the query argument passed on the command line.
	// Hint: call client.CreateEmbedding(query)
	fmt.Printf("\nClassifying: '%s'\n", query)

	var queryEmbedding []float32
	_ = queryEmbedding

	// TODO G: Call CosineSimilarity between query and each centroid
	// Compare the query embedding to each category centroid.
	// Hint: range over centroids, call math.CosineSimilarity(queryEmbedding, centroid)
	type categoryScore struct {
		category   string
		similarity float32
	}

	var scores []categoryScore
	_ = scores

	// TODO H: Find and print the best category
	// Sort scores descending by similarity, print all, then declare the winner.
	// Hint: use sort.Slice, then print each score and scores[0] as the result

	// Suppress unused import warnings — remove these lines as you implement the TODOs
	_ = os.ReadFile
	_ = json.Unmarshal
	_ = math.CosineSimilarity
	_ = sort.Slice

	return nil
}

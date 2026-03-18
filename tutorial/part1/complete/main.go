package main

import (
	"embedtutorial/cmd"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	provider string
	model    string
)

var rootCmd = &cobra.Command{
	Use:   "embed",
	Short: "Embeddings CLI for semantic analysis",
	Long: `A command-line tool for working with vector embeddings.

Supports semantic search, similarity analysis, anomaly detection,
and visualization of high-dimensional embedding spaces.`,
}

func init() {
	// Add persistent flags for provider and model
	rootCmd.PersistentFlags().StringVar(&provider, "provider", "openai", "Embedding provider (openai or ollama)")
	rootCmd.PersistentFlags().StringVar(&model, "model", "text-embedding-3-small", "Embedding model name")

	// Register commands
	cmd.Register(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

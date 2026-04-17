package main

import (
	"embedtutorial/cmd"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "embed",
	Short: "Semantic code search over a Go codebase",
	Long: `A command-line tool for semantic search over Go source code.

Indexes function bodies using text embeddings, then retrieves relevant
functions by natural language query. Supports enriched chunks, namespace
filtering, and RRF hybrid search.`,
}

func init() {
	rootCmd.PersistentFlags().String("provider", "openai", "Embedding provider (openai or ollama)")
	rootCmd.PersistentFlags().String("model", "text-embedding-3-small", "Embedding model name")

	cmd.Register(rootCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

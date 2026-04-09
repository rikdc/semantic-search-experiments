package cmd

import (
	"context"
	"fmt"
	"os"

	chromem "github.com/philippgille/chromem-go"
	"github.com/spf13/cobra"

	"part6/indexer"
	"part6/shared/embedder"
)

const collectionName = "functions"

var (
	flagProvider  string
	flagModel     string
	flagDir       string
	flagDBPath    string
	flagTop       int
	flagIndex     bool
	flagBenchmark bool
)

// benchmarkCases maps natural-language queries to the expected top-1 function name.
// Cases marked with a comment are expected to fail at the baseline — they are the
// motivation for Part 7's signal-cleaning techniques.
var benchmarkCases = []struct {
	query    string
	expected string
	note     string
}{
	{"how do I connect to the database", "Connect", ""},
	{"check if a token is valid", "ValidateToken", ""},
	{"verify the user has permission", "Authorize", ""},
	{"write a JSON response", "WriteJSON", ""},
	{"handle an incoming HTTP request", "HandleRequest", ""},
	{"refresh an expired token", "RefreshToken", ""},
	{"store a value in the cache", "Set", ""},
	{"get all rows from a query", "FetchAll", ""},
	{"run a command", "Execute", "ambiguous: db.Execute and cli.Execute score similarly"},
	{"how long before a request times out", "getTimeout", "short function, weak signal"},
}

func NewSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Index Go source functions and search them by natural language",
		Long: `embed search indexes Go function bodies into a local vector store,
then retrieves the most semantically similar functions for a natural-language query.

To build the index:
  embed search --index

To search the index:
  embed search "how do I connect to the database"

To measure baseline accuracy:
  embed search --benchmark

The default --dir points to corpus/internal relative to the repo root.
Pass --dir to index a different codebase.`,
		RunE: runSearch,
	}

	cmd.Flags().StringVar(&flagProvider, "provider", "openai", "embedding provider: openai or ollama")
	cmd.Flags().StringVar(&flagModel, "model", "text-embedding-3-small", "embedding model name")
	cmd.Flags().StringVar(&flagDir, "dir", "../../../corpus/internal", "directory of Go source files to index")
	cmd.Flags().StringVar(&flagDBPath, "db", "./index", "path to persist the vector index")
	cmd.Flags().IntVar(&flagTop, "top", 3, "number of results to return")
	cmd.Flags().BoolVar(&flagIndex, "index", false, "build or rebuild the index before searching")
	cmd.Flags().BoolVar(&flagBenchmark, "benchmark", false, "run the benchmark suite and report top-1 accuracy")

	return cmd
}

func runSearch(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	emb, err := embedder.NewClient(flagProvider, flagModel)
	if err != nil {
		return err
	}

	// Wrap our embedder as a chromem EmbeddingFunc so the collection can
	// embed query text during Search without a separate call.
	embeddingFunc := chromem.EmbeddingFunc(func(ctx context.Context, text string) ([]float32, error) {
		return emb.CreateEmbedding(ctx, text)
	})

	db, err := chromem.NewPersistentDB(flagDBPath, false)
	if err != nil {
		return fmt.Errorf("opening vector store at %s: %w", flagDBPath, err)
	}

	collection, err := db.GetOrCreateCollection(collectionName, nil, embeddingFunc)
	if err != nil {
		return fmt.Errorf("getting collection: %w", err)
	}

	if flagIndex {
		fmt.Printf("Indexing %s\n", flagDir)
		if err := indexer.BuildIndex(ctx, flagDir, collection, emb); err != nil {
			return fmt.Errorf("building index: %w", err)
		}
		fmt.Printf("\nIndexed %d functions. Stored to %s.\n\n", collection.Count(), flagDBPath)
	}

	if collection.Count() == 0 {
		return fmt.Errorf("index is empty — run with --index to build it first")
	}

	if flagBenchmark {
		return runBenchmark(ctx, collection)
	}

	if len(args) == 0 {
		if flagIndex {
			return nil
		}
		return fmt.Errorf("provide a search query, or use --index to build the index")
	}

	return runQuery(ctx, args[0], collection)
}

func runQuery(ctx context.Context, query string, collection *chromem.Collection) error {
	results, err := indexer.Search(ctx, query, flagTop, collection)
	if err != nil {
		return err
	}

	fmt.Printf("Query: %q\n\n", query)
	for i, r := range results {
		fmt.Printf("%d. %-36s :: %-20s sim: %.4f\n",
			i+1, r.Metadata["file"], r.Metadata["func"], r.Similarity)
	}
	return nil
}

func runBenchmark(ctx context.Context, collection *chromem.Collection) error {
	fmt.Printf("Running benchmark (%d queries)...\n\n", len(benchmarkCases))

	correct := 0
	for _, tc := range benchmarkCases {
		results, err := indexer.Search(ctx, tc.query, 1, collection)
		if err != nil {
			return fmt.Errorf("benchmark query %q: %w", tc.query, err)
		}

		hit := len(results) > 0 && results[0].Metadata["func"] == tc.expected
		if hit {
			correct++
		}

		status := "✓"
		if !hit {
			status = "✗"
		}

		top := "(no results)"
		if len(results) > 0 {
			top = fmt.Sprintf("%s :: %s (%.4f)",
				results[0].Metadata["file"], results[0].Metadata["func"], results[0].Similarity)
		}

		fmt.Printf("%s %-50q → %s\n", status, tc.query, top)
		if tc.note != "" {
			fmt.Printf("   note: %s\n", tc.note)
		}
		if !hit {
			fmt.Printf("   expected: %s\n", tc.expected)
		}
	}

	pct := float64(correct) / float64(len(benchmarkCases)) * 100
	fmt.Printf("\nTop-1 accuracy: %d/%d (%.0f%%)\n", correct, len(benchmarkCases), pct)

	if _, err := os.Stat("./index"); err == nil {
		fmt.Println("\nRun with --index to rebuild if you've changed the source files.")
	}

	return nil
}

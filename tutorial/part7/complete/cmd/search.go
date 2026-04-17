package cmd

import (
	"context"
	"fmt"

	chromem "github.com/philippgille/chromem-go"
	"github.com/spf13/cobra"

	"embedtutorial/benchmark"
	"embedtutorial/indexer"
	"embedtutorial/search"
	"embedtutorial/shared/embedder"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Index and search Go source code by natural language",
	Long: `search indexes a Go codebase and retrieves functions by natural language query.

Use --index to build or rebuild the index before searching. Use --enrich to
prepend file and package context before embedding (fixes name-collision cases).
Use --package to restrict results to a single package directory. Use --hybrid
to combine dense and keyword search with RRF.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSearch,
}

func Register(root *cobra.Command) {
	root.AddCommand(searchCmd)
}

func init() {
	searchCmd.Flags().Bool("index", false, "Build or rebuild the index before searching")
	searchCmd.Flags().Bool("enrich", false, "Prepend file/package context to each function before embedding")
	searchCmd.Flags().Int("top", 3, "Number of results to return")
	searchCmd.Flags().String("package", "", "Restrict results to this package directory (e.g. cli)")
	searchCmd.Flags().Bool("hybrid", false, "Combine dense and keyword search with RRF")
	searchCmd.Flags().Bool("benchmark", false, "Run the benchmark query suite")
	searchCmd.Flags().String("dir", "../../corpus/internal", "Directory of Go source to index")
	searchCmd.Flags().String("index-dir", "./index", "Directory for the persistent index")
}

func runSearch(cmd *cobra.Command, args []string) error {
	provider, _ := cmd.Flags().GetString("provider")
	model, _ := cmd.Flags().GetString("model")
	doIndex, _ := cmd.Flags().GetBool("index")
	enrich, _ := cmd.Flags().GetBool("enrich")
	topN, _ := cmd.Flags().GetInt("top")
	pkg, _ := cmd.Flags().GetString("package")
	hybrid, _ := cmd.Flags().GetBool("hybrid")
	doBenchmark, _ := cmd.Flags().GetBool("benchmark")
	dir, _ := cmd.Flags().GetString("dir")
	indexDir, _ := cmd.Flags().GetString("index-dir")

	// Resolve provider/model from root persistent flags.
	if p, err := cmd.Root().PersistentFlags().GetString("provider"); err == nil && p != "" {
		provider = p
	}
	if m, err := cmd.Root().PersistentFlags().GetString("model"); err == nil && m != "" {
		model = m
	}

	var client embedder.Client
	switch provider {
	case "openai":
		client = embedder.NewOpenAIClient(model)
	case "ollama":
		client = embedder.NewOllamaClient(model)
	default:
		return fmt.Errorf("unsupported provider: %s (use 'openai' or 'ollama')", provider)
	}

	if doIndex {
		enrich_label := ""
		if enrich {
			enrich_label = " (enriched)"
		}
		fmt.Printf("Indexing %s%s via %s (%s)...\n", dir, enrich_label, provider, model)
		if err := indexer.BuildIndex(dir, indexDir, client, enrich); err != nil {
			return fmt.Errorf("building index: %w", err)
		}
		fmt.Printf("Index saved to %s\n", indexDir)
	}

	// If there is nothing else to do, return.
	query := ""
	if len(args) > 0 {
		query = args[0]
	}
	if query == "" && !doBenchmark {
		if !doIndex {
			return fmt.Errorf("provide a query, --benchmark, or --index")
		}
		return nil
	}

	// Load the index and docs.
	col, err := indexer.OpenCollection(indexDir)
	if err != nil {
		return err
	}
	docs, err := indexer.LoadDocs(indexDir)
	if err != nil {
		return err
	}

	if doBenchmark {
		fmt.Printf("Benchmark: %d queries, top-1 accuracy\n", len(benchmark.Queries))
		benchmark.Run(col, docs, client, hybrid)
		return nil
	}

	// Single query.
	return runQuery(cmd.Context(), col, docs, client, query, pkg, topN, hybrid)
}

func runQuery(
	ctx context.Context,
	col *chromem.Collection,
	docs []indexer.Document,
	client embedder.Client,
	query, pkg string,
	topN int,
	hybrid bool,
) error {
	queryVec, err := client.CreateEmbedding(query)
	if err != nil {
		return fmt.Errorf("embedding query: %w", err)
	}

	// Filter docs by package for keyword search and result filtering.
	filtered := search.FilterByPackage(docs, pkg)

	// Build a lookup map.
	docMap := make(map[string]indexer.Document, len(docs))
	for _, d := range docs {
		docMap[d.ID] = d
	}

	if hybrid {
		return runHybrid(ctx, col, filtered, docMap, client, queryVec, query, topN)
	}
	return runDense(ctx, col, filtered, queryVec, topN)
}

func runDense(
	ctx context.Context,
	col *chromem.Collection,
	filtered []indexer.Document,
	queryVec []float32,
	topN int,
) error {
	// Build a set of allowed IDs when package filtering is active.
	var allowSet map[string]bool
	if len(filtered) < col.Count() {
		allowSet = make(map[string]bool, len(filtered))
		for _, d := range filtered {
			allowSet[d.ID] = true
		}
	}

	// Query enough results to survive the filter.
	n := col.Count()
	if n == 0 {
		return fmt.Errorf("index is empty")
	}
	results, err := col.QueryEmbedding(ctx, queryVec, n, nil, nil)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	rank := 1
	for _, r := range results {
		if allowSet != nil && !allowSet[r.ID] {
			continue
		}
		fmt.Printf("%d. %-36s :: %-20s sim: %.4f\n", rank, r.Metadata["file"], r.Metadata["func"], r.Similarity)
		rank++
		if rank > topN {
			break
		}
	}
	return nil
}

func runHybrid(
	ctx context.Context,
	col *chromem.Collection,
	filtered []indexer.Document,
	docMap map[string]indexer.Document,
	client embedder.Client,
	queryVec []float32,
	query string,
	topN int,
) error {
	n := col.Count()
	if n == 0 {
		return fmt.Errorf("index is empty")
	}
	results, err := col.QueryEmbedding(ctx, queryVec, n, nil, nil)
	if err != nil {
		return fmt.Errorf("dense query: %w", err)
	}

	// Build filtered ID set.
	var allowSet map[string]bool
	if len(filtered) < n {
		allowSet = make(map[string]bool, len(filtered))
		for _, d := range filtered {
			allowSet[d.ID] = true
		}
	}

	// Dense ranks (filtered).
	denseRanks := make([]string, 0, len(results))
	for _, r := range results {
		if allowSet == nil || allowSet[r.ID] {
			denseRanks = append(denseRanks, r.ID)
		}
	}

	// Keyword ranks (already filtered, since we pass filtered docs).
	keywordRanks := search.KeywordRanks(query, filtered)

	fused := search.FuseRRF(denseRanks, keywordRanks, docMap)

	for i, r := range fused {
		if i >= topN {
			break
		}
		fmt.Printf("%d. %-36s :: %-20s rrf: %.5f\n", i+1, r.File, r.Func, r.Score)
	}

	return nil
}

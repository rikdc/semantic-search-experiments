package benchmark

import (
	"context"
	"fmt"

	chromem "github.com/philippgille/chromem-go"
	"embedtutorial/indexer"
	"embedtutorial/search"
	"embedtutorial/shared/embedder"
)

// Query pairs a natural-language query with the expected top-1 function ID.
type Query struct {
	Text       string
	ExpectedID string // e.g. "cli/root.go::Execute"
}

// Queries is the benchmark suite: 18 natural-language queries, one per
// expected function. The IDs use the relative path format stored by BuildIndex.
var Queries = []Query{
	{"route an HTTP request to the right handler", "api/handler.go::HandleRequest"},
	{"write a JSON response to an HTTP client", "api/handler.go::WriteJSON"},
	{"return an error as JSON", "api/handler.go::WriteError"},
	{"extract and validate the bearer token from a request", "auth/middleware.go::Authenticate"},
	{"check whether a user has the required permission", "auth/middleware.go::Authorize"},
	{"verify a token signature and check it has not expired", "auth/token.go::ValidateToken"},
	{"issue a new token with a fresh expiry", "auth/token.go::RefreshToken"},
	{"invalidate a token immediately", "auth/token.go::RevokeToken"},
	{"create an in-memory key value cache", "cache/store.go::NewStore"},
	{"store a value with a TTL", "cache/store.go::Set"},
	{"look up a key in the cache", "cache/store.go::Get"},
	{"remove one entry from the cache", "cache/store.go::Delete"},
	{"clear all entries from the cache", "cache/store.go::Flush"},
	{"parse command line arguments into a struct", "cli/flags.go::ParseFlags"},
	{"how long before a request times out", "cli/flags.go::getTimeout"},
	{"run a command", "cli/root.go::Execute"},
	{"how do I connect to the database", "db/connection.go::Connect"},
	{"check if the database is still reachable", "db/connection.go::Ping"},
}

// Run executes all benchmark queries against col and docs. If hybrid is true,
// it uses RRF to combine dense and keyword results; otherwise dense only.
// It prints each result and returns the number of correct top-1 hits.
func Run(col *chromem.Collection, docs []indexer.Document, client embedder.Client, hybrid bool) int {
	ctx := context.Background()

	// Build a lookup map for keyword search and result decoration.
	docMap := make(map[string]indexer.Document, len(docs))
	for _, d := range docs {
		docMap[d.ID] = d
	}

	correct := 0
	fmt.Printf("\n%-55s %-40s %s\n", "Query", "Got (top-1)", "Expected")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────────────────────────────────────────────")

	for _, q := range Queries {
		topID, err := topResult(ctx, col, docs, docMap, client, q.Text, hybrid)
		if err != nil {
			fmt.Printf("  error: %v\n", err)
			continue
		}

		mark := "✗"
		if topID == q.ExpectedID {
			mark = "✓"
			correct++
		}
		fmt.Printf("%s %-53s %-40s %s\n", mark, q.Text, topID, q.ExpectedID)
	}

	total := len(Queries)
	pct := 100 * correct / total
	fmt.Printf("\nTop-1 accuracy: %d / %d (%d%%)\n", correct, total, pct)
	return correct
}

func topResult(
	ctx context.Context,
	col *chromem.Collection,
	docs []indexer.Document,
	docMap map[string]indexer.Document,
	client embedder.Client,
	query string,
	hybrid bool,
) (string, error) {
	queryVec, err := client.CreateEmbedding(query)
	if err != nil {
		return "", fmt.Errorf("embedding query: %w", err)
	}

	n := col.Count()
	if n == 0 {
		return "", fmt.Errorf("collection is empty")
	}

	results, err := col.QueryEmbedding(ctx, queryVec, n, nil, nil)
	if err != nil {
		return "", fmt.Errorf("dense query: %w", err)
	}

	if !hybrid {
		if len(results) == 0 {
			return "", fmt.Errorf("no results")
		}
		return results[0].ID, nil
	}

	// Build dense rank list.
	denseRanks := make([]string, len(results))
	for i, r := range results {
		denseRanks[i] = r.ID
	}

	// Build keyword rank list.
	keywordRanks := search.KeywordRanks(query, docs)

	// Fuse.
	fused := search.FuseRRF(denseRanks, keywordRanks, docMap)
	if len(fused) == 0 {
		return "", fmt.Errorf("no fused results")
	}
	return fused[0].ID, nil
}

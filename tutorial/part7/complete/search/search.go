package search

import (
	"sort"
	"strings"
	"unicode"

	"embedtutorial/indexer"
)

// RankedResult holds a document's position in a fused ranking.
type RankedResult struct {
	ID    string
	Score float64
	File  string
	Func  string
}

// filterByPackage returns only those documents whose File path starts with pkg.
// If pkg is empty, all documents are returned unchanged.
func filterByPackage(docs []indexer.Document, pkg string) []indexer.Document {
	if pkg == "" {
		return docs
	}
	var filtered []indexer.Document
	for _, d := range docs {
		if strings.HasPrefix(d.File, pkg+"/") || strings.HasPrefix(d.File, pkg) {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

// tokenize splits s into lowercase tokens, splitting on non-alphanumeric
// characters and on camelCase boundaries (lowercase letter followed by
// uppercase letter). This ensures that "getTimeout" produces ["get", "timeout"]
// and a query containing "timeout" matches the function name.
func tokenize(s string) []string {
	var parts []string
	current := strings.Builder{}
	runes := []rune(strings.ToLower(s))

	for i, r := range []rune(s) {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			if current.Len() > 0 {
				parts = append(parts, strings.ToLower(current.String()))
				current.Reset()
			}
			continue
		}
		// Split on camelCase boundary: lowercase followed by uppercase.
		if i > 0 && unicode.IsUpper(r) && unicode.IsLower([]rune(s)[i-1]) {
			if current.Len() > 0 {
				parts = append(parts, strings.ToLower(current.String()))
				current.Reset()
			}
		}
		current.WriteRune(runes[i])
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	return parts
}

// tokenizeSet returns the unique tokens of s as a set.
func tokenizeSet(s string) map[string]bool {
	tokens := tokenize(s)
	set := make(map[string]bool, len(tokens))
	for _, t := range tokens {
		set[t] = true
	}
	return set
}

// keywordScore returns the fraction of query tokens that appear in content.
// Repeated query tokens count multiple times in both numerator and denominator.
func keywordScore(query, content string) float32 {
	queryTokens := tokenize(query)
	docTokens := tokenizeSet(content)

	var matches int
	for _, qt := range queryTokens {
		if docTokens[qt] {
			matches++
		}
	}
	if len(queryTokens) == 0 {
		return 0
	}
	return float32(matches) / float32(len(queryTokens))
}

const rrfK = 60

// fuseRRF combines two ranked lists of document IDs using Reciprocal Rank
// Fusion. Each document's score is the sum of 1/(k+rank) over both lists.
// Results are sorted by score descending, with ID as a stable tiebreaker.
// The docs map is used to populate File and Func fields in the result.
func fuseRRF(denseRanks, keywordRanks []string, docs map[string]indexer.Document) []RankedResult {
	scores := make(map[string]float64)

	for i, id := range denseRanks {
		scores[id] += 1.0 / float64(rrfK+i+1)
	}
	for i, id := range keywordRanks {
		scores[id] += 1.0 / float64(rrfK+i+1)
	}

	results := make([]RankedResult, 0, len(scores))
	for id, score := range scores {
		r := RankedResult{ID: id, Score: score}
		if d, ok := docs[id]; ok {
			r.File = d.File
			r.Func = d.Func
		}
		results = append(results, r)
	}

	// Primary: score descending. Secondary: ID ascending for stable tie-breaking.
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].ID < results[j].ID
	})

	return results
}

// KeywordRanks returns all docs sorted by keyword overlap with query, highest
// first. Returns a slice of document IDs in rank order.
func KeywordRanks(query string, docs []indexer.Document) []string {
	type scored struct {
		id    string
		score float32
	}
	results := make([]scored, len(docs))
	for i, d := range docs {
		results[i] = scored{id: d.ID, score: keywordScore(query, d.Content)}
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].score != results[j].score {
			return results[i].score > results[j].score
		}
		return results[i].id < results[j].id
	})
	ids := make([]string, len(results))
	for i, r := range results {
		ids[i] = r.id
	}
	return ids
}

// FilterByPackage is the exported wrapper for use in cmd/.
func FilterByPackage(docs []indexer.Document, pkg string) []indexer.Document {
	return filterByPackage(docs, pkg)
}

// FuseRRF is the exported wrapper for use in cmd/.
func FuseRRF(denseRanks, keywordRanks []string, docs map[string]indexer.Document) []RankedResult {
	return fuseRRF(denseRanks, keywordRanks, docs)
}

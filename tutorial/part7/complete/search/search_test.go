package search

import (
	"math"
	"testing"

	"embedtutorial/indexer"
)

func TestTokenize_CamelCase(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{"getTimeout", []string{"get", "timeout"}},
		{"ValidateToken", []string{"validate", "token"}},
		{"WriteJSON", []string{"write", "json"}},
		{"ParseFlags", []string{"parse", "flags"}},
		{"handleRequest", []string{"handle", "request"}},
		{"timeout", []string{"timeout"}},
	}
	for _, tc := range cases {
		got := tokenize(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("tokenize(%q) = %v, want %v", tc.input, got, tc.want)
			continue
		}
		for i := range tc.want {
			if got[i] != tc.want[i] {
				t.Errorf("tokenize(%q)[%d] = %q, want %q", tc.input, i, got[i], tc.want[i])
			}
		}
	}
}

func TestTokenize_Punctuation(t *testing.T) {
	got := tokenize("get_timeout")
	want := []string{"get", "timeout"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Errorf("tokenize(\"get_timeout\") = %v, want %v", got, want)
	}
}

// TestFuseRRF_ConcreteExample uses the concrete case from the blog post:
//
//	Dense:   [getTimeout, cacheSet, dbPing]
//	Keyword: [getTimeout, dbPing,   cacheSet]
//
// Expected RRF scores:
//
//	getTimeout: 1/61 + 1/61 = 0.03279...
//	dbPing:     1/63 + 1/62 = 0.03200...
//	cacheSet:   1/62 + 1/63 = 0.03200...  (same, tied on score, sorted by ID)
func TestFuseRRF_ConcreteExample(t *testing.T) {
	dense := []string{"getTimeout", "cacheSet", "dbPing"}
	keyword := []string{"getTimeout", "dbPing", "cacheSet"}

	docs := map[string]indexer.Document{
		"getTimeout": {ID: "getTimeout", File: "cli/flags.go", Func: "getTimeout"},
		"cacheSet":   {ID: "cacheSet", File: "cache/store.go", Func: "Set"},
		"dbPing":     {ID: "dbPing", File: "db/connection.go", Func: "Ping"},
	}

	results := fuseRRF(dense, keyword, docs)

	if len(results) != 3 {
		t.Fatalf("want 3 results, got %d", len(results))
	}

	// getTimeout must be first
	if results[0].ID != "getTimeout" {
		t.Errorf("results[0].ID = %q, want %q", results[0].ID, "getTimeout")
	}

	wantTop := 2.0 / float64(rrfK+1)
	if math.Abs(results[0].Score-wantTop) > 1e-9 {
		t.Errorf("results[0].Score = %v, want %v", results[0].Score, wantTop)
	}

	// dbPing and cacheSet tie — sorted by ID ascending: "cacheSet" < "dbPing"
	if results[1].ID != "cacheSet" {
		t.Errorf("results[1].ID = %q, want %q (tie broken by ID)", results[1].ID, "cacheSet")
	}
	if results[2].ID != "dbPing" {
		t.Errorf("results[2].ID = %q, want %q (tie broken by ID)", results[2].ID, "dbPing")
	}

	wantTied := 1.0/float64(rrfK+2) + 1.0/float64(rrfK+3)
	if math.Abs(results[1].Score-wantTied) > 1e-9 {
		t.Errorf("results[1].Score = %v, want %v", results[1].Score, wantTied)
	}
}

func TestFuseRRF_SingleRetriever(t *testing.T) {
	dense := []string{"a", "b", "c"}
	docs := map[string]indexer.Document{
		"a": {ID: "a"}, "b": {ID: "b"}, "c": {ID: "c"},
	}
	results := fuseRRF(dense, nil, docs)
	if results[0].ID != "a" {
		t.Errorf("results[0].ID = %q, want %q", results[0].ID, "a")
	}
}

func TestFilterByPackage(t *testing.T) {
	docs := []indexer.Document{
		{ID: "1", File: "cli/root.go", PkgDir: "cli"},
		{ID: "2", File: "cli/flags.go", PkgDir: "cli"},
		{ID: "3", File: "db/query.go", PkgDir: "db"},
	}

	got := filterByPackage(docs, "cli")
	if len(got) != 2 {
		t.Fatalf("want 2, got %d", len(got))
	}
	for _, d := range got {
		if d.PkgDir != "cli" {
			t.Errorf("unexpected doc: %+v", d)
		}
	}
}

func TestFilterByPackage_Empty(t *testing.T) {
	docs := []indexer.Document{{ID: "1", File: "cli/root.go"}}
	got := filterByPackage(docs, "")
	if len(got) != 1 {
		t.Errorf("empty pkg filter should return all docs")
	}
}

func TestKeywordScore_ExactMatch(t *testing.T) {
	score := keywordScore("timeout", "func getTimeout() time.Duration { return 30 * time.Second }")
	if score <= 0 {
		t.Errorf("expected positive score for matching token, got %v", score)
	}
}

func TestKeywordScore_NoMatch(t *testing.T) {
	score := keywordScore("timeout", "func Set(key string, value any) {}")
	if score != 0 {
		t.Errorf("expected zero score for non-matching content, got %v", score)
	}
}

package indexer_test

import (
	"testing"

	"part6/indexer"
)

// handcraftedFile is a small Go source with three functions: two exported, one unexported.
// The expected output is deterministic — use it to verify ParseFunctions before
// touching any embedding API.
const handcraftedFile = `package example

import "fmt"

// Greet returns a greeting for the given name.
func Greet(name string) string {
	return fmt.Sprintf("hello, %s", name)
}

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

// helper is an unexported utility used internally.
func helper() bool {
	return true
}
`

func TestParseFunctions_Count(t *testing.T) {
	docs, err := indexer.ParseFunctions("example.go", []byte(handcraftedFile))
	if err != nil {
		t.Fatalf("ParseFunctions returned error: %v", err)
	}

	if len(docs) != 3 {
		t.Errorf("expected 3 documents, got %d", len(docs))
	}
}

func TestParseFunctions_FunctionNames(t *testing.T) {
	docs, err := indexer.ParseFunctions("example.go", []byte(handcraftedFile))
	if err != nil {
		t.Fatalf("ParseFunctions returned error: %v", err)
	}

	wantNames := []string{"Greet", "Add", "helper"}
	for i, want := range wantNames {
		if i >= len(docs) {
			t.Errorf("missing document at index %d (expected func %q)", i, want)
			continue
		}
		if docs[i].Func != want {
			t.Errorf("docs[%d].Func = %q, want %q", i, docs[i].Func, want)
		}
	}
}

func TestParseFunctions_IDs(t *testing.T) {
	docs, err := indexer.ParseFunctions("example.go", []byte(handcraftedFile))
	if err != nil {
		t.Fatalf("ParseFunctions returned error: %v", err)
	}

	for _, doc := range docs {
		wantID := "example.go::" + doc.Func
		if doc.ID != wantID {
			t.Errorf("doc.ID = %q, want %q", doc.ID, wantID)
		}
	}
}

func TestParseFunctions_ContentContainsFuncKeyword(t *testing.T) {
	docs, err := indexer.ParseFunctions("example.go", []byte(handcraftedFile))
	if err != nil {
		t.Fatalf("ParseFunctions returned error: %v", err)
	}

	for _, doc := range docs {
		if len(doc.Content) == 0 {
			t.Errorf("doc %q has empty Content", doc.ID)
		}
		// The content should start with "func"
		if len(doc.Content) < 4 || doc.Content[:4] != "func" {
			t.Errorf("doc %q Content does not start with 'func': %q", doc.ID, doc.Content[:min(20, len(doc.Content))])
		}
	}
}

func TestParseFunctions_EmptyFile(t *testing.T) {
	src := []byte("package empty\n")
	docs, err := indexer.ParseFunctions("empty.go", []byte(src))
	if err != nil {
		t.Fatalf("ParseFunctions returned error on empty file: %v", err)
	}
	if len(docs) != 0 {
		t.Errorf("expected 0 documents for file with no functions, got %d", len(docs))
	}
}

func TestParseFunctions_InvalidGo(t *testing.T) {
	src := []byte("this is not valid go code {{{")
	_, err := indexer.ParseFunctions("bad.go", src)
	if err == nil {
		t.Error("expected an error for invalid Go source, got nil")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

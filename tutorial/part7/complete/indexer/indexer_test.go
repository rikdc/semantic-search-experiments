package indexer

import (
	"testing"
)

// handcrafted source with three functions: two exported, one unexported.
const testSrc = `package example

// Greet returns a greeting string.
func Greet(name string) string {
	return "Hello, " + name
}

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

func helper() {}
`

func TestParseFunctions_Count(t *testing.T) {
	docs, err := ParseFunctions("testfile.go", []byte(testSrc), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 3 {
		t.Fatalf("want 3 documents, got %d", len(docs))
	}
}

func TestParseFunctions_FunctionNames(t *testing.T) {
	docs, err := ParseFunctions("testfile.go", []byte(testSrc), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{"Greet", "Add", "helper"}
	for i, w := range want {
		if docs[i].Func != w {
			t.Errorf("doc[%d].Func = %q, want %q", i, docs[i].Func, w)
		}
	}
}

func TestParseFunctions_IDs(t *testing.T) {
	docs, err := ParseFunctions("testfile.go", []byte(testSrc), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []string{
		"testfile.go::Greet",
		"testfile.go::Add",
		"testfile.go::helper",
	}
	for i, w := range want {
		if docs[i].ID != w {
			t.Errorf("doc[%d].ID = %q, want %q", i, docs[i].ID, w)
		}
	}
}

func TestParseFunctions_ContentContainsFuncKeyword(t *testing.T) {
	docs, err := ParseFunctions("testfile.go", []byte(testSrc), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, doc := range docs {
		if len(doc.Content) < 4 || doc.Content[:4] != "func" {
			t.Errorf("doc %q Content does not start with 'func': %q", doc.ID, doc.Content[:min(20, len(doc.Content))])
		}
	}
}

func TestParseFunctions_EmptyFile(t *testing.T) {
	src := "package example\n"
	docs, err := ParseFunctions("empty.go", []byte(src), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 0 {
		t.Errorf("want 0 documents, got %d", len(docs))
	}
}

func TestParseFunctions_InvalidGo(t *testing.T) {
	src := "this is not go code {"
	_, err := ParseFunctions("bad.go", []byte(src), false)
	if err == nil {
		t.Fatal("expected error for invalid Go source, got nil")
	}
}

func TestEnrichContent(t *testing.T) {
	body := "func Execute() error {\n\treturn nil\n}"

	got := enrichContent("cli/root.go", "cli", "", "Execute", body)
	want := "File: cli/root.go | Package: cli | Function: Execute\n\n" + body
	if got != want {
		t.Errorf("enrichContent without receiver:\ngot:  %q\nwant: %q", got, want)
	}

	got = enrichContent("db/store.go", "db", "*Store", "Set", body)
	want = "File: db/store.go | Package: db | Function: (*Store) Set\n\n" + body
	if got != want {
		t.Errorf("enrichContent with receiver:\ngot:  %q\nwant: %q", got, want)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Package indexer walks a Go source tree, extracts function declarations,
// embeds them, and stores the results in a chromem-go collection for search.
package indexer

import (
	"context"
	"fmt"

	chromem "github.com/philippgille/chromem-go"

	"part6/shared/embedder"
)

// Document represents a single indexed function.
type Document struct {
	ID      string // "path/to/file.go::FunctionName"
	Content string // full function source text
	File    string // relative file path
	Func    string // function name
}

// ParseFunctions extracts all function declarations from a Go source file.
// Each function body becomes one Document. Both exported and unexported functions
// are included.
//
// TODO: implement ParseFunctions
//
// You will need these imports:
//   - "go/ast"
//   - "go/parser"
//   - "go/token"
//
// Steps:
//  1. Create a token.NewFileSet() to track source positions.
//  2. Call parser.ParseFile(fset, path, src, parser.ParseComments) to get an *ast.File.
//     Return the error if parsing fails.
//  3. Range over f.Decls. Use a type assertion to select *ast.FuncDecl values only.
//  4. For each function, compute byte offsets with:
//     start := fset.Position(fn.Pos()).Offset
//     end   := fset.Position(fn.End()).Offset
//  5. Build a Document with:
//     ID:      fmt.Sprintf("%s::%s", path, fn.Name.Name)
//     Content: string(src[start:end])
//     File:    path
//     Func:    fn.Name.Name
func ParseFunctions(path string, src []byte) ([]Document, error) {
	panic("not implemented: ParseFunctions")
}

// BuildIndex walks dir, extracts all Go functions, embeds each one, and stores
// the results in collection. Progress is printed to stdout as each file is processed.
//
// TODO: implement BuildIndex
//
// You will need these imports:
//   - "os"
//   - "path/filepath"
//   - "strings"
//
// Steps:
//  1. Use filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error)
//     to visit every file. Skip directories and files that don't end in ".go".
//  2. Read each file with os.ReadFile(path).
//  3. Call ParseFunctions(path, src) to extract its Documents.
//     Skip files that return an error (some generated files won't parse cleanly).
//  4. For each Document, call emb.CreateEmbedding(ctx, doc.Content) to get a []float32.
//  5. Store the result with collection.AddDocuments:
//     collection.AddDocuments(ctx, []chromem.Document{{
//         ID:        doc.ID,
//         Content:   doc.Content,
//         Embedding: vec,
//         Metadata:  map[string]string{"file": doc.File, "func": doc.Func},
//     }}, 1)
func BuildIndex(ctx context.Context, dir string, collection *chromem.Collection, emb embedder.Embedder) error {
	panic("not implemented: BuildIndex")
}

// Search queries the collection for the top-k functions most similar to queryText.
// The collection's embedding function embeds the query automatically.
func Search(ctx context.Context, queryText string, topK int, collection *chromem.Collection) ([]chromem.Result, error) {
	results, err := collection.Query(ctx, queryText, topK, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("querying collection: %w", err)
	}
	return results, nil
}

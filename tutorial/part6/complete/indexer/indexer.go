// Package indexer walks a Go source tree, extracts function declarations,
// embeds them, and stores the results in a chromem-go collection for search.
package indexer

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

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
// are included — code search is most useful when it covers the full codebase.
func ParseFunctions(path string, src []byte) ([]Document, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	var docs []Document
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		start := fset.Position(fn.Pos()).Offset
		end := fset.Position(fn.End()).Offset
		body := string(src[start:end])

		docs = append(docs, Document{
			ID:      fmt.Sprintf("%s::%s", path, fn.Name.Name),
			Content: body,
			File:    path,
			Func:    fn.Name.Name,
		})
	}
	return docs, nil
}

// BuildIndex walks dir, extracts all Go functions, embeds each one, and stores
// the results in collection. Progress is printed to stdout as each file is processed.
func BuildIndex(ctx context.Context, dir string, collection *chromem.Collection, emb embedder.Embedder) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		docs, err := ParseFunctions(path, src)
		if err != nil {
			// Skip files that don't parse cleanly (generated code, build tags, etc.)
			fmt.Printf("  skipping %s: %v\n", path, err)
			return nil
		}
		if len(docs) == 0 {
			return nil
		}

		// Print the file and its functions as we go.
		names := make([]string, len(docs))
		for i, d := range docs {
			names[i] = d.Func
		}
		fmt.Printf("  %-36s → %s\n", path, strings.Join(names, ", "))

		// Embed and store each function.
		for _, doc := range docs {
			vec, err := emb.CreateEmbedding(ctx, doc.Content)
			if err != nil {
				return fmt.Errorf("embedding %s: %w", doc.ID, err)
			}

			err = collection.AddDocuments(ctx, []chromem.Document{
				{
					ID:        doc.ID,
					Content:   doc.Content,
					Embedding: vec,
					Metadata: map[string]string{
						"file": doc.File,
						"func": doc.Func,
					},
				},
			}, 1)
			if err != nil {
				return fmt.Errorf("storing %s: %w", doc.ID, err)
			}
		}

		return nil
	})
}

// Search queries the collection for the top-k functions most similar to queryText.
// The collection's embedding function is used to embed the query before searching.
func Search(ctx context.Context, queryText string, topK int, collection *chromem.Collection) ([]chromem.Result, error) {
	results, err := collection.Query(ctx, queryText, topK, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("querying collection: %w", err)
	}
	return results, nil
}

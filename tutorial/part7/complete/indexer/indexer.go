package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	chromem "github.com/philippgille/chromem-go"
	"embedtutorial/shared/embedder"
)

// Document represents a single indexed function.
type Document struct {
	ID       string `json:"id"`
	Content  string `json:"content"` // the text that was embedded
	File     string `json:"file"`    // relative path, e.g. "cli/root.go"
	PkgDir   string `json:"pkg_dir"` // first path component, e.g. "cli"
	Func     string `json:"func"`
	Receiver string `json:"receiver"`
}

// enrichContent builds a synthetic string that prepends file, package, and
// function context to the raw body. This shifts the embedding away from the
// bare function text and toward the domain the function belongs to.
func enrichContent(file, pkg, receiver, name, body string) string {
	if receiver != "" {
		return fmt.Sprintf(
			"File: %s | Package: %s | Function: (%s) %s\n\n%s",
			file, pkg, receiver, name, body,
		)
	}
	return fmt.Sprintf(
		"File: %s | Package: %s | Function: %s\n\n%s",
		file, pkg, name, body,
	)
}

// ParseFunctions parses src (Go source at path) and returns one Document per
// function declaration. If enrich is true, each document's Content is the
// enriched synthetic string; otherwise it is the raw function body text.
func ParseFunctions(path string, src []byte, enrich bool) ([]Document, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	pkgDir := strings.SplitN(path, "/", 2)[0]

	var docs []Document
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		start := fset.Position(fn.Pos()).Offset
		end := fset.Position(fn.End()).Offset
		body := string(src[start:end])

		var receiver string
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			field := fn.Recv.List[0]
			switch t := field.Type.(type) {
			case *ast.StarExpr:
				if ident, ok := t.X.(*ast.Ident); ok {
					receiver = "*" + ident.Name
				}
			case *ast.Ident:
				receiver = t.Name
			}
		}

		content := body
		if enrich {
			content = enrichContent(path, pkgDir, receiver, fn.Name.Name, body)
		}

		docs = append(docs, Document{
			ID:       fmt.Sprintf("%s::%s", path, fn.Name.Name),
			Content:  content,
			File:     path,
			PkgDir:   pkgDir,
			Func:     fn.Name.Name,
			Receiver: receiver,
		})
	}
	return docs, nil
}

// BuildIndex walks dir, parses each .go file, embeds each function, and stores
// the results in a chromem-go collection persisted at indexDir. A docs.json
// file is also written to indexDir for later keyword search.
func BuildIndex(dir, indexDir string, client embedder.Client, enrich bool) error {
	db, err := chromem.NewPersistentDB(indexDir, false)
	if err != nil {
		return fmt.Errorf("opening index at %s: %w", indexDir, err)
	}

	// Drop and recreate the collection so re-indexing is clean.
	db.DeleteCollection("functions")
	col, err := db.GetOrCreateCollection("functions", nil, nil)
	if err != nil {
		return fmt.Errorf("creating collection: %w", err)
	}

	ctx := context.Background()
	var allDocs []Document

	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
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

		// Store relative paths so IDs are portable.
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = path
		}
		// Normalise to forward slashes.
		relPath = filepath.ToSlash(relPath)

		docs, err := ParseFunctions(relPath, src, enrich)
		if err != nil {
			fmt.Printf("  skipping %s: %v\n", relPath, err)
			return nil
		}
		if len(docs) == 0 {
			return nil
		}

		names := make([]string, len(docs))
		for i, d := range docs {
			names[i] = d.Func
		}
		fmt.Printf("  %-40s → %s\n", relPath, strings.Join(names, ", "))

		for _, doc := range docs {
			vec, err := client.CreateEmbedding(doc.Content)
			if err != nil {
				return fmt.Errorf("embedding %s: %w", doc.ID, err)
			}

			err = col.AddDocument(ctx, chromem.Document{
				ID:        doc.ID,
				Content:   doc.Content,
				Embedding: vec,
				Metadata: map[string]string{
					"file":    doc.File,
					"func":    doc.Func,
					"pkg_dir": doc.PkgDir,
				},
			})
			if err != nil {
				return fmt.Errorf("storing %s: %w", doc.ID, err)
			}

			allDocs = append(allDocs, doc)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return saveDocs(indexDir, allDocs)
}

// saveDocs writes all documents to indexDir/docs.json for keyword search.
func saveDocs(indexDir string, docs []Document) error {
	if err := os.MkdirAll(indexDir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(docs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(indexDir, "docs.json"), data, 0644)
}

// LoadDocs reads the docs.json written by BuildIndex.
func LoadDocs(indexDir string) ([]Document, error) {
	data, err := os.ReadFile(filepath.Join(indexDir, "docs.json"))
	if err != nil {
		return nil, fmt.Errorf("loading docs (run --index first): %w", err)
	}
	var docs []Document
	if err := json.Unmarshal(data, &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

// OpenCollection opens an existing chromem-go collection from indexDir.
func OpenCollection(indexDir string) (*chromem.Collection, error) {
	db, err := chromem.NewPersistentDB(indexDir, false)
	if err != nil {
		return nil, fmt.Errorf("opening index: %w", err)
	}
	col := db.GetCollection("functions", nil)
	if col == nil {
		return nil, fmt.Errorf("collection not found (run --index first)")
	}
	return col, nil
}

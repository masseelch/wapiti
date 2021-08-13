package wapiti

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/masseelch/wapiti/wapiti/config"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

type Wapiti struct {
	cfg   *config.Config
	graph *gen.Graph
	fset  *token.FileSet
	ast   *ast.Package
}

func New(cfg *config.Config) (*Wapiti, error) {
	g, err := entc.LoadGraph(cfg.SchemaPath, &gen.Config{})
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	tree, err := parser.ParseDir(fset, cfg.SchemaPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return &Wapiti{
		cfg:   cfg,
		graph: g,
		fset:  fset,
		ast:   tree[filepath.Base(g.Schema)],
	}, nil
}

// Reload reloads the ent graph from the file system and re-parses the ast.
func (w *Wapiti) Reload() error {
	n, err := New(w.cfg)
	if err != nil {
		return err
	}
	w.graph = n.graph
	w.fset = n.fset
	w.ast = n.ast
	return nil
}

// Run runs the interactive cli, asking questions how to change the schema.
func (w *Wapiti) Run() error {
	n, err := w.SelectNode()
	if err != nil {
		return err
	}
	if n == nil {
		// fmt.Printf(greenBG("\t\t\t\nSuccess\n\t\t\t"))
	}
	fmt.Println(n)
	return nil
}

package wapiti

import (
	"entgo.io/ent/entc/load"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/masseelch/wapiti/wapiti/config"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

type Wapiti struct {
	cfg  *config.Config
	spec *load.SchemaSpec
	fset *token.FileSet
	ast  *ast.Package
}

func New(cfg *config.Config) (*Wapiti, error) {
	spec, err := (&load.Config{Path: cfg.SchemaPath}).Load()
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	tree, err := parser.ParseDir(fset, cfg.SchemaPath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return &Wapiti{
		cfg:  cfg,
		spec: spec,
		fset: fset,
		ast:  tree[filepath.Base(spec.PkgPath)],
	}, nil
}

// Reload reloads the ent spec from the file system and re-parses the ast.
func (w *Wapiti) Reload() error {
	n, err := New(w.cfg)
	if err != nil {
		return err
	}
	w.spec = n.spec
	w.fset = n.fset
	w.ast = n.ast
	return nil
}

// Run runs the interactive cli, asking questions how to change the schema.
func (w *Wapiti) Run() error {
	// Select the node to edit.
	n, err := w.SelectNode()
	if err != nil {
		return fmt.Errorf(aurora.Red("ERROR: %w").String(), err)
	}
	if n == nil {
		fmt.Println(aurora.Green("Success!").Bold())
	}
	for {
		f, err := w.NewField(n)
		if err != nil {
			return err
		}
		if f == nil {
			break
		}
		fmt.Println("TODO: Add message here")
	}
	return nil
}

package wapiti

import (
	"bytes"
	"entgo.io/ent/entc/load"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const schemaCreatedFormat = `
created: %s

Schema generated! Now let's add some fields and edges!
You can always add more fields and edges later manually or by re-running this command.

`

var schemaTpl = template.Must(template.New("schema").Parse(`package schema

import "entgo.io/ent"

// {{ . }} holds the schema definition for the {{ . }} entity.
type {{ . }} struct {
	ent.Schema
}
`))

var fieldsTpl = template.Must(template.New("fields").Parse(`
// Fields of the {{ . }}.
func ({{ . }}) Fields() []ent.Field {
	return nil
}
`))

// CreateSchema creates a new schema with the given name and writes it to file. Calls Reload afterwards.
func (w *Wapiti) CreateSchema(name string) error {
	b := new(bytes.Buffer)
	if err := schemaTpl.Execute(b, name); err != nil {
		return fmt.Errorf("executing template %s: %w", name, err)
	}
	f := filepath.Join(w.cfg.SchemaPath, strings.ToLower(name+".go"))
	if err := ioutil.WriteFile(f, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing file %s: %w", f, err)
	}
	fmt.Printf(schemaCreatedFormat, aurora.Cyan(f))
	return w.Reload()
}

// AddField adds the field to the schema. Calls Reload afterwards.
func (w *Wapiti) AddField(f *Field) (*load.Field, error) {
	// Find 'Fields()'-method of our schema.
	file, fieldsDecl := w.fieldsMethod(f.Schema)
	// If there is no 'Fields()'-method add one.
	if fieldsDecl == nil {
		// Get the file our schema resides in.
		file = w.file(f.Schema)
		// Add the 'Fields()'-method.
		b := bytes.NewBuffer([]byte("package main\n"))
		if err := fieldsTpl.Execute(b, f.Schema.Name); err != nil {
			return nil, err
		}
		expr, err := parser.ParseFile(token.NewFileSet(), "", b.String(), parser.ParseComments)
		if err != nil {
			return nil, err
		}
		file.Decls = append(file.Decls, expr.Decls...)
		printer.Fprint(os.Stdout, w.fset, file)
	} else {
		fmt.Println("fields already exist")
	}
	return nil, nil
}

// schemaNode extracts the ast.FuncDecl and the ast.File it was found in of the 'Fields()'-method for the load.Schema.
func (w *Wapiti) fieldsMethod(s *load.Schema) (*ast.File, *ast.FuncDecl) {
	var decl *ast.FuncDecl
	var file *ast.File
	for _, f := range w.ast.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			if fn, ok := n.(*ast.FuncDecl); ok {
				if fn.Name.Name == "Fields" && fn.Recv != nil && len(fn.Recv.List) == 1 {
					if r, ok := fn.Recv.List[0].Type.(*ast.Ident); ok && r.Name == s.Name {
						file = f
						decl = fn
						return false
					}
				}
			}
			return true
		})
	}
	return file, decl
}

// schemaNode extracts the ast.File the given load.Schema resides in.
func (w *Wapiti) file(s *load.Schema) *ast.File {
	for _, f := range w.ast.Files {
		for _, decl := range f.Decls {
			if genDecl, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range genDecl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						if _, ok := typeSpec.Type.(*ast.StructType); ok {
							if typeSpec.Name.Name == s.Name {
								return f
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (w *Wapiti) printNode(node interface{}) {
	printer.Fprint(os.Stdout, w.fset, node)
}

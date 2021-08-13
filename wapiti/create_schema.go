package wapiti

import (
	"bytes"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

const schemaCreatedFormat = `
created: %s

Schema generated! Now let's add some fields and edges!
You can always add more fields and edges later manually or by re-running this command.

`

var tpl = template.Must(template.New("schema").Parse(`package schema

import "entgo.io/ent"

// {{ . }} holds the schema definition for the {{ . }} entity.
type {{ . }} struct {
	graph.Schema
}
`))

// CreateSchema creates a new schema with the given name and writes it to file. Reload the Graph afterwards.
func (w *Wapiti) CreateSchema(name string) error {
	b := new(bytes.Buffer)
	if err := tpl.Execute(b, name); err != nil {
		return fmt.Errorf("executing template %s: %w", name, err)
	}
	f := filepath.Join(w.cfg.SchemaPath, strings.ToLower(name+".go"))
	if err := ioutil.WriteFile(f, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing file %s: %w", f, err)
	}
	fmt.Printf(schemaCreatedFormat, aurora.Cyan(f))
	return w.Reload()
}

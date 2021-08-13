package wapiti

import (
	"entgo.io/ent/entc/gen"
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/masseelch/wapiti/wapiti/sillyname"
	"regexp"
	"strings"
)

const prefix = "> "

var nameRgx = regexp.MustCompile("[A-Z][A-Za-z]*")

func (w *Wapiti) SelectNode() (*gen.Type, error) {
	// Suggest already created schemas.
	var s []prompt.Suggest
	for _, n := range w.graph.Nodes {
		s = append(s, prompt.Suggest{Text: n.Name, Description: fmt.Sprintf("Edit the %s node", n.Name)})
	}
	n := ask(s, "Name of the schema to create or update (e.g. %s):", aurora.Yellow(sillyname.Schema()))
	// If the user sent us nothing stop the execution.
	if n == "" {
		return nil, nil
	}
	// Validate the name.
	if !nameRgx.MatchString(n) {
		return nil, errors.New("schema names must begin with uppercase")
	}
	// If there is no schema of this name yet create it.
	t := w.LookupNode(n)
	if t == nil {
		// Create a new schema.
		if err := w.CreateSchema(n); err != nil {
			return nil, err
		}
		t = w.LookupNode(n)
	}
	return t, nil
}

// LookupNode looks for a node with the name and returns it. nil if no such node exists.
func (w *Wapiti) LookupNode(n string) *gen.Type {
	for _, t := range w.graph.Nodes {
		if t.Name == n {
			return t
		}
	}
	return nil
}

// ask asks the user the given question and returns the answer. Ensures the answer has a fresh line.
func ask(s []prompt.Suggest, question string, args ...interface{}) string {
	if !strings.HasSuffix(question, "\n") {
		question += "\n"
	}
	fmt.Printf(aurora.Sprintf(aurora.Cyan(question), args...))
	return prompt.Input(prefix, func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	})
}

package wapiti

import (
	"entgo.io/ent/entc/load"
	"errors"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/logrusorgru/aurora/v3"
	"github.com/masseelch/wapiti/wapiti/sillyname"
	"regexp"
	"strings"
)

const (
	prefix      = "> "
	defaultType = "string"
)

var (
	nodeNameRgx  = regexp.MustCompile("[A-Z][A-Za-z]*")
	fieldNameRgx = regexp.MustCompile("[A-Za-z][A-Za-z1-9_-]*")
	types        = []prompt.Suggest{
		{Text: "int"}, {Text: "uint"},
		{Text: "int8"}, {Text: "int16"}, {Text: "int32"}, {Text: "int64"},
		{Text: "uint8"}, {Text: "uint16"}, {Text: "uint32"}, {Text: "uint64"},
		{Text: "float"}, {Text: "float32"},
		{Text: "bool"},
		{Text: "string"}, {Text: "text"},
		{Text: "time"},
		{Text: "uuid"},
		{Text: "[]byte"},
		{Text: "json"},
		{Text: "enum"},
		{Text: "other"},
	}
	yesNo = []prompt.Suggest{{Text: "yes"}, {Text: "no"}}
)

// SelectNode asks the user what node to edit. If an unknown name is given a new schema with that name is created and
// returned.
func (w *Wapiti) SelectNode() (*load.Schema, error) {
	// Suggest already created schemas.
	var sgst []prompt.Suggest
	for _, n := range w.spec.Schemas {
		sgst = append(sgst, prompt.Suggest{Text: n.Name, Description: fmt.Sprintf("Edit the %s node", n.Name)})
	}
	name := ask(sgst, "Name of the schema to create or update (e.g. %s):", aurora.Yellow(sillyname.Schema()))
	// If the user sent us nothing stop the execution.
	if name == "" {
		return nil, nil
	}
	// Validate the name.
	if !nodeNameRgx.MatchString(name) {
		return nil, errors.New("schema names must begin with uppercase and contain only letters")
	}
	// If there is no schema of this name yet, create it.
	s := w.LookupNode(name)
	if s == nil {
		// Create a new schema.
		if err := w.CreateSchema(name); err != nil {
			return nil, err
		}
		s = w.LookupNode(name)
	}
	return s, nil
}

// LookupNode looks for a node with the name and returns it. nil if no such node exists.
func (w *Wapiti) LookupNode(n string) *load.Schema {
	for _, s := range w.spec.Schemas {
		if s.Name == n {
			return s
		}
	}
	return nil
}

type Field struct {
	Schema    *load.Schema
	Name      string
	Type      string
	Optional  bool
	Nillable  bool
	Immutable bool
}

// NewField asks the user what field to add to the given node. Returns nil, nil if the user wants to stop adding fields.
func (w *Wapiti) NewField(s *load.Schema) (*load.Field, error) {
	f := &Field{Schema: s}
	// Ask for the field name.
	f.Name = ask(nil, "Name of the field to add (press <return> to stop adding fields):")
	// If the user sent us nothing stop the execution.
	if f.Name == "" {
		return nil, nil
	}
	// Validate the name.
	if !fieldNameRgx.MatchString(f.Name) {
		return nil, errors.New("field names must begin with a letter and only contain alphanumeric characters, underscores and dashes")
	}
	// Ask for the field type.
	f.Type = ask(types, "Field type [%s]:", aurora.Yellow(defaultType))
	if f.Type == "" {
		f.Type = defaultType
	}
	// TODO: Maybe only ask "Do you want set any other option?" - if so wizard them to all the rest (optional, nil, immutable, storage_key, struct_tag)
	// Ask if the field is optional.
	f.Optional = "no" == ask(yesNo, "Is this field required on creation (optional) (yes/no) [%s]", aurora.Yellow("yes"))
	// Ask if the field is nillable.
	f.Nillable = "yes" == ask(yesNo, "Can this field be nil (yes/no) (nillable) [%s]", aurora.Yellow("no"))
	// Ask if the field is immutable.
	f.Immutable = "yes" == ask(yesNo, "Can this field be updated after creation (immutable) (yes/no) [%s]", aurora.Yellow("no"))
	// Add the field to the schema.
	return w.AddField(f)
}

// ask asks the user the given question and returns the answer. Ensures question and prompt have a fresh line.
func ask(s []prompt.Suggest, question string, args ...interface{}) string {
	if !strings.HasPrefix(question, "\n") {
		question = "\n" + question
	}
	if !strings.HasSuffix(question, "\n") {
		question += "\n"
	}
	fmt.Printf(aurora.Sprintf(aurora.Cyan(question), args...))
	return prompt.Input(prefix, func(d prompt.Document) []prompt.Suggest {
		return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
	})
}

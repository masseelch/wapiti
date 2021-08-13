package schema

import "entgo.io/ent"

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

func (Pet) Fields() []ent.Field {
	return nil
}

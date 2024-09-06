package schema

import "github.com/actord/actord/pkg/process/logic"

type Trigger struct {
	Name       string        `hcl:"name,label"`
	EventType  string        `hcl:"event_type"`
	Logic      logic.Program `hcl:"logic,block"`
	Transition *string       `hcl:"transition"`
}

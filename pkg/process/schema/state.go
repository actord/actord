package schema

import (
	"github.com/actord/actord/pkg/process/execontext"
	"github.com/actord/actord/pkg/process/logic"
)

type State struct {
	Name  string        `hcl:",label"`
	Logic logic.Program `hcl:"logic,block"`
}

func (s *State) Execute(ctx *execontext.ExecutionContext) error {
	// TODO: maybe Execute should return next tick instruction instead of modifying context

	ctx.CurrentState = s.Name
	return s.Logic.Execute(ctx)
}

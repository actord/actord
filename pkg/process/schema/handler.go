package schema

import (
	"github.com/actord/actord/pkg/process/execontext"
	"github.com/actord/actord/pkg/process/logic"
)

type Handler struct {
	Name  string        `hcl:",label"`
	Logic logic.Program `hcl:"logic,block"`
}

func (s *Handler) Execute(ctx *execontext.ExecutionContext) error {

	return s.Logic.Execute(ctx)
}

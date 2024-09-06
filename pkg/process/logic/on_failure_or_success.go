package logic

import (
	"errors"

	"github.com/actord/actord/pkg/process/execontext"
)

type OnFailure struct {
	Transition *string `hcl:"transition"`
	Exception  *string `hcl:"exception"`
}

func (o OnFailure) Execute(ctx *execontext.ExecutionContext) error {
	if o.Transition != nil {
		ctx.CurrentState = *o.Transition
		ctx.ShouldTransit = true
		return nil
	}
	if o.Exception != nil {
		return errors.New(*o.Exception)
	}
	return errors.New("exception or transition must be set")
}

type OnSuccess struct {
	Transition *string `hcl:"transition"`
	Exception  *string `hcl:"exception"`
}

func (o OnSuccess) Execute(ctx *execontext.ExecutionContext) error {
	if o.Transition != nil {
		ctx.CurrentState = *o.Transition
		ctx.ShouldTransit = true
		return nil
	}
	if o.Exception != nil {
		return errors.New(*o.Exception)
	}
	return errors.New("exception or transition must be set")
}

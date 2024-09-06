package logic

import (
	"errors"

	"github.com/actord/actord/pkg/process/execontext"
)

type Await struct {
	Events    []AwaitEvent        `hcl:"event,block"`
	OnUnknown AwaitEventOnUnknown `hcl:"on_unknown,block"`
}

func (a *Await) Execute(ctx *execontext.ExecutionContext) error {
	if !ctx.AwaitEvent {
		ctx.AwaitEvent = true
		return nil
	}
	if !ctx.HasEventData() {
		return errors.New("no event data")
	}

	event := ctx.GetEvent()
	if event == nil {
		return errors.New("event not provided")
	}

	for _, e := range a.Events {
		if e.Name == event.Type.Name {
			return e.Execute(ctx)
		}
	}

	return a.OnUnknown.Execute(ctx)
}

type AwaitEvent struct {
	Name       string  `hcl:"name,label"`
	Transition *string `hcl:"transition"`
	Handler    *string `hcl:"handler"`
}

func (a AwaitEvent) Execute(ctx *execontext.ExecutionContext) error {
	ctx.AwaitEvent = false
	if a.Transition != nil {
		ctx.CurrentState = *a.Transition
		ctx.ShouldTransit = true
		return nil
	}

	if a.Handler != nil {
		ctx.CurrentHandler = *a.Handler
		ctx.ShouldHandle = true
		return nil
	}

	// there is no error - just unblock the execution and let the state machine to transit to next logic block
	return nil
}

type AwaitEventOnUnknown struct {
	Exception  *string `hcl:"exception"`
	Transition *string `hcl:"transition"`
}

func (a AwaitEventOnUnknown) Execute(ctx *execontext.ExecutionContext) error {
	if a.Exception != nil {
		return errors.New(*a.Exception)
	}
	if a.Transition != nil {
		ctx.CurrentState = *a.Transition
		ctx.ShouldTransit = true
		return nil
	}
	return errors.New("do not know what to do with event")
}

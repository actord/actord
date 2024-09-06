package schema

import (
	"fmt"

	"github.com/actord/actord/pkg/process/execontext"
)

type Process struct {
	States   []State   `hcl:"state,block"`
	Handlers []Handler `hcl:"handler,block"`
}

func (p *Process) Execute(ctx *execontext.ExecutionContext) error {
	for range 1000 {
		var err error
		if ctx.CurrentHandler != "" {
			currentHandler := ctx.CurrentHandler
			if err := p.ExecuteHandle(ctx); err != nil {
				return fmt.Errorf("failed to execute handler '%s': %w", currentHandler, err)
			}
			if currentHandler != ctx.CurrentHandler {
				continue
			} else {
				ctx.CurrentHandler = ""
				ctx.ShouldHandle = false
			}
			if ctx.AwaitEvent {
				return fmt.Errorf("await event in handler is not allowed")
			}
		}

		err = p.ExecuteState(ctx)
		if err != nil {
			return err
		}
		if ctx.ShouldTransit {
			ctx.ShouldTransit = false
			continue
		}

		if ctx.AwaitEvent {
			return nil
		}
	}

	// todo: maybe here we should block process until fix the issue in program logic
	return fmt.Errorf("max recursion depth reached")
}

func (p *Process) ExecuteState(ctx *execontext.ExecutionContext) error {
	stateIndex := -1
	for i, s := range p.States {
		if s.Name == ctx.CurrentState {
			stateIndex = i
			break
		}
	}
	if stateIndex == -1 {
		return fmt.Errorf("state '%s' is not found in process", ctx.CurrentState)
	}

	return p.States[stateIndex].Execute(ctx)
}

func (p *Process) ExecuteHandle(ctx *execontext.ExecutionContext) error {
	handlerIndex := -1
	for i, h := range p.Handlers {
		if h.Name == ctx.CurrentHandler {
			handlerIndex = i
			break
		}
	}
	if handlerIndex == -1 {
		return fmt.Errorf("handler '%s' is not found in process", ctx.CurrentHandler)
	}

	return p.Handlers[handlerIndex].Execute(ctx)
}

package logic

import (
	"fmt"

	"github.com/actord/actord/pkg/process/execontext"
)

type Program []Logic

func (p Program) Execute(ctx *execontext.ExecutionContext) error {
	if p == nil {
		return fmt.Errorf("program is empty in state '%s'", ctx.CurrentState)
	}
	ctx.ShouldTransit = false // if transition happen - cleanup this flag
	var startFromIndex int
	if ctx.CurrentLoginIndex == -1 {
		startFromIndex = 0
	} else {
		startFromIndex = ctx.CurrentLoginIndex
	}
	if startFromIndex >= len(p) {
		return fmt.Errorf("current logic index is out of range")
	}
	for i, l := range p[startFromIndex:] {
		ctx.CurrentLoginIndex = i

		err := l.Execute(ctx)
		if err != nil {
			return err
		}
		if ctx.ShouldTransit {
			ctx.CurrentLoginIndex = -1
			return nil
		}
		if ctx.AwaitEvent {
			return nil
		}
	}

	ctx.CurrentLoginIndex = -1
	return nil
}

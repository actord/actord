package executor

import (
	"fmt"

	"github.com/actord/actord/pkg/actor"
	"github.com/actord/actord/pkg/process/execontext"
	"github.com/actord/actord/pkg/process/schema"
	"github.com/actord/actord/pkg/process/typedef"
)

func (e *Executor) SendEvent(p *schema.Schema, act *actor.Actor, eventTypeName string, eventData []byte) error {
	eventType := p.Events.Find(eventTypeName)
	if eventType == nil {
		return fmt.Errorf("event type not found")
	}
	event, err := typedef.NewTypedJSONFromBytes(*eventType, eventData)
	if err != nil {
		return fmt.Errorf("failed to parse event data: %w", err)
	}
	ctx, err := execontext.NewExecutionContext(p.Types, act, event)
	if err != nil {
		return fmt.Errorf("failed to create execution context: %w", err)
	}

	if err := p.Process.Execute(ctx); err != nil {
		return fmt.Errorf("failed to execute process: %w", err)
	}

	act.State = ctx.CurrentState
	act.LogicIndex = ctx.CurrentLoginIndex
	act.AwaitEvent = ctx.AwaitEvent
	act.Data = ctx.GetActorData()

	return nil
}

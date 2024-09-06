package executor

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"github.com/actord/actord/pkg/actor"
	"github.com/actord/actord/pkg/process/execontext"
	"github.com/actord/actord/pkg/process/schema"
	"github.com/actord/actord/pkg/process/typedef"
)

func (e *Executor) Trigger(p *schema.Schema, triggerName string, eventData []byte) (*actor.Actor, error) {
	var trigger schema.Trigger
	var triggerFound bool
	for _, t := range p.Triggers {
		if t.Name == triggerName {
			trigger = t
			triggerFound = true
			break
		}
	}
	if !triggerFound {
		return nil, ErrTriggerNotFound
	}

	var eventType typedef.Type
	var eventTypeFound bool
	for _, e := range p.Events {
		if e.Name == trigger.EventType {
			eventType = e
			eventTypeFound = true
			break
		}
	}
	if !eventTypeFound {
		return nil, ErrEventTypeNotFound
	}
	log.Println("EVENT TYPE", eventType)

	event, err := typedef.NewTypedJSONFromBytes(eventType, eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse event data: %w", err)
	}

	ctx, err := execontext.NewExecutionContext(p.Types, nil, event)
	if err != nil {
		return nil, err
	}

	if err := trigger.Logic.Execute(ctx); err != nil {
		return nil, err
	}

	if !ctx.ShouldTransit {
		if trigger.Transition == nil {
			return nil, errors.New("no default transition defined in trigger")
		}
		ctx.CurrentState = *trigger.Transition
	}
	if ctx.AwaitEvent {
		return nil, errors.New("await not allowed in trigger program")
	}

	if err := p.Process.Execute(ctx); err != nil {
		return nil, err
	}

	act := &actor.Actor{
		ID:    uuid.New().String(),
		State: ctx.CurrentState,

		LogicIndex:   ctx.CurrentLoginIndex,
		AwaitEvent:   ctx.AwaitEvent,
		Data:         ctx.GetActorData(),
		DataTypeName: ctx.GetActorType().Name,
	}

	return act, nil
}

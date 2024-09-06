package execontext

import (
	"fmt"
	"strings"

	"github.com/valyala/fastjson"

	"github.com/actord/actord/pkg/actor"
	"github.com/actord/actord/pkg/process/typedef"
)

type ExecutionContext struct {
	actor *typedef.TypedJSON
	event *typedef.TypedJSON
	temp  *typedef.TypedJSON

	CurrentState      string
	ShouldTransit     bool // indicates that the state should be changed now
	CurrentHandler    string
	ShouldHandle      bool // indicates that the handler should be executed now
	CurrentLoginIndex int
	AwaitEvent        bool
}

func NewExecutionContext(
	types typedef.Collection,
	act *actor.Actor,
	event *typedef.TypedJSON,
) (ctx *ExecutionContext, err error) {
	if act != nil {
		ctx = &ExecutionContext{
			CurrentState:      act.State,
			CurrentLoginIndex: act.LogicIndex,
			AwaitEvent:        act.AwaitEvent,
		}
		actorJSON, err := fastjson.ParseBytes(act.Data)
		if err != nil {
			return nil, fmt.Errorf("broken actor: %w", err)
		}
		actorType := types.Find(act.DataTypeName)
		if actorType == nil {
			return nil, fmt.Errorf("actor type not found")
		}
		ctx.actor = typedef.NewTypedJSON(*actorType, actorJSON)
	} else {
		ctx = &ExecutionContext{}
		emptyActorJSON, _ := fastjson.ParseBytes([]byte(`{}`))
		actorType := types.Find("actor")
		if actorType == nil {
			return nil, fmt.Errorf("actor type not found")
		}
		ctx.actor = typedef.NewTypedJSON(*actorType, emptyActorJSON)
	}

	if err := ctx.actor.Validate(); err != nil {
		return nil, fmt.Errorf("actor json validation failed: %w", err)
	}

	ctx.event = event
	if ctx.event != nil {
		if err := ctx.event.Validate(); err != nil {
			return nil, fmt.Errorf("event json validation failed: %w", err)
		}
	}

	ctx.MakeCleanTemp()

	return ctx, nil
}

func (ctx *ExecutionContext) GetActorData() []byte {
	return ctx.actor.Marshal()
}

func (ctx *ExecutionContext) GetActorType() typedef.Type {
	return ctx.actor.Type
}

func (ctx *ExecutionContext) GetActor() *typedef.TypedJSON {
	return ctx.actor
}

func (ctx *ExecutionContext) GetEvent() *typedef.TypedJSON {
	return ctx.event
}

func (ctx *ExecutionContext) Get(key string) (*fastjson.Value, error) {
	source, key, err := ctx.getSource(key)
	if err != nil {
		return nil, err
	}

	// todo slice key
	return source.Get(key)
}

func (ctx *ExecutionContext) Set(key string, value *fastjson.Value) error {
	source, key, err := ctx.getSource(key)
	if err != nil {
		return err
	}

	return source.Set(key, value)
}

func (ctx *ExecutionContext) RemoveEventData() {
	ctx.event = nil
}

func (ctx *ExecutionContext) HasEventData() bool {
	return ctx.event != nil
}

func (ctx *ExecutionContext) MakeCleanTemp() {
	temp, err := typedef.NewTypedJSONFromBytes(typedef.Type{
		Name:   "",
		Any:    true,
		Fields: []typedef.Field{},
	}, []byte(`{}`))
	if err != nil {
		panic(fmt.Errorf("failed to create temporary json: %w", err))
	}
	ctx.temp = temp
}

func (ctx *ExecutionContext) getSource(key string) (*typedef.TypedJSON, string, error) {
	var source *typedef.TypedJSON
	var err error

	keyParts := strings.Split(key, ".")
	prefix := keyParts[0]

	switch prefix {
	case "event":
		source = ctx.event
	case "actor":
		source = ctx.actor
	case "temp":
		source = ctx.temp
	}

	if source == nil {
		err = fmt.Errorf("no valid source found")
	} else {
		key = strings.Join(strings.Split(key, ".")[1:], ".")
	}
	return source, key, err
}

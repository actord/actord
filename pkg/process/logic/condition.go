package logic

import (
	"bytes"

	"github.com/actord/actord/pkg/process/execontext"
)

type Condition struct {
	Key       string     `hcl:",label"`
	Equals    []string   `hcl:"equals,optional"`
	OnFailure OnFailure  `hcl:"on_failure,block"`
	OnSuccess *OnSuccess `hcl:"on_success,block"`
}

func (c *Condition) Execute(ctx *execontext.ExecutionContext) error {
	val, err := ctx.Get(c.Key)
	if err != nil {
		return err
	}
	valBytes := val.MarshalTo(nil)
	for _, key := range c.Equals {
		compareWith, err := ctx.Get(key)
		if err != nil {
			panic(err)
		}

		if !bytes.Equal(valBytes, compareWith.MarshalTo(nil)) {
			return c.OnFailure.Execute(ctx)
		}
	}

	if c.OnSuccess != nil {
		return c.OnSuccess.Execute(ctx)
	}
	return nil
}

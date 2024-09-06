package logic

import "github.com/actord/actord/pkg/process/execontext"

type Logic struct {
	Label       *string `hcl:"label"`
	Description *string `hcl:"description"`

	Condition      *Condition      `hcl:"condition,block"`
	Await          *Await          `hcl:"await,block"`
	VerifyPassword *VerifyPassword `hcl:"verify_password,block"`
	Set            []Set           `hcl:"set,block"`
	Reply          *Reply          `hcl:"reply,block"`
	Transition     *string         `hcl:"transition"`
}

func (l *Logic) Execute(ctx *execontext.ExecutionContext) error {
	if l.Condition != nil {
		if err := l.Condition.Execute(ctx); err != nil {
			return err
		}
		if ctx.ShouldTransit {
			return nil
		}
	}

	if l.Await != nil {
		if err := l.Await.Execute(ctx); err != nil {
			return err
		}
		if ctx.ShouldTransit {
			return nil
		}
		if ctx.AwaitEvent {
			return nil
		}
	}

	if l.VerifyPassword != nil {
		if err := l.VerifyPassword.Execute(ctx); err != nil {
			return err
		}
	}

	if len(l.Set) > 0 {
		for _, s := range l.Set {
			if err := s.Execute(ctx); err != nil {
				return err
			}
		}
	}

	if l.Reply != nil {
		if err := l.Reply.Execute(ctx); err != nil {
			return err
		}
		// todo: send reply
	}

	if l.Transition != nil {
		ctx.CurrentState = *l.Transition
		ctx.ShouldTransit = true
	}

	return nil
}

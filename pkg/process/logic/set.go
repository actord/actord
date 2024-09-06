package logic

import (
	"fmt"

	"github.com/valyala/fastjson"
	"golang.org/x/crypto/bcrypt"

	"github.com/actord/actord/pkg/process/execontext"
)

type Set struct {
	Key string `hcl:",label"`

	Copy *string `hcl:"copy"`
	Json *string `hcl:"json"`

	Strconv *[]string `hcl:"strconv"`
}

func (s Set) Execute(ctx *execontext.ExecutionContext) error {
	var newValue *fastjson.Value
	if s.Copy != nil {
		val, err := ctx.Get(*s.Copy)
		if err != nil {
			return err
		}
		newValue = val
	} else if s.Json != nil {
		val, err := fastjson.Parse(*s.Json)
		if err != nil {
			return err
		}
		newValue = val
	} else {
		panic("nothing to set")
	}

	if s.Strconv != nil {
		for _, process := range *s.Strconv {
			switch process {
			case "password_hash":
				hashed, err := bcrypt.GenerateFromPassword(newValue.GetStringBytes(), bcrypt.DefaultCost)
				if err != nil {
					return fmt.Errorf("failed to hash password: %w", err)
				}
				newValue = fastjson.MustParse(`"` + string(hashed) + `"`)
			default:
				return fmt.Errorf("unknown process: %s", process)
			}
		}
	}

	if err := ctx.Set(s.Key, newValue); err != nil {
		return err
	}

	return nil
}

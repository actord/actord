package logic

import (
	"fmt"

	"github.com/valyala/fastjson"
	"golang.org/x/crypto/bcrypt"

	"github.com/actord/actord/pkg/process/execontext"
)

type VerifyPassword struct {
	Password string `hcl:"password,label"`
	Hash     string `hcl:"hash"`
	SetTo    string `hcl:"set_to"`
}

func (v VerifyPassword) Execute(ctx *execontext.ExecutionContext) error {
	password, err := ctx.Get(v.Password)
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	hash, err := ctx.Get(v.Hash)
	if err != nil {
		return fmt.Errorf("failed to get hash: %w", err)
	}

	isMatchedJSON := "false"
	err = bcrypt.CompareHashAndPassword(hash.GetStringBytes(), password.GetStringBytes())
	if err == nil {
		isMatchedJSON = "true"
	}

	if err := ctx.Set(v.SetTo, fastjson.MustParse(isMatchedJSON)); err != nil {
		return fmt.Errorf("failed to set value: %w", err)
	}

	return nil
}

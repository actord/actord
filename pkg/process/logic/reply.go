package logic

import (
	"fmt"
	"log"

	"github.com/actord/actord/pkg/process/execontext"
)

type Reply struct {
	From string `hcl:"from"`
}

func (v Reply) Execute(ctx *execontext.ExecutionContext) error {
	reply, err := ctx.Get(v.From)
	if err != nil {
		return fmt.Errorf("failed to get reply: %w", err)
	}
	log.Println("REPLY:", reply)
	return nil
}

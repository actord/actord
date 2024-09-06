package schema

import (
	"github.com/actord/actord/pkg/process/typedef"
)

type Schema struct {
	Types    typedef.Collection `hcl:"type,block"`
	Events   typedef.Collection `hcl:"event,block"`
	Triggers []Trigger          `hcl:"trigger,block"`
	Process  Process            `hcl:"process,block"`
}

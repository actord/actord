package executor

import (
	"errors"
)

var (
	ErrTriggerNotFound   = errors.New("trigger not found")
	ErrEventTypeNotFound = errors.New("event type not found")
)

type Executor struct {
}

func New() *Executor {
	return &Executor{}
}

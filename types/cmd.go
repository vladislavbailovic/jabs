package types

import "context"

type Executable interface {
	Execute() (string, error)
}

type Scriptable interface {
	GetScript() string
}

type Runnable interface {
	Executable
	Scriptable
}

type ActionType int
type ActionState int

const (
	STATE_INIT ActionState = iota
	STATE_RUN
	STATE_DONE
)

type Action interface {
	Run()
	Output() chan string
	State() chan ActionState
}

type SubcommandType int
type Subcommand interface {
	Action
	Init(ctx context.Context) context.Context
}

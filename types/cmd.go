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
type Action interface {
	Run()
	Output() chan string
}

type SubcommandType int
type Subcommand interface {
	Action
	Init(ctx context.Context) context.Context
}

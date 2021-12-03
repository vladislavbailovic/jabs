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

type Action interface {
	Run()
}

type Subcommand interface {
	Action
	Init(ctx context.Context) context.Context
}

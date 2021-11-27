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

type Subcommand interface {
	Executable
	Init(ctx context.Context) context.Context
}

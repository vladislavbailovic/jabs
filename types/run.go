package types

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

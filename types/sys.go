package types

type Executable interface {
	Execute() (string, error)
}

type Scriptable interface {
	GetScript() string
}

type Instruction interface {
	Executable
	Scriptable
}

type FilenamePattern string
type PathPattern string

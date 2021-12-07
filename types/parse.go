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

type Rule struct {
	Name      string
	Observes  []Instruction
	DependsOn []string
	Tasks     []Instruction
}

type Macro struct {
	Name  string
	Value string
}

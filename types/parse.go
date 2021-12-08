package types

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

package types

type RuleName string

type Rule struct {
	Name      RuleName
	Observes  []Instruction
	DependsOn []RuleName
	Tasks     []Instruction
}

type Macro struct {
	Name  string
	Value string
}

package parse

import (
	"jabs/dbg"
	"jabs/types"
	"strings"
)

type Limit int

const (
	LIMIT_MACRO_EXPANSION_PASS Limit = 1000
)

type Task string

type Preprocessor struct {
	Macros map[string]types.Macro
	Rules  map[string]types.Rule
}

func NewPreprocessor(file string) Preprocessor {
	parser := NewParser(file)
	preprocessor := Preprocessor{}
	preprocessor.initMacros(parser.MacroDefinitions)
	preprocessor.initRules(parser.RuleDefinitions)
	return preprocessor
}

func (p *Preprocessor) initMacros(dfns []MacroDefinition) {
	mds := NewMacroDefinitions(dfns)
	p.Macros = mds.Macros
}

func (p *Preprocessor) initRules(dfns []RuleDefinition) {
	rules := map[string]types.Rule{}
	for _, dfn := range dfns {
		tasks := []types.Instruction{}
		for _, item := range dfn.Tasks {
			tasks = append(tasks, NewCommand(p.expand(item)))
		}
		dependencies := []string{}
		for _, item := range dfn.DependsOn {
			dependencies = append(dependencies, p.expand(item))
		}
		observes := []types.Instruction{}
		for _, obs := range dfn.Observes {
			observes = append(observes, NewCommand(p.expand(obs)))
		}
		name := p.expand(dfn.Name)
		rules[name] = types.Rule{
			Name:      name,
			DependsOn: dependencies,
			Observes:  observes,
			Tasks:     tasks,
		}
	}
	p.Rules = rules
}

func (p Preprocessor) expand(subj string) string {
	result := subj
	for i := Limit(0); i < LIMIT_MACRO_EXPANSION_PASS; i++ {
		expanded := false
		for name, macro := range p.Macros {
			key := GetExpansionKey(name)
			if strings.Contains(macro.Value, key) {
				dbg.Warning("Direct recursion encountered for %s", key)
				continue
			}
			expanded = strings.Contains(result, key)
			if !expanded {
				continue
			}
			result = strings.Replace(result, key, macro.Value, -1)
		}
		if !expanded {
			break
		}
	}
	return result
}

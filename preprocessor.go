package main

import (
	"strings"
)

type Limit int

const (
	LIMIT_MACRO_EXPANSION_PASS Limit = 1000
)

type Macro struct {
	Name  string
	Value string
}

type Task string

type Rule struct {
	Name      string
	Observes  []string
	DependsOn []string
	Tasks     []string
}

type Preprocessor struct {
	Macros map[string]Macro
	Rules  map[string]Rule
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
	rules := map[string]Rule{}
	for _, dfn := range dfns {
		tasks := []string{}
		for _, item := range dfn.Tasks {
			tasks = append(tasks, p.expand(item))
		}
		dependencies := []string{}
		for _, item := range dfn.DependsOn {
			dependencies = append(dependencies, p.expand(item))
		}
		observes := []string{}
		for _, obs := range dfn.Observes {
			observes = append(observes, p.expand(obs))
		}
		name := p.expand(dfn.Name)
		rules[name] = Rule{
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
				Warning("Direct recursion encountered for %s", key)
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

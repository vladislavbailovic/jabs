package main

import "strings"

type Preprocessor struct {
	Macros map[string]string
	Rules  map[string]RuleDefinition
}

func NewPreprocessor(file string) Preprocessor {
	parser := NewParser(file)
	preprocessor := Preprocessor{}
	preprocessor.initMacroDefinitions(parser.MacroDefinitions)
	preprocessor.initRules(parser.RuleDefinitions)
	return preprocessor
}

func (p *Preprocessor) initMacroDefinitions(dfns []MacroDefinition) {
	macros := map[string]string{}
	for _, dfn := range dfns {
		value := dfn.Value
		if "" == value {
			value = dfn.Command
		}
		expanded := expandMacroDfns(value, dfns)
		if dfn.Command != "" {
			expanded = "/bin/bash -c '" + expanded + "'"
		}
		macros[dfn.Name] = expanded
	}
	p.Macros = macros
}

func (p *Preprocessor) initRules(dfns []RuleDefinition) {
	rules := map[string]RuleDefinition{}
	for _, dfn := range dfns {
		tasks := []string{}
		for _, item := range dfn.Tasks {
			tasks = append(tasks, p.expand(item))
		}
		dependencies := []string{}
		for _, item := range dfn.DependsOn {
			dependencies = append(dependencies, p.expand(item))
		}
		name := p.expand(dfn.Name)
		rules[name] = RuleDefinition{
			Name:      name,
			DependsOn: dependencies,
			Tasks:     tasks,
		}
	}
	p.Rules = rules
}

func (p Preprocessor) expand(subj string) string {
	result := subj
	for i := 0; i < 1000; i++ {
		expanded := false
		for name, value := range p.Macros {
			key := getExpansionKey(name)
			expanded = strings.Contains(result, key)
			if !expanded {
				continue
			}
			result = strings.Replace(result, key, value, -1)
		}
		if !expanded {
			break
		}
	}
	return result
}

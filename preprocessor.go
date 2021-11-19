package main

import (
	"fmt"
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
	macros := map[string]Macro{}
	dfns = append(GetSystemMacroDefinitions(), dfns...)
	for _, dfn := range dfns {
		value := dfn.Value
		if "" == value {
			value = dfn.Command
		}
		expanded := expandMacroDfns(value, dfn.Name, dfns)
		if dfn.Command != "" {
			cmd := NewExecutable(expanded)
			out, err := cmd.Execute()
			if err != nil {
				panic(err)
			}
			expanded = out
		}
		macros[dfn.Name] = Macro{Name: dfn.Name, Value: expanded}
	}
	p.Macros = macros
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
			key := getExpansionKey(name)
			if strings.Contains(macro.Value, key) {
				// @TODO warn about recursion
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

func getExpansionKey(what string) string {
	return fmt.Sprintf("${{%s}}", what)
}

func expandMacroDfns(subj string, in string, dfns []MacroDefinition) string {
	result := subj
	for i := Limit(0); i < LIMIT_MACRO_EXPANSION_PASS; i++ {
		expanded := false
		for _, macro := range dfns {
			if macro.Name == in {
				// @TODO warn about recursion
				break
			}
			key := getExpansionKey(macro.Name)
			expanded = strings.Contains(result, key)
			if !expanded {
				continue
			}
			value := macro.Value
			if "" == value {
				value = macro.Command
			}
			result = strings.Replace(result, key, value, -1)
		}
		if !expanded {
			break
		}
	}
	return result
}

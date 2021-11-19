package main

import (
	"fmt"
	"strings"
)

type Macro struct {
	Name  string
	Value string
}

type Task string

type Rule struct {
	Name      string
	Type      string
	DependsOn []string // @TODO these should actually be other rules
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
	for _, dfn := range dfns {
		value := dfn.Value
		if "" == value {
			value = dfn.Command
		}
		expanded := expandMacroDfns(value, dfns)
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
		name := p.expand(dfn.Name)
		rules[name] = Rule{
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
		for name, macro := range p.Macros {
			key := getExpansionKey(name)
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

func expandMacroDfns(subj string, dfns []MacroDefinition) string {
	result := subj
	for i := 0; i < 1000; i++ {
		expanded := false
		for _, macro := range dfns {
			key := getExpansionKey(macro.Name)
			expanded = strings.Contains(result, key)
			if !expanded {
				continue
			}
			value := macro.Value
			if "" == value {
				// @TODO execute this???
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

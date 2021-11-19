package main

import (
	"fmt"
	"strings"
)

type MacroDefinition struct {
	Name    string `yaml:"Name"`
	Command string `yaml:"Command"`
	Value   string `yaml:"Value"`
}

type RuleDefinition struct {
	Name      string   `yaml:"Name"`
	Type      string   `yaml:"Type"`
	DependsOn []string `yaml:"DependsOn"`
	Tasks     []string `yaml:"Tasks"`
}

func main() {
	p := NewPreprocessor("./examples/first.yml")
	fmt.Println(p.Rules)
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

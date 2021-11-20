package main

import (
	"bytes"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type DefinitionType string

const (
	DFN_MACRO DefinitionType = "Macro"
	DFN_RULE                 = "Rule"
)

type MacroDefinition struct {
	Name    string `yaml:"Name"`
	Command string `yaml:"Command"`
	Value   string `yaml:"Value"`
}

type RuleDefinition struct {
	Name      string   `yaml:"Name"`
	Observes  []string `yaml:"Observes"`
	DependsOn []string `yaml:"DependsOn"`
	Tasks     []string `yaml:"Tasks"`
}

type Parser struct {
	MacroDefinitions []MacroDefinition
	RuleDefinitions  []RuleDefinition
}

func NewParser(file string) Parser {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic("Can't read file: " + file)
	}
	return Parser{
		MacroDefinitions: parseMacros(data),
		RuleDefinitions:  parseRules(data),
	}
}

func parseMacros(data []byte) []MacroDefinition {
	dcd := yaml.NewDecoder(bytes.NewReader(data))
	items := []MacroDefinition{}
	for {
		instance := make(map[DefinitionType]MacroDefinition)
		err := dcd.Decode(&instance)
		if err != nil {
			if err != io.EOF {
				Warning("error: %v", err)
			}
			break
		}

		item, ok := instance[DFN_MACRO]
		if ok {
			items = append(items, item)
		}
	}
	return items
}

func parseRules(data []byte) []RuleDefinition {
	dcd := yaml.NewDecoder(bytes.NewReader(data))
	items := []RuleDefinition{}
	for {
		instance := make(map[DefinitionType]RuleDefinition)
		err := dcd.Decode(&instance)
		if err != nil {
			if err != io.EOF {
				Warning("error: %v", err)
			}
			break
		}

		item, ok := instance[DFN_RULE]
		if ok {
			items = append(items, item)
			continue
		}
	}
	return items
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

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
	instance := make(map[string]MacroDefinition)
	items := []MacroDefinition{}
	for {
		err := dcd.Decode(&instance)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("error: %v\n", err)
			}
			break
		}

		item, ok := instance["Macro"]
		if ok {
			items = append(items, item)
			continue
		}
	}
	return items
}

func parseRules(data []byte) []RuleDefinition {
	dcd := yaml.NewDecoder(bytes.NewReader(data))
	instance := make(map[string]RuleDefinition)
	items := []RuleDefinition{}
	for {
		err := dcd.Decode(&instance)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("error: %v\n", err)
			}
			break
		}

		item, ok := instance["Rule"]
		if ok {
			items = append(items, item)
			continue
		}
	}
	return items
}

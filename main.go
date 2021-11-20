package main

import (
	"fmt"
	"strings"
)

func main() {
	p := NewPreprocessor("./examples/first.yml")
	es := NewEvaluationStack("root", p.Rules)
	printStack(es)
}

func printStack(es EvaluationStack) {
	out := []string{"#!/bin/bash", ""}

	for idx, rule := range es.stack {
		out = append(out, fmt.Sprintf("# Rule %d -- %s", idx, rule.Name))
		for _, task := range rule.Tasks {
			cmd := NewScriptable(task)
			out = append(out, cmd.GetScript())
		}
		out = append(out, "")
	}
	Info("\n" + strings.Join(out[:], "\n"))
}

func executeStack(es EvaluationStack) {
	Debug("Stack")
	Debug("--------------------")
	for _, rl := range es.stack {
		Debug("\t- %s", rl.Name)
		for i, task := range rl.Tasks {
			cmd := NewExecutable(task)
			out, err := cmd.Execute()
			if err != nil {
				Error("%v", err)
			}
			Debug("\t\t%d) %v", i+1, out)
		}
	}
}

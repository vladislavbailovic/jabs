package main

import (
	"context"
	"fmt"
	"strings"
)

func main() {
	ctx := context.WithValue(context.Background(), OPT_VERBOSITY, int(LOG_INFO))
	InitOptions(ctx)
	p := NewPreprocessor("./examples/self.yml")
	es := NewEvaluationStack("cover:html", p.Rules)
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

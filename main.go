package main

import (
	"context"
	"fmt"
	"jabs/dbg"
	"jabs/options"
	"strings"
)

func main() {
	ctx := ApplyEnvironment(context.Background())

	//ctx = context.WithValue(ctx, OPT_VERBOSITY, int(LOG_INFO))
	options.InitOptions(ctx)

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
	dbg.Info("\n" + strings.Join(out[:], "\n"))
}

func executeStack(es EvaluationStack) {
	dbg.Debug("Stack")
	dbg.Debug("--------------------")
	for _, rl := range es.stack {
		dbg.Debug("\t- %s", rl.Name)
		for i, task := range rl.Tasks {
			cmd := NewExecutable(task)
			out, err := cmd.Execute()
			if err != nil {
				dbg.Error("%v", err)
			}
			dbg.Debug("\t\t%d) %v", i+1, out)
		}
	}
}

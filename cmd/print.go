package cmd

import (
	"fmt"
	"jabs/dbg"
	"jabs/parse"
	"strings"
)

func Print(file string, root string) {
	timer := dbg.GetTimer()

	p := parse.NewPreprocessor(file)
	timer.Lap("preprocess")
	es := parse.NewEvaluationStack(root, p.Rules)
	timer.Lap("parse")
	printStack(es)
}

func printStack(es parse.EvaluationStack) {
	out := []string{"#!/bin/bash", ""}

	for idx, rule := range es.GetStack() {
		out = append(out, fmt.Sprintf("# Rule %d -- %s", idx, rule.Name))
		for _, task := range rule.Tasks {
			out = append(out, task.GetScript())
		}
		out = append(out, "")
	}
	dbg.Info("\n" + strings.Join(out[:], "\n"))
}

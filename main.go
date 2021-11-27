package main

import (
	"context"
	"fmt"
	"jabs/dbg"
	"jabs/opts"
	"strings"
)

func main() {
	timer := dbg.GetTimer()
	ctx := ApplyEnvironment(context.Background())

	//ctx = context.WithValue(ctx, OPT_VERBOSITY, int(LOG_INFO))
	opts.InitOptions(ctx)
	timer.Lap("boot")

	p := NewPreprocessor("./examples/self.yml")
	timer.Lap("preprocess")
	es := NewEvaluationStack("cover:html", p.Rules)
	timer.Lap("parse")

	printStack(es)
	// executeStack(es)
	timer.Lap("print")

	dbg.Debug("duration: %dms", timer.Duration()/dbg.TIME_MS)
	for name, time := range timer.GetLaps() {
		dbg.Debug("\t%s: %dms", name, time/dbg.TIME_MS)
	}
}

func printStack(es EvaluationStack) {
	out := []string{"#!/bin/bash", ""}

	for idx, rule := range es.stack {
		out = append(out, fmt.Sprintf("# Rule %d -- %s", idx, rule.Name))
		for _, task := range rule.Tasks {
			out = append(out, task.GetScript())
		}
		out = append(out, "")
	}
	dbg.Info("\n" + strings.Join(out[:], "\n"))
}

func executeStack(es EvaluationStack) {
	dbg.Debug("Stack")
	dbg.Debug("--------------------")
	timer := dbg.NewStopwatch()
	for _, rl := range es.stack {
		dbg.Debug("\t- %s", rl.Name)
		for i, task := range rl.Tasks {
			out, err := task.Execute()
			if err != nil {
				dbg.Error("%v", err)
			}
			dbg.Debug("\t\t%d) %v", i+1, out)
			timer.Lap(fmt.Sprintf("Rule %s :: Task %d", rl.Name, i))
		}
	}
	dbg.Debug("--------------------")
	dbg.Info("--------------------")
	dbg.Info("Execution times")
	for name, time := range timer.GetLaps() {
		dbg.Info("\t%s: %dms", name, time/dbg.TIME_MS)
	}
}

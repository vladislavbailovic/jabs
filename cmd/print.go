package cmd

import (
	"context"
	"flag"
	"fmt"
	"jabs/dbg"
	"jabs/opts"
	"jabs/parse"
	"strings"
)

type PrintSubcommand struct {
	PrintAction
	fs       *flag.FlagSet
	whatever *string
}

func NewPrintSubcommand(fs *flag.FlagSet) *PrintSubcommand {
	ps := PrintSubcommand{fs: fs}
	ps.whatever = fs.String("whatever", "this is whatever", "whatever arg")
	return &ps
}

func (ps *PrintSubcommand) Init(ctx context.Context) context.Context {
	// ...
	// privates are now populated with flags
	// so init context and return it

	// dbg.Debug("WHATEVER: [%s]", *ps.whatever)
	return ctx
}

type PrintAction struct{}

func (pa PrintAction) Run() {
	timer := dbg.GetTimer()
	options := opts.GetOptions()

	p := parse.NewPreprocessor(options.Path)
	timer.Lap("preprocess")
	es := parse.NewEvaluationStack(options.Root, p.Rules)
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

package cmd

import (
	"context"
	"flag"
	"fmt"
	"jabs/dbg"
	"jabs/parse"
	"strings"
)

type PrintSubcommand struct {
	flags *flag.FlagSet

	root string
	file string
}

func NewPrintSubcommand(file string, root string) *PrintSubcommand {
	ps := PrintSubcommand{
		flags: flag.NewFlagSet("print", flag.ContinueOnError),
		file:  file,
		root:  root,
	}
	return &ps
}

func (ps *PrintSubcommand) Init(ctx context.Context) context.Context {
	ps.flags.Parse(flag.Args())
	// ...
	// privates are now populated with flags
	// so init context and return it
	return ctx
}

func (ps PrintSubcommand) Execute() (string, error) {
	timer := dbg.GetTimer()

	p := parse.NewPreprocessor(ps.file)
	timer.Lap("preprocess")
	es := parse.NewEvaluationStack(ps.root, p.Rules)
	timer.Lap("parse")
	printStack(es)
	return "", nil
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

package cmd

import (
	"context"
	"flag"
	"fmt"
	"jabs/dbg"
	"jabs/parse"
)

type RunSubcommand struct {
	flags *flag.FlagSet

	root string
	file string
}

func NewRunSubcommand(file string, root string) *RunSubcommand {
	rs := RunSubcommand{
		flags: flag.NewFlagSet("print", flag.ContinueOnError),
		file:  file,
		root:  root,
	}
	return &rs
}

func (rs *RunSubcommand) Init(ctx context.Context) context.Context {
	rs.flags.Parse(flag.Args())
	// ...
	// privates are now populated with flags
	// so init context and return it
	return ctx
}

func (rs RunSubcommand) Execute() (string, error) {
	timer := dbg.GetTimer()

	p := parse.NewPreprocessor(rs.file)
	timer.Lap("preprocess")
	es := parse.NewEvaluationStack(rs.root, p.Rules)
	timer.Lap("parse")
	executeStack(es)
	return "", nil
}

func executeStack(es parse.EvaluationStack) {
	dbg.Debug("Stack")
	dbg.Debug("--------------------")
	timer := dbg.NewStopwatch()
	for _, rl := range es.GetStack() {
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

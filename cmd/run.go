package cmd

import (
	"context"
	"flag"
	"fmt"
	"jabs/dbg"
	"jabs/opts"
	"jabs/parse"
)

type RunSubcommand struct {
	RunAction
	fs *flag.FlagSet
}

func NewRunSubcommand(fs *flag.FlagSet) *RunSubcommand {
	rs := RunSubcommand{fs: fs}
	rs.out = make(chan string)
	return &rs
}

func (rs *RunSubcommand) Init(ctx context.Context) context.Context {
	// ...
	// privates are now populated with flags
	// so init context and return it
	return ctx
}

type RunAction struct {
	out chan string
}

func (a RunAction) Output() chan string {
	return a.out
}

func (ra RunAction) Run() {
	timer := dbg.GetTimer()
	options := opts.GetOptions()

	p := parse.NewPreprocessor(options.Path)
	timer.Lap("preprocess")
	es := parse.NewEvaluationStack(options.Root, p.Rules)
	timer.Lap("parse")
	ra.executeStack(es)
}

func (ra RunAction) executeStack(es parse.EvaluationStack) {
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
			ra.out <- out
			dbg.Info("\t\t%d) %v", i+1, out)
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

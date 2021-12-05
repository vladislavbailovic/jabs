package cmd

import (
	"context"
	"flag"
	"fmt"
	"jabs/dbg"
	"jabs/opts"
	"jabs/parse"
	"jabs/types"
)

type RunSubcommand struct {
	RunAction
	fs *flag.FlagSet
}

func NewRunSubcommand(fs *flag.FlagSet) *RunSubcommand {
	rs := RunSubcommand{fs: fs}
	rs.out = make(chan string)
	rs.state = make(chan types.ActionState)
	return &rs
}

func (rs *RunSubcommand) Init(ctx context.Context) context.Context {
	rs.state <- types.STATE_INIT
	// ...
	// privates are now populated with flags
	// so init context and return it
	return ctx
}

type RunAction struct {
	out   chan string
	state chan types.ActionState
}

func (a RunAction) Output() chan string {
	return a.out
}
func (a RunAction) State() chan types.ActionState {
	return a.state
}

func (ra RunAction) Run() {
	ra.state <- types.STATE_RUN
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
		dbg.Debug("\t- Running task: [%s]", rl.Name)
		for i, task := range rl.Tasks {
			out, err := task.Execute()
			if err != nil {
				dbg.Info("\tTask errored out: %s %d", rl.Name, i)
				dbg.Error("%v", err)
			}
			ra.out <- out
			timer.Lap(fmt.Sprintf("Rule %s :: Task %d", rl.Name, i))
		}
	}
	dbg.Info("--------------------")
	dbg.Info("Execution times")
	for name, time := range timer.GetLaps() {
		dbg.Info("\t%s: %dms", name, time/dbg.TIME_MS)
	}
	ra.state <- types.STATE_DONE
}

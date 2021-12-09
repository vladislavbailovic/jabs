package cmd

import (
	"context"
	"flag"
	"fmt"
	"jabs/dbg"
	"jabs/opts"
	"jabs/parse"
	"jabs/sys"
	"jabs/types"
	"strings"
)

type PrintAction struct {
	out   chan string
	state chan types.ActionState
}

type PrintSubcommand struct {
	PrintAction
	fs                *flag.FlagSet
	includeConditions *bool
}

func NewPrintSubcommand(fs *flag.FlagSet) *PrintSubcommand {
	ps := PrintSubcommand{fs: fs}
	ps.out = make(chan string)
	ps.state = make(chan types.ActionState)
	ps.includeConditions = fs.Bool("conds", false, "Include rule conditions in output")
	return &ps
}

func (ps *PrintSubcommand) Init(ctx context.Context) context.Context {
	ps.state <- types.STATE_INIT
	ctx = context.WithValue(ctx, opts.OPT_CONDITIONS, *ps.includeConditions)
	return ctx
}

func (a PrintAction) Output() chan string {
	return a.out
}
func (a PrintAction) State() chan types.ActionState {
	return a.state
}

func (pa PrintAction) Run() {
	pa.state <- types.STATE_RUN
	timer := dbg.GetTimer()
	options := opts.GetOptions()

	p := parse.NewPreprocessor(options.Path)
	timer.Lap("preprocess")
	es := parse.NewEvaluationStack(options.Root, p.Rules)
	timer.Lap("parse")
	pa.printStack(es)
}

func (pa PrintAction) printStack(es parse.EvaluationStack) {
	out := []string{sys.Shebang(), ""}

	for _, rule := range es.GetStack() {
		out = append(out, pa.printRule(rule)...)
	}
	pa.out <- strings.Join(out, "\n")
	pa.state <- types.STATE_DONE
}

func (pa PrintAction) printRule(rule types.Rule) []string {
	out := []string{}
	options := opts.GetOptions()
	out = append(out, fmt.Sprintf("# Rule [%s] ---", rule.Name))
	dbg.Info("Printing %s", rule.Name)
	if options.Conditions {
		dbg.Info("\tEmitting rule conditions for %s", rule.Name)
		for _, obs := range rule.Observes {
			out = append(out, obs.GetScript())
		}
	}
	dbg.Info("\tEmitting tasks for %s", rule.Name)
	for _, task := range rule.Tasks {
		out = append(out, task.GetScript())
	}
	out = append(out, "")
	return out
}

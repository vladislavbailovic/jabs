package cmd

import (
	"fmt"
	"jabs/dbg"
	"jabs/parse"
)

func Execute(file string, root string) {
	timer := dbg.GetTimer()

	p := parse.NewPreprocessor(file)
	timer.Lap("preprocess")
	es := parse.NewEvaluationStack(root, p.Rules)
	timer.Lap("parse")
	executeStack(es)
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

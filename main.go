package main

import (
	"context"
	"jabs/cmd"
	"jabs/dbg"
	"jabs/opts"
)

func main() {
	timer := dbg.GetTimer()
	ctx := ApplyEnvironment(context.Background())

	//ctx = context.WithValue(ctx, OPT_VERBOSITY, int(LOG_INFO))
	opts.InitOptions(ctx)
	timer.Lap("boot")

	cmd.Print("examples/self.yml", "cover:html")
	// cmd.Execute("./examples/self.yml", "cover:html")
	timer.Lap("print")

	dbg.Debug("duration: %dms", timer.Duration()/dbg.TIME_MS)
	for name, time := range timer.GetLaps() {
		dbg.Debug("\t%s: %dms", name, time/dbg.TIME_MS)
	}
}

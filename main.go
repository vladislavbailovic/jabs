package main

import (
	"context"
	"flag"
	"fmt"
	"jabs/cmd"
	"jabs/dbg"
	"jabs/opts"
	"jabs/types"
	"os"
)

func main() {
	timer := dbg.GetTimer()
	ctx := ApplyEnvironment(context.Background())

	flags := flag.NewFlagSet("main", flag.ContinueOnError)
	file := flags.String("f", "./examples/self.yml", "File to process")
	// file := flag.String("f", "./examples/self.yml", "File to process")
	flag.StringVar(file, "f", "./examples/self.yml", "File to process")

	var subcmdtype string
	if len(os.Args) >= 2 {
		subcmdtype = os.Args[1]
	} else {
		subcmdtype = ""
	}

	position := 2
	switch subcmdtype {
	case "run":
	case "print":
	default:
		subcmdtype = "print"
		fmt.Println("print subcommand")
		position = 0
	}
	flags.Parse(os.Args[position:])
	flag.Parse()

	//ctx = context.WithValue(ctx, OPT_VERBOSITY, int(LOG_INFO))

	var root string
	if position == 0 {
		if len(flag.Args()) > 0 {
			root = flag.Args()[len(flag.Args())-1]
		}
	} else {
		if len(flags.Args()) > 0 {
			root = flags.Args()[len(flags.Args())-1]
		}
	}
	if "" == root {
		root = "cover:html"
	}

	var subcmd types.Subcommand
	switch subcmdtype {
	case "run":
		subcmd = cmd.NewRunSubcommand(*file, root)
	case "print":
		subcmd = cmd.NewPrintSubcommand(*file, root)
	}

	ctx = subcmd.Init(ctx)
	opts.InitOptions(ctx)
	timer.Lap("boot")

	subcmd.Execute()
	timer.Lap("subcommand")

	dbg.Debug("duration: %dms", timer.Duration()/dbg.TIME_MS)
	for name, time := range timer.GetLaps() {
		dbg.Debug("\t%s: %dms", name, time/dbg.TIME_MS)
	}
}

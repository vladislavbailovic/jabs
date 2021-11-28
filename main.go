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

	fs := flag.NewFlagSet("main", flag.ContinueOnError)
	file := fs.String("f", "./examples/self.yml", "File to process")

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
		position = 1
		subcmdtype = "print"
		fmt.Println("print subcommand")
	}

	var subcmd types.Subcommand
	switch subcmdtype {
	case "run":
		subcmd = cmd.NewRunSubcommand(fs)
	case "print":
		subcmd = cmd.NewPrintSubcommand(fs)
	}
	fs.Parse(os.Args[position:])

	var root string
	if len(fs.Args()) > 0 {
		root = fs.Args()[len(fs.Args())-1]
	}
	if "" == root {
		root = "cover:html"
	}

	ctx = context.WithValue(ctx, opts.OPT_PATH, *file)
	ctx = context.WithValue(ctx, opts.OPT_ROOT, root)

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

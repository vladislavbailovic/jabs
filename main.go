package main

import (
	"context"
	"flag"
	"jabs/cmd"
	"jabs/dbg"
	"jabs/opts"
	"os"
)

func usage(fs *flag.FlagSet) {
	dbg.Out("Usage: jabs <SUBCOMMAND> <FLAGS> <RULE>\n")
	dbg.Out("Subcommands:\n")
	dbg.Out("  - print\n")
	dbg.Out("  - run\n")
	dbg.Out("  - watch\n")
	dbg.Out("The watch subcommand will read a list of files to watch from STDIN\n")
	dbg.Out("Flags:\n")
	fs.PrintDefaults()
	os.Exit(0)
}

func main() {
	timer := dbg.GetTimer()
	ctx := ApplyEnvironment(context.Background())

	fs := flag.NewFlagSet("main", flag.ContinueOnError)
	file := fs.String("f", "./examples/self.yml", "File to process")
	force := fs.Bool("force", false, "Force-run (do not stop at recoverable errors)")
	help := fs.Bool("h", false, "Show help")

	var wantedSubcommand string
	position := 2
	if len(os.Args) >= 2 {
		wantedSubcommand = os.Args[1]
	} else {
		position = 1
		wantedSubcommand = ""
	}

	subcmd := cmd.NewSubcommand(cmd.SubcommandType(wantedSubcommand), fs)

	fs.Parse(os.Args[position:])

	if *help {
		usage(fs)
	}

	var root string
	if len(fs.Args()) > 0 {
		root = fs.Args()[len(fs.Args())-1]
	}
	if "" == root {
		root = "cover:html"
	}

	ctx = context.WithValue(ctx, opts.OPT_PATH, *file)
	ctx = context.WithValue(ctx, opts.OPT_ROOT, root)
	ctx = context.WithValue(ctx, opts.OPT_FORCE, *force)

	ctx = subcmd.Init(ctx)

	opts.InitOptions(ctx)
	timer.Lap("boot")

	subcmd.Run()
	timer.Lap("subcommand")

	dbg.Debug("duration: %dms", timer.Duration()/dbg.TIME_MS)
	for name, time := range timer.GetLaps() {
		dbg.Debug("\t%s: %dms", name, time/dbg.TIME_MS)
	}
}

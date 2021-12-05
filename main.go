package main

import (
	"context"
	"flag"
	"jabs/cmd"
	"jabs/dbg"
	"jabs/opts"
	"jabs/out"
	"jabs/types"
	"os"
)

func usage(fs *flag.FlagSet) {
	out.Cli.Out("Usage: jabs <SUBCOMMAND> <FLAGS> <RULE>")
	out.Cli.Out("Subcommands:")
	out.Cli.Out("  - print")
	out.Cli.Out("  - run")
	out.Cli.Out("  - watch")
	out.Cli.Out("The watch subcommand will read a list of files to watch from STDIN")
	out.Cli.Out("Flags:")
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
		wantedSubcommand = ""
	}
	subcmdType := cmd.SubcommandType(wantedSubcommand)
	if cmd.SUBCMD_DEFAULT == subcmdType {
		position = 1
	}

	subcmd := cmd.NewSubcommand(subcmdType, fs)

	done := make(chan bool)
	go func() {
		for {
			select {
			case state := <-subcmd.State():
				switch state {
				case types.STATE_INIT:
					dbg.Info("--- Init ---")
				case types.STATE_RUN:
					dbg.Info("--- Run ---")
				case types.STATE_DONE:
					dbg.Info("--- Done ---")
					done <- true
				}
			case output := <-subcmd.Output():
				dbg.Debug("--- Output ---")
				out.Cli.Out(output)
				dbg.Debug("--- End of output ---")
			}
		}
	}()

	fs.Parse(os.Args[position:])

	if *help {
		usage(fs)
		os.Exit(0)
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
	dbg.Debug("\t%v subcommand (%s)", subcmdType, wantedSubcommand)
	dbg.Debug("\tfile: %#v, root: %#v", *file, root)

	go func() {
		for val := range done {
			if !val {
				continue
			}
			timer.Lap("subcommand")

			dbg.Debug("duration: %dms", timer.Duration()/dbg.TIME_MS)
			for name, time := range timer.GetLaps() {
				dbg.Debug("\t%s: %dms", name, time/dbg.TIME_MS)
			}
			os.Exit(0)
		}
	}()

	subcmd.Run()
	<-done
}

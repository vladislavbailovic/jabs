package main

import (
	"context"
	"flag"
	"fmt"
	"jabs/sys"
	"jabs/cmd"
	"jabs/dbg"
	"jabs/opts"
	"jabs/out"
	"jabs/types"
	"os"
)

const (
	DEFAULT_FILE string = "./examples/self.yml"
	DEFAULT_RULE        = "cover:html"
)

func usage(fs *flag.FlagSet) {
	isDefault := func(action string) string {
		if cmd.SubcommandType(action) == cmd.DefaultSubcommand() {
			return "(default)"
		}
		return ""
	}
	out.Cli.Out("Usage: jabs <SUBCOMMAND> <FLAGS> <RULE>")
	out.Cli.Out("Subcommands:")
	out.Cli.Out(fmt.Sprintf("  - print %s", isDefault("print")))
	out.Cli.Out(fmt.Sprintf("  - run %s", isDefault("run")))
	out.Cli.Out(fmt.Sprintf("  - watch %s", isDefault("watch")))
	out.Cli.Out("The watch subcommand will read a list of files to watch from STDIN")
	out.Cli.Out("Flags:")
	fs.PrintDefaults()
	os.Exit(0)
}

func initContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, opts.OPT_VERBOSITY, int(types.LOG_WARNING))
	return ctx
}

func main() {
	timer := dbg.GetTimer()
	ctx := sys.ApplyEnvironment(initContext())

	fs := flag.NewFlagSet("main", flag.ContinueOnError)
	file := fs.String("f", DEFAULT_FILE, "File to process")
	force := fs.Bool("force", false, "Force-run (do not stop at recoverable errors)")

	notice := fs.Bool("v", false, "Verbose output (verbosity level: notice")
	info := fs.Bool("vv", false, "Verbose output (verbosity level: info")
	debug := fs.Bool("vvv", false, "Verbose output (verbosity level: debug")

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
				case types.STATE_RUN:
				case types.STATE_DONE:
					done <- true
				}
			case output := <-subcmd.Output():
				out.Cli.Out(output)
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
		root = DEFAULT_RULE
	}

	ctx = context.WithValue(ctx, opts.OPT_PATH, *file)
	ctx = context.WithValue(ctx, opts.OPT_ROOT, root)
	ctx = context.WithValue(ctx, opts.OPT_FORCE, *force)

	if *notice {
		ctx = context.WithValue(ctx, opts.OPT_VERBOSITY, int(types.LOG_NOTICE))
	}
	if *info {
		ctx = context.WithValue(ctx, opts.OPT_VERBOSITY, int(types.LOG_INFO))
	}
	if *debug {
		ctx = context.WithValue(ctx, opts.OPT_VERBOSITY, int(types.LOG_DEBUG))
	}

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

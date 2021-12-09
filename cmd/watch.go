package cmd

import (
	"context"
	"flag"
	"jabs/dbg"
	"jabs/opts"
	"jabs/sys"
	"jabs/types"
	"time"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

const SLEEPYTIME time.Duration = 5000

type WatchSubcommand struct {
	fs     *flag.FlagSet
	out    chan string
	state  chan types.ActionState

	action *string
	dir *string
	filter *string
}

func NewWatchSubcommand(fs *flag.FlagSet) *WatchSubcommand {
	ws := WatchSubcommand{fs: fs}
	ws.state = make(chan types.ActionState)
	ws.out = make(chan string)

	ws.action = fs.String("action",
		"print", "Action to perform on resource change")
	ws.dir = fs.String("dir",
		"./*", "Directory to watch")
	ws.filter = fs.String("filter",
		"*", "Pattern filter for watched files")

	return &ws
}

func (s WatchSubcommand) Output() chan string {
	return s.out
}
func (s WatchSubcommand) State() chan types.ActionState {
	return s.state
}

func (ws WatchSubcommand) Init(ctx context.Context) context.Context {
	ws.state <- types.STATE_INIT
	ctx = context.WithValue(ctx, opts.OPT_ACTION, int(ActionType(*ws.action)))
	ctx = context.WithValue(ctx, opts.OPT_DIRECTORY, string(types.PathPattern(*ws.dir)))
	ctx = context.WithValue(ctx, opts.OPT_FILTER, string(types.FilenamePattern(*ws.filter)))
	return ctx
}

func (ws WatchSubcommand) Run() {
	ws.state <- types.STATE_RUN
	options := opts.GetOptions()

	dbg.Info("File from %s", options.Path)
	dbg.Info("Rule is %s", options.Root)

	sources := sys.NewDirlist(options.Directory)
	dbg.Info("Sources: %v", sources)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		dbg.FatalError("%v", err)
	}
	defer watcher.Close()

	action := NewAction(options.Action)
	filter := sys.NewFilter(options.Filter)

	go func() {
		for {
			select {
			case state := <-action.State():
				switch state {
				case types.STATE_INIT:
					dbg.Info("-- Action Init --")
				case types.STATE_RUN:
					dbg.Info("-- Action Run --")
				case types.STATE_DONE:
					dbg.Info("-- Action Done --")
				}
			case output := <-action.Output():
				ws.out <- output
			}
		}
	}()

	// Rate-limit the runs
	limiter := time.Tick(SLEEPYTIME * time.Millisecond)
	var triggeringAction bool
	go func() {
		for {
			<-limiter
			triggeringAction = false
		}
	}()

	go func() {
		for {
			select {

			case event := <-watcher.Events:
				if triggeringAction {
					dbg.Debug("Waiting for action cooldown: %dms", SLEEPYTIME)
					continue
				}
				fname := filepath.Base(event.Name)
				if !filter.Matches(fname) {
					dbg.Debug("Event source does not match filter: %s",
						event.Name)
					continue
				}
				switch {
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					watcher.Remove(event.Name)
					continue
				default:
					dbg.Debug("---- %v on %v ----", event.Op, event.Name)
					ws.state <- types.STATE_RUN
					action.Run()
					triggeringAction = true
				}

			case err := <-watcher.Errors:
				dbg.FatalError("%v", err)
				continue
			}
		}
	}()

	done := make(chan bool)
	for _, fp := range sources.Read() {
		err := watcher.Add(fp)
		if err != nil {
			dbg.FatalError("%v", err)
		}
	}
	<-done
}

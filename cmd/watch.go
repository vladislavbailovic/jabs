package cmd

import (
	"bufio"
	"context"
	"flag"
	"jabs/dbg"
	"jabs/opts"
	"jabs/types"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const SLEEPYTIME time.Duration = 500

type WatchSubcommand struct {
	fs     *flag.FlagSet
	action *string
	out    chan string
	state  chan types.ActionState
}

func NewWatchSubcommand(fs *flag.FlagSet) *WatchSubcommand {
	ws := WatchSubcommand{fs: fs}
	ws.action = fs.String("action", "print", "Action to perform on resource change")
	ws.out = make(chan string)
	ws.state = make(chan types.ActionState)
	return &ws
}

func (s WatchSubcommand) Output() chan string {
	return s.out
}
func (s WatchSubcommand) State() chan types.ActionState {
	return s.state
}

func (ws *WatchSubcommand) Init(ctx context.Context) context.Context {
	ws.state <- types.STATE_INIT
	// ...
	// privates are now populated with flags
	// so init context and return it

	// dbg.Debug("WHATEVER: [%s]", *ws.whatever)
	return ctx
}

func (ws WatchSubcommand) Run() {
	ws.state <- types.STATE_RUN
	var sources []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sources = append(sources, strings.TrimSpace(scanner.Text()))
	}
	options := opts.GetOptions()
	dbg.Info("File from %s", options.Path)
	dbg.Info("Rule is %s", options.Root)
	dbg.Info("%#v", sources)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		dbg.FatalError("%v", err)
	}
	defer watcher.Close()

	action := NewAction(ActionType(*ws.action))

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

	go func() {
		for {
			select {

			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					dbg.Notice("re-adding %s", event.Name)
					watcher.Remove(event.Name)
					time.Sleep(time.Millisecond * SLEEPYTIME)
					err := watcher.Add(event.Name)
					if err != nil {
						dbg.FatalError("%v", err)
					}
					time.Sleep(time.Millisecond * SLEEPYTIME)
					continue
				default:
					dbg.Debug("---- %v on %v ----", event.Op, event.Name)
					time.Sleep(time.Millisecond * SLEEPYTIME)
					action.Run()
					ws.state <- types.STATE_RUN
				}

			case err := <-watcher.Errors:
				dbg.FatalError("%v", err)
				continue
			}
		}
	}()

	done := make(chan bool)
	for _, fp := range sources {
		err := watcher.Add(fp)
		if err != nil {
			dbg.FatalError("%v", err)
		}
	}
	<-done
}

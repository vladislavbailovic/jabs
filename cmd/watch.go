package cmd

import (
	"bufio"
	"context"
	"flag"
	"jabs/dbg"
	"jabs/opts"
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
}

func NewWatchSubcommand(fs *flag.FlagSet) *WatchSubcommand {
	ws := WatchSubcommand{fs: fs}
	ws.action = fs.String("action", "print", "Action to perform on resource change")
	ws.out = make(chan string)
	return &ws
}

func (s WatchSubcommand) Output() chan string {
	return s.out
}

func (ws *WatchSubcommand) Init(ctx context.Context) context.Context {
	// ...
	// privates are now populated with flags
	// so init context and return it

	// dbg.Debug("WHATEVER: [%s]", *ws.whatever)
	return ctx
}

func (ws WatchSubcommand) Run() {
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
		for event := range action.Output() {
			ws.out <- event
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

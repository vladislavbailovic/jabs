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
}

func NewWatchSubcommand(fs *flag.FlagSet) *WatchSubcommand {
	ws := WatchSubcommand{fs: fs}
	ws.action = fs.String("action", "print", "Action to perform on resource change")
	return &ws
}

func (ws *WatchSubcommand) Init(ctx context.Context) context.Context {
	// ...
	// privates are now populated with flags
	// so init context and return it

	// dbg.Debug("WHATEVER: [%s]", *ws.whatever)
	return ctx
}

func (ws WatchSubcommand) Execute() (string, error) {
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

	Action := Print
	switch *ws.action {
	case "print":
		Action = Print
	case "run":
		Action = Run
	}

	go func() {
		for {
			select {

			case event := <-watcher.Events:
				switch {
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					dbg.Warning("re-adding %s", event.Name)
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
					Action()
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

	return "", nil
}

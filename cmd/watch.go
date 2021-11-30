package cmd

import (
	"bufio"
	"context"
	"flag"
	"jabs/dbg"
	"jabs/opts"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

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

	var Action func()
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
				case event.Op&fsnotify.Write == fsnotify.Write:
					log.Printf("Write:  %s: %s", event.Op, event.Name)
					Action()
				case event.Op&fsnotify.Create == fsnotify.Create:
					log.Printf("Create: %s: %s", event.Op, event.Name)
					Action()
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					log.Printf("Remove: %s: %s", event.Op, event.Name)
					Action()
					watcher.Add(event.Name) // @TODO: handle error
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					log.Printf("Rename: %s: %s", event.Op, event.Name)
					Action()
					// @TODO: readd?
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					log.Printf("Chmod:  %s: %s", event.Op, event.Name)
					Action()
				}
			case err := <-watcher.Errors:
				log.Printf("ERROR: %v", err)
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

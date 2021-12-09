package cmd

import (
	"context"
	"flag"
	"jabs/dbg"
	"jabs/opts"
	"jabs/types"
	"os"
	"time"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

const SLEEPYTIME time.Duration = 500

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

	sources := getDirs()
	dbg.Info("Sources: %v", sources)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		dbg.FatalError("%v", err)
	}
	defer watcher.Close()

	action := NewAction(options.Action)
	filter := FilenameFilter{ options.Filter }

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
				fname := filepath.Base(event.Name)
				if !filter.Matches(fname) {
					dbg.Warning("Event source does not match filter: %s", event.Name)
					continue
				}
				switch {
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					dbg.Notice("re-adding %s", event.Name)
					watcher.Remove(event.Name)
					time.Sleep(time.Millisecond * SLEEPYTIME)
					// err := watcher.Add(event.Name)
					// if err != nil {
					// 	dbg.FatalError("%v", err)
					// }
					// time.Sleep(time.Millisecond * SLEEPYTIME)
					continue
				default:
					dbg.Debug("---- %v on %v ----", event.Op, event.Name)
					time.Sleep(time.Millisecond * SLEEPYTIME)
					ws.state <- types.STATE_RUN
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


// @TODO refactor into FS: dirlist and filter structs

func getDirs() []string {
	options := opts.GetOptions()
	root := string(options.Directory)
	if root == "" {
		root = "../*"
	}
	return getSubdirs(root, []string{})
}

func getSubdirs(root string, dirs []string) []string {
	dirFilter := DirFilter{}
	paths, err := filepath.Glob(root)
	if err != nil {
		dbg.FatalError("Unable to get directories (%s): %v", root, err)
	}

	for _, path := range paths {
		if filepath.Base(path)[0:1] == "." {
			continue
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			dbg.FatalError("Unable to determine path for %s: %v", path, err)
		}

		if dirFilter.Matches(abs) {
			dirs = append(dirs, abs)
			dirs = getSubdirs(abs + string(os.PathSeparator) + "*", dirs)
		}
	}

	return dirs
}

type Filter interface {
	Matches(path string) bool
}

type FilenameFilter struct {
	pattern types.FilenamePattern
}

func (pf FilenameFilter)Matches (what string) bool {
	isMatch, err := filepath.Match(string(pf.pattern), what)
	if err != nil {
		return false
	}
	return isMatch
}

type DirFilter struct {}
func (df DirFilter)Matches (what string) bool {
	stat, err := os.Stat(what)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

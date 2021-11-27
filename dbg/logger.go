package dbg

import (
	"fmt"
	"jabs/opts"
	"jabs/types"
	"os"
)

func Log(lvl types.LogLevel, msg string) {
	options := opts.GetOptions()
	if lvl < options.Verbosity {
		return
	}
	fmt.Printf("[%s] %s\n",
		types.LOG_LEVELS[lvl], msg)
}

func Debug(msg string, args ...interface{}) {
	Log(types.LOG_DEBUG, fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...interface{}) {
	Log(types.LOG_INFO, fmt.Sprintf(msg, args...))
}

func Notice(msg string, args ...interface{}) {
	Log(types.LOG_NOTICE, fmt.Sprintf(msg, args...))
}

func Warning(msg string, args ...interface{}) {
	Log(types.LOG_WARNING, fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...interface{}) {
	Log(types.LOG_ERROR, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

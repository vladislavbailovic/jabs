package dbg

import (
	"fmt"
	"jabs/options"
	"jabs/types"
	"os"
)

var levelNames = map[types.LogLevel]string{
	types.LOG_DEBUG:   "debug",
	types.LOG_INFO:    "info",
	types.LOG_NOTICE:  "notice",
	types.LOG_WARNING: "warning",
	types.LOG_ERROR:   "error",
}

func Log(lvl types.LogLevel, msg string) {
	opts := options.GetOptions()
	if lvl < opts.Verbosity {
		return
	}
	fmt.Printf("[%s] %s\n",
		levelNames[lvl], msg)
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

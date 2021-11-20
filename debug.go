package main

import (
	"fmt"
	"os"
)

type LogLevel int

const (
	LOG_DEBUG LogLevel = iota
	LOG_INFO
	LOG_NOTICE
	LOG_WARNING
	LOG_ERROR
)

var levelNames = map[LogLevel]string{
	LOG_DEBUG:   "debug",
	LOG_INFO:    "info",
	LOG_NOTICE:  "notice",
	LOG_WARNING: "warning",
	LOG_ERROR:   "error",
}

func Log(lvl LogLevel, msg string) {
	opts := GetOptions()
	if lvl < opts.Verbosity {
		return
	}
	fmt.Printf("[%s] %s\n",
		levelNames[lvl], msg)
}

func Debug(msg string, args ...interface{}) {
	Log(LOG_DEBUG, fmt.Sprintf(msg, args...))
}

func Info(msg string, args ...interface{}) {
	Log(LOG_INFO, fmt.Sprintf(msg, args...))
}

func Notice(msg string, args ...interface{}) {
	Log(LOG_NOTICE, fmt.Sprintf(msg, args...))
}

func Warning(msg string, args ...interface{}) {
	Log(LOG_WARNING, fmt.Sprintf(msg, args...))
}

func Error(msg string, args ...interface{}) {
	Log(LOG_ERROR, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

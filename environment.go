package main

import (
	"context"
	"jabs/options"
	"os"
	"strconv"
)

var envmap = map[options.OptionKey]string{
	options.OPT_VERBOSITY: "JABS_LOG_LEVEL",
}

func ApplyEnvironment(ctx context.Context) context.Context {
	logLevel, err := strconv.Atoi(os.Getenv(envmap[options.OPT_VERBOSITY]))
	if err != nil {
		logLevel = 0
	}
	ctx = context.WithValue(ctx, options.OPT_VERBOSITY, logLevel)
	return ctx
}

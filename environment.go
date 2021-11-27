package main

import (
	"context"
	"jabs/opts"
	"os"
	"strconv"
)

var envmap = map[opts.OptionKey]string{
	opts.OPT_VERBOSITY: "JABS_LOG_LEVEL",
}

func ApplyEnvironment(ctx context.Context) context.Context {
	logLevel, err := strconv.Atoi(os.Getenv(envmap[opts.OPT_VERBOSITY]))
	if err != nil {
		logLevel = 0
	}
	ctx = context.WithValue(ctx, opts.OPT_VERBOSITY, logLevel)
	return ctx
}

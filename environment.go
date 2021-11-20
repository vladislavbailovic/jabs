package main

import (
	"context"
	"os"
	"strconv"
)

var envmap = map[OptionKey]string{
	OPT_VERBOSITY: "JABS_LOG_LEVEL",
}

func ApplyEnvironment(ctx context.Context) context.Context {
	logLevel, err := strconv.Atoi(os.Getenv(envmap[OPT_VERBOSITY]))
	if err != nil {
		logLevel = 0
	}
	ctx = context.WithValue(ctx, OPT_VERBOSITY, logLevel)
	return ctx
}

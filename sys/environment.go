package sys

import (
	"context"
	"jabs/opts"
	"os"
	"strconv"
)

type envVar string

const (
	ENV_VERBOSITY envVar = "JABS_LOG_LEVEL"
)

var envmap = map[opts.OptionKey]envVar{
	opts.OPT_VERBOSITY: ENV_VERBOSITY,
}

func getEnv(opt opts.OptionKey) string {
	key := envmap[opt]
	return os.Getenv(string(key))
}

func setEnv(opt opts.OptionKey, val string) {
	key := envmap[opt]
	os.Setenv(string(key), val)
}

func ApplyEnvironment(ctx context.Context) context.Context {
	logLevel, err := strconv.Atoi(getEnv(opts.OPT_VERBOSITY))
	if err == nil {
		ctx = context.WithValue(ctx, opts.OPT_VERBOSITY, logLevel)
	}
	return ctx
}

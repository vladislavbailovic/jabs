package opts

import (
	"context"
	"jabs/types"
)

type OptionKey int

const (
	OPT_WATCH OptionKey = iota
	OPT_ROOT
	OPT_PATH
	OPT_FORCE
	OPT_VERBOSITY

	// Print subcommand
	OPT_CONDITIONS
)

type Options struct {
	Watch     bool
	Root      types.RuleName
	Path      string
	Force     bool
	Verbosity types.LogLevel

	// Print subcommand
	Conditions bool
}

var _options Options

func getBoolean(ctx context.Context, key OptionKey) bool {
	val := ctx.Value(key)
	if val == nil {
		return false
	}
	return val.(bool)
}

func getString(ctx context.Context, key OptionKey) string {
	val := ctx.Value(key)
	if val == nil {
		return ""
	}
	return val.(string)
}

func getInt(ctx context.Context, key OptionKey) int {
	val := ctx.Value(key)
	if val == nil {
		return 0
	}
	return val.(int)
}

func InitOptions(ctx context.Context) Options {
	if (Options{}) != _options {
		return _options
	}
	// @TODO add defaults
	_options = Options{
		Watch:     getBoolean(ctx, OPT_WATCH),
		Root:      types.RuleName(getString(ctx, OPT_ROOT)),
		Path:      getString(ctx, OPT_PATH),
		Force:     getBoolean(ctx, OPT_FORCE),
		Verbosity: types.LogLevel(getInt(ctx, OPT_VERBOSITY)),

		// Print subcommand
		Conditions: getBoolean(ctx, OPT_CONDITIONS),
	}
	return _options
}

func GetOptions() Options {
	if (Options{}) == _options {
		return InitOptions(context.TODO())
	}
	return _options
}

func resetOptions() {
	_options = Options{}
}

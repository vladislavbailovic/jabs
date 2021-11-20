package main

import "context"

type OptionKey int

const (
	OPT_WATCH OptionKey = iota
	OPT_ROOT
	OPT_PATH
	OPT_VERBOSITY
)

type Options struct {
	Watch     bool
	Root      string
	Path      string
	Verbosity LogLevel
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

func NewOptions(ctx context.Context) Options {
	if (Options{}) != _options {
		return _options
	}
	// @TODO add defaults
	_options = Options{
		Watch:     getBoolean(ctx, OPT_WATCH),
		Root:      getString(ctx, OPT_ROOT),
		Path:      getString(ctx, OPT_PATH),
		Verbosity: LogLevel(getInt(ctx, OPT_VERBOSITY)),
	}
	return _options
}

func GetOptions() Options {
	if (Options{}) == _options {
		return NewOptions(context.TODO())
	}
	return _options
}

func resetOptions() {
	_options = Options{}
}

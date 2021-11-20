package main

import "context"

type OptionKey int

const (
	OPT_WATCH OptionKey = iota
	OPT_ROOT
	OPT_PATH
)

type Options struct {
	Watch bool
	Root  string
	Path  string
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

func NewOptions(ctx context.Context) Options {
	if (Options{}) != _options {
		return _options
	}
	// @TODO add defaults
	_options = Options{
		Watch: getBoolean(ctx, OPT_WATCH),
		Root:  getString(ctx, OPT_ROOT),
		Path:  getString(ctx, OPT_PATH),
	}
	return _options
}

func resetOptions() {
	_options = Options{}
}

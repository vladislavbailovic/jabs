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

func NewOptions(ctx context.Context) Options {
	if (Options{}) != _options {
		return _options
	}
	_options = Options{
		Watch: ctx.Value(OPT_WATCH).(bool),
		Root:  ctx.Value(OPT_ROOT).(string),
		Path:  ctx.Value(OPT_PATH).(string),
	}
	return _options
}

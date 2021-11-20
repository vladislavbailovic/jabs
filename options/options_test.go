package options

import (
	"context"
	"testing"
)

func Test_Options(t *testing.T) {
	ctx := context.TODO()
	ctx = context.WithValue(ctx, OPT_WATCH, true)
	ctx = context.WithValue(ctx, OPT_ROOT, "root")

	ctx1 := context.WithValue(ctx, OPT_PATH, "./examples/first.yml")
	opts1 := InitOptions(ctx1)

	ctx2 := context.WithValue(ctx, OPT_PATH, "./examples/self.yml")
	opts2 := InitOptions(ctx2)
	if opts1.Path != opts2.Path {
		t.Fatalf("Options will only be loaded once")
	}
	resetOptions()
}

func Test_Options_Defaults(t *testing.T) {
	ctx := context.TODO()
	opts := InitOptions(ctx)

	if "" != opts.Root {
		t.Fatalf("defaults are not yet implemented")
	}
	resetOptions()
}

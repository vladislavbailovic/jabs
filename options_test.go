package main

import (
	"context"
	"testing"
)

func Test_Options(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, OPT_WATCH, true)
	ctx = context.WithValue(ctx, OPT_ROOT, "root")

	ctx1 := context.WithValue(ctx, OPT_PATH, "./examples/first.yml")
	opts1 := NewOptions(ctx1)

	ctx2 := context.WithValue(ctx, OPT_PATH, "./examples/self.yml")
	opts2 := NewOptions(ctx2)
	if opts1.Path != opts2.Path {
		t.Fatalf("Options will only be loaded once")
	}
}

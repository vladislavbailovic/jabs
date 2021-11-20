package main

import (
	"context"
	"fmt"
	"jabs/options"
	"jabs/types"
	"os"
	"testing"
)

func Test_Environment(t *testing.T) {
	ctx := ApplyEnvironment(context.TODO())
	level := ctx.Value(options.OPT_VERBOSITY)
	if nil == level {
		t.Fatalf("Expected to have initial log level")
	}

	os.Setenv(envmap[options.OPT_VERBOSITY], fmt.Sprintf("%v", types.LOG_ERROR))
	ctx = ApplyEnvironment(context.TODO())
	l2 := ctx.Value(options.OPT_VERBOSITY)
	if l2 == nil {
		t.Fatalf("Expected to have set log level")
	}
	if types.LogLevel(l2.(int)) != types.LOG_ERROR {
		t.Fatalf("Expected %v (error), got %v", types.LOG_ERROR, l2)
	}
	os.Setenv(envmap[options.OPT_VERBOSITY], "")
}

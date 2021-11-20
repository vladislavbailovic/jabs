package main

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func Test_Environment(t *testing.T) {
	ctx := ApplyEnvironment(context.TODO())
	level := ctx.Value(OPT_VERBOSITY)
	if nil == level {
		t.Fatalf("Expected to have initial log level")
	}

	os.Setenv(envmap[OPT_VERBOSITY], fmt.Sprintf("%v", LOG_ERROR))
	ctx = ApplyEnvironment(context.TODO())
	l2 := ctx.Value(OPT_VERBOSITY)
	if l2 == nil {
		t.Fatalf("Expected to have set log level")
	}
	if LogLevel(l2.(int)) != LOG_ERROR {
		t.Fatalf("Expected %v (error), got %v", LOG_ERROR, l2)
	}
	os.Setenv(envmap[OPT_VERBOSITY], "")
}

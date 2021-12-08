package sys

import (
	"context"
	"fmt"
	"jabs/opts"
	"jabs/types"
	"testing"
)

func Test_Environment(t *testing.T) {
	ctx := ApplyEnvironment(context.TODO())

	setEnv(opts.OPT_VERBOSITY, fmt.Sprintf("%v", types.LOG_ERROR))
	ctx = ApplyEnvironment(context.TODO())
	l2 := ctx.Value(opts.OPT_VERBOSITY)
	if l2 == nil {
		t.Fatalf("Expected to have set log level")
	}
	if types.LogLevel(l2.(int)) != types.LOG_ERROR {
		t.Fatalf("Expected %v (error), got %v", types.LOG_ERROR, l2)
	}
	setEnv(opts.OPT_VERBOSITY, "")
}

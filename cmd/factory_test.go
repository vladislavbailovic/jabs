package cmd

import (
	"flag"
	"jabs/types"
	"testing"
)

func Test_ActionType(t *testing.T) {
	var kind types.ActionType

	kind = ActionType("print")
	if kind != ACTION_PRINT {
		t.Log(kind)
		t.Fatalf("expected print action type")
	}

	kind = ActionType("run")
	if kind != ACTION_RUN {
		t.Log(kind)
		t.Fatalf("expected run action type")
	}

	kind = ActionType("whatever")
	if kind != ACTION_DEFAULT {
		t.Log(kind)
		t.Fatalf("expected default action type")
	}
}

func Test_NewAction(t *testing.T) {
	var action types.Action
	var ok bool

	action = NewAction(ACTION_PRINT)
	_, ok = action.(PrintAction)
	if !ok {
		t.Log(action)
		t.Fatalf("expected print action request to work properly")
	}

	action = NewAction(ACTION_RUN)
	_, ok = action.(RunAction)
	if !ok {
		t.Log(action)
		t.Fatalf("expected run action request to work properly")
	}

	action = NewAction(ACTION_DEFAULT)
	_, ok = action.(PrintAction)
	if !ok {
		t.Log(action)
		t.Fatalf("expected print action to be default")
	}
}

func Test_SubcommandType(t *testing.T) {
	if SUBCMD_PRINT != SubcommandType("print") {
		t.Fatalf("expected print subcommand type")
	}
	if SUBCMD_RUN != SubcommandType("run") {
		t.Fatalf("expected run subcommand type")
	}
	if SUBCMD_WATCH != SubcommandType("watch") {
		t.Fatalf("expected watch subcommand type")
	}
	if SUBCMD_DEFAULT != SubcommandType("whatever") {
		t.Fatalf("expected default subcommand type")
	}
}

func Test_NewSubcommand(t *testing.T) {
	var subcmd types.Subcommand
	var ok bool
	fs := flag.NewFlagSet("jabs", flag.ContinueOnError)

	subcmd = NewSubcommand(SUBCMD_PRINT, fs)
	_, ok = subcmd.(*PrintSubcommand)
	if !ok {
		t.Log(subcmd)
		t.Fatalf("expected new print subcommand")
	}

	subcmd = NewSubcommand(SUBCMD_RUN, fs)
	_, ok = subcmd.(*RunSubcommand)
	if !ok {
		t.Log(subcmd)
		t.Fatalf("expected new run subcommand")
	}

	subcmd = NewSubcommand(SUBCMD_WATCH, fs)
	_, ok = subcmd.(*WatchSubcommand)
	if !ok {
		t.Log(subcmd)
		t.Fatalf("expected new watch subcommand")
	}

	fs2 := flag.NewFlagSet("jabs", flag.ContinueOnError)
	subcmd = NewSubcommand(SUBCMD_DEFAULT, fs2)
	_, ok = subcmd.(*PrintSubcommand)
	if !ok {
		t.Log(subcmd)
		t.Fatalf("expected new print subcommand as default")
	}
}

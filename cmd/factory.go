package cmd

import (
	"flag"
	"jabs/types"
)

const (
	ACTION_DEFAULT types.ActionType = iota
	ACTION_PRINT
	ACTION_RUN
)

func ActionType(which string) types.ActionType {
	switch which {
	case "run":
		return ACTION_RUN
	case "print":
		return ACTION_PRINT
	default:
		return ACTION_DEFAULT
	}
}

func NewAction(kind types.ActionType) types.Action {
	switch kind {
	case ACTION_DEFAULT:
		return PrintAction{}
	case ACTION_PRINT:
		return PrintAction{}
	case ACTION_RUN:
		return RunAction{}
	}
	return nil
}

const (
	SUBCMD_DEFAULT types.SubcommandType = iota
	SUBCMD_RUN
	SUBCMD_PRINT
	SUBCMD_WATCH
)

func SubcommandType(which string) types.SubcommandType {
	switch which {
	case "run":
		return SUBCMD_RUN
	case "print":
		return SUBCMD_PRINT
	case "watch":
		return SUBCMD_WATCH
	default:
		return SUBCMD_DEFAULT
	}
}

func NewSubcommand(kind types.SubcommandType, fs *flag.FlagSet) types.Subcommand {
	switch kind {
	case SUBCMD_DEFAULT:
		return NewPrintSubcommand(fs)
	case SUBCMD_PRINT:
		return NewPrintSubcommand(fs)
	case SUBCMD_RUN:
		return NewRunSubcommand(fs)
	case SUBCMD_WATCH:
		return NewWatchSubcommand(fs)
	}
	return nil
}
package parse

import (
	"fmt"
	"strings"
	"time"
)

const (
	MACRO_OPEN  string = "${{"
	MACRO_CLOSE string = "}}"
)

func GetSystemMacroDefinitions() []MacroDefinition {
	now := time.Now().Unix()
	sys := []MacroDefinition{
		MacroDefinition{Name: "System:LastRuntime", Value: fmt.Sprintf("%d", now)},
		MacroDefinition{Name: "System:Now", Value: fmt.Sprintf("%d", now)},
	}
	return sys
}

func GetExpansionKey(what string) string {
	return MACRO_OPEN + what + MACRO_CLOSE
}

func ContainsMacros(subject string) bool {
	return strings.Contains(subject, MACRO_OPEN)
}

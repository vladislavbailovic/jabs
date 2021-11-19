package main

import (
	"fmt"
	"time"
)

func GetSystemMacroDefinitions() []MacroDefinition {
	now := time.Now().Unix()
	sys := []MacroDefinition{
		MacroDefinition{Name: "System:LastRuntime", Value: fmt.Sprintf("%d", now)},
		MacroDefinition{Name: "System:Now", Value: fmt.Sprintf("%d", now)},
	}
	return sys
}

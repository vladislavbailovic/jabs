package main

import (
	"jabs/dbg"
	"strings"
)

type MacroDefinitions struct {
	Dfns   []MacroDefinition
	Macros map[string]Macro
}

func NewMacroDefinitions(dfns []MacroDefinition) *MacroDefinitions {
	md := MacroDefinitions{dfns, map[string]Macro{}}
	md.preprocess()
	return &md
}

func (md *MacroDefinitions) preprocess() int {
	if len(md.Dfns) == 0 {
		return 0
	}
	for i := Limit(0); i < LIMIT_MACRO_EXPANSION_PASS; i++ {
		initialCount := len(md.Dfns)
		if initialCount == 0 {
			// We're all done here
			break
		}
		md.convertSimple()
		md.convertValue()
		md.convertShellcode()
		currentCount := len(md.Dfns)
		if initialCount == currentCount && currentCount > 0 {
			dbg.Warning("Some macro definitions could not be expanded: (%d)", currentCount)
			dbg.Debug("--------------------")
			for idx, dfn := range md.Dfns {
				dbg.Debug("\t%d) %s", idx, dfn.Name)
			}
			dbg.Debug("--------------------")
			break
		}
	}
	return len(md.Dfns)
}

func (md *MacroDefinitions) convertSimple() {
	for i := Limit(0); i < LIMIT_MACRO_EXPANSION_PASS; i++ {
		initialCount := len(md.Dfns)
		if initialCount == 0 {
			// We're all done here
			break
		}
		for idx, dfn := range md.Dfns {
			value := dfn.Value
			if "" == value {
				value = dfn.Command
			}
			if ContainsMacros(value) {
				continue
			}
			if "" == value {
				value = execute(value)
			}
			dbg.Debug("Processed any dfn '%s' fully, adding it as a macro", dfn.Name)
			md.Macros[dfn.Name] = Macro{dfn.Name, value}
			md.Dfns = append(md.Dfns[:idx], md.Dfns[idx+1:]...)
			break
		}
		currentCount := len(md.Dfns)
		if initialCount == currentCount {
			break
		}
	}
}

func (md *MacroDefinitions) convertValue() {
	for _, macro := range md.Macros {
		for idx, dfn := range md.Dfns {
			// Not a value macro?
			if "" == dfn.Value {
				continue
			}

			// We've already preprocessed this definition entirely
			// Add it as a macro and remove definition
			if !ContainsMacros(md.Dfns[idx].Value) {
				dbg.Debug("Processed dfn '%s' by value, adding it as a macro", dfn.Name)
				md.Macros[dfn.Name] = Macro{dfn.Name, dfn.Value}
				md.Dfns = append(md.Dfns[:idx], md.Dfns[idx+1:]...)
				break
			}

			key := GetExpansionKey(macro.Name)
			if !strings.Contains(dfn.Value, key) {
				continue
			}

			md.Dfns[idx].Value = strings.Replace(dfn.Value, key, macro.Value, -1)

			// We've already preprocessed this definition entirely
			// Add it as a macro and remove definition
			if !ContainsMacros(md.Dfns[idx].Value) {
				dbg.Debug("Processed dfn '%s' by value, adding it as a macro", dfn.Name)
				md.Macros[dfn.Name] = Macro{dfn.Name, dfn.Value}
				md.Dfns = append(md.Dfns[:idx], md.Dfns[idx+1:]...)
				break
			}
		}
	}
}

func (md *MacroDefinitions) convertShellcode() {
	for _, macro := range md.Macros {
		for idx, dfn := range md.Dfns {
			// Not a command macro?
			if "" == dfn.Command {
				continue
			}
			dbg.Debug("\tchecking command %s", dfn.Name)

			// We've already preprocessed this definition entirely
			// Add it as a macro and remove definition
			if !ContainsMacros(md.Dfns[idx].Command) {
				dbg.Debug("Processed dfn '%s' by cmd, adding it as a macro", dfn.Name)
				md.Macros[dfn.Name] = Macro{dfn.Name, execute(md.Dfns[idx].Command)}
				md.Dfns = append(md.Dfns[:idx], md.Dfns[idx+1:]...)
				break
			}

			key := GetExpansionKey(macro.Name)
			if !strings.Contains(dfn.Command, key) {
				dbg.Debug("\t\tcommand does not have %s", key)
				continue
			}
			dbg.Debug("\t\tcommand DOES have %s, processing", key)

			md.Dfns[idx].Command = strings.Replace(dfn.Command, key, macro.Value, -1)
			dbg.Debug("\t\tresult: %s", md.Dfns[idx].Command)

			// We've already preprocessed this definition entirely
			// Add it as a macro and remove definition
			if !ContainsMacros(md.Dfns[idx].Command) {
				dbg.Debug("Processed dfn '%s' by cmd, adding it as a macro", dfn.Name)
				md.Macros[dfn.Name] = Macro{dfn.Name, execute(md.Dfns[idx].Command)}
				md.Dfns = append(md.Dfns[:idx], md.Dfns[idx+1:]...)
				break
			}
		}
	}
}

func execute(what string) string {
	cmd := NewExecutable(what)
	out, err := cmd.Execute()
	if err != nil {
		dbg.Error("Unable to run command %s: %v", what, err)
	}
	return out
}

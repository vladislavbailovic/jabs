package main

import (
	"fmt"
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
			// We've done nothing
			// @TODO we have unexpanded macros - might wanna warn
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
			if containsMacros(value) {
				continue
			}
			if "" == value {
				value = execute(value)
			}
			fmt.Printf("Processed any dfn '%s' fully, adding it as a macro\n", dfn.Name)
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
			if !containsMacros(md.Dfns[idx].Value) {
				fmt.Printf("Processed dfn '%s' by value, adding it as a macro\n", dfn.Name)
				md.Macros[dfn.Name] = Macro{dfn.Name, dfn.Value}
				md.Dfns = append(md.Dfns[:idx], md.Dfns[idx+1:]...)
				break
			}

			key := getExpansionKey(macro.Name)
			if !strings.Contains(dfn.Value, key) {
				continue
			}

			md.Dfns[idx].Value = strings.Replace(dfn.Value, key, macro.Value, -1)

			// We've already preprocessed this definition entirely
			// Add it as a macro and remove definition
			if !containsMacros(md.Dfns[idx].Value) {
				fmt.Printf("Processed dfn '%s' by value, adding it as a macro\n", dfn.Name)
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
			fmt.Println("\tchecking command", dfn.Name)

			// We've already preprocessed this definition entirely
			// Add it as a macro and remove definition
			if !containsMacros(md.Dfns[idx].Command) {
				fmt.Printf("Processed dfn '%s' by cmd, adding it as a macro\n", dfn.Name)
				md.Macros[dfn.Name] = Macro{dfn.Name, execute(md.Dfns[idx].Command)}
				md.Dfns = append(md.Dfns[:idx], md.Dfns[idx+1:]...)
				break
			}

			key := getExpansionKey(macro.Name)
			if !strings.Contains(dfn.Command, key) {
				fmt.Println("\t\tcommand does not have", key)
				continue
			}
			fmt.Println("\t\tcommand DOES have", key, "processing")

			md.Dfns[idx].Command = strings.Replace(dfn.Command, key, macro.Value, -1)
			fmt.Println("\t\tresult", md.Dfns[idx].Command)

			// We've already preprocessed this definition entirely
			// Add it as a macro and remove definition
			if !containsMacros(md.Dfns[idx].Command) {
				fmt.Printf("Processed dfn '%s' by cmd, adding it as a macro\n", dfn.Name)
				md.Macros[dfn.Name] = Macro{dfn.Name, execute(md.Dfns[idx].Command)}
				md.Dfns = append(md.Dfns[:idx], md.Dfns[idx+1:]...)
				break
			}
		}
	}
}

func containsMacros(where string) bool {
	return strings.Contains(where, "${{")
}

func execute(what string) string {
	cmd := NewExecutable(what)
	out, err := cmd.Execute()
	if err != nil {
		fmt.Printf("ERROR: unable to run command %s\n", what)
		panic(err)
	}
	return out
}
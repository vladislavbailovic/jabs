package main

import (
	"fmt"
	"testing"
)

func Test_MacroDefinitions_First(t *testing.T) {
	p := NewParser("./examples/first.yml")
	dfns := NewMacroDefinitions(p.MacroDefinitions)

	if len(dfns.Dfns) > 1 {
		fmt.Println("Macros")
		fmt.Println("--------------------")
		for name, macro := range dfns.Macros {
			fmt.Printf("\t- Macro %s: [%s]\n", name, macro.Value)
		}
		fmt.Println("--------------------")

		fmt.Println("Leftovers")
		fmt.Println("--------------------")
		for idx, dfn := range dfns.Dfns {
			fmt.Printf("\t%d) Dfn %s: [%s]\n", idx+1, dfn.Name, dfn.Value)
		}
		fmt.Println("--------------------")

		t.Fatalf("expected exactly 1 unprocessed definition, got %d", len(dfns.Dfns))
	}
}

func Test_MacroDefinitions_Self(t *testing.T) {
	p := NewParser("./examples/self.yml")
	dfns := NewMacroDefinitions(p.MacroDefinitions)

	if len(dfns.Dfns) > 0 {
		fmt.Println("Macros")
		fmt.Println("--------------------")
		for name, macro := range dfns.Macros {
			fmt.Printf("\t- Macro %s: [%s]\n", name, macro.Value)
		}
		fmt.Println("--------------------")

		fmt.Println("Leftovers")
		fmt.Println("--------------------")
		for idx, dfn := range dfns.Dfns {
			fmt.Printf("\t%d) Dfn %s: [%s]\n", idx+1, dfn.Name, dfn.Value)
		}
		fmt.Println("--------------------")

		t.Fatalf("expected no unprocessed definitions, got %d", len(dfns.Dfns))
	}
}

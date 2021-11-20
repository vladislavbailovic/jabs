package main

import (
	"testing"
)

func Test_MacroDefinitions_First(t *testing.T) {
	p := NewParser("./examples/first.yml")
	dfns := NewMacroDefinitions(p.MacroDefinitions)

	if len(dfns.Dfns) > 1 {
		Debug("Macros")
		Debug("--------------------")
		for name, macro := range dfns.Macros {
			Debug("\t- Macro %s: [%s]", name, macro.Value)
		}
		Debug("--------------------")

		Debug("Leftovers")
		Debug("--------------------")
		for idx, dfn := range dfns.Dfns {
			Debug("\t%d) Dfn %s: [%s]", idx+1, dfn.Name, dfn.Value)
		}
		Debug("--------------------")

		t.Fatalf("expected exactly 1 unprocessed definition, got %d", len(dfns.Dfns))
	}
}

func Test_MacroDefinitions_Self(t *testing.T) {
	p := NewParser("./examples/self.yml")
	dfns := NewMacroDefinitions(p.MacroDefinitions)

	if len(dfns.Dfns) > 0 {
		Debug("Macros")
		Debug("--------------------")
		for name, macro := range dfns.Macros {
			Debug("\t- Macro %s: [%s]", name, macro.Value)
		}
		Debug("--------------------")

		Debug("Leftovers")
		Debug("--------------------")
		for idx, dfn := range dfns.Dfns {
			Debug("\t%d) Dfn %s: [%s]", idx+1, dfn.Name, dfn.Value)
		}
		Debug("--------------------")

		t.Fatalf("expected no unprocessed definitions, got %d", len(dfns.Dfns))
	}
}

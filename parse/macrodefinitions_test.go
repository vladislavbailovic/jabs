package parse

import (
	"jabs/dbg"
	"testing"
)

func Test_MacroDefinitions_First(t *testing.T) {
	p := NewParser("../examples/first.yml")
	dfns := NewMacroDefinitions(p.MacroDefinitions)

	if len(dfns.Dfns) > 1 {
		dbg.Debug("Macros")
		dbg.Debug("--------------------")
		for name, macro := range dfns.Macros {
			dbg.Debug("\t- Macro %s: [%s]", name, macro.Value)
		}
		dbg.Debug("--------------------")

		dbg.Debug("Leftovers")
		dbg.Debug("--------------------")
		for idx, dfn := range dfns.Dfns {
			dbg.Debug("\t%d) Dfn %s: [%s]", idx+1, dfn.Name, dfn.Value)
		}
		dbg.Debug("--------------------")

		t.Fatalf("expected exactly 1 unprocessed definition, got %d", len(dfns.Dfns))
	}
}

func Test_MacroDefinitions_Self(t *testing.T) {
	p := NewParser("../examples/self.yml")
	dfns := NewMacroDefinitions(p.MacroDefinitions)

	if len(dfns.Dfns) > 0 {
		dbg.Debug("Macros")
		dbg.Debug("--------------------")
		for name, macro := range dfns.Macros {
			dbg.Debug("\t- Macro %s: [%s]", name, macro.Value)
		}
		dbg.Debug("--------------------")

		dbg.Debug("Leftovers")
		dbg.Debug("--------------------")
		for idx, dfn := range dfns.Dfns {
			dbg.Debug("\t%d) Dfn %s: [%s]", idx+1, dfn.Name, dfn.Value)
		}
		dbg.Debug("--------------------")

		t.Fatalf("expected no unprocessed definitions, got %d", len(dfns.Dfns))
	}
}

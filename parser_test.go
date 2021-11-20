package main

import (
	"jabs/dbg"
	"testing"
)

func Test_Parser(t *testing.T) {
	p := NewParser("./examples/first.yml")
	expectedMacroDfns := 5
	expectedRuleDfns := 7
	if len(p.MacroDefinitions) != expectedMacroDfns {
		dbg.Debug("Macros")
		dbg.Debug("--------------------")
		for i, dfn := range p.MacroDefinitions {
			dbg.Debug("\t%d) ${{%s}}: [%v]", i, dfn.Name, dfn.Value)
		}
		t.Fatalf("expected %d macro dfns, got %d", expectedMacroDfns, len(p.MacroDefinitions))
	}

	if len(p.RuleDefinitions) != expectedRuleDfns {
		dbg.Debug("Rules")
		dbg.Debug("--------------------")
		for i, dfn := range p.RuleDefinitions {
			dbg.Debug("\t%d) ${{%s}}", i, dfn.Name)
		}
		t.Fatalf("expected %d rule dfns, got %d", expectedRuleDfns, len(p.RuleDefinitions))
	}
}

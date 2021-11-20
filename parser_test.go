package main

import (
	"testing"
)

func Test_Parser(t *testing.T) {
	p := NewParser("./examples/first.yml")
	expectedMacroDfns := 5
	expectedRuleDfns := 7
	if len(p.MacroDefinitions) != expectedMacroDfns {
		Debug("Macros")
		Debug("--------------------")
		for i, dfn := range p.MacroDefinitions {
			Debug("\t%d) ${{%s}}: [%v]", i, dfn.Name, dfn.Value)
		}
		t.Fatalf("expected %d macro dfns, got %d", expectedMacroDfns, len(p.MacroDefinitions))
	}

	if len(p.RuleDefinitions) != expectedRuleDfns {
		Debug("Rules")
		Debug("--------------------")
		for i, dfn := range p.RuleDefinitions {
			Debug("\t%d) ${{%s}}", i, dfn.Name)
		}
		t.Fatalf("expected %d rule dfns, got %d", expectedRuleDfns, len(p.RuleDefinitions))
	}
}

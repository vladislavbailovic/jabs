package main

import (
	"fmt"
	"testing"
)

func Test_Parser(t *testing.T) {
	p := NewParser("./examples/first.yml")
	expectedMacroDfns := 5
	expectedRuleDfns := 5
	if len(p.MacroDefinitions) != expectedMacroDfns {
		fmt.Println("Macros")
		fmt.Println("--------------------")
		for i, dfn := range p.MacroDefinitions {
			fmt.Printf("\t%d) ${{%s}}: [%v]\n", i, dfn.Name, dfn.Value)
		}
		t.Fatalf("expected %d macro dfns, got %d", expectedMacroDfns, len(p.MacroDefinitions))
	}

	if len(p.RuleDefinitions) != expectedRuleDfns {
		fmt.Println("Rules")
		fmt.Println("--------------------")
		for i, dfn := range p.RuleDefinitions {
			fmt.Printf("\t%d) ${{%s}}\n", i, dfn.Name)
		}
		t.Fatalf("expected %d rule dfns, got %d", expectedRuleDfns, len(p.RuleDefinitions))
	}
}

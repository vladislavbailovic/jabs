package main

import (
	"testing"
)

func Test_Preprocessor(t *testing.T) {
	p := NewPreprocessor("./examples/first.yml")
	expectedMacros := 4
	actualMacros := 0
	for name, macro := range p.Macros {
		if len(macro.Value) > 100 {
			t.Fatalf("macro value length too long for: %s (%d)", name, len(macro.Value))
		}
		actualMacros += 1
	}
	if expectedMacros != actualMacros {
		Debug("Macros")
		Debug("--------------------")
		for name, _ := range p.Macros {
			Debug("%s", name)
		}
		Debug("--------------------")
		t.Fatalf("expected %d macros, but got %d instead", expectedMacros, actualMacros)
	}
}

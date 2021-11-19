package main

import "testing"

func Test_Preprocessor(t *testing.T) {
	p := NewPreprocessor("./examples/first.yml")
	for name, macro := range p.Macros {
		if len(macro.Value) > 100 {
			t.Fatalf("macro value length too long for: %s (%d)", name, len(macro.Value))
		}
	}
}

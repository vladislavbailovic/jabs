package main

import (
	"fmt"
	"testing"
)

func Test_Preprocessor(t *testing.T) {
	p := NewPreprocessor("./examples/first.yml")
	expectedMacros := 7
	actualMacros := 0
	for name, macro := range p.Macros {
		if len(macro.Value) > 100 {
			t.Fatalf("macro value length too long for: %s (%d)", name, len(macro.Value))
		}
		actualMacros += 1
	}
	if expectedMacros != actualMacros {
		fmt.Println("Macros")
		fmt.Println("--------------------")
		for name, _ := range p.Macros {
			fmt.Printf("%s\n", name)
		}
		fmt.Println("--------------------")
		t.Fatalf("expected %d macros, but got %d instead", expectedMacros, actualMacros)
	}
}

package main

import (
	"testing"
)

func Test_EvaluationStack_Dependent(t *testing.T) {
	p := NewPreprocessor("./examples/first.yml")
	es := NewEvaluationStack("root", p.Rules)
	expectedRules := 4
	if len(es.stack) != expectedRules {
		Debug("Stack")
		Debug("--------------------")
		for _, rl := range es.stack {
			Debug("\t- %s", rl.Name)
		}
		Debug("--------------------")
		t.Fatalf("expected %d rules, but got %d", expectedRules, len(es.stack))
	}
}

func Test_EvaluationStack_Standalone(t *testing.T) {
	p := NewPreprocessor("./examples/first.yml")
	es := NewEvaluationStack("Standalone rule", p.Rules)
	expectedRules := 1
	if len(es.stack) != expectedRules {
		Debug("Stack")
		Debug("--------------------")
		for _, rl := range es.stack {
			Debug("\t- %s", rl.Name)
		}
		Debug("--------------------")
		t.Fatalf("expected %d rules, but got %d", expectedRules, len(es.stack))
	}
}

func Test_EvaluationStack_Observable(t *testing.T) {
	p := NewPreprocessor("./examples/first.yml")
	es := NewEvaluationStack("Wants subtask with failing observable", p.Rules)
	expectedRules := 2
	if len(es.stack) != expectedRules {
		Debug("Stack")
		Debug("--------------------")
		for _, rl := range es.stack {
			Debug("\t- %s", rl.Name)
		}
		Debug("--------------------")
		t.Fatalf("expected %d rules, but got %d", expectedRules, len(es.stack))
	}
}

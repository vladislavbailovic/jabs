package main

import "jabs/dbg"

type EvaluationStack struct {
	root  string
	rules map[string]Rule
	stack []Rule
}

func NewEvaluationStack(root string, rules map[string]Rule) EvaluationStack {
	es := EvaluationStack{root, rules, []Rule{}}
	es.init()
	return es
}

func (es *EvaluationStack) init() {
	es.stack = es.getSubstack(es.root, []Rule{})
}

func (es *EvaluationStack) getSubstack(root string, stack []Rule) []Rule {
	rule, ok := es.rules[root]
	if !ok {
		dbg.Error("No such rule: %s", root)
	}

	// Test rule observables state - should we even consider this rule?
	if len(rule.Observes) > 0 {
		observedChange := false
		for _, observable := range rule.Observes {
			cmd := NewExecutable(observable)
			_, err := cmd.Execute()
			if err != nil {
				observedChange = true
			}
		}
		if !observedChange {
			return stack
		}
	}

	for _, dependency := range rule.DependsOn {
		stack = es.getSubstack(dependency, stack)
	}

	stack = append(stack, rule)
	return stack
}

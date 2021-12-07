package parse

import (
	"jabs/dbg"
	"jabs/types"
)

type EvaluationStack struct {
	root  string
	rules map[string]types.Rule
	stack []types.Rule
}

func NewEvaluationStack(root string, rules map[string]types.Rule) EvaluationStack {
	es := EvaluationStack{root, rules, []types.Rule{}}
	es.init()
	return es
}

func (es *EvaluationStack) init() {
	es.stack = es.getSubstack(es.root, []types.Rule{})
}

func (es EvaluationStack) GetStack() []types.Rule {
	return es.stack
}

func (es *EvaluationStack) getSubstack(root string, stack []types.Rule) []types.Rule {
	rule, ok := es.rules[root]
	if !ok {
		dbg.FatalError("No such rule: %s", root)
	}

	// Test rule observables state - should we even consider this rule?
	if len(rule.Observes) > 0 {
		observedChange := false
		for _, observable := range rule.Observes {
			_, err := observable.Execute()
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

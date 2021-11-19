package main

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
		panic("No such rule: " + root)
	}

	for _, dependency := range rule.DependsOn {
		stack = es.getSubstack(dependency, stack)
	}

	stack = append(stack, rule)
	return stack
}

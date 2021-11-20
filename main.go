package main

import (
	"fmt"
	"strings"
)

// Kinda print
func main() {
	p := NewPreprocessor("./examples/self.yml")
	es := NewEvaluationStack("cover:html", p.Rules)
	out := []string{"#!/bin/bash", ""}

	for idx, rule := range es.stack {
		out = append(out, fmt.Sprintf("# Rule %d -- %s", idx, rule.Name))
		for _, task := range rule.Tasks {
			cmd := NewScriptable(task)
			out = append(out, cmd.GetScript())
		}
		out = append(out, "")
	}
	Info("\n" + strings.Join(out[:], "\n"))
}

// Kinda execute
// func main() {
// 	p := NewPreprocessor("./examples/first.yml")

// 	Debug("Macros")
// 	Debug("--------------------")
// 	for name, value := range p.Macros {
// 		Debug("\t- ${{%s}}: [%v]", name, value)
// 	}

// 	Debug("Rules")
// 	Debug("--------------------")
// 	for name, rule := range p.Rules {
// 		Debug("\t- %s: %v", name, rule)
// 	}

// 	stack := NewEvaluationStack("root", p.Rules)
// 	Debug("Stack")
// 	Debug("--------------------")
// 	for _, rl := range stack.stack {
// 		Debug("\t- %s", rl.Name)
// 		for i, task := range rl.Tasks {
// 			cmd := NewExecutable(task)
// 			out, err := cmd.Execute()
// 			if err != nil {
// 				Error("%v", err)
// 			}
// 			Debug("\t\t%d) %v", i+1, out)
// 		}
// 	}

// 	// @TODO execute or compile stack
// }

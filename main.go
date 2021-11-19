package main

import (
	"fmt"
	"strings"
)

// Kinda print
func main() {
	p := NewPreprocessor("./examples/self.yml")
	es := NewEvaluationStack("build", p.Rules)
	out := []string{"#!/bin/bash", ""}

	for idx, rule := range es.stack {
		out = append(out, fmt.Sprintf("# Rule %d -- %s", idx, rule.Name))
		for _, task := range rule.Tasks {
			cmd := NewScriptable(task)
			out = append(out, cmd.GetScript())
		}
		out = append(out, "")
	}
	fmt.Println(strings.Join(out[:], "\n"))
}

// Kinda execute
// func main() {
// 	p := NewPreprocessor("./examples/first.yml")

// 	fmt.Println("Macros")
// 	fmt.Println("--------------------")
// 	for name, value := range p.Macros {
// 		fmt.Printf("\t- ${{%s}}: [%v]\n", name, value)
// 	}

// 	fmt.Println("Rules")
// 	fmt.Println("--------------------")
// 	for name, rule := range p.Rules {
// 		fmt.Printf("\t- %s: %v\n", name, rule)
// 	}

// 	stack := NewEvaluationStack("root", p.Rules)
// 	fmt.Println("Stack")
// 	fmt.Println("--------------------")
// 	for _, rl := range stack.stack {
// 		fmt.Printf("\t- %s\n", rl.Name)
// 		for i, task := range rl.Tasks {
// 			cmd := NewExecutable(task)
// 			out, err := cmd.Execute()
// 			if err != nil {
// 				panic(err)
// 			}
// 			fmt.Printf("\t\t%d) %v\n", i+1, out)
// 		}
// 	}

// 	// @TODO execute or compile stack
// }

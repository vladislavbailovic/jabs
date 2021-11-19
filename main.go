package main

import (
	"fmt"
)

func main() {
	p := NewPreprocessor("./examples/first.yml")

	fmt.Println("Macros")
	fmt.Println("--------------------")
	for name, value := range p.Macros {
		fmt.Printf("\t- ${{%s}}: [%v]\n", name, value)
	}

	fmt.Println("Rules")
	fmt.Println("--------------------")
	for name, rule := range p.Rules {
		fmt.Printf("\t- %s: %v\n", name, rule)
	}

	stack := NewEvaluationStack("root", p.Rules)
	fmt.Println("Stack")
	fmt.Println("--------------------")
	for _, rl := range stack.stack {
		fmt.Printf("\t- %s\n", rl.Name)
		for i, task := range rl.Tasks {
			cmd := NewExecutable(task)
			out, err := cmd.Execute()
			if err != nil {
				panic(err)
			}
			fmt.Printf("\t\t%d) %v\n", i+1, out)
		}
	}

	// @TODO execute or compile stack
}

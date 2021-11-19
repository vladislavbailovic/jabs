package main

import (
	"fmt"
)

func main() {
	p := NewPreprocessor("./examples/first.yml")
	fmt.Println(p.Rules)

	// @TODO prepare proper stack
	// @TODO execute or compile stack
}

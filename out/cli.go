package out

import (
	"fmt"
	"os"
)

type Output struct {
	Stdout *os.File
	Stderr *os.File
}

func (out Output) Out(msg string, args ...interface{}) {
	fmt.Println(msg)
}

func (out Output) Err(msg string, args ...interface{}) {
	fmt.Fprintf(out.Stderr, fmt.Sprintf(msg, args...))
}

var Cli Output = Output{
	Stdout: os.Stdout,
	Stderr: os.Stderr,
}

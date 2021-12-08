package sys

import "fmt"

func Shell() string {
	return "/bin/bash"
}

func Shebang() string {
	return fmt.Sprintf("#!%s", Shell())
}

func ShellCommandParam() string {
	return "-c"
}

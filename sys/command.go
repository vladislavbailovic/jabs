package sys

import (
	"fmt"
	"io/ioutil"
	"jabs/dbg"
	"jabs/types"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	source string
}

type Cmdlet struct {
	Command
}

type Scriptlet struct {
	Command
}

func NewCommand(cmd string) types.Instruction {
	command := Command{cmd}
	if len(strings.Split(cmd, "\n")) > 1 {
		return &Scriptlet{command}
	}
	return &Cmdlet{command}
}

func (c *Cmdlet) Execute() (string, error) {
	args := []string{ShellCommandParam(), c.source}
	out, err := exec.Command(Shell(), args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (c Cmdlet) GetScript() string {
	return c.source
}

func (s *Scriptlet) Execute() (string, error) {
	file := "./tmp"
	err := ioutil.WriteFile(file, []byte(s.source), 0744)
	if err != nil {
		dbg.FatalError("could not write file %s: %v", file, err)
	}
	out, err := exec.Command(file).Output()
	os.Remove(file)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s Scriptlet) GetScript() string {
	file := TempFilePath()
	ret := []string{
		fmt.Sprintf("cat <<'EOF' > %s", file),
	}
	ret = append(ret, strings.Split(s.source, "\n")...)
	ret = append(ret, "EOF")
	ret = append(ret, fmt.Sprintf("chmod u+x %s && %s", file, file))
	ret = append(ret, fmt.Sprintf("rm %s", file))
	return strings.Join(ret, "\n")
}

package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Executable interface {
	Execute() (string, error)
}

type Command struct {
	source string
}

type Cmdlet struct {
	Command
}

type Scriptlet struct {
	Command
}

func NewExecutable(cmd string) Executable {
	command := Command{cmd}
	if len(strings.Split(cmd, "\n")) > 1 {
		return &Scriptlet{command}
	}
	return &Cmdlet{command}
}

func (c *Cmdlet) Execute() (string, error) {
	args := []string{"-c", c.source}
	out, err := exec.Command("/bin/sh", args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *Scriptlet) Execute() (string, error) {
	file := "/data/Projects/geek/jabs/tmp"
	err := ioutil.WriteFile(file, []byte(s.source), 0744)
	if err != nil {
		panic("could not write file")
	}
	out, err := exec.Command(file).Output()
	os.Remove(file)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

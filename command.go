package main

import (
	"github.com/gonuts/go-shlex"
	"os/exec"
)

type runner interface {
	Run(*exec.Cmd) error
}

type prepareFunc func(*exec.Cmd)

func (f prepareFunc) Run(cmd *exec.Cmd) error {
	f(cmd)
	return cmd.Run()
}

func runCommand(text string, r runner) error {
	command := Render("", text)
	args, err := shlex.Split(command)
	if err != nil {
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	return r.Run(cmd)
}

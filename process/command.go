package process

import (
	"os/exec"

	"code.google.com/p/go-shlex"
	"github.com/yosisa/pave/template"
)

type Command struct {
	Template    string
	PrepareFunc func(*exec.Cmd)
	Cmd         *exec.Cmd
}

func NewCommand(cmd string) *Command {
	return &Command{Template: cmd}
}

func (c *Command) Start() error {
	command := template.Render("", c.Template)
	args, err := shlex.Split(command)
	if err != nil {
		return err
	}

	c.Cmd = exec.Command(args[0], args[1:]...)
	if c.PrepareFunc != nil {
		c.PrepareFunc(c.Cmd)
	}

	return c.Cmd.Start()
}

func (c *Command) Wait() error {
	return c.Cmd.Wait()
}

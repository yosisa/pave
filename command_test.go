package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestRunCommand(t *testing.T) {
	w := new(bytes.Buffer)
	err := runCommand(`echo -n {{env "USER"}}`, prepareFunc(func(cmd *exec.Cmd) {
		cmd.Stdout = w
	}))
	assert.Nil(t, err)
	assert.Equal(t, os.Getenv("USER"), w.String())
}

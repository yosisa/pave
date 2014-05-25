package main

import (
	"bytes"
	"fmt"
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

func TestRunComplexCommand(t *testing.T) {
	w := new(bytes.Buffer)
	cmd := `echo -n {{if env "USER"}}--user={{env "USER"}}{{end}}
                {{if env "DO_NOT_MATCH"}}--group={{env "USER"}}{{end}}
                --listen={{ipv4 "127.0"}}:{{env "PORT" "80"}}`
	expected := fmt.Sprintf("--user=%s --listen=127.0.0.1:80", os.Getenv("USER"))

	err := runCommand(cmd, prepareFunc(func(cmd *exec.Cmd) {
		cmd.Stdout = w
	}))
	assert.Nil(t, err)
	assert.Equal(t, expected, w.String())
}

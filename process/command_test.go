package process

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCommand(t *testing.T) {
	w := new(bytes.Buffer)
	cmd := NewCommand(`echo -n {{env "USER"}}`)
	cmd.PrepareFunc = func(cmd *exec.Cmd) {
		cmd.Stdout = w
	}

	err := cmd.Start()
	assert.Nil(t, err)
	err = cmd.Wait()
	assert.Nil(t, err)
	assert.Equal(t, os.Getenv("USER"), w.String())
}

func TestRunComplexCommand(t *testing.T) {
	w := new(bytes.Buffer)
	cmd := NewCommand(`echo -n {{if env "USER"}}--user={{env "USER"}}{{end}}
                           {{if env "DO_NOT_MATCH"}}--group={{env "USER"}}{{end}}
                           --listen={{ipv4 "127.0"}}:{{env "PORT" "80"}}`)
	cmd.PrepareFunc = func(cmd *exec.Cmd) {
		cmd.Stdout = w
	}

	err := cmd.Start()
	assert.Nil(t, err)
	err = cmd.Wait()
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("--user=%s --listen=127.0.0.1:80", os.Getenv("USER")), w.String())
}

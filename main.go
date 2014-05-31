package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
	"os/exec"
	"sync"
)

var opts struct {
	Files    []string `short:"f" long:"file" description:"Files to be rendered"`
	Commands []string `short:"c" long:"command" description:"Commands to be executed"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	for _, f := range opts.Files {
		if err := NewTemplate(f).Execute(); err != nil {
			fmt.Println(err)
		}
	}

	var wg sync.WaitGroup
	for _, command := range opts.Commands {
		wg.Add(1)
		go func(command string) {
			runCommand(command, prepareFunc(func(cmd *exec.Cmd) {
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
			}))
			wg.Done()
		}(command)
	}

	wg.Wait()
}

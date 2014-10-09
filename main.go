package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/yosisa/pave/process"
	"github.com/yosisa/pave/template"
)

var opts struct {
	Files       []string                `short:"f" long:"file" description:"Files to be rendered"`
	Commands    []string                `short:"c" long:"command" description:"Commands to be executed"`
	Strategy    process.RestartStrategy `short:"r" long:"restart" description:"Restart strategy (none|always|error)"`
	RestartWait time.Duration           `short:"w" long:"restart-wait" description:"Duration for restarting"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	for _, f := range opts.Files {
		if err := template.NewTemplate(f).Execute(); err != nil {
			fmt.Println(err)
		}
	}

	if len(opts.Commands) > 0 {
		pm := process.NewProcessManager(opts.Strategy, opts.RestartWait)
		for _, command := range opts.Commands {
			pm.Add(process.New(command, func(cmd *exec.Cmd) {
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
			}, nil))
		}
		pm.Run()
	}
}

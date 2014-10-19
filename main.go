package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
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

	if len(opts.Commands) == 0 {
		return
	}

	pm := process.NewProcessManager(opts.Strategy, opts.RestartWait)
	for _, command := range opts.Commands {
		pm.Add(process.New(command, func(cmd *exec.Cmd) {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}, nil))
	}
	go signalHandler(pm)
	go childCollector()
	pm.Run()
}

func signalHandler(pm *process.ProcessManager) {
	sigC := make(chan os.Signal)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for sig := range sigC {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			pm.Stop()
			signal.Stop(sigC)
			close(sigC)
		default:
			pm.SignalAll(sig)
		}
	}
}

func childCollector() {
	var ws syscall.WaitStatus
	tick := time.Tick(10 * time.Second)
	for _ = range tick {
		var err error
		for err == nil {
			_, err = syscall.Wait4(-1, &ws, 0, nil)
		}
	}
}

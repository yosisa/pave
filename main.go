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

var Version string

var opts struct {
	Files       []string                `short:"f" long:"file" description:"Files to be rendered"`
	Commands    []string                `short:"c" long:"command" description:"Commands to be executed"`
	Strategy    process.RestartStrategy `short:"r" long:"restart" description:"Restart strategy (none|always|error)"`
	RestartWait time.Duration           `short:"w" long:"restart-wait" description:"Duration for restarting"`
	Forever     bool                    `long:"forever" description:"Run forever even if no child process exists"`
	Version     bool                    `long:"version" description:"Show version"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}
	if opts.Version {
		fmt.Println("pave", Version)
		return
	}

	for _, f := range opts.Files {
		if err := template.NewTemplate(f).Execute(); err != nil {
			fmt.Println(err)
		}
	}

	pm := process.NewProcessManager(opts.Strategy, opts.RestartWait)
	for _, command := range opts.Commands {
		pm.Add(process.New(command, func(cmd *exec.Cmd) {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}, nil))
	}

	c := make(chan bool)
	go signalHandler(pm, c)
	go childCollector()
	pm.Run()
	if opts.Forever {
		<-c
	}
}

func signalHandler(pm *process.ProcessManager, c chan bool) {
	sigC := make(chan os.Signal)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	for sig := range sigC {
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			pm.Stop()
			signal.Stop(sigC)
			close(sigC)
			close(c)
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

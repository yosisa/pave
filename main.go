package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
	"sync"
)

func main() {
	app := cli.NewApp()
	app.Name = "pave"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{"file, f", &cli.StringSlice{}, "description"},
		cli.StringSliceFlag{"command, c", &cli.StringSlice{}, "description"},
	}
	app.Action = realMain

	app.Run(os.Args)
}

func realMain(c *cli.Context) {
	for _, f := range c.StringSlice("file") {
		if err := NewTemplate(f).Execute(); err != nil {
			fmt.Println(err)
		}
	}

	var wg sync.WaitGroup
	for _, command := range c.StringSlice("command") {
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

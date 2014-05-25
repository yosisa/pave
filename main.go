package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"os/exec"
)

func main() {
	app := cli.NewApp()
	app.Name = "pave"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{"file, f", &cli.StringSlice{}, "description"},
		cli.StringFlag{"command, c", "", "description"},
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

	if command := c.String("command"); command != "" {
		runCommand(command, prepareFunc(func(cmd *exec.Cmd) {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}))
	}
}

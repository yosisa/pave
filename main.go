package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "pave"
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{"file, f", &cli.StringSlice{}, "description"},
	}
	app.Action = realMain

	app.Run(os.Args)
}

func realMain(c *cli.Context) {
	for _, f := range c.StringSlice("file") {
		println(NewTemplate(f).Execute())
	}
}

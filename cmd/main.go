package main

import (
	"gopkg.in/urfave/cli.v1"
	"os"
)

var (
	appName = "Bitstamp Robot"
	version = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = version
	app.Copyright = "(c) 2018 - Ricardo Pinto"
	app.Usage = "Bitstamp Robot"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "robot, r",
			Usage: "Runs in robot mode (as in a crontab)",
		},
		cli.BoolFlag{
			Name:  "balance, b",
			Usage: `Checks account balance`,
		},
		cli.BoolFlag{
			Name:  "check, c",
			Usage: "Checks the orders in place",
		},
	}

	app.Action = Serve

	app.Run(os.Args)
}

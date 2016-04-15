package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "swift-provisionning"
	app.Usage = "fight the loneliness!"
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}

	app.Commands = []cli.Command{
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Report status Of Disks",
			Action: Status,
		},
		{
			Name:    "complete",
			Aliases: []string{"c"},
			Usage:   "complete a task on the list",
			Action: func(c *cli.Context) {
				println("completed task: ", c.Args().First())
			},
		},
		{
			Name:    "template",
			Aliases: []string{"r"},
			Usage:   "options for task templates",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new template",
					Action: func(c *cli.Context) {
						println("new task template: ", c.Args().First())
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing template",
					Action: func(c *cli.Context) {
						println("removed task template: ", c.Args().First())
					},
				},
			},
		},
	}

	app.Run(os.Args)
}

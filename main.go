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
			Action:  Status,
		},
		{
			Name:    "rename",
			Aliases: []string{"r"},
			Usage:   "Rename Vdisks conforming to cloudwatt spec",
			Action:  RenameVdisks,
			Flags: []cli.Flag{
				cli.BoolTFlag{
					Name:  "dry, d",
					Usage: "Add --dry=false",
				},
			},
		},
	}

	app.Run(os.Args)
}

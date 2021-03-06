package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "swift-provisionning"
	app.Usage = "fight the loneliness!"
	app.Action = func(c *cli.Context) {
		println("Hello friend!")
	}

	app.EnableBashCompletion = true

	app.Commands = []cli.Command{
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Report status Of Disks",
			Action:  Status,
		},
		{
			Name:    "rename-vdisk",
			Aliases: []string{"r"},
			Usage:   "Rename Vdisks conforming to cloudwatt spec",
			Action:  RenameVdisks,
			Flags: []cli.Flag{
				cli.BoolTFlag{
					Name:  "dry, d",
					Usage: "Add --dry=false to actually rename stuff",
				},
				cli.BoolFlag{
					Name:  "yolo",
					Usage: "DANGER : Add --yolo to force renaming Vdisk with already valid pattern",
				},
				cli.IntFlag{
					Name:  "maxhdd",
					Usage: "Maximum number of data HDD to use",
					Value: -1,
				},
				cli.IntFlag{
					Name:  "maxssd",
					Usage: "Maximum number of data SSD to use",
					Value: -1,
				},
			},
		},
		{
			Name:   "add-labels",
			Usage:  "Add xfs labels to fs matching Vdisk Name",
			Action: AddLabels,
			Flags: []cli.Flag{
				cli.BoolTFlag{
					Name:  "dry, d",
					Usage: "Add --dry=false to actually rename stuff",
				},
			},
		},
	}

	app.Run(os.Args)
}

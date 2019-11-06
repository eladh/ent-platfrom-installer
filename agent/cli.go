package main

import (
	"agent/installers"
	"github.com/urfave/cli"
	"log"
	"os"
	"sort"
)

func main() {
	app := cli.NewApp()
	app.Name = "Platform Installer"
	app.Usage = "make infra as code E+ automation"
	app.Version = "1.0.0"

	var setupInfoLocation string

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "config",
			Usage: "Load setupInfo configuration from yaml `FILE`",
			Destination: &setupInfoLocation,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "Install platform",
			Action:  func(c *cli.Context) error {
				if setupInfoLocation == "" {
					return cli.NewExitError("you must supply setupInfo yaml file by using --config=`FILE`", 86)
				}
				installers.InstallPlatform(setupInfoLocation)
				return nil
			},
		},
		{
			Name:    "uninstall",
			Aliases: []string{"u"},
			Usage:   "Uninstall platform",
			Action:  func(c *cli.Context) error {
				if setupInfoLocation == "" {
					return cli.NewExitError("you must supply setupInfo yaml file by using --config=`FILE`", 86)
				}
				installers.UninstallPlatform(setupInfoLocation)
				return nil
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
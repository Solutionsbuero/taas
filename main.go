package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/solutionsbuero/ttrn/src"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ttrn",
		Usage: "remote trains-as-a-servie server part",
		Authors: []*cli.Author{
			{
				Name:  "72nd",
				Email: "msg@frg72.com",
			},
		},
		Action: func(c *cli.Context) error {
			path := getArgument(c)
			cfg := ttrn.OpenConfig(path)
			ttrn.Run(cfg)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "new",
				Usage: "creates a new config file",
				Action: func(c *cli.Context) error {
					path := getArgument(c)
					cfg := ttrn.DefaultConfig()
					cfg.SaveConfig(path)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("cli error, %s", err)
	}
}

// getArgument tries to get the first positional argument given by the user. Or fatals with a error
// message.
func getArgument(c *cli.Context) string {
	if c.Args().Len() != 1 {
		logrus.Fatalf("one positional argument needed (path to config file)")
	}
	return c.Args().Get(0)
}

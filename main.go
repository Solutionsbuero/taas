package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "ttrn",
		Usage: "remote trains-as-a-servie server part",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "72nd",
				Email: "msg@frg72.com",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalf("cli error, %s", err)
	}
}

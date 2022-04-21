package main

import (
	"ROUKIEN/rundeck-activity-monitor/cmd"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

func main() {
	run(os.Args)
}

func run(args []string) {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.yml",
				Usage:   "path to the configuration file to use",
			},
		},
		Commands: []*cli.Command{
			cmd.NewScrapeCmd(),
			cmd.NewDatabaseCmd(),
			cmd.NewServeCmd(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

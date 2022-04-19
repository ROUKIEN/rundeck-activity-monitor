package main

import (
	"ROUKIEN/rundeck-activity-monitor/cmd"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
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

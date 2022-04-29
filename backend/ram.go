package main

import (
	"ROUKIEN/rundeck-activity-monitor/cmd"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/urfave/cli/v2"
)

func main() {
	run(os.Args)
}

func run(args []string) {
	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	logLevel, err := logrus.ParseLevel(os.Getenv("RAM_LOG_LEVEL"))
	if err == nil {
		logrus.SetLevel(logLevel)
	}

	logrus.SetFormatter(formatter)
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

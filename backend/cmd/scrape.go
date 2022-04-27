package cmd

import (
	"ROUKIEN/rundeck-activity-monitor/config"
	"ROUKIEN/rundeck-activity-monitor/database"
	"ROUKIEN/rundeck-activity-monitor/rundeck"
	"ROUKIEN/rundeck-activity-monitor/rundeck/spec"
	"bufio"
	"database/sql"
	"math/rand"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func NewScrapeCmd() *cli.Command {
	return &cli.Command{
		Name:   "scrape",
		Usage:  "scrape rundeck executions that ended in the defined timeframe across all defined instances",
		Action: scrapeExecute,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "begin",
				Value: "2022-01-01T00:00:00Z",
				Usage: "begin date to scrape jobs",
			},
			&cli.StringFlag{
				Name:  "end",
				Value: "2022-01-01T00:00:00Z",
				Usage: "end date to scrape jobs",
			},
			&cli.StringFlag{
				Name:    "newer-than",
				Aliases: []string{"nt"},
				Usage:   "Scrape job executions whose completion date is newer than NEWER-THAN",
			},
		},
	}
}

func scrapeExecute(c *cli.Context) error {
	f, err := os.Open(c.String("config"))
	if err != nil {
		return err
	}
	conf, err := config.Parse(bufio.NewReader(f))
	if err != nil {
		return err
	}

	db, err := database.Db()
	if err != nil {
		return err
	}

	opts := map[string]string{
		"begin":      c.String("begin"),
		"end":        c.String("end"),
		"newer-than": c.String("newer-than"),
	}
	scrapeOptions, err := spec.NewScrapeOptions(opts)
	if err != nil {
		return err
	}

	log.Infof("there are %d instances to scrape", len(conf.Instances))

	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	instanceExecutionsChannel := make(chan *config.ScrapedExecution)

	for instance_label, instance := range conf.Instances {
		wg.Add(1)
		log.Infof("Scraping %s", instance_label)

		go func(i config.RundeckInstance, il string, so *spec.ScrapeOptions) {
			defer wg.Done()
			err := scrapeInstanceExecutions(i, il, so, instanceExecutionsChannel)
			if err != nil {
				log.Error(err)
			}
		}(instance, instance_label, scrapeOptions)
	}

	go func() {
		wg.Wait()
		log.Info("Done.")
		close(instanceExecutionsChannel)
	}()

	for execution := range instanceExecutionsChannel {
		if err := handleExecutionRecording(db, execution.Instance, execution.Execution); err != nil {
			log.Errorf("[%s]failed to save execution #%d: %s", execution.Instance, execution.Execution.ID, err.Error())
		}
	}

	return nil
}

func scrapeInstanceExecutions(instance config.RundeckInstance, instanceLabel string, so *spec.ScrapeOptions, ch chan<- *config.ScrapedExecution) error {
	client := rundeck.NewRundeckClient(instance.Url, instance.Token, instance.ApiVersion, time.Duration(instance.Timeout)*time.Millisecond)
	projects, err := client.ListProjects()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, project := range projects {
		wg.Add(1)
		go func(p spec.Project) {

			defer wg.Done()
			execCh := client.ListProjectExecutions(p.Name, so)
			i := 0
			for execution := range execCh {
				i = i + 1
				se := config.ScrapedExecution{
					Execution: &execution,
					Instance:  instanceLabel,
				}
				ch <- &se
			}

			log.WithFields(log.Fields{
				"instance": instanceLabel,
				"project":  p.Name,
			}).Infof("scraped %d executions", i)
		}(project)
	}

	wg.Wait()

	log.Infof("Scraping of %s is over.", instanceLabel)

	return nil
}

func handleExecutionRecording(db *sql.DB, instance_name string, e *spec.Execution) error {
	executionInDB, err := database.FindExecution(db, instance_name, e)
	if err != nil {
		return err
	}

	if executionInDB == nil {
		log.WithFields(log.Fields{
			"instance": instance_name,
			"project":  e.Project,
			"job":      e.Job.ID,
		}).Debugf("Will save execution #%d", e.ID)
		return database.SaveExecution(db, instance_name, e)
	}

	return nil
}

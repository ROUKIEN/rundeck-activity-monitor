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
		Usage:  "scrape rundeck instances",
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
				Usage: "begin date to scrape jobs",
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

	layout := "2006-01-02T15:04:05.000Z"
	begin, err := time.Parse(layout, c.String("begin"))
	if err != nil {
		return err
	}
	end, err := time.Parse(layout, c.String("end"))
	if err != nil {
		return err
	}

	log.Infof("there are %d instances to scrape", len(conf.Instances))

	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	instanceExecutionsChannel := make(chan *ScrapedExecution)

	for instance_label, instance := range conf.Instances {
		wg.Add(1)
		log.Infof("Scraping %s", instance_label)

		go func(i config.RundeckInstance, il string, b time.Time, e time.Time) {
			defer wg.Done()
			executionsChan, err := scrapeInstanceExecutions(i, il, b, e)
			if err != nil {
				log.Error(err)
				// fmt.Printf("%s\n", err.Error())
			}

			for execution := range executionsChan {
				instanceExecutionsChannel <- execution
			}
		}(instance, instance_label, begin, end)
	}

	go func() {
		wg.Wait()
		close(instanceExecutionsChannel)
	}()

	for execution := range instanceExecutionsChannel {
		if err := handleExecutionRecording(db, execution.Instance, execution.Execution); err != nil {
			log.Errorf("[%s]failed to save execution #%d: %s", execution.Instance, execution.Execution.ID, err.Error())
		}
	}

	return nil
}

type ScrapedExecution struct {
	Execution *spec.Execution
	Instance  string
}

func scrapeInstanceExecutions(instance config.RundeckInstance, instanceLabel string, begin time.Time, end time.Time) (chan *ScrapedExecution, error) {
	client := rundeck.NewRundeckClient(instance.Url, instance.Token, instance.ApiVersion, time.Duration(instance.Timeout)*time.Millisecond)
	projects, err := client.ListProjects()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	currChan := make(chan *ScrapedExecution)
	for _, project := range projects {
		wg.Add(1)
		go func(client *rundeck.Rundeck, p spec.Project) {
			defer wg.Done()
			executions, _ := client.ListProjectExecutions(p.Name, begin, end)
			for _, execution := range executions {
				se := ScrapedExecution{
					Execution: &execution,
					Instance:  instanceLabel,
				}
				currChan <- &se
			}
			log.WithFields(log.Fields{
				"instance": instanceLabel,
				"project":  p.Name,
			}).Debugf("scraped %d executions", len(executions))
		}(client, project)
	}

	go func() {
		wg.Wait()
		close(currChan)

		log.Info("Scraping is over.")
	}()

	return currChan, nil
}

func handleExecutionRecording(db *sql.DB, instance_name string, e *spec.Execution) error {
	executionInDB, err := database.FindExecution(db, instance_name, e)
	if err != nil {
		return err
	}

	if executionInDB == nil {
		return database.SaveExecution(db, instance_name, e)
	}

	return nil
}

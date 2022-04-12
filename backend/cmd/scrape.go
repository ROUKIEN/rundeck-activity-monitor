package cmd

import (
	"ROUKIEN/rundeck-activity-monitor/config"
	"ROUKIEN/rundeck-activity-monitor/rundeck"
	"ROUKIEN/rundeck-activity-monitor/rundeck/spec"
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

func NewScrapeCmd() *cli.Command {
	return &cli.Command{
		Name:    "scrape",
		Aliases: []string{"s"},
		Usage:   "scrape rundeck instances",
		Action:  scrapeExecute,
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
	f, err := os.Open("config.yml")
	if err != nil {
		return err
	}
	conf, err := config.Parse(bufio.NewReader(f))
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

	fmt.Printf("there are %d instances to scrape\n", len(conf.Instances))

	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())
	instanceExecutionsChannel := make(chan *spec.Execution)

	for _, instance := range conf.Instances {
		wg.Add(1)

		go func(i config.RundeckInstance, b time.Time, e time.Time) {
			defer wg.Done()
			fmt.Printf("Scraping %s...\n", i.Url)
			executionsChan, err := scrapeInstanceExecutions(i, b, e)
			if err != nil {
				fmt.Printf("%s\n", err.Error())
			}

			for execution := range executionsChan {
				instanceExecutionsChannel <- execution
			}
		}(instance, begin, end)
	}

	go func() {
		wg.Wait()
		close(instanceExecutionsChannel)
	}()

	allExecutions := make([]spec.Execution, 0)
	for execution := range instanceExecutionsChannel {
		allExecutions = append(allExecutions, *execution)
	}

	fmt.Printf("%d executions over all instances\n", len(allExecutions))

	return nil
}

func scrapeInstanceExecutions(instance config.RundeckInstance, begin time.Time, end time.Time) (chan *spec.Execution, error) {
	client := rundeck.NewRundeckClient(instance.Url, instance.Token, instance.ApiVersion, time.Duration(instance.Timeout)*time.Millisecond)
	projects, err := client.ListProjects()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	currChan := make(chan *spec.Execution)
	for _, project := range projects {
		wg.Add(1)

		go func(client *rundeck.Rundeck, p spec.Project) {
			defer wg.Done()
			executions, _ := client.ListProjectExecutions(p.Name, begin, end)
			for _, execution := range executions {
				currChan <- &execution
			}
			fmt.Printf("[%s][%s] %d executions\n", client.Url, p.Name, len(executions))
		}(client, project)
	}

	go func() {
		wg.Wait()
		close(currChan)
	}()

	return currChan, nil
}

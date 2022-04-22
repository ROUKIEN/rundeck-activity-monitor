package config

import "ROUKIEN/rundeck-activity-monitor/rundeck/spec"

type ScrapedExecution struct {
	Execution *spec.Execution
	Instance  string
}

package spec

import "time"

type Execution struct {
	ID            int         `json:"id"`
	Href          string      `json:"href"`
	Permalink     string      `json:"permalink"`
	Status        string      `json:"status"`
	Project       string      `json:"project"`
	ExecutionType string      `json:"executionType"`
	User          string      `json:"user"`
	DateStarted   RundeckDate `json:"date-started"`
	DateEnded     RundeckDate `json:"date-ended"`
	Job           Job         `json:"job"`

	Description string `json:"description"`
}

type RundeckDate struct {
	UnixTime int    `json:"unixtime"`
	Date     string `json:"date"`
}

func NewRundeckDate(datefrom time.Time) RundeckDate {
	return RundeckDate{
		UnixTime: int(datefrom.UTC().UnixMicro()),
		Date:     datefrom.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

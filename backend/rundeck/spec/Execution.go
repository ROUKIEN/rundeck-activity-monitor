package spec

import (
	"errors"
	"time"
)

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

type ScrapeOptions struct {
	Begin     *time.Time
	End       *time.Time
	NewerThan *time.Time
}

func NewScrapeOptions(opts map[string]string) (*ScrapeOptions, error) {
	newerThan := opts["newer-than"]
	if newerThan != "" {
		duration, err := time.ParseDuration(newerThan)
		if err != nil {
			return nil, err
		}
		newTime := time.Now().Add(-duration)
		return &ScrapeOptions{
			NewerThan: &newTime,
		}, nil
	}

	layout := "2006-01-02T15:04:05.000Z"

	var htr ScrapeOptions
	begin := opts["begin"]
	if begin != "" {
		begin, err := time.Parse(layout, begin)
		if err == nil {
			htr.Begin = &begin
		}
	}
	end := opts["end"]
	if end != "" {
		end, err := time.Parse(layout, end)
		if err == nil {
			htr.End = &end
		}
	}

	if end == "" && begin == "" {
		return nil, errors.New("at least one scraping argument must be specified")
	}

	return &htr, nil
}

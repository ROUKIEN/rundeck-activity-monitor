package spec

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewScrapeOptions(t *testing.T) {

	opts := map[string]string{
		"begin":      "2022-04-23T00:00:00.000Z",
		"end":        "2022-05-23T00:00:00.000Z",
		"newer-than": "",
	}

	so, err := NewScrapeOptions(opts)
	assert.Nil(t, err)
	assert.Equal(t, time.Date(2022, time.Month(4), 23, 0, 0, 0, 0, time.UTC), *so.Begin)
	assert.Equal(t, time.Date(2022, time.Month(5), 23, 0, 0, 0, 0, time.UTC), *so.End)
}

func TestNewScrapeOptionsNewerThan(t *testing.T) {

	opts := map[string]string{
		"begin":      "2022-04-23T00:00:00.000Z",
		"end":        "2022-05-23T00:00:00.000Z",
		"newer-than": "24h",
	}

	expectedDate := time.Now().Add(-24 * time.Hour)

	so, err := NewScrapeOptions(opts)
	assert.Nil(t, err)
	assert.Nil(t, so.Begin)
	assert.Nil(t, so.End)
	assert.Equal(t, expectedDate.UTC().Hour(), so.NewerThan.UTC().Hour())
}

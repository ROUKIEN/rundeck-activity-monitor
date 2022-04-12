package config

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleConfig = `---

instances:
  dev:
    url: http://rundeck.dev
    token: my.T0k3n
    apiversion: 38
    timeout: 5000
  stg:
    url: http://rundeck.stg
    token: my.T0k3333333n
    apiversion: 38
    timeout: 5000
`

func TestParse(t *testing.T) {
	config, err := Parse(strings.NewReader(exampleConfig))
	assert.Nil(t, err)
	assert.Len(t, config.Instances, 2)
}

package rundeck

import (
	"ROUKIEN/rundeck-activity-monitor/rundeck/spec"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListProjects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.String(), "/api/38/projects")

		projects := []spec.Project{
			{
				Name:        "project1",
				Description: "the first project",
				Url:         "http://rundeck.dev/api/38/project/project1",
				Label:       "",
				Created:     time.Now(),
			},
			{
				Name:        "project2",
				Description: "the second project",
				Url:         "http://rundeck.dev/api/38/project/project2",
				Label:       "",
				Created:     time.Now(),
			},
		}
		projects_json, _ := json.Marshal(projects)
		str := string(projects_json)

		rw.Write([]byte(str))
	}))
	defer server.Close()

	rd := NewRundeckClient(server.URL, "myAP1.T0k3n", 38, 5*time.Second)
	projects, err := rd.ListProjects()
	assert.Nil(t, err)
	assert.Len(t, projects, 2)
}

func TestListProjectExecutions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, req.URL.Path, "/api/38/project/my_project/executions")
		query := req.URL.Query()

		executions := generateDummyExecutions(64)
		var offset int
		var err error
		if query.Has("offset") {
			offset, err = strconv.Atoi(query.Get("offset"))
			assert.Nil(t, err)
		} else {
			offset = 0
		}
		pagedResponse := newPagedResponse(&executions, offset, 20)

		projects_json, _ := json.Marshal(pagedResponse)
		str := string(projects_json)

		rw.Write([]byte(str))
	}))
	defer server.Close()

	rd := NewRundeckClient(server.URL, "myAP1.T0k3n", 38, 0*time.Second)
	executions, err := rd.ListProjectExecutions("my_project", time.Date(2022, 01, 01, 12, 0, 0, 0, time.UTC), time.Date(2022, 01, 01, 13, 0, 0, 0, time.UTC))
	assert.Nil(t, err)
	assert.Len(t, executions, 64)
}

func dummyExecution(i int) spec.Execution {
	start := time.Date(2022, 01, 01, 12, 1, 1, 0, time.UTC)
	end := time.Date(2022, 01, 01, 12, 1, 2, 0, time.UTC)
	return spec.Execution{
		ID:            i,
		Href:          fmt.Sprintf("http://rundeck.dev/api/38/execution/%d", i),
		Permalink:     fmt.Sprintf("http://rundeck.dev/project/my_project/execution/show/%d", i),
		Status:        "succeeded",
		Project:       "my_project",
		ExecutionType: "scheduled",
		User:          "admin",
		DateStarted:   spec.NewRundeckDate(start),
		DateEnded:     spec.NewRundeckDate(end),
		Job: spec.Job{
			ID:              fmt.Sprintf("1462746-47462642-%d", i),
			AverageDuration: 1234,
			Name:            "My job",
			Group:           "",
			Project:         "my_project",
			Description:     "my awesome job",
			Href:            fmt.Sprintf("http://rundeck.dev/api/38/job/1462746-47462642-%d", i),
			Permalink:       fmt.Sprintf("http://rundeck.dev/project/my_project/job/show/1462746-47462642-%d", i),
		},
		Description: "foo",
	}
}

func generateDummyExecutions(number int) []spec.Execution {
	executions := make([]spec.Execution, 0)
	for i := 0; i < number; i++ {
		executions = append(executions, dummyExecution(i))
	}

	return executions
}

func newPagedResponse(ex *[]spec.Execution, offset int, max int) executionsResponse {
	var last int
	if offset+max < len(*ex) {
		last = offset + max
	} else {
		last = len(*ex)
	}
	sliced := (*ex)[offset:last]

	return executionsResponse{
		Paging: rundeckPaging{
			Count:  len(sliced),
			Total:  len(*ex),
			Offset: offset,
			Max:    max,
		},
		Executions: sliced,
	}
}

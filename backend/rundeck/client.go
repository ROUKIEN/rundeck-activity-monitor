package rundeck

import (
	"ROUKIEN/rundeck-activity-monitor/rundeck/spec"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Rundeck struct {
	Client     *http.Client
	Url        string
	ApiVersion int
}

type RundeckHeaderTransport struct {
	T     http.RoundTripper
	Token string
}

func (rht *RundeckHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Rundeck-Auth-Token", rht.Token)
	return rht.T.RoundTrip(req)
}

func newRundeckHeaderTransport(token string) *RundeckHeaderTransport {
	return &RundeckHeaderTransport{
		T:     http.DefaultTransport,
		Token: token,
	}
}

func NewRundeckClient(url string, token string, apiversion int, timeout time.Duration) *Rundeck {

	client := &http.Client{
		Timeout:   timeout,
		Transport: newRundeckHeaderTransport(token),
	}

	rd := &Rundeck{
		Client:     client,
		ApiVersion: apiversion,
		Url:        url,
	}

	return rd
}

func (rd *Rundeck) ListProjects() ([]spec.Project, error) {
	resp, err := rd.Client.Get(fmt.Sprintf("%s/api/%d/projects", rd.Url, rd.ApiVersion))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	projectsResp := []spec.Project{}

	if err := json.NewDecoder(resp.Body).Decode(&projectsResp); err != nil {
		return nil, err
	}

	return projectsResp, nil
}

type rundeckPaging struct {
	Count  int `json:"count"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
	Max    int `json:"max"`
}

type executionsResponse struct {
	Paging     rundeckPaging    `json:"paging"`
	Executions []spec.Execution `json:"executions"`
}

func (rd *Rundeck) ListProjectExecutions(project string, begin time.Time, end time.Time) ([]spec.Execution, error) {
	max := 20
	offset := 0

	executions := make([]spec.Execution, 0)
	for {
		url := fmt.Sprintf("%s/api/%d/project/%s/executions?begin=%s&end=%s&max=%d&offset=%d", rd.Url, rd.ApiVersion, project, begin.UTC().Format("2006-01-02T15:04:05Z"), end.UTC().Format("2006-01-02T15:04:05Z"), max, offset)
		resp, err := rd.Client.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		executionsResp := executionsResponse{}

		if err := json.NewDecoder(resp.Body).Decode(&executionsResp); err != nil {
			return nil, err
		}

		executions = append(executions, executionsResp.Executions...)

		if executionsResp.Paging.Count+executionsResp.Paging.Offset == executionsResp.Paging.Total {
			break
		}

		offset += executionsResp.Paging.Max
	}

	return executions, nil
}

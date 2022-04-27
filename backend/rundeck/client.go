package rundeck

import (
	"ROUKIEN/rundeck-activity-monitor/rundeck/spec"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
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
	req.Header.Add("User-Agent", "Rundeck-Activity-Monitor/1.0")
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

func (rd *Rundeck) ListProjectExecutions(project string, so *spec.ScrapeOptions) <-chan spec.Execution {
	// fmt.Printf("%s\n", so)
	ch := make(chan spec.Execution)
	max := 20
	offset := 0
	go func() {
		for {
			base, err := url.Parse(fmt.Sprintf("%s/api/%d/project/%s/executions", rd.Url, rd.ApiVersion, project))
			if err != nil {
				log.Error(err)
				break
			}
			params := url.Values{}
			params.Add("max", strconv.Itoa(max))
			params.Add("offset", strconv.Itoa(offset))

			if so.NewerThan != nil {
				since := so.NewerThan.UTC().UnixMilli()
				params.Add("begin", strconv.FormatInt(since, 10))
			} else {
				begin := so.Begin.UTC().UnixMilli()
				params.Add("begin", strconv.FormatInt(begin, 10))
				end := so.End.UTC().UnixMilli()
				params.Add("end", strconv.FormatInt(end, 10))
			}

			base.RawQuery = params.Encode()

			log.Trace(base.String())

			resp, err := rd.Client.Get(base.String())
			if err != nil {
				log.Error(err)
				break
			}
			defer resp.Body.Close()

			executionsResp := executionsResponse{}

			if err := json.NewDecoder(resp.Body).Decode(&executionsResp); err != nil {
				log.Error(err)
			}

			for _, execution := range executionsResp.Executions {
				ch <- execution
			}

			if executionsResp.Paging.Count+executionsResp.Paging.Offset == executionsResp.Paging.Total {
				break
			}

			offset += executionsResp.Paging.Max
		}

		close(ch)
	}()

	return ch
}

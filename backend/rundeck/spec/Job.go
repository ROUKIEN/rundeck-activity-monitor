package spec

type Job struct {
	ID              string `json:"id"`
	AverageDuration int    `json:"averageDuration"`
	Name            string `json:"name"`
	Group           string `json:"group"`
	Project         string `json:"project"`
	Description     string `json:"description"`
	Href            string `json:"href"`
	Permalink       string `json:"permalink"`
}

package spec

import "time"

type Project struct {
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	Label       string    `json:"label"`
	Created     time.Time `json:"created"`
}

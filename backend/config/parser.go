package config

import (
	"io"

	"gopkg.in/yaml.v3"
)

type RundeckInstance struct {
	Url        string `yaml:"url"`
	Token      string `yaml:"token"`
	ApiVersion int    `yaml:"apiversion"`
	Timeout    int    `yaml:"timeout"`
}

type RAMConfig struct {
	Instances map[string]RundeckInstance `yaml:"instances"`
}

func Parse(r io.Reader) (*RAMConfig, error) {
	config := RAMConfig{}
	buff, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(buff, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

func LoadConf(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := Config{}

	err = yaml.Unmarshal(file, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

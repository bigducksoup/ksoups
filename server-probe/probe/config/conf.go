package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Name        string `yaml:"name"`
	CenterAddr  string `yaml:"center-addr"`
	ScriptsPath string `yaml:"scripts-path"`
	DbPath      string `yaml:"db-path"`
}

var Conf Config = Config{}

func LoadConf(path string) error {

	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &Conf)
	if err != nil {
		return err
	}

	if Conf.ScriptsPath == "" {
		pwd, _ := os.Getwd()
		Conf.ScriptsPath = pwd
	}

	return nil

}

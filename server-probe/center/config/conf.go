package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Api    `yaml:"api"`
	Center `yaml:"center"`
}

type Api struct {
	Port     string `yaml:"port"`
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
}

type Center struct {
	Port string `yaml:"port"`
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

	return nil
}

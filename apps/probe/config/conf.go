package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Name       string `yaml:"name"`
	CenterAddr string `yaml:"center-addr"`
	ScriptPath string `yaml:"script-path"`
	PublicKey  string `yaml:"public-key"`
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

	if Conf.ScriptPath == "" {
		pwd, _ := os.Getwd()
		Conf.ScriptPath = pwd + "/scripts"
	}
	afterPropertySet()

	return nil

}

func afterPropertySet() {

	if Conf.CenterAddr == "" {
		panic("center-addr is not set")
	}

	if Conf.Name == "" {
		panic("name is not set")
	}

	if Conf.PublicKey == "" {
		panic("you should have a public key to identify yourself")
	}

	info, err := os.Stat(Conf.ScriptPath)

	if err != nil {
		err := os.MkdirAll(Conf.ScriptPath, 0o755)
		if err != nil {
			panic(err)
		}
		return
	}

	if !info.IsDir() {
		panic("script-path is not a directory")
	}
}

package config

type Config struct {
	WebApi `yaml:"web-api"`
	Center `yaml:"center"`
}

type WebApi struct {
	Port     string `yaml:"port"`
	Account  string `yaml:"account"`
	Password string `yaml:"password"`
}

type Center struct {
	Port int `yaml:"port"`
}

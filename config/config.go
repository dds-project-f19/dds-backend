package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Addr        string `yaml:"addr"`
	DSN         string `yaml:"dsn"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
}

var config *Config

func Load(path string) error {
	result, err := ioutil.ReadFile(path)
	if err != nil {
		// LOAD DEFAULT
		config = GetDefault()
	}

	return yaml.Unmarshal(result, &config)
}

func Get() *Config {
	return config
}

func GetDefault() *Config {
	return &Config{
		Addr:        ":9000",
		DSN:         "root:ddspassword14882@tcp(127.0.0.1:3306)/ddstest?charset=utf8&parseTime=True&loc=Local", // TODO: refactor
		MaxIdleConn: 100,
	}
}

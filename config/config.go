package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Addr        string `yaml:"addr"`
	DBLogin     string `yaml:"dblogin"`
	DBPassword  string `yaml:"dbpassword"`
	DBAddress   string `yaml:"dbaddress"`
	DBPort      string `yaml:"dbport"`
	DBName      string `yaml:"dbname"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
}

var config *Config

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.DBLogin, c.DBPassword, c.DBAddress, c.DBPort, c.DBName)
}

func Load(path string) error {
	result, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(result, &config)
}

func Get() *Config {
	return config
}

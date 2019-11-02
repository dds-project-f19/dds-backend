package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
)

type DBConfig struct {
	DBAddress  string `yaml:"dbaddress"`
	DBLogin    string `yaml:"dblogin"`
	DBPassword string `yaml:"dbpassword"`
	DBPort     string `yaml:"dbport"`
	DBName     string `yaml:"dbname"`
}

type GeneralConfig struct {
	Address     string `yaml:"addr"`
	MaxIdleConn int    `yaml:"max_idle_conn"`
}

var dbConfigArgNames = [...]string{"dbaddress", "dblogin", "dbpassword", "dbport", "dbname"}

func GetDefaultDBConfig() DBConfig {
	return DBConfig{
		DBAddress:  "127.0.0.1",
		DBLogin:    "root",
		DBPassword: "ddspassword14882",
		DBPort:     "3306",
		DBName:     "ddstest",
	}
}

func GetDefaultGeneralConfig() GeneralConfig {
	return GeneralConfig{
		Address:     ":9000",
		MaxIdleConn: 100,
	}
}

func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.DBLogin, c.DBPassword, c.DBAddress, c.DBPort, c.DBName)
}

// TODO: consider deletion
func LoadConfigFromFile(path string) (DBConfig, error) {
	result := DBConfig{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return result, err
	}
	err = yaml.Unmarshal(file, &result)
	return result, err
}

// Extract configuration files from command line arguments
func LoadConfigFromCmdArgs() DBConfig {
	configTemplate := GetDefaultDBConfig()
	conf := reflect.ValueOf(&configTemplate)
	for i := 0; i < conf.Elem().NumField(); i++ {
		flag.StringVar(
			conf.Elem().Field(i).Addr().Interface().(*string),
			dbConfigArgNames[i],
			conf.Elem().Field(i).Interface().(string),
			"")
	}
	flag.Parse()
	return configTemplate
}

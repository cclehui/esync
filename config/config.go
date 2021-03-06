package config

import (
	"io/ioutil"

	daoongorm "github.com/cclehui/dao-on-gorm"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var Conf *Config

type Config struct {
	Mysql *daoongorm.DBClientConfig `yaml:"mysql"`
	Redis *Redis                    `yaml:"redis"`
}

func InitConfigFromFile(filePath string) *Config {
	configData, err := (&Config{}).decodeFromFile(filePath)
	if err != nil {
		panic(err)
	}

	return configData
}

func (c *Config) decodeFromFile(filePath string) (*Config, error) {
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(fileData, c)
	if err != nil {
		return c, errors.WithStack(err)
	}

	return c, nil
}

type Redis struct {
	Server   string `yaml:"server"`   // "xxxxx:6379"
	Password string `yaml:"password"` // "wxxxxxxx"
}

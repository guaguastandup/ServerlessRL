package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
}

func FromYaml(yamlBytes []byte, objAddr interface{}) {
	err := yaml.Unmarshal(yamlBytes, objAddr)
	if err != nil {
		log.Panicf("Failed to unmarshal YAML %#v to Object: %v", yamlBytes, err)
	}
}

func NewConfig(configPath string) *Config {
	c := &Config{}
	yamlBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	FromYaml(yamlBytes, c)
	return c
}

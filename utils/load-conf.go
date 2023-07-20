package utils

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConf() (*Config, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	// Substitute environment variables
	yamlBytes := []byte(os.ExpandEnv(string(yamlFile)))

	var config Config
	err = yaml.Unmarshal(yamlBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB      DBConfig      `yaml:"db"`
	Handler HandlerConfig `yaml:"handler"`
	Models  []string      `yaml:"models"`
	Repos   []string      `yaml:"repositories"`
}

type DBConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type HandlerConfig struct {
	Path     string                   `yaml:"path"`
	Prefix   string                   `yaml:"prefix"`
	Versions map[string]VersionConfig `yaml:"versions"`
}

type VersionConfig struct {
	Prefix string                 `yaml:"prefix"`
	Routes map[string]RouteConfig `yaml:"routes"`
}

type RouteConfig struct {
	Method  string `yaml:"method"`
	Path    string `yaml:"path"`
	Handler string `yaml:"handler"`
}

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

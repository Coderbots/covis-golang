package helpers

import (
	"os"

	"gopkg.in/yaml.v2"
)

const configPath = "helpers/config.yml"

type Config struct {
	CovidRepo struct {
		Url string `yaml:"url"`
	} `yaml:"covidRepo"`
}

var AppConfig *Config

func ReadConfig() {
	fileHandler, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer fileHandler.Close()
	decoder := yaml.NewDecoder(fileHandler)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		panic(err)
	}
}

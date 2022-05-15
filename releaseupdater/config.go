package releaseupdater

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Update struct {
	Repo  string // name of the repo
	File  string // reference to file in repo
	Regex string // regex for searching the to-be-replaced item
}

type Action struct {
	From      string
	Component string
	Update    []Update
}

type Config struct {
	Actions []Action
}

func GetConfig(configFilename string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
		return nil, err
	}
	return getConfigFrom(yamlFile)
}

func getConfigFrom(content []byte) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal(content, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c, nil
}

package releaseupdater

import (
	"io/ioutil"
	"log"

	"github.com/ktrysmt/go-bitbucket"
	"github.com/mpapenbr/go-probot/probot"
	"gopkg.in/yaml.v3"
)

type Context struct {
	Config          *Config
	ProbotCtx       *probot.Context
	BitbucketClient *bitbucket.Client
}

type Update struct {
	RepoType string   `yaml:"repoType"`
	Repo     string   // name of the repo
	Files    []string // reference to files in repo
	Regex    string   // regex for searching the to-be-replaced item
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

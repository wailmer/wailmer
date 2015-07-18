package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/fsouza/go-dockerclient"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
)

var (
	flagEndpoint = flag.String("endpoint", "", "Docker endpoint. Can be unix:///my_path or tcp://IP:PORT")
)

type job struct {
	Name   string        `json:"name",yaml:"name"`
	Config docker.Config `json:"config",yaml:"config"`
}

type config struct {
	Name    string        `json:"name",yaml:"name"`
	Config  docker.Config `json:"config",yaml:"config"`
	Version string        `json:"Version",yaml:"Version"`
	Jobs    []job         `json:"jobs",yaml:"jobs"`
}

func (c *config) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func main() {
	//	var help bool

	//	var helpText = /*help text omitted for readability, shown in output instead.*/

	flag.Parse()
	var client *docker.Client
	if *flagEndpoint != "" {
		var err error
		if client, err = docker.NewClient(*flagEndpoint); err != nil {
			log.Fatal(err)
		}
	} else {
		var err error
		if client, err = docker.NewClientFromEnv(); err != nil {
			log.Fatal(err)
		}
	}
	data, err := ioutil.ReadFile("wailmer.yml")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	if err := config.Parse(data); err != nil {
		log.Fatal(err)
	}
	for j, _ := range config.Jobs {
		if err := mergo.Merge(&config.Jobs[j].Config, config.Config); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("%+v\n", config)
}

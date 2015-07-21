package main

import (
	"errors"
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
	Version string        `json:"version",yaml:"version"`
	Jobs    []job         `json:"jobs",yaml:"jobs"`
}

func (c *config) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func findJob(name string, config config) (job, error) {
	for _, job := range config.Jobs {
		if job.Name == name {
			return job, nil
		}
	}
	return job{}, errors.New("No job found")
}

func main() {
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
	for _, job := range flag.Args() {
		log.Print(job)
		j, err := findJob(job, config)
		if err != nil {
			log.Fatal(err)
		}
		opts := docker.CreateContainerOptions{Name: j.Name, Config: &j.Config}
		container, err := client.CreateContainer(opts)
		if err != nil {
			log.Fatal(err)
		}
		log.Print(container)
	}

}

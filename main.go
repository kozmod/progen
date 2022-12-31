package main

import (
	"flag"
	"log"

	"github.com/kozmod/progen/interanl/config"
	"github.com/kozmod/progen/interanl/entity"
)

var (
	configPath = flag.String("c", "", "config file path")
)

func main() {
	flag.Parse()

	if *configPath == entity.Empty {
		log.Fatal("config file is not set")
	}

	rawConfig, err := config.MustPreprocessConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	conf, err := config.UnmarshalYamlConfig(rawConfig)
	if err != nil {
		log.Fatal(err)
	}

	order, err := config.YamlRootNodesOrder(rawConfig)
	if err != nil {
		log.Fatal(err)
	}

	procChain, err := config.MustConfigureProcChain(conf, order)
	if err != nil {
		log.Fatal(err)
	}

	err = procChain.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

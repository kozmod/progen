package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/factory"
)

const (
	defaultConfigFilePath = "progen.yml"
)

var (
	flagConfigPath = flag.String("f", "",
		fmt.Sprintf("configuration file path (if not set, default configuration is used: %s)", defaultConfigFilePath))
	flagVerbose = flag.Bool("v", false, "verbose output")
)

func main() {
	flag.Parse()

	logger, err := factory.NewLogger(*flagVerbose)
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	if *flagConfigPath == entity.Empty {
		*flagConfigPath = defaultConfigFilePath
		logger.Infof("default configuration is used: %s", defaultConfigFilePath)
	}

	rawConfig, templateData, err := config.PreprocessRawConfigData(*flagConfigPath)
	if err != nil {
		logger.Fatalf("preprocess raw config: %v", err)
	}
	logger.Infof("config to apply:\n%s", string(rawConfig))

	conf, err := config.UnmarshalYamlConfig(rawConfig)
	if err != nil {
		logger.Fatalf("unmarshal config: %v", err)
	}

	order, err := config.YamlRootNodesOrder(rawConfig)
	if err != nil {
		logger.Fatalf("define tag order: %v", err)
	}

	procChain, err := factory.NewProcChain(conf, order, templateData, logger)
	if err != nil {
		logger.Fatalf("create processors chain: %v", err)
	}

	_ = procChain

	err = procChain.Exec()
	if err != nil {
		logger.Fatalf("execute chain: %v", err)
	}
}

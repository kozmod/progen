package main

import (
	"flag"
	"fmt"
	"github.com/kozmod/progen/internal"
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/factory"
)

//goland:noinspection SpellCheckingInspection
const (
	defaultConfigFilePath = "progen.yml"
)

var (
	flagConfigPath = flag.String(
		"f",
		entity.Empty,
		fmt.Sprintf("configuration file path (default file: %s)", defaultConfigFilePath))
	flagVerbose = flag.Bool(
		"v",
		false,
		"verbose output")
	flagDryRun = flag.Bool(
		"dr",
		false,
		"dry run mode (to verbose output should be combine with`-v`)")
	flagVersion = flag.Bool(
		"version",
		false,
		"output version")
)

func main() {
	flag.Parse()

	if flagVersion != nil && *flagVersion {
		fmt.Println(internal.GetVersion())
		return
	}

	logger, err := factory.NewLogger(*flagVerbose)
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	defer func(start time.Time) {
		logger.Infof("execution time: %v", time.Since(start))
	}(time.Now())

	if *flagConfigPath == entity.Empty {
		*flagConfigPath = defaultConfigFilePath
		logger.Infof("default configuration is used: %s", defaultConfigFilePath)
	}

	data, err := os.ReadFile(*flagConfigPath)
	if err != nil {
		logger.Fatalf("read config: %v", err)
	}

	rawConfig, templateData, err := config.PreprocessRawConfigData(*flagConfigPath, data)
	if err != nil {
		logger.Fatalf("preprocess raw config: %v", err)
	}

	var (
		eg errgroup.Group

		conf  config.Config
		files map[string][]config.File
	)

	eg.Go(func() error {
		c, err := config.UnmarshalYamlConfig(rawConfig)
		if err != nil {
			return fmt.Errorf("unmarshal config: %w", err)
		}
		conf = c
		return nil
	})

	eg.Go(func() error {
		f, err := config.UnmarshalYamlFiles(data)
		if err != nil {
			return fmt.Errorf("unmarshal file config: %w", err)
		}
		files = f
		return nil
	})

	if err = eg.Wait(); err != nil {
		logger.Fatalf("prepare config: %v", err)
	}

	conf, err = config.PrepareFiles(conf, files)
	if err != nil {
		logger.Fatalf("prepare config files section: %v", err)
	}

	procChain, err := factory.NewProcChain(conf, templateData, logger, *flagDryRun)
	if err != nil {
		logger.Fatalf("create processors chain: %v", err)
	}

	err = procChain.Exec()
	if err != nil {
		logger.Fatalf("execute chain: %v", err)
	}
}

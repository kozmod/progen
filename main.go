package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/kozmod/progen/internal"
	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/factory"
	"github.com/kozmod/progen/internal/flag"
)

func main() {

	flags := flag.Parse()

	if flags.Version {
		fmt.Println(internal.GetVersion())
		return
	}

	logger, err := factory.NewLogger(flags.Verbose)
	if err != nil {
		log.Fatalf("create logger: %v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	{
		if err := os.Chdir(flags.AWD); err != nil {
			logger.Fatalf("changes the application working directory: %v", err)
		}
		awd, err := os.Getwd()
		if err != nil {
			logger.Fatalf("get the application working directory: %v", err)
		}
		logger.Infof("application working directory: %s", awd)
	}

	defer func(start time.Time) {
		logger.Infof("execution time: %v", time.Since(start))
	}(time.Now())

	logger.Infof("configuration read: %s", flags.FileLocationMessage())

	data, err := config.NewConfigReader(flags).Read()
	if err != nil {
		logger.Fatalf("read config: %v", err)
	}

	rawConfig, templateData, err := config.NewRawPreprocessor(
		flags.ConfigPath,
		flags.TemplateVars.Vars,
		entity.TemplateFnsMap).
		Process(data)
	if err != nil {
		logger.Fatalf("preprocess raw config: %v", err)
	}

	var (
		eg errgroup.Group

		conf  config.Config
		files map[string][]config.File

		tagFilter = entity.NewRegexpChain(flags.Skip...)
	)

	eg.Go(func() error {
		c, err := config.NewYamlConfigUnmarshaler(tagFilter, logger).Unmarshal(rawConfig)
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

	procChain, err := factory.NewProcChain(conf, templateData, logger, flags.DryRun)
	if err != nil {
		logger.Fatalf("create processors chain: %v", err)
	}

	err = procChain.Exec()
	if err != nil {
		logger.Fatalf("execute chain: %v", err)
	}
}

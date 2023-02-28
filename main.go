package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

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
		log.Fatalf("create logger: %+v", err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	{
		if err = os.Chdir(flags.AWD); err != nil {
			logger.Fatalf("changes the application working directory: %+v", err)
		}

		var awd string
		awd, err = os.Getwd()
		if err != nil {
			logger.Fatalf("get the application working directory: %+v", err)
		}
		logger.Infof("application working directory: %s", awd)
	}

	defer func(start time.Time) {
		logger.Infof("execution time: %v", time.Since(start))
	}(time.Now())

	logger.Infof("configuration file: %s", flags.FileLocationMessage())

	data, err := config.NewConfigReader(flags).Read()
	if err != nil {
		logger.Fatalf("read config: %+v", err)
	}

	rawConfig, templateData, err := config.NewRawPreprocessor(
		flags.ConfigPath,
		flags.TemplateVars.Vars,
		entity.TemplateFnsMap,
		[]string{flags.MissingKey.String()},
	).
		Process(data)
	if err != nil {
		logger.Fatalf("preprocess raw config: %+v", err)
	}

	var (
		eg errgroup.Group

		conf      config.Config
		tagFilter = entity.NewRegexpChain(flags.Skip...)
	)

	eg.Go(func() error {
		conf, err = config.NewYamlConfigUnmarshaler(tagFilter, logger).Unmarshal(rawConfig)
		if err != nil {
			return xerrors.Errorf("unmarshal config: %w", err)
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		logger.Fatalf("prepare config: %+v", err)
	}

	procChain, err := factory.NewExecutorChain(
		conf,
		templateData,
		[]string{flags.MissingKey.String()},
		logger,
		flags.PreprocessFiles,
		flags.DryRun)
	if err != nil {
		logger.Fatalf("create processors chain: %+v", err)
	}

	err = procChain.Exec()
	if err != nil {
		logger.Fatalf("execute chain: %+v", err)
	}
}

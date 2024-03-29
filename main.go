package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	logFatalSuffixFn := func(s string) string {
		const pv, v = "%+v", "%v"
		if flags.PrintErrorStackTrace {
			return s + pv
		}
		return s + v
	}

	logger, err := factory.NewLogger(flags.Verbose)
	if err != nil {
		log.Fatalf(logFatalSuffixFn("create logger: "), err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	{
		if err = os.Chdir(flags.AWD); err != nil {
			logger.Fatalf(logFatalSuffixFn("changes the application working directory: "), xerrors.Errorf("%w", err))
		}

		var awd string
		awd, err = os.Getwd()
		if err != nil {
			logger.Fatalf(logFatalSuffixFn("get the application working directory: "), xerrors.Errorf("%w", err))
		}
		logger.Infof("application working directory: %s", awd)
	}

	defer func(start time.Time) {
		logger.Infof("execution time: %v", time.Since(start))
	}(time.Now())

	logger.Infof("configuration file: %s", flags.FileLocationMessage())

	data, err := config.NewConfigReader(flags).Read()
	if err != nil {
		logger.Fatalf(logFatalSuffixFn("read config: "), err)
	}

	rawConfig, templateData, err := config.NewRawPreprocessor(
		flags.ConfigPath,
		flags.TemplateVars.Vars,
		entity.TemplateFnsMap,
		[]string{flags.MissingKey.String()},
	).Process(data)
	if err != nil {
		logger.Fatalf(logFatalSuffixFn("preprocess raw config: "), err)
	}

	if flags.PrintProcessedConfig {
		logger.ForceInfof("preprocessed config:\n%s", string(rawConfig))
	}

	var (
		conf config.Config
	)
	conf, err = config.NewYamlConfigUnmarshaler().Unmarshal(rawConfig)
	if err != nil {
		logger.Fatalf(logFatalSuffixFn("unmarshal config: "), err)
	}

	if err = conf.Validate(); err != nil {
		logger.Fatalf(logFatalSuffixFn("validate config: "), err)
	}

	procChain, err := factory.NewExecutorChain(
		conf,
		factory.NewActionFilter(
			flags.Skip,
			flags.Group,
			conf.Settings.Groups.GroupByAction(),
			conf.Settings.Groups.ManualActions(),
			logger,
		),
		templateData,
		[]string{flags.MissingKey.String()},
		logger,
		flags.PreprocessFiles,
		flags.DryRun)
	if err != nil {
		logger.Fatalf(logFatalSuffixFn("create processors chain: "), err)
	}

	err = procChain.Exec()
	if err != nil {
		logger.Fatalf(logFatalSuffixFn("execute chain: "), err)
	}
}

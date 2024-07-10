package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kozmod/progen/internal/exec"
	"github.com/kozmod/progen/internal/factory"
	"log"
	"os"
	"time"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal"
	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/flag"
)

func main() {

	flags := flag.Parse()

	if flags.Version {
		fmt.Println(internal.GetVersion())
		return
	}

	logFatalSuffixFn := entity.NewAppendVPlusOrV(flags.PrintErrorStackTrace)

	logger, err := factory.NewLogger(flags.Verbose)
	if err != nil {
		log.Fatalf(logFatalSuffixFn("create logger: "), err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	{
		if err = os.Chdir(flags.AWD); err != nil {
			logger.Errorf(logFatalSuffixFn("changes the application working directory: "), xerrors.Errorf("%w", err))
			return
		}

		var awd string
		awd, err = os.Getwd()
		if err != nil {
			logger.Errorf(logFatalSuffixFn("get the application working directory: "), xerrors.Errorf("%w", err))
			return
		}
		logger.Infof("application working directory: %s", awd)
	}

	defer func(start time.Time) {
		logger.Infof("execution time: %v", time.Since(start))
	}(time.Now())

	logger.Infof("configuration file: %s", flags.FileLocationMessage())

	data, err := config.NewConfigReader(flags).Read()
	if err != nil {
		logger.Errorf(logFatalSuffixFn("read config: "), err)
		return
	}

	rawConfig, templateData, err := config.NewRawPreprocessor(
		flags.ConfigPath,
		flags.TemplateVars.Vars,
		entity.TemplateFnsMap,
		[]string{flags.MissingKey.String()},
	).Process(data)
	if err != nil {
		logger.Errorf(logFatalSuffixFn("preprocess raw config: "), err)
		return
	}

	if flags.PrintProcessedConfig {
		logger.ForceInfof("preprocessed config:\n%s", string(rawConfig))
	}

	var (
		conf config.Config
	)
	conf, err = config.NewYamlConfigUnmarshaler().Unmarshal(rawConfig)
	if err != nil {
		logger.Errorf(logFatalSuffixFn("unmarshal config: "), err)
		return
	}

	if err = conf.Validate(); err != nil {
		logger.Errorf(logFatalSuffixFn("validate config: "), err)
		return
	}

	var (
		actionFilter = factory.NewActionFilter(
			flags.Skip,
			flags.Group,
			conf.Settings.Groups.GroupByAction(),
			conf.Settings.Groups.ManualActions(),
			logger,
		)
		templateOptions = []string{flags.MissingKey.String()}
		preprocessors   = &exec.Preprocessors{}
	)

	procChain, err := factory.NewExecutorChainFactory(
		logger,
		flags.DryRun,
		func(executors []entity.Executor) entity.Executor {
			return exec.NewPreprocessingChain(preprocessors, executors)
		},
		factory.NewExecutorBuilderFactory(
			conf.DirActions(),
			factory.NewMkdirExecutor,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			conf.RmActions(),
			factory.NewRmExecutor,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			conf.CommandActions(),
			factory.NewRunCommandExecutor,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			conf.FilesActions(),
			factory.NewPreprocessorsFileExecutorFactory(
				templateData,
				templateOptions,
				flags.PreprocessFiles,
				preprocessors,
				func(logger entity.Logger) *resty.Client {
					return factory.NewHTTPClient(conf.Settings.HTTP, logger)
				},
			).Create,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			conf.FsActions(),
			factory.NewFsExecFactory(
				templateData,
				templateOptions,
			).Create,
			actionFilter,
		),
	).Create()
	if err != nil {
		logger.Errorf(logFatalSuffixFn("create processors chain: "), err)
		return
	}

	err = procChain.Exec()
	if err != nil {
		logger.Errorf(logFatalSuffixFn("execute chain: "), err)
		return
	}
}

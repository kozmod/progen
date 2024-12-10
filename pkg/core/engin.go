package core

import (
	"sync"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
	"github.com/kozmod/progen/internal/factory"
)

type Engin struct {
	mx       sync.RWMutex
	files    []entity.Action[[]entity.UndefinedFile]
	cmd      []entity.Action[[]entity.Command]
	dirs     []entity.Action[[]string]
	fsModify []entity.Action[[]string]
	fsSave   []entity.Action[[]entity.TargetFs]
	rm       []entity.Action[[]string]

	logger entity.Logger
}

// AddActions adds action to the [Engin].
func (e *Engin) AddActions(actions ...action) *Engin {
	e.mx.Lock()
	defer e.mx.Unlock()
	for _, action := range actions {
		action.add(e)
	}
	return e
}

// Run all actions.
func (e *Engin) Run(config *Config) error {
	if config == nil {
		config = &Config{}
	}

	var (
		logFatalSuffixFn = entity.NewAppendVPlusOrV(config.PrintErrorStackTrace)
		actionFilter     factory.DummyActionFilter
	)

	e.mx.RLock()
	defer e.mx.RUnlock()

	var (
		logger entity.Logger
		err    error
	)
	if e.logger == nil {
		logger, err = factory.NewLogger(config.Verbose)
		if err != nil {
			return xerrors.Errorf("failed to initialize logger: %w", err)
		}
	}

	procChain, err := factory.NewExecutorChainFactory(
		logger,
		config.DryRun,
		func(executors []entity.Executor) entity.Executor {
			return exec.NewChain(executors)
		},
		factory.NewExecutorBuilderFactory(
			e.dirs,
			factory.NewMkdirExecutor,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			e.fsModify,
			factory.NewFsModifyExecFactory(
				config.TemplateVars.Vars,
				[]string{config.MissingKey.String()},
			).Create,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			e.fsSave,
			factory.NewFsSaveExecFactory(
				config.TemplateVars.Vars,
				[]string{config.MissingKey.String()},
			).Create,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			e.rm,
			factory.NewRmExecutor,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			e.files,
			factory.NewFileExecutorFactory(
				config.TemplateVars.Vars,
				[]string{config.MissingKey.String()},
			).Create,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			e.cmd,
			factory.NewRunCommandExecutor,
			actionFilter,
		),
	).Create()

	if err != nil {
		//goland:noinspection ALL
		logger.Errorf(logFatalSuffixFn("create processors chain: "), err)
		return err
	}

	err = procChain.Exec()
	if err != nil {
		//goland:noinspection ALL
		logger.Errorf(logFatalSuffixFn("execute chain: "), err)
		return err
	}
	return nil
}

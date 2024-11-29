package core

import (
	"sync"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
	"github.com/kozmod/progen/internal/factory"
	"github.com/kozmod/progen/internal/flag"
)

type Engin struct {
	mx    sync.RWMutex
	files []entity.Action[[]entity.UndefinedFile]
	cmd   []entity.Action[[]entity.Command]
	dirs  []entity.Action[[]string]
	fs    []entity.Action[[]string]
	rm    []entity.Action[[]string]

	logger entity.Logger
}

func (e *Engin) AddActions(actions ...Action) *Engin {
	e.mx.Lock()
	defer e.mx.Unlock()
	for _, action := range actions {
		action.add(e)
	}
	return e
}

func (e *Engin) Run() error {
	var (
		flags = flag.ParseDefault()

		logFatalSuffixFn = entity.NewAppendVPlusOrV(flags.PrintErrorStackTrace)
		actionFilter     factory.DummyActionFilter
	)

	e.mx.RLock()
	defer e.mx.RUnlock()

	var logger entity.Logger
	if e.logger == nil {
		l, err := factory.NewLogger(flags.Verbose)
		if err != nil {
			return xerrors.Errorf("failed to initialize logger: %w", err)
		}
		logger = l
	}

	procChain, err := factory.NewExecutorChainFactory(
		logger,
		flags.DryRun,
		func(executors []entity.Executor) entity.Executor {
			return exec.NewChain(executors)
		},
		factory.NewExecutorBuilderFactory(
			e.dirs,
			factory.NewMkdirExecutor,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			e.fs,
			factory.NewFsExecFactory(
				flags.TemplateVars.Vars,
				[]string{flags.MissingKey.String()},
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
				flags.TemplateVars.Vars,
				[]string{flags.MissingKey.String()},
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

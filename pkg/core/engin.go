package core

import (
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
	"github.com/kozmod/progen/internal/factory"
	"golang.org/x/xerrors"
	"sync"
)

type Engin struct {
	mx    sync.RWMutex
	files []entity.Action[[]entity.UndefinedFile]
	cmd   []entity.Action[[]entity.Command]
	dirs  []entity.Action[[]string]
	fs    []entity.Action[[]string]
	rm    []entity.Action[[]string]

	templateVars    map[string]any
	templateOptions []string

	logger entity.Logger
}

func (e *Engin) AddTemplateVars(vars map[string]any) *Engin {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.templateVars = entity.MergeKeys(e.templateVars, vars)
	return e
}

func (e *Engin) AddFilesActions(a ...entity.Action[[]File]) *Engin {
	e.mx.Lock()
	defer e.mx.Unlock()
	actions := toEntityActions(func(f File) entity.UndefinedFile {
		return entity.UndefinedFile{
			Path: f.Path,
			Data: &f.Data,
		}
	}, a...)
	e.files = append(e.files, actions...)
	return e
}

func (e *Engin) AddExecuteCommandActions(a ...entity.Action[[]Cmd]) *Engin {
	e.mx.Lock()
	defer e.mx.Unlock()
	commands := toEntityActions(func(cmd Cmd) entity.Command {
		return entity.Command(cmd)
	}, a...)
	e.cmd = append(e.cmd, commands...)
	return e
}

func (e *Engin) AddCreateDirsActions(a ...entity.Action[[]string]) *Engin {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.dirs = append(e.dirs, a...)
	return e
}

func (e *Engin) AddRmActions(a ...entity.Action[[]string]) *Engin {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.rm = append(e.rm, a...)
	return e
}

func (e *Engin) AddFsActions(a ...entity.Action[[]string]) *Engin {
	e.mx.Lock()
	defer e.mx.Unlock()
	e.fs = append(e.fs, a...)
	return e
}

func (e *Engin) Run() error {
	var (
		logFatalSuffixFn = entity.NewAppendVPlusOrV(true) //todo
		actionFilter     factory.DummyActionFilter
	)

	e.mx.RLock()
	defer e.mx.RUnlock()

	var logger entity.Logger
	if e.logger == nil {
		l, err := factory.NewLogger(true) //todo
		if err != nil {
			return xerrors.Errorf("failed to initialize logger: %w", err)
		}
		logger = l
	}

	procChain, err := factory.NewExecutorChainFactory(
		logger,
		true, //todo
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
				e.templateVars,
				e.templateOptions,
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
				e.templateVars,
				e.templateOptions,
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
		logger.Errorf(logFatalSuffixFn("create processors chain: "), err)
		return err
	}

	err = procChain.Exec()
	if err != nil {
		logger.Errorf(logFatalSuffixFn("execute chain: "), err)
		return err
	}
	return nil
}

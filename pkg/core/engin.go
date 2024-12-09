package core

import (
	"flag"
	"log"
	"os"
	"sync"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
	"github.com/kozmod/progen/internal/factory"
	engineFlags "github.com/kozmod/progen/internal/flag"
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

	flags engineFlags.DefaultFlags
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
func (e *Engin) Run() error {
	var (
		args    = os.Args
		flagSet = flag.NewFlagSet(args[0], flag.ExitOnError)
	)
	err := e.flags.Parse(flagSet, args)
	if err != nil {
		log.Printf("parse flags failed: %v", err)
		return err
	}

	var (
		logFatalSuffixFn = entity.NewAppendVPlusOrV(e.flags.PrintErrorStackTrace)
		actionFilter     factory.DummyActionFilter
	)

	e.mx.RLock()
	defer e.mx.RUnlock()

	var logger entity.Logger
	if e.logger == nil {
		logger, err = factory.NewLogger(e.flags.Verbose)
		if err != nil {
			return xerrors.Errorf("failed to initialize logger: %w", err)
		}
	}

	procChain, err := factory.NewExecutorChainFactory(
		logger,
		e.flags.DryRun,
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
				e.flags.TemplateVars.Vars,
				[]string{e.flags.MissingKey.String()},
			).Create,
			actionFilter,
		),
		factory.NewExecutorBuilderFactory(
			e.fsSave,
			factory.NewFsSaveExecFactory(
				e.flags.TemplateVars.Vars,
				[]string{e.flags.MissingKey.String()},
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
				e.flags.TemplateVars.Vars,
				[]string{e.flags.MissingKey.String()},
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

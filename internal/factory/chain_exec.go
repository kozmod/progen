package factory

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

type (
	executorBuilderFactory interface {
		Create(logger entity.Logger, dryRun bool) []entity.ExecutorBuilder
	}
)

type ExecutorChainFactory struct {
	logger entity.Logger
	dryRun bool

	executorBuilderFactories []executorBuilderFactory
	createFn                 func([]entity.Executor) entity.Executor
}

func NewExecutorChainFactory(
	logger entity.Logger,
	dryRun bool,
	createFn func([]entity.Executor) entity.Executor,
	executorBuilderFactories ...executorBuilderFactory,

) *ExecutorChainFactory {
	return &ExecutorChainFactory{
		createFn:                 createFn,
		logger:                   logger,
		dryRun:                   dryRun,
		executorBuilderFactories: executorBuilderFactories,
	}
}

func (f ExecutorChainFactory) Create() (entity.Executor, error) {
	var (
		allBuilders []entity.ExecutorBuilder
	)
	for _, factory := range f.executorBuilderFactories {
		builder := factory.Create(f.logger, f.dryRun)
		allBuilders = append(allBuilders, builder...)
	}

	sort.Slice(allBuilders, func(i, j int) bool {
		return allBuilders[i].Priority < allBuilders[j].Priority
	})

	actionNames := make([]string, len(allBuilders))
	for i, builder := range allBuilders {
		actionNames[i] = fmt.Sprintf("'%d':'%s'", builder.Priority, builder.Action)
	}
	f.logger.Infof("action is going to be execute ('priopiry':'name')[%s]", strings.Join(actionNames, ","))

	executors := make([]entity.Executor, 0, len(allBuilders))
	for _, builder := range allBuilders {
		e, err := builder.ProcFn()
		if err != nil {
			return nil, xerrors.Errorf("configure executor [%s]: %w", builder.Action, err)
		}
		if e == nil {
			continue
		}
		executors = append(executors, e)
	}

	return f.createFn(executors), nil
}

type (
	actionValConsumer[T any] func(vals []T, logger entity.Logger, dryRun bool) (entity.Executor, error)
)

type ExecutorBuilderFactory[T any] struct {
	actionsSupplier   []entity.Action[[]T]
	actionValConsumer actionValConsumer[T]
	actionFilter      entity.ActionFilter
}

func NewExecutorBuilderFactory[T any](
	actionSupplier []entity.Action[[]T],
	actionValConsumer actionValConsumer[T],
	actionFilter entity.ActionFilter,
) *ExecutorBuilderFactory[T] {
	return &ExecutorBuilderFactory[T]{
		actionsSupplier:   actionSupplier,
		actionValConsumer: actionValConsumer,
		actionFilter:      actionFilter}
}

func (y ExecutorBuilderFactory[T]) Create(logger entity.Logger, dryRun bool) []entity.ExecutorBuilder {
	var (
		actions  = y.actionsSupplier
		builders = make([]entity.ExecutorBuilder, 0, len(actions))
	)
	for _, action := range actions {
		var (
			a    = action
			name = a.Name
		)
		if !y.actionFilter.MatchString(name) {
			continue
		}
		builders = append(builders,
			entity.ExecutorBuilder{
				Action:   name,
				Priority: a.Priority,
				ProcFn: func() (entity.Executor, error) {
					executor, err := y.actionValConsumer(a.Val, logger, dryRun)
					return executor, err
				},
			})
	}
	return builders
}

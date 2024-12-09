package exec

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

type FileSystemModifyStrategy struct {
	logger         entity.Logger
	strategiesFn   func(paths map[string]string) []entity.FileStrategy
	templateProcFn func() entity.TemplateProc
	dirExecutorFn  func(dirs []string) entity.Executor
	fileExecutorFn func(producers []entity.FileProducer, strategies []entity.FileStrategy) entity.Executor
	removeAllFn    func(path string) error
}

func NewFileSystemModifyStrategy(
	templateData,
	templateFns map[string]any,
	templateOptions []string,
	logger entity.Logger) *FileSystemModifyStrategy {
	return &FileSystemModifyStrategy{
		logger: logger,
		strategiesFn: func(paths map[string]string) []entity.FileStrategy {
			return []entity.FileStrategy{
				NewTemplateFileStrategy(templateData, templateFns, templateOptions),
				NewReplacePathFileStrategy(paths),
				NewSaveFileStrategy(logger),
			}
		},
		templateProcFn: func() entity.TemplateProc {
			return entity.NewTemplateProc(templateData, templateFns, templateOptions)
		},
		dirExecutorFn: func(dirs []string) entity.Executor {
			return NewDirExecutor(dirs, []entity.DirStrategy{NewMkdirAllStrategy(logger)})
		},
		fileExecutorFn: func(producers []entity.FileProducer, strategies []entity.FileStrategy) entity.Executor {
			return NewFilesExecutor(producers, strategies)
		},
		removeAllFn: os.RemoveAll,
	}
}

func (e *FileSystemModifyStrategy) Apply(dir string) (string, error) {
	type (
		Entity struct {
			Path  string
			IsDir bool
		}
	)
	var (
		entitySet = make(map[string]Entity)
		filePaths = make(map[string]string)
	)

	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return err
		}
		if dir == path {
			return err
		}
		entPath, err := e.templateProcFn().Process(path, path)
		if err != nil {
			return xerrors.Errorf("fs modify: process template to path [%s]: %w", path, err)
		}
		filePaths[path] = entPath
		entitySet[path] = Entity{
			Path:  entPath,
			IsDir: info.IsDir(),
		}
		return err
	})
	if err != nil {
		return entity.Empty, xerrors.Errorf("fs modify: walk dir [%s]: %w", dir, err)
	}

	var (
		dirs          = make([]string, 0, len(entitySet))
		fileProducers = make([]entity.FileProducer, 0, len(entitySet))
	)

	for old, ent := range entitySet {
		if ent.IsDir {
			dirs = append(dirs, ent.Path)
		}

		if !ent.IsDir {
			file := entity.LocalFile{
				FileInfo:  entity.NewFileInfo(ent.Path),
				LocalPath: old,
			}
			fileProducers = append(fileProducers, NewLocalProducer(file))
		}
	}

	dirExec := e.dirExecutorFn(dirs)
	err = dirExec.Exec()
	if err != nil {
		return entity.Empty, xerrors.Errorf("fs modify: dirs execute: %w", err)
	}

	fileExec := e.fileExecutorFn(fileProducers, e.strategiesFn(filePaths))
	err = fileExec.Exec()
	if err != nil {
		return entity.Empty, xerrors.Errorf("fs modify: files execute: %w", err)
	}

	for old, ent := range entitySet {
		if ent.Path == old {
			continue
		}
		err = e.removeAllFn(old)
		if err != nil {
			return entity.Empty, xerrors.Errorf("fs modify: remove old: %w", err)
		}
		e.logger.Infof("fs modify: remove: %s", old)
	}

	return dir, nil
}

type DryRunFileSystemModifyStrategy struct {
	logger entity.Logger
}

func NewDryRunFileSystemModifyStrategy(logger entity.Logger) *DryRunFileSystemModifyStrategy {
	return &DryRunFileSystemModifyStrategy{
		logger: logger,
	}
}

func (e *DryRunFileSystemModifyStrategy) Apply(dir string) (string, error) {
	e.logger.Infof("fs modify: dir execute: %s", dir)
	return dir, nil
}

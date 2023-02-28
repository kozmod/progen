package exec

import (
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

type FileSystemStrategy struct {
	logger         entity.Logger
	strategiesFn   func(paths map[string]string) []entity.FileStrategy
	templateProcFn func() entity.TemplateProc
	dirExecutorFn  func(dirs []string) entity.Executor
	fileExecutorFn func(producers []entity.FileProducer, strategies []entity.FileStrategy) entity.Executor
	removeAllFn    func(path string) error
}

func NewFileSystemStrategy(
	templateData,
	templateFns map[string]any,
	templateOptions []string,
	logger entity.Logger) *FileSystemStrategy {
	return &FileSystemStrategy{
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

func (e *FileSystemStrategy) Apply(dir string) (string, error) {
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
			return xerrors.Errorf("fs: process template to path [%s]: %w", path, err)
		}
		filePaths[path] = entPath
		entitySet[path] = Entity{
			Path:  entPath,
			IsDir: info.IsDir(),
		}
		return err
	})
	if err != nil {
		return entity.Empty, xerrors.Errorf("fs: walk dir [%s]: %w", dir, err)
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
		return entity.Empty, xerrors.Errorf("fs: dirs execute: %w", err)
	}

	fileExec := e.fileExecutorFn(fileProducers, e.strategiesFn(filePaths))
	err = fileExec.Exec()
	if err != nil {
		return entity.Empty, xerrors.Errorf("fs: files execute: %w", err)
	}

	for old, ent := range entitySet {
		if ent.Path == old {
			continue
		}
		err = e.removeAllFn(old)
		if err != nil {
			return entity.Empty, xerrors.Errorf("fs: remove old: %w", err)
		}
		e.logger.Infof("fs: remove: %s", old)
	}

	return dir, nil
}

type DryRunFileSystemStrategy struct {
	logger entity.Logger
}

func NewDryRunFileSystemStrategy(logger entity.Logger) *DryRunFileSystemStrategy {
	return &DryRunFileSystemStrategy{
		logger: logger,
	}
}

func (e *DryRunFileSystemStrategy) Apply(dir string) (string, error) {
	e.logger.Infof("fs: dir execute: %s", dir)
	return dir, nil
}

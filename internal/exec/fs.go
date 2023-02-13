package exec

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kozmod/progen/internal/entity"
)

type FileSystemProc struct {
	logger         entity.Logger
	processorsFn   func(paths map[string]string) []entity.FileProc
	templateProcFn func() entity.TemplateProc
	dirExecutorFn  func(dirs []string) entity.Executor
	fileExecutorFn func(producers []entity.FileProducer, processors []entity.FileProc) entity.Executor
	removeAllFn    func(path string) error
}

func NewFileSystemProc(
	templateData,
	templateFns map[string]any,
	templateOptions []string,
	logger entity.Logger) *FileSystemProc {
	return &FileSystemProc{
		logger: logger,
		processorsFn: func(paths map[string]string) []entity.FileProc {
			return []entity.FileProc{
				NewTemplateFileProc(templateData, templateFns, templateOptions),
				NewReplacePathFileProc(paths),
				NewSaveFileProc(logger),
			}
		},
		templateProcFn: func() entity.TemplateProc {
			return entity.NewTemplateProc(templateData, templateFns, templateOptions)
		},
		dirExecutorFn: func(dirs []string) entity.Executor {
			return NewDirExecutor(dirs, []entity.DirProc{NewMkdirAllProc(logger)})
		},
		fileExecutorFn: func(producers []entity.FileProducer, processors []entity.FileProc) entity.Executor {
			return NewFilesExecutor(producers, processors)
		},
		removeAllFn: os.RemoveAll,
	}
}

func (e *FileSystemProc) Process(dir string) (string, error) {
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
			return fmt.Errorf("fs: process template to path [%s]: %w", path, err)
		}
		filePaths[path] = entPath
		entitySet[path] = Entity{
			Path:  entPath,
			IsDir: info.IsDir(),
		}
		return err
	})
	if err != nil {
		return entity.Empty, fmt.Errorf("fs: walk dir [%s]: %w", dir, err)
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
		return entity.Empty, fmt.Errorf("fs: dirs execute: %w", err)
	}

	fileExec := e.fileExecutorFn(fileProducers, e.processorsFn(filePaths))
	err = fileExec.Exec()
	if err != nil {
		return entity.Empty, fmt.Errorf("fs: files execute: %w", err)
	}

	for old, ent := range entitySet {
		if ent.Path == old {
			continue
		}
		err = e.removeAllFn(old)
		if err != nil {
			return entity.Empty, fmt.Errorf("fs: remove old: %w", err)
		}
		e.logger.Infof("fs: remove: %s", old)
	}

	return dir, nil
}

type DryRunFileSystemProc struct {
	logger entity.Logger
}

func NewDryRunFileSystemProc(logger entity.Logger) *DryRunFileSystemProc {
	return &DryRunFileSystemProc{
		logger: logger,
	}
}

func (e *DryRunFileSystemProc) Process(dir string) (string, error) {
	e.logger.Infof("fs: dir execute: %s", dir)
	return dir, nil
}

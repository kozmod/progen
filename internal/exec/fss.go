package exec

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

type FileSystemSaveStrategy struct {
	fs             fs.FS
	logger         entity.Logger
	strategiesFn   func() []entity.FileStrategy
	templateProcFn func() entity.TemplateProc
	dirExecutorFn  func(dirs []string) entity.Executor
	fileExecutorFn func(producers []entity.FileProducer, strategies []entity.FileStrategy) entity.Executor
	removeAllFn    func(path string) error
}

func NewFileSystemSaveStrategy(
	fs fs.FS,
	templateData,
	templateFns map[string]any,
	templateOptions []string,
	logger entity.Logger) *FileSystemSaveStrategy {
	return &FileSystemSaveStrategy{
		fs:     fs,
		logger: logger,
		strategiesFn: func() []entity.FileStrategy {
			return []entity.FileStrategy{
				NewTemplateFileStrategy(templateData, templateFns, templateOptions),
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

func (e *FileSystemSaveStrategy) Apply(targetDir string) (string, error) {
	var (
		dirs          []string
		fileProducers []entity.FileProducer
		root          = entity.Dot
	)

	err := fs.WalkDir(e.fs, root, func(path string, info fs.DirEntry, err error) error {
		if info == nil {
			return err
		}

		entPath, err := e.templateProcFn().Process(path, path)
		if err != nil {
			return xerrors.Errorf("fs save: process template to path [%s]: %w", path, err)
		}

		entPath = filepath.Join(targetDir, entPath)

		var (
			srcFile fs.File
			data    []byte
		)
		switch {
		case info.IsDir():
			dirs = append(dirs, entPath)
		default:
			srcFile, err = e.fs.Open(path)
			if err != nil {
				return fmt.Errorf("fs save: open fs file [%s]: %v", path, err)
			}
			defer func() {
				_ = srcFile.Close()
			}()

			data, err = io.ReadAll(srcFile)
			if err != nil {
				return fmt.Errorf("fs save: read fs file [%s]: %v", path, err)
			}

			fileProducers = append(fileProducers,
				NewDummyProducer(
					entity.DataFile{
						FileInfo: entity.NewFileInfo(entPath),
						Data:     data,
					},
				),
			)
		}
		return err
	})
	if err != nil {
		return entity.Empty, xerrors.Errorf("fs save: walk dir: %w", err)
	}

	dirExec := e.dirExecutorFn(dirs)
	err = dirExec.Exec()
	if err != nil {
		return entity.Empty, xerrors.Errorf("fs save: dirs execute: %w", err)
	}

	fileExec := e.fileExecutorFn(fileProducers, e.strategiesFn())
	err = fileExec.Exec()
	if err != nil {
		return entity.Empty, xerrors.Errorf("fs save: files execute: %w", err)
	}

	return targetDir, nil
}

type DryRunFileSystemSaveStrategy struct {
	logger entity.Logger
}

func NewDryRunFileSystemSaveStrategy(logger entity.Logger) *DryRunFileSystemSaveStrategy {
	return &DryRunFileSystemSaveStrategy{
		logger: logger,
	}
}

func (e *DryRunFileSystemSaveStrategy) Apply(dir string) (string, error) {
	e.logger.Infof("fs save: dir execute: %s", dir)
	return dir, nil
}

package proc

import (
	"container/list"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/kozmod/progen/internal/entity"
)

type MkdirAllProc struct {
	defaultPerm os.FileMode
	dirs        []entity.Dir
	logger      entity.Logger
}

func NewMkdirAllProc(dirs []entity.Dir, logger entity.Logger) *MkdirAllProc {
	return &MkdirAllProc{
		defaultPerm: os.ModePerm,
		dirs:        dirs,
		logger:      logger,
	}
}

func (p *MkdirAllProc) Exec() error {
	for _, dir := range p.dirs {
		err := os.MkdirAll(dir.Path, p.defaultPerm)
		if err != nil {
			return fmt.Errorf("create dir [%s:%v]: %w", dir.Path, dir.Perm, err)
		}

		if p.defaultPerm != dir.Perm {
			l := list.New()
			err = filepath.WalkDir(dir.Path, func(path string, d fs.DirEntry, err error) error {
				if d.IsDir() {
					l.PushFront(path)
				}
				return err
			})
			if err != nil {
				return fmt.Errorf("walk dir [%s]: %w", dir.Path, err)
			}

			for e := l.Front(); e != nil; e = e.Next() {
				path := e.Value.(string)
				err = os.Chmod(path, dir.Perm)
				if err != nil {
					return fmt.Errorf("chmod dir [path:%s, perm:%v]: %w", path, dir.Perm, err)
				}
			}

		}

		p.logger.Infof("dir created [perm: %v]: %s ", dir.Perm, dir.Path)
	}
	return nil
}

type DryRunMkdirAllProc struct {
	fileMode os.FileMode
	dirs     []entity.Dir
	logger   entity.Logger
}

func NewDryRunMkdirAllProc(dirs []entity.Dir, logger entity.Logger) *DryRunMkdirAllProc {
	return &DryRunMkdirAllProc{
		fileMode: os.ModePerm,
		dirs:     dirs,
		logger:   logger,
	}
}

func (p *DryRunMkdirAllProc) Exec() error {
	for _, dir := range p.dirs {
		p.logger.Infof("dir created [perm: %v]: %s ", dir.Perm, dir.Path)
	}
	return nil
}

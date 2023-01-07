package proc

import (
	"fmt"
	"os"

	"github.com/kozmod/progen/internal/entity"
)

type MkdirAllProc struct {
	fileMode os.FileMode
	dirs     []string
	logger   entity.Logger
}

func NewMkdirAllProc(dirs []string, logger entity.Logger) *MkdirAllProc {
	return &MkdirAllProc{
		fileMode: os.ModePerm,
		dirs:     dirs,
		logger:   logger,
	}
}

func (p *MkdirAllProc) Exec() error {
	for _, dir := range p.dirs {
		err := os.MkdirAll(dir, p.fileMode)
		if err != nil {
			return fmt.Errorf("create dir [%s]: %w", dir, err)
		}
		p.logger.Infof("dir created: %s", dir)
	}
	return nil
}

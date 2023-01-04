package factory

import (
	"path/filepath"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewFileProc(conf config.Config, logger entity.Logger) (proc.Proc, error) {
	if len(conf.Files) == 0 {
		return nil, nil
	}

	files := make([]entity.File, 0, len(conf.Files))
	for _, f := range conf.Files {
		tmpl := entity.File{
			Name: filepath.Base(f.Path),
			Path: filepath.Dir(f.Path),
			Data: []byte(f.Data),
		}
		files = append(files, tmpl)
	}

	return proc.NewFileProc(files, logger), nil
}

package config

import (
	"path/filepath"

	"github.com/kozmod/progen/interanl/entity"
	"github.com/kozmod/progen/interanl/proc"
)

func MustConfigureWriteFileProc(conf Config) (proc.Proc, error) {
	if len(conf.Files) == 0 {
		return nil, nil
	}
	files := toGeneratorFiles(conf.Files)
	return proc.NewWriteFileProc(files), nil
}

func toGeneratorFiles(files []File) []entity.File {
	gFiles := make([]entity.File, 0, len(files))
	for _, file := range files {
		gFiles = append(gFiles, toGeneratorFile(file))
	}
	return gFiles
}

func toGeneratorFile(f File) entity.File {
	return entity.File{
		Name: filepath.Base(f.Path),
		Path: filepath.Dir(f.Path),
		Data: []byte(f.Template),
	}
}

package factory

import (
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewFileProc(conf config.Config, templateData map[string]any, logger entity.Logger) (proc.Proc, error) {
	if len(conf.Files) == 0 {
		return nil, nil
	}

	producers := make([]entity.FileProducer, 0, len(conf.Files))
	for _, f := range conf.Files {

		var (
			producer entity.FileProducer
			client   *resty.Client
		)

		switch {
		case f.Data != nil:
			file := entity.File{
				Name: filepath.Base(f.Path),
				Path: filepath.Dir(f.Path),
				Data: []byte(*f.Data),
				// always `false` because template preprocess on raw config preprocessing step
				Template: false,
			}
			producer = proc.NewStoredProducer(file)
		case f.Get != nil:
			file := entity.RemoteFile{
				Name:     filepath.Base(f.Path),
				Path:     filepath.Dir(f.Path),
				URL:      f.Get.URL,
				Headers:  f.Get.Headers,
				Template: f.Template,
			}

			if client == nil {
				client = NewHTTPClient(conf.HTTP)
			}

			producer = proc.NewRemoteProducer(file, client)
		case f.Local != nil:
			file := entity.LocalFile{
				Name:      filepath.Base(f.Path),
				Path:      filepath.Dir(f.Path),
				LocalPath: *f.Local,
				Template:  f.Template,
			}
			producer = proc.NewLocalProducer(file)

		default:
			return nil, fmt.Errorf("create file processor: one of `data`, `get`, `local` must not be empty")
		}

		producers = append(producers, producer)
	}

	return proc.NewFileProc(producers, templateData, logger), nil
}

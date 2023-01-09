package factory

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-resty/resty/v2"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewFileProc(conf config.Config, logger entity.Logger) (proc.Proc, error) {
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
			}
			producer = proc.NewStoredProducer(file)
		case f.Get != nil:
			file := entity.RemoteFile{
				Name:    filepath.Base(f.Path),
				Path:    filepath.Dir(f.Path),
				URL:     f.Get.URL.String(),
				Headers: f.Get.Headers,
			}

			if client == nil {
				client = resty.NewWithClient(&http.Client{})
			}

			producer = proc.NewRemoteProducer(file, client)
		default:
			return nil, fmt.Errorf("create file processor: `data` or `get` must not be empty")
		}

		producers = append(producers, producer)
	}

	return proc.NewFileProc(producers, logger), nil
}

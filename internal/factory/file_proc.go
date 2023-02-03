package factory

import (
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewFileProc(
	files []config.File,
	http *config.HTTPClient,
	templateData map[string]any,
	logger entity.Logger,
	dryRun bool,
) (proc.Proc, error) {
	if len(files) == 0 {
		logger.Infof("`files` section is empty")
		return nil, nil
	}

	producers := make([]entity.FileProducer, 0, len(files))

	var client *resty.Client
	for _, f := range files {
		var (
			tmpl = entity.Template{
				Name:     filepath.Base(f.Path),
				Path:     filepath.Dir(f.Path),
				ExecTmpl: !f.ExecTmplSkip,
			}
		)

		var producer entity.FileProducer
		switch {
		case f.Data != nil:
			file := entity.DataFile{
				Template: tmpl,
				Data:     []byte(*f.Data),
			}
			producer = proc.NewStoredProducer(file)
		case f.Get != nil:
			file := entity.RemoteFile{
				Template: tmpl,
				HTTPClientParams: entity.HTTPClientParams{
					URL:         f.Get.URL,
					Headers:     f.Get.Headers,
					QueryParams: f.Get.QueryParams,
				},
			}

			if client == nil {
				client = NewHTTPClient(http, logger)
			}

			producer = proc.NewRemoteProducer(file, client)
		case f.Local != nil:
			file := entity.LocalFile{
				Template:  tmpl,
				LocalPath: *f.Local,
			}
			producer = proc.NewLocalProducer(file)

		default:
			return nil, fmt.Errorf("create file processor: one of `data`, `get`, `local` must not be empty")
		}

		producers = append(producers, producer)
	}

	if dryRun {
		return proc.NewDryRunFileProc(producers, templateData, logger), nil
	}

	return proc.NewFileProc(producers, templateData, entity.TemplateFnsMap, logger), nil
}

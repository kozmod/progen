package factory

import (
	"fmt"
	"path/filepath"

	"github.com/go-resty/resty/v2"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewFileExecutor(
	files []config.File,
	http *config.HTTPClient,
	templateData map[string]any,
	logger entity.Logger,
	preload,
	dryRun bool,
	templateOptions []string,
) (entity.Executor, []entity.Preprocessor, error) {
	if len(files) == 0 {
		logger.Infof("`files` section is empty")
		return nil, nil, nil
	}

	producers := make([]entity.FileProducer, 0, len(files))

	var client *resty.Client
	for _, f := range files {
		var (
			tmpl = entity.FileInfo{
				Name: filepath.Base(f.Path),
				Dir:  filepath.Dir(f.Path),
			}
		)

		var producer entity.FileProducer
		switch {
		case f.Data != nil:
			file := entity.DataFile{
				FileInfo: tmpl,
				Data:     []byte(*f.Data),
			}
			producer = proc.NewDummyProducer(file)
		case f.Get != nil:
			file := entity.RemoteFile{
				FileInfo: tmpl,
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
				FileInfo:  tmpl,
				LocalPath: *f.Local,
			}
			producer = proc.NewLocalProducer(file)

		default:
			return nil, nil, fmt.Errorf("create file processor: one of `data`, `get`, `local` must not be empty")
		}

		producers = append(producers, producer)
	}

	var preprocessors []entity.Preprocessor
	if preload {
		preloadProducers := make([]entity.FileProducer, 0, len(producers))
		preloader := proc.NewPreloadProducer(producers, logger)
		for i := 0; i < len(producers); i++ {
			preloadProducers = append(preloadProducers, preloader)
		}
		producers = preloadProducers
		preprocessors = append(preprocessors, preloader)
	}

	processors := []entity.FileProc{proc.NewTemplateFileProc(templateData, entity.TemplateFnsMap, templateOptions)}

	switch {
	case dryRun:
		processors = append(processors, proc.NewDryRunFileProc(logger))
	default:
		processors = append(processors, proc.NewSaveFileProc(logger))
	}
	executor := proc.NewFilesExecutor(producers, processors)

	return executor, preprocessors, nil
}

package factory

import (
	resty "github.com/go-resty/resty/v2"
	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

func NewFileExecutor(
	files []config.File,
	http *config.HTTPClient,
	templateData map[string]any,
	templateOptions []string,
	logger entity.Logger,
	preprocess,
	dryRun bool,
) (entity.Executor, []entity.Preprocessor, error) {
	if len(files) == 0 {
		logger.Infof("`files` section is empty")
		return nil, nil, nil
	}

	producers := make([]entity.FileProducer, 0, len(files))

	var client *resty.Client
	for _, f := range files {
		var (
			tmpl = entity.NewFileInfo(f.Path)
		)

		var producer entity.FileProducer
		switch {
		case f.Data != nil:
			file := entity.DataFile{
				FileInfo: tmpl,
				Data:     []byte(*f.Data),
			}
			producer = exec.NewDummyProducer(file)
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

			producer = exec.NewRemoteProducer(file, client)
		case f.Local != nil:
			file := entity.LocalFile{
				FileInfo:  tmpl,
				LocalPath: *f.Local,
			}
			producer = exec.NewLocalProducer(file)

		default:
			return nil, nil, xerrors.Errorf("build file executor: one of `data`, `get`, `local` must not be empty")
		}

		producers = append(producers, producer)
	}

	var preprocessors []entity.Preprocessor
	if preprocess {
		preloadProducers := make([]entity.FileProducer, 0, len(producers))
		preloader := exec.NewPreloadProducer(producers, logger)
		for i := 0; i < len(producers); i++ {
			preloadProducers = append(preloadProducers, preloader)
		}
		producers = preloadProducers
		preprocessors = append(preprocessors, preloader)
	}

	strategies := []entity.FileStrategy{exec.NewTemplateFileStrategy(templateData, entity.TemplateFnsMap, templateOptions)}

	switch {
	case dryRun:
		strategies = append(strategies, exec.NewDryRunFileStrategy(logger))
	default:
		strategies = append(strategies, exec.NewSaveFileStrategy(logger))
	}
	executor := exec.NewFilesExecutor(producers, strategies)

	return executor, preprocessors, nil
}

package factory

import (
	"golang.org/x/xerrors"

	resty "github.com/go-resty/resty/v2"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

type FileExecutorFactory struct {
	templateData    map[string]any
	templateOptions []string
}

func NewFileExecutorFactory(
	templateData map[string]any,
	templateOptions []string,
) *FileExecutorFactory {
	return &FileExecutorFactory{
		templateData:    templateData,
		templateOptions: templateOptions,
	}
}

func (ff *FileExecutorFactory) Create(files []entity.UndefinedFile, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(files) == 0 {
		logger.Infof("`files` section is empty")
		return nil, nil
	}

	producers := make([]entity.FileProducer, 0, len(files))
	for _, f := range files {
		file := entity.DataFile{
			FileInfo: entity.NewFileInfo(f.Path),
			Data:     *f.Data,
		}
		producer := exec.NewDummyProducer(file)
		producers = append(producers, producer)
	}

	strategies := []entity.FileStrategy{exec.NewTemplateFileStrategy(ff.templateData, entity.TemplateFnsMap, ff.templateOptions)}

	switch {
	case dryRun:
		strategies = append(strategies, exec.NewDryRunFileStrategy(logger))
	default:
		strategies = append(strategies, exec.NewSaveFileStrategy(logger))
	}
	executor := exec.NewFilesExecutor(producers, strategies)

	return executor, nil
}

type PreprocessorsFileExecutorFactory struct {
	templateData    map[string]any
	templateOptions []string

	preprocess         bool
	preprocessors      *exec.Preprocessors
	httpClientSupplier func(logger entity.Logger) *resty.Client
}

func NewPreprocessorsFileExecutorFactory(
	templateData map[string]any,
	templateOptions []string,
	preprocess bool,
	preprocessors *exec.Preprocessors,
	httpClientSupplier func(logger entity.Logger) *resty.Client,
) *PreprocessorsFileExecutorFactory {
	return &PreprocessorsFileExecutorFactory{
		templateData:       templateData,
		templateOptions:    templateOptions,
		preprocess:         preprocess,
		preprocessors:      preprocessors,
		httpClientSupplier: httpClientSupplier,
	}
}

func (ff *PreprocessorsFileExecutorFactory) Create(files []entity.UndefinedFile, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(files) == 0 {
		logger.Infof("`files` section is empty")
		return nil, nil
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
				Data:     *f.Data,
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
				client = ff.httpClientSupplier(logger)
			}

			producer = exec.NewRemoteProducer(file, client)
		case f.Local != nil:
			file := entity.LocalFile{
				FileInfo:  tmpl,
				LocalPath: *f.Local,
			}
			producer = exec.NewLocalProducer(file)

		default:
			return nil, xerrors.Errorf("build file executor from config: one of `data`, `get`, `local` must not be empty")
		}

		producers = append(producers, producer)
	}

	if preprocess := ff.preprocess; preprocess {
		if ff.preprocessors == nil {
			return nil, xerrors.Errorf("creating file executor - preprocesso is nil [preprocess:%v]", preprocess)
		}
		preloadProducers := make([]entity.FileProducer, 0, len(producers))
		preloader := exec.NewPreloadProducer(producers, logger)
		for i := 0; i < len(producers); i++ {
			preloadProducers = append(preloadProducers, preloader)
		}
		producers = preloadProducers
		ff.preprocessors.Add(preloader)
	}

	strategies := []entity.FileStrategy{exec.NewTemplateFileStrategy(ff.templateData, entity.TemplateFnsMap, ff.templateOptions)}

	switch {
	case dryRun:
		strategies = append(strategies, exec.NewDryRunFileStrategy(logger))
	default:
		strategies = append(strategies, exec.NewSaveFileStrategy(logger))
	}
	executor := exec.NewFilesExecutor(producers, strategies)

	return executor, nil
}

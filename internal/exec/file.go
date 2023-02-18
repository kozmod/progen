package exec

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"sync"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"

	"github.com/kozmod/progen/internal/entity"
)

type FilesExecutor struct {
	producers  []entity.FileProducer
	processors []entity.FileProc
}

func NewFilesExecutor(producers []entity.FileProducer, processors []entity.FileProc) *FilesExecutor {
	return &FilesExecutor{
		producers:  producers,
		processors: processors,
	}
}

func (e *FilesExecutor) Exec() error {
	for _, producer := range e.producers {
		file, err := producer.Get()
		if err != nil {
			return fmt.Errorf("execute file: get file: %w", err)
		}

		for _, processor := range e.processors {
			file, err = processor.Process(file)
			if err != nil {
				return fmt.Errorf("execute file: process file: %w", err)
			}
		}
	}
	return nil
}

type TemplateFileProc struct {
	templateProcFn func() entity.TemplateProc
}

func NewTemplateFileProc(templateData, templateFns map[string]any, templateOptions []string) *TemplateFileProc {
	return &TemplateFileProc{
		templateProcFn: func() entity.TemplateProc {
			return entity.NewTemplateProc(templateData, templateFns, templateOptions)
		},
	}
}

func (p *TemplateFileProc) Process(file entity.DataFile) (entity.DataFile, error) {
	filePath := file.Path()

	data, err := p.templateProcFn().Process(filePath, string(file.Data))
	if err != nil {
		return file, fmt.Errorf("process file template: %w", err)
	}
	file.Data = []byte(data)
	return file, nil
}

type SaveFileProc struct {
	fileMode os.FileMode
	logger   entity.Logger
}

func NewSaveFileProc(logger entity.Logger) *SaveFileProc {
	return &SaveFileProc{
		fileMode: os.ModePerm,
		logger:   logger,
	}
}

func (p *SaveFileProc) Process(file entity.DataFile) (entity.DataFile, error) {
	fileDir := file.Dir()
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		err = os.MkdirAll(fileDir, p.fileMode)
		if err != nil {
			return file, fmt.Errorf("save file: create file dir [%s]: %w", fileDir, err)
		}
	}

	filePath := file.Path()
	err := os.WriteFile(filePath, file.Data, p.fileMode)
	if err != nil {
		return file, fmt.Errorf("save file: write file [%s]: %w", file.Name(), err)
	}
	p.logger.Infof("file saved: %s", filePath)
	return file, nil
}

type ReplacePathFileProc struct {
	// paths contains old and new Files path to replace
	paths map[string]string
}

func NewReplacePathFileProc(paths map[string]string) *ReplacePathFileProc {
	return &ReplacePathFileProc{
		paths: paths,
	}
}

func (p *ReplacePathFileProc) Process(file entity.DataFile) (entity.DataFile, error) {
	newPath, ok := p.paths[file.Path()]
	if ok {
		file = entity.DataFile{
			FileInfo: entity.NewFileInfo(newPath),
			Data:     file.Data,
		}
	}
	return file, nil
}

type DryRunFileProc struct {
	logger entity.Logger
}

func NewDryRunFileProc(logger entity.Logger) *DryRunFileProc {
	return &DryRunFileProc{
		logger: logger,
	}
}

func (p *DryRunFileProc) Process(file entity.DataFile) (entity.DataFile, error) {
	fileDir := file.Dir()
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		p.logger.Infof("save file: create dir [%s] to store file [%s]", fileDir, file.Name)
	}

	filePath := file.Path()
	p.logger.Infof("file saved [path: %s]:\n%s", filePath, string(file.Data))
	return file, nil
}

type PreloadProducer struct {
	mx        sync.Mutex
	producers []entity.FileProducer
	logger    entity.Logger
}

func NewPreloadProducer(producers []entity.FileProducer, logger entity.Logger) *PreloadProducer {
	return &PreloadProducer{
		producers: producers,
		logger:    logger,
	}
}

func (p *PreloadProducer) Process() error {
	type OrderedFile struct {
		index int
		entity.DataFile
	}

	p.mx.Lock()
	defer p.mx.Unlock()

	var (
		files = make([]OrderedFile, 0, len(p.producers))
		fChan = make(chan OrderedFile, len(p.producers))

		eg, ctx = errgroup.WithContext(context.Background())
	)

	for i, p := range p.producers {
		index, producer := i, p
		eg.Go(func() error {
			if ctx.Err() != nil {
				return nil
			}
			file, err := producer.Get()
			if err != nil {
				return fmt.Errorf("preload file [%d]: %w", index, err)
			}
			fChan <- OrderedFile{index: index, DataFile: file}
			return nil
		})
	}

	err := eg.Wait()
	close(fChan)
	if err != nil {
		return err
	}

	for file := range fChan {
		files = append(files, file)
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].index < files[j].index
	})

	dummyProducers := make([]entity.FileProducer, 0, len(files))
	for _, file := range files {
		dummyProducers = append(dummyProducers, NewDummyProducer(file.DataFile))
		p.logger.Infof("file process: %s", file.Path())
	}

	p.producers = dummyProducers
	return nil
}

func (p *PreloadProducer) Get() (entity.DataFile, error) {
	p.mx.Lock()

	if len(p.producers) == 0 {
		return entity.DataFile{}, fmt.Errorf("process files list is empty")
	}
	producer := p.producers[0]
	p.producers = p.producers[1:]
	p.mx.Unlock()

	file, err := producer.Get()
	if err != nil {
		return entity.DataFile{}, fmt.Errorf("process files: get: %w", err)
	}

	return file, nil
}

type DummyProducer struct {
	file entity.DataFile
}

func NewDummyProducer(file entity.DataFile) *DummyProducer {
	return &DummyProducer{
		file: file,
	}
}

func (p *DummyProducer) Get() (entity.DataFile, error) {
	return p.file, nil
}

type LocalProducer struct {
	file entity.LocalFile
}

func NewLocalProducer(file entity.LocalFile) *LocalProducer {
	return &LocalProducer{
		file: file,
	}
}

func (p *LocalProducer) Get() (entity.DataFile, error) {
	data, err := os.ReadFile(p.file.LocalPath)
	if err != nil {
		return entity.DataFile{}, fmt.Errorf("read local: %w", err)
	}
	return entity.DataFile{
		FileInfo: p.file.FileInfo,
		Data:     data,
	}, nil
}

type RemoteProducer struct {
	client *resty.Client
	file   entity.RemoteFile
}

func NewRemoteProducer(file entity.RemoteFile, client *resty.Client) *RemoteProducer {
	return &RemoteProducer{
		file:   file,
		client: client,
	}
}

func (p *RemoteProducer) Get() (entity.DataFile, error) {
	var (
		url = p.file.URL
	)

	rq := p.client.R().
		SetHeaders(p.file.Headers).
		SetQueryParams(p.file.QueryParams)

	rs, err := rq.Get(url)
	if err != nil {
		return entity.DataFile{}, fmt.Errorf("get [%s]: %w", url, err)
	}

	statusCode := rs.StatusCode()
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		return entity.DataFile{
			FileInfo: p.file.FileInfo,
			Data:     rs.Body(),
		}, nil

	}
	return entity.DataFile{}, fmt.Errorf("get [%s]: ststus [%d]: response status is not in the 2xx range", url, statusCode)
}

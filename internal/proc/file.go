package proc

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/kozmod/progen/internal/entity"
)

type SaveFilesProc struct {
	fileMode  os.FileMode
	producers []entity.FileProducer
	logger    entity.Logger
}

func NewSaveFilesProc(producers []entity.FileProducer, logger entity.Logger) *SaveFilesProc {
	return &SaveFilesProc{
		fileMode:  os.ModePerm,
		producers: producers,

		logger: logger,
	}
}

func (p *SaveFilesProc) Exec() error {
	for _, producer := range p.producers {
		file, err := producer.Get()
		if err != nil {
			return fmt.Errorf("save file: get file to write: %w", err)
		}

		fileDir := file.Path
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			err = os.MkdirAll(fileDir, p.fileMode)
			if err != nil {
				return fmt.Errorf("save file: create file dir [%s]: %w", fileDir, err)
			}
		}

		filePath := path.Join(file.Path, file.Name)
		err = os.WriteFile(filePath, file.Data, p.fileMode)
		if err != nil {
			return fmt.Errorf("save file: create file [%s]: %w", file.Name, err)
		}
		p.logger.Infof("file saved [template: %v]: %s", file.ExecTmpl, filePath)
	}
	return nil
}

type DryRunFilesProc struct {
	fileMode  os.FileMode
	producers []entity.FileProducer
	logger    entity.Logger
	executor  templateExecutor
}

func NewDryRunFilesProc(producers []entity.FileProducer, templateData map[string]any, logger entity.Logger) *DryRunFilesProc {
	return &DryRunFilesProc{
		fileMode:  os.ModePerm,
		producers: producers,
		executor:  templateExecutor{templateData: templateData},
		logger:    logger,
	}
}

func (p *DryRunFilesProc) Exec() error {
	for _, producer := range p.producers {
		file, err := producer.Get()
		if err != nil {
			return fmt.Errorf("process file: get file to write: %w", err)
		}

		fileDir := file.Path
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			p.logger.Infof("process file: create dir [%s] to store file [%s]", fileDir, file.Name)
		}

		filePath := path.Join(file.Path, file.Name)
		p.logger.Infof("file created [template: %v, path: %s]:\n%s", file.ExecTmpl, filePath, string(file.Data))
	}
	return nil
}

type FilesPreProc struct {
	mx        sync.Mutex
	producers []entity.FileProducer
}

func NewFilesPreProc(producers []entity.FileProducer) *FilesPreProc {
	return &FilesPreProc{
		producers: producers,
	}
}

func (p *FilesPreProc) Exec() error {
	p.mx.Lock()
	defer p.mx.Unlock()

	processed := make([]entity.FileProducer, 0, len(p.producers))
	for i, producer := range p.producers {
		file, err := producer.Get()
		if err != nil {
			return fmt.Errorf("pre process file [%d]: %w", i, err)
		}
		processed = append(processed, &StoredProducer{file: file})
	}
	p.producers = processed
	return nil
}

func (p *FilesPreProc) Get() (entity.DataFile, error) {
	p.mx.Lock()
	defer p.mx.Unlock()
	if len(p.producers) == 0 {
		return entity.DataFile{}, fmt.Errorf("get preprocessed: empty")
	}
	producer := p.producers[0]
	p.producers = p.producers[1:]

	return producer.Get()
}

type TemplateExecProducerDecorator struct {
	producer entity.FileProducer
	executor templateExecutor
}

func NewTemplateExecProducer(producer entity.FileProducer, templateData, templateFns map[string]any) *TemplateExecProducerDecorator {
	return &TemplateExecProducerDecorator{
		producer: producer,
		executor: templateExecutor{
			templateData: templateData,
			templateFns:  templateFns,
		},
	}
}

func (p *TemplateExecProducerDecorator) Get() (entity.DataFile, error) {
	file, err := p.producer.Get()
	if err != nil {
		return file, fmt.Errorf("process file: get file to write: %w", err)
	}

	filePath := path.Join(file.Path, file.Name)
	if file.ExecTmpl {
		data, err := p.executor.Exec(filePath, file.Data)
		if err != nil {
			return file, fmt.Errorf("process process: %w", err)
		}
		file.Data = data
	}
	return file, nil
}

type StoredProducer struct {
	file entity.DataFile
}

func NewStoredProducer(file entity.DataFile) *StoredProducer {
	return &StoredProducer{
		file: file,
	}
}

func (p *StoredProducer) Get() (entity.DataFile, error) {
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
		Template: p.file.Template,
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
			Template: p.file.Template,
			Data:     rs.Body(),
		}, nil

	}
	return entity.DataFile{}, fmt.Errorf("get [%s]: ststus [%d]: response status is not in the 2xx range", url, statusCode)
}

package proc

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/go-resty/resty/v2"

	"github.com/kozmod/progen/internal/entity"
)

type FileProc struct {
	defaultPerm  os.FileMode
	producers    []entity.FileProducer
	logger       entity.Logger
	templateData map[string]any
	executor     templateExecutor
}

func NewFileProc(producers []entity.FileProducer, templateData map[string]any, logger entity.Logger) *FileProc {
	return &FileProc{
		defaultPerm:  os.ModePerm,
		producers:    producers,
		templateData: templateData,
		executor:     templateExecutor{templateData: templateData},
		logger:       logger,
	}
}

func (p *FileProc) Exec() error {
	for _, producer := range p.producers {
		file, err := producer.Get()
		if err != nil {
			return fmt.Errorf("process file: get file to write: %w", err)
		}

		fileDir := file.Path
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			err = os.MkdirAll(fileDir, p.defaultPerm)
			if err != nil {
				return fmt.Errorf("process file: create file dir [%s]: %w", fileDir, err)
			}
		}

		filePath := path.Join(file.Path, file.Name)

		if file.ExecTmpl {
			if file.ExecTmpl {
				data, err := p.executor.Exec(filePath, file.Data)
				if err != nil {
					return fmt.Errorf("process file: %w", err)
				}
				file.Data = data
			}
		}

		err = os.WriteFile(filePath, file.Data, file.Perm)
		if err != nil {
			return fmt.Errorf("process file: create file [%s]: %w", file.Name, err)
		}
		p.logger.Infof("file created [template: %v, perm: %v]: %s", file.ExecTmpl, file.Perm, filePath)
	}
	return nil
}

type DryRunFileProc struct {
	fileMode  os.FileMode
	producers []entity.FileProducer
	logger    entity.Logger
	executor  templateExecutor
}

func NewDryRunFileProc(producers []entity.FileProducer, templateData map[string]any, logger entity.Logger) *DryRunFileProc {
	return &DryRunFileProc{
		fileMode:  os.ModePerm,
		producers: producers,
		executor:  templateExecutor{templateData: templateData},
		logger:    logger,
	}
}

func (p *DryRunFileProc) Exec() error {
	for _, producer := range p.producers {
		file, err := producer.Get()
		if err != nil {
			return fmt.Errorf("process file: get file to write: %w", err)
		}

		fileDir := file.Path
		if _, err = os.Stat(fileDir); os.IsNotExist(err) {
			p.logger.Infof("process file: create dir to store file [%s]: %s", file.Name, fileDir)
		}

		filePath := path.Join(file.Path, file.Name)

		if file.ExecTmpl {
			data, err := p.executor.Exec(filePath, file.Data)
			if err != nil {
				return fmt.Errorf("process file: %w", err)
			}
			file.Data = data
		}
		p.logger.Infof("file created [template: %v, perm: %v]: %s\n%s", file.ExecTmpl, file.Perm, filePath, string(file.Data))
	}
	return nil
}

type StoredProducer struct {
	file *entity.DataFile
}

func NewStoredProducer(file entity.DataFile) *StoredProducer {
	return &StoredProducer{
		file: &file,
	}
}

func (p *StoredProducer) Get() (*entity.DataFile, error) {
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

func (p *LocalProducer) Get() (*entity.DataFile, error) {
	data, err := os.ReadFile(p.file.LocalPath)
	if err != nil {
		return nil, fmt.Errorf("read local: %w", err)
	}
	return &entity.DataFile{
		File:     p.file.File,
		Data:     data,
		ExecTmpl: p.file.ExecTmpl,
	}, nil
}

type RemoteProducer struct {
	client *resty.Client
	file   *entity.RemoteFile
}

func NewRemoteProducer(file entity.RemoteFile, client *resty.Client) *RemoteProducer {
	return &RemoteProducer{
		file:   &file,
		client: client,
	}
}

func (p *RemoteProducer) Get() (*entity.DataFile, error) {
	var (
		url = p.file.URL
	)

	rq := p.client.R().SetHeaders(p.file.Headers)

	rs, err := rq.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get [%s]: %w", url, err)
	}

	statusCode := rs.StatusCode()
	if statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices {
		return &entity.DataFile{
			File:     p.file.File,
			Data:     rs.Body(),
			ExecTmpl: p.file.ExecTmpl,
		}, nil

	}
	return nil, fmt.Errorf("get [%s]: ststus [%d]: response status is not in the 2xx range", url, statusCode)
}

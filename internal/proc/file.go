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
	fileMode  os.FileMode
	producers []entity.FileProducer
	logger    entity.Logger
}

func NewFileProc(producers []entity.FileProducer, logger entity.Logger) *FileProc {
	return &FileProc{
		fileMode:  os.ModePerm,
		producers: producers,
		logger:    logger,
	}
}

func (p *FileProc) Exec() error {
	for _, producer := range p.producers {
		file, err := producer.Get()
		if err != nil {
			return fmt.Errorf("get file to write: %w", err)
		}

		fileDir := file.Path
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			err = os.MkdirAll(fileDir, p.fileMode)
			if err != nil {
				return fmt.Errorf("create template dir [%s]: %w", fileDir, err)
			}
		}

		filePath := path.Join(file.Path, file.Name)
		err = os.WriteFile(filePath, file.Data, p.fileMode)
		if err != nil {
			return fmt.Errorf("create file [%s]: %w", file.Name, err)
		}
		p.logger.Infof("file created: %s", filePath)
	}
	return nil
}

type StoredProducer struct {
	file *entity.File
}

func NewStoredProducer(file entity.File) *StoredProducer {
	return &StoredProducer{
		file: &file,
	}
}

func (p *StoredProducer) Get() (*entity.File, error) {
	return p.file, nil
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

func (p *RemoteProducer) Get() (*entity.File, error) {
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
		return &entity.File{
			Name: p.file.Name,
			Path: p.file.Path,
			Data: rs.Body(),
		}, nil

	}
	return nil, fmt.Errorf("get [%s]: ststus [%d]: response status is not in the 2xx range", url, statusCode)
}

package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kozmod/progen/internal/entity"
)

type Reader struct {
	reader io.Reader
	path   string
}

func NewConfigReader(f entity.Flags) *Reader {
	if f.ReadStdin {
		return &Reader{
			reader: os.Stdin,
		}
	}
	return &Reader{
		path: f.ConfigPath,
	}

}

func (r *Reader) Read() ([]byte, error) {
	if r.reader == nil {
		data, err := os.ReadFile(r.path)
		if err != nil {
			return nil, fmt.Errorf("read config from file: %w", err)
		}
		return data, nil
	}

	reader := bufio.NewReader(r.reader)

	var data []byte
	for {
		line, err := reader.ReadBytes(entity.NewLine)
		switch {
		case errors.Is(err, io.EOF):
			return data, nil
		case err != nil:
			return nil, fmt.Errorf("read config stdin: %w", err)
		default:
			data = append(data, line...)
		}
	}
}
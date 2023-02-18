package config

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/kozmod/progen/internal/flag"

	"github.com/kozmod/progen/internal/entity"
)

var (
	newLIne = entity.NewLine[0]
)

type Reader struct {
	reader io.Reader
	path   string
}

func NewConfigReader(f flag.Flags) *Reader {
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
			return nil, fmt.Errorf("config file: %w", err)
		}
		return data, nil
	}

	reader := bufio.NewReader(r.reader)

	var data []byte
	for {
		line, err := reader.ReadBytes(newLIne)
		switch {
		case errors.Is(err, io.EOF):
			return data, nil
		case err != nil:
			return nil, fmt.Errorf("stdin: %w", err)
		default:
			data = append(data, line...)
		}
	}
}

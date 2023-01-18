package config

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/kozmod/progen/internal/entity"
)

const (
	TagDirs  = "dirs"
	TagFiles = "files"
	TagCmd   = "cmd"
)

type Config struct {
	HTTP  *HTTPClient `yaml:"http"`
	Dirs  []string    `yaml:"dirs,flow"`
	Files []File      `yaml:"files,flow"`
	Cmd   []string    `yaml:"cmd,flow"`
}

type HTTPClient struct {
	BaseURL AddrURL           `yaml:"base_url"`
	Headers map[string]string `yaml:"headers"`
	Debug   bool              `yaml:"debug"`
}

type File struct {
	Path         string  `yaml:"path"`
	Data         *string `yaml:"data"`
	Get          *Get    `yaml:"get"`
	Local        *string `yaml:"local"`
	ExecTmplSkip bool    `yaml:"tmpl_skip"`
}

type Get struct {
	Headers map[string]string `yaml:"headers"`
	URL     string            `yaml:"url"`
}

type AddrURL struct {
	*url.URL `yaml:",inline"`
}

func (addr *AddrURL) UnmarshalYAML(unmarshal func(any) error) error {
	var raw string
	if err := unmarshal(&raw); err != nil {
		return err
	}
	u, err := url.Parse(raw)
	if err != nil {
		return fmt.Errorf("unmarshal url: %w", err)
	}
	addr.URL = u
	return nil
}

func UnmarshalYamlConfig(in []byte) (Config, error) {
	var conf Config
	err := yaml.Unmarshal(in, &conf)
	if err != nil {
		return conf, err
	}

	for i, file := range conf.Files {
		err = ValidateFile(file)
		if err != nil {
			return conf, fmt.Errorf("validate config: files: %d [%s]: %w", i, file.Path, err)
		}
	}

	return conf, nil
}

func UnmarshalYamlFiles(in []byte) ([]File, error) {
	var files struct {
		Files []File `yaml:"files,flow"`
	}
	err := yaml.Unmarshal(in, &files)
	if err != nil {
		return nil, err
	}
	return files.Files, nil
}

func ValidateFile(file File) error {
	notNil := notNilValues(file.Get, file.Data, file.Local)
	switch {
	case notNil == 0:
		return fmt.Errorf("files: `get`, `data`, `local` tags - only one can be present")
	case notNil > 1:
		return fmt.Errorf("files: `get`, `data`, `local` - all are empty")
	case strings.TrimSpace(file.Path) == entity.Empty:
		return fmt.Errorf("files: save `path` are empty")
	}
	return nil
}

func notNilValues(values ...any) int {
	counter := 0
	for _, value := range values {
		val := reflect.ValueOf(value)
		if val.Kind() == reflect.Ptr && val.IsNil() {
			continue
		}
		counter++
	}
	return counter
}

package config

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strconv"
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
	Files []File      `yaml:"files,flow"`
	Cmd   []string    `yaml:"cmd,flow"`
	Dirs  []Dir       `yaml:"dirs,flow"`
}

type HTTPClient struct {
	BaseURL AddrURL           `yaml:"base_url"`
	Headers map[string]string `yaml:"headers"`
	Debug   bool              `yaml:"debug"`
}

type Dir struct {
	Path string `yaml:"path"`
	Perm *Perm  `yaml:"perm"`
}

func (dir *Dir) UnmarshalYAML(unmarshal func(any) error) error {
	var raw string
	err := unmarshal(&raw)
	if err == nil {
		*dir = Dir{
			Perm: nil,
			Path: raw,
		}
		return nil
	}

	if _, ok := err.(*yaml.TypeError); !ok {
		return fmt.Errorf("unmarshal dir: string: %w", err)
	}

	type alias Dir
	var val alias
	err = unmarshal(&val)
	if err != nil {
		return fmt.Errorf("unmarshal dir: struct: %w", err)
	}
	*dir = Dir(val)
	return nil
}

type File struct {
	Path         string  `yaml:"path"`
	Data         *string `yaml:"data"`
	Get          *Get    `yaml:"get"`
	Local        *string `yaml:"local"`
	ExecTmplSkip bool    `yaml:"tmpl_skip"`
	Perm         *Perm   `yaml:"perm"`
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

type Perm struct {
	os.FileMode `yaml:",inline"`
}

func (perm *Perm) UnmarshalYAML(unmarshal func(any) error) error {
	var raw string
	if err := unmarshal(&raw); err != nil {
		return fmt.Errorf("unmarshal perm: %w", err)
	}
	fm, err := strconv.ParseUint(raw, 0, 32)
	if err != nil {
		return fmt.Errorf("unmarshal perm unit: %w", err)
	}

	*perm = Perm{FileMode: os.FileMode(fm)}
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

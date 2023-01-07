package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kozmod/progen/internal/entity"

	"gopkg.in/yaml.v3"
)

const (
	TagDirs  = "dirs"
	TagFiles = "files"
	TagCmd   = "cmd"
)

type Config struct {
	Dirs  []string `yaml:"dirs,flow"`
	Files []File   `yaml:"files,flow"`
	Cmd   []string `yaml:"cmd,flow"`
}

type File struct {
	Path string  `yaml:"path"`
	Data *string `yaml:"data"`
	Get  *Get    `yaml:"get"`
}

type Get struct {
	Headers map[string]string `yaml:"headers"`
	URL     AddrURL           `yaml:"url"`
}

type AddrURL struct {
	*url.URL
}

func (addr *AddrURL) UnmarshalYAML(unmarshal func(interface{}) error) error {
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
		return conf, fmt.Errorf("unmarshal config: %w", err)
	}

	for i, file := range conf.Files {
		err = ValidateFile(file)
		if err != nil {
			return conf, fmt.Errorf("unmarshal config: validate files: %d [%s]: %w", i, file.Path, err)
		}
	}

	return conf, nil
}

func ValidateFile(file File) error {
	switch {
	case file.Get != nil && file.Data != nil:
		return fmt.Errorf("files: `get` and `data` tags can not be present both")
	case file.Get == nil && file.Data == nil:
		return fmt.Errorf("files: `get` and `data` are empty both")
	case strings.TrimSpace(file.Path) == entity.Empty:
		return fmt.Errorf("files: `path` is empty")
	}
	return nil
}

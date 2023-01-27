package config

import (
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/kozmod/progen/internal/entity"
)

const (
	TagDirs      = "dirs"
	TagFiles     = "files"
	TagCmd       = "cmd"
	SettingsHTTP = "settings"
)

type Config struct {
	Settings Settings            `yaml:"settings"`
	Dirs     []Section[[]string] `yaml:"dirs,flow"`
	Files    []Section[[]File]   `yaml:"files,flow"`
	Cmd      []Section[[]string] `yaml:"cmd,flow"`
}

type Settings struct {
	HTTP *HTTPClient `yaml:"http"`
}

type HTTPClient struct {
	HTTPClientParams `yaml:",inline"`
	BaseURL          AddrURL `yaml:"base_url"`
	Debug            bool    `yaml:"debug"`
}

type Section[T any] struct {
	Line int32
	Tag  string
	Val  T
}

type File struct {
	Path         string  `yaml:"path"`
	Data         *string `yaml:"data"`
	Get          *Get    `yaml:"get"`
	Local        *string `yaml:"local"`
	ExecTmplSkip bool    `yaml:"tmpl_skip"`
}

type Get struct {
	HTTPClientParams `yaml:",inline"`
	URL              string `yaml:"url"`
}

type HTTPClientParams struct {
	Headers     map[string]string `yaml:"headers"`
	QueryParams map[string]string `yaml:"query_params"`
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

func UnmarshalYamlFiles(rawConfig []byte) (map[string][]File, error) {
	var (
		fileByIndex = make(map[string][]File)
		rootTags    map[string]yaml.Node
	)

	err := yaml.Unmarshal(rawConfig, &rootTags)
	if err != nil {
		return nil, fmt.Errorf("unmarshal url: %w", err)
	}

	for tag, node := range rootTags {
		switch {
		case strings.Index(tag, TagFiles) == 0:
			var files []File
			err = node.Decode(&files)
			if err != nil {
				return nil, fmt.Errorf("unmarshal files config [%s]: %w", tag, err)
			}
			for _, file := range files {
				err = ValidateFile(file)
				if err != nil {
					return nil, fmt.Errorf("validate config: files [tag: %s, path: %s]: %w", tag, file.Path, err)
				}
			}
			fileByIndex[tag] = files
		}
	}
	return fileByIndex, nil
}

func ValidateFile(file File) error {
	notNil := entity.NotNilValues(file.Get, file.Data, file.Local)
	switch {
	case notNil == 0:
		return fmt.Errorf("files: `get`, `data`, `local` - all are empty")
	case notNil > 1:
		return fmt.Errorf("files: sections `get`, `data`, `local` tags - only one can be present")
	case strings.TrimSpace(file.Path) == entity.Empty:
		return fmt.Errorf("files: save `path` are empty")
	}
	return nil
}

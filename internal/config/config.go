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
	TagHTTP  = "http"
)

type Config struct {
	HTTP  *HTTPClient         `yaml:"http"`
	Dirs  []Section[[]string] `yaml:"dirs,flow"`
	Files []Section[[]File]   `yaml:"files,flow"`
	Cmd   []Section[[]string] `yaml:"cmd,flow"`
}

type HTTPClient struct {
	BaseURL AddrURL           `yaml:"base_url"`
	Headers map[string]string `yaml:"headers"`
	Debug   bool              `yaml:"debug"`
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

func UnmarshalYamlConfig(rawConfig []byte) (Config, error) {
	var (
		conf     Config
		rootTags map[string]yaml.Node
	)

	if err := yaml.Unmarshal(rawConfig, &rootTags); err != nil {
		return conf, fmt.Errorf("unmarshal url: %w", err)
	}

	for tag, node := range rootTags {
		var err error

		switch {
		case tag == TagHTTP:
			var client HTTPClient
			err = node.Decode(&client)
			conf.HTTP = &client
		case strings.Index(tag, TagDirs) == 0:
			dir := Section[[]string]{Line: int32(node.Line), Tag: tag}
			err = node.Decode(&dir.Val)
			conf.Dirs = append(conf.Dirs, dir)
		case strings.Index(tag, TagFiles) == 0:
			files := Section[[]File]{Line: int32(node.Line), Tag: tag}
			err = node.Decode(&files.Val)
			conf.Files = append(conf.Files, files)
		case strings.Index(tag, TagCmd) == 0:
			cmd := Section[[]string]{Line: int32(node.Line), Tag: tag}
			err = node.Decode(&cmd.Val)
			conf.Cmd = append(conf.Cmd, cmd)
		}

		if err != nil {
			return conf, fmt.Errorf("unmarshal tag [%s]: %w", tag, err)
		}
	}

	for i, files := range conf.Files {
		for _, file := range files.Val {
			err := ValidateFile(file)
			if err != nil {
				return conf, fmt.Errorf("validate config: files: %d [%s]: %w", i, file.Path, err)
			}
		}
	}

	if files, dirs, cmd := len(conf.Files), len(conf.Dirs), len(conf.Cmd); files == 0 && dirs == 0 && cmd == 0 {
		return conf, fmt.Errorf(
			"validate config: config not contains executable actions [dirs: %d, files: %d, cms: %d]",
			dirs, files, cmd)
	}

	return conf, nil
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

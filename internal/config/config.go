package config

import (
	"fmt"

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
	Path string `yaml:"path"`
	Data string `yaml:"data"`
}

func UnmarshalYamlConfig(in []byte) (Config, error) {
	var conf Config
	err := yaml.Unmarshal(in, &conf)
	if err != nil {
		return conf, fmt.Errorf("config.unmarshalConfig: %w", err)
	}

	return conf, nil
}

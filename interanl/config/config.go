package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	tagDirs  = "dirs"
	tagFiles = "files"
	tagCmd   = "cmd"
)

type VarsConfig struct {
	Vars []string `yaml:"vars,flow"`
}

func UnmarshalYamlVarsConfig(data []byte) (VarsConfig, error) {
	var conf VarsConfig
	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return conf, fmt.Errorf("unmarshal config: %w", err)
	}

	return conf, nil
}

type Config struct {
	Dirs  []string `yaml:"dirs,flow"`
	Files []File   `yaml:"files,flow"`
	Cmd   []string `yaml:"cmd,flow"`
}

type File struct {
	Path     string `yaml:"path"`
	Template string `yaml:"template"`
}

func UnmarshalYamlConfig(in []byte) (Config, error) {
	var conf Config
	err := yaml.Unmarshal(in, &conf)
	if err != nil {
		return conf, fmt.Errorf("config.unmarshalConfig: %w", err)
	}

	return conf, nil
}

package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/kozmod/progen/internal/entity"
)

type YamlUnmarshaler struct {
	tagFilter *entity.RegexpChain
	logger    entity.Logger
}

func NewYamlConfigUnmarshaler(tagFilter *entity.RegexpChain, logger entity.Logger) *YamlUnmarshaler {
	return &YamlUnmarshaler{
		tagFilter: tagFilter,
		logger:    logger,
	}
}

func (u *YamlUnmarshaler) Unmarshal(rawConfig []byte) (Config, error) {
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
		case tag == SettingsHTTP:
			var settings Settings
			err = node.Decode(&settings)
			conf.Settings = settings
		case u.tagFilter != nil && u.tagFilter.MatchString(tag):
			u.logger.Infof("action tag will be skipped: %s", tag)
			continue
		case strings.Index(tag, TagDirs) == 0:
			conf.Dirs, err = decode(conf.Dirs, node, tag)
		case strings.Index(tag, TagFiles) == 0:
			conf.Files, err = decode(conf.Files, node, tag)
		case strings.Index(tag, TagCmd) == 0:
			conf.Cmd, err = decode(conf.Cmd, node, tag)
		case strings.Index(tag, TagFS) == 0:
			conf.FS, err = decode(conf.FS, node, tag)
		case strings.Index(tag, TagScripts) == 0:
			conf.Scripts, err = decode(conf.Scripts, node, tag)
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

	return conf, validateConfigSections(conf)
}

func decode[T any](target []Section[T], node yaml.Node, tag string) ([]Section[T], error) {
	section := Section[T]{Line: int32(node.Line), Tag: tag}
	err := node.Decode(&section.Val)
	target = append(target, section)
	return target, err
}

func validateConfigSections(conf Config) error {
	var (
		files   = len(conf.Files)
		dirs    = len(conf.Dirs)
		cmd     = len(conf.Cmd)
		fs      = len(conf.FS)
		scripts = len(conf.Scripts)
	)
	if files == 0 && dirs == 0 && cmd == 0 && fs == 0 && scripts == 0 {
		return fmt.Errorf(
			"validate config: config not contains executable actions [dirs: %d, files: %d, cms: %d, scripts: %d]",
			dirs, files, cmd, scripts)
	}
	return nil
}

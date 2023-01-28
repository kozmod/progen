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

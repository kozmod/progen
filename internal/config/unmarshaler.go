package config

import (
	"strings"

	"golang.org/x/xerrors"
	yaml "gopkg.in/yaml.v3"

	"github.com/kozmod/progen/internal/entity"
)

var (
	ErrCommandEmpty = xerrors.Errorf("command declaration is empty")
)

type YamlUnmarshaler struct{}

func NewYamlConfigUnmarshaler() *YamlUnmarshaler {
	return &YamlUnmarshaler{}
}

func (u *YamlUnmarshaler) Unmarshal(rawConfig []byte) (Config, error) {
	var (
		conf     Config
		rootTags map[string]yaml.Node
	)

	if err := yaml.Unmarshal(rawConfig, &rootTags); err != nil {
		return conf, xerrors.Errorf("unmarshal url: %w", err)
	}

	for tag, node := range rootTags {
		var err error

		switch {
		case tag == SettingsHTTP:
			var settings Settings
			err = node.Decode(&settings)
			conf.Settings = settings
		case strings.Index(tag, TagDirs) == 0:
			conf.Dirs, err = decode(conf.Dirs, node, tag)
		case strings.Index(tag, TagRm) == 0:
			conf.Rm, err = decode(conf.Rm, node, tag)
		case strings.Index(tag, TagFiles) == 0:
			conf.Files, err = decode(conf.Files, node, tag)
		case strings.Index(tag, TagCmd) == 0:
			conf.Cmd, err = decode(conf.Cmd, node, tag)
		case strings.Index(tag, TagFS) == 0:
			conf.FS, err = decode(conf.FS, node, tag)
		}

		if err != nil {
			return conf, xerrors.Errorf("unmarshal tag [%s]: %w", tag, err)
		}
	}

	return conf, nil
}

func decode[T any](target []Section[T], node yaml.Node, tag string) ([]Section[T], error) {
	section := Section[T]{Line: int32(node.Line), Tag: tag}
	err := node.Decode(&section.Val)
	target = append(target, section)
	return target, err
}

func commandFromString(cmd string) (Command, error) {
	var (
		splitCmd = strings.Split(cmd, entity.Space)
		command  = make([]string, 0, len(splitCmd))
	)

	for _, val := range splitCmd {
		if trimmed := strings.TrimSpace(val); val != entity.Empty {
			command = append(command, trimmed)
		}
	}

	if len(command) == 0 {
		return Command{}, ErrCommandEmpty
	}

	return Command{
		Exec: command[0],
		Args: command[1:],
	}, nil
}

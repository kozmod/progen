package config

import (
	"net/url"
	"strings"

	"golang.org/x/xerrors"

	"github.com/kozmod/progen/internal/entity"
)

const (
	TagDirs      = "dirs"
	TagFiles     = "files"
	TagCmd       = "cmd"
	TagFS        = "fs"
	SettingsHTTP = "settings"
)

type Config struct {
	Settings Settings             `yaml:"settings"`
	Dirs     []Section[[]string]  `yaml:"dirs,flow"`
	Files    []Section[[]File]    `yaml:"files,flow"`
	Cmd      []Section[[]Command] `yaml:"cmd,flow"`
	FS       []Section[[]string]  `yaml:"fs,flow"`
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
	Path  string  `yaml:"path"`
	Data  *string `yaml:"data"`
	Get   *Get    `yaml:"get"`
	Local *string `yaml:"local"`
}

type Command struct {
	Dir  string   `yaml:"dir"`
	Exec string   `yaml:"exec"`
	Args []string `yaml:"args,flow"`
}

func (c *Command) UnmarshalYAML(unmarshal func(any) error) error {
	var raw string
	if err := unmarshal(&raw); err == nil {
		*c, err = commandFromString(raw)
		if err != nil {
			return err
		}
		return nil
	}
	type alias Command
	var cmd alias
	if err := unmarshal(&cmd); err != nil {
		return err
	}
	*c = (Command)(cmd)
	return nil
}

type Script struct {
	Dir    string   `yaml:"dir"`
	Exec   string   `yaml:"exec"`
	Args   []string `yaml:"args,flow"`
	Script string   `yaml:"script,flow"`
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
		return xerrors.Errorf("unmarshal url: %w", err)
	}
	addr.URL = u
	return nil
}

func ValidateFile(file File) error {
	notNil := entity.NotNilValues(file.Get, file.Data, file.Local)
	switch {
	case notNil == 0:
		return xerrors.Errorf("files: `get`, `data`, `local` - all are empty")
	case notNil > 1:
		return xerrors.Errorf("files: sections `get`, `data`, `local` tags - only one can be present")
	case strings.TrimSpace(file.Path) == entity.Empty:
		return xerrors.Errorf("files: save `path` are empty")
	}
	return nil
}

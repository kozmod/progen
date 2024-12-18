package config

import (
	"net/url"
	"strings"

	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"

	"github.com/kozmod/progen/internal/entity"
)

const (
	TagDirs      = "dirs"
	TagRm        = "rm"
	TagFiles     = "files"
	TagCmd       = "cmd"
	TagFS        = "fs"
	SettingsHTTP = "settings"
)

type Config struct {
	Settings Settings             `yaml:"settings"`
	Dirs     []Section[[]string]  `yaml:"dirs,flow"`
	Rm       []Section[[]string]  `yaml:"rm,flow"`
	Files    []Section[[]File]    `yaml:"files,flow"`
	Cmd      []Section[[]Command] `yaml:"cmd,flow"`
	FS       []Section[[]string]  `yaml:"fs,flow"`
}

func (c Config) CommandActions() []entity.Action[[]entity.Command] {
	return toActionsSlice(c.Cmd, func(cmd Command) entity.Command {
		return entity.Command{
			Cmd:  cmd.Exec,
			Args: cmd.Args,
			Dir:  cmd.Dir,
		}
	})
}

func (c Config) FilesActions() []entity.Action[[]entity.UndefinedFile] {
	return toActionsSlice(c.Files, func(file File) entity.UndefinedFile {
		uFile := entity.UndefinedFile{
			Path: file.Path,
		}
		if file.Data != nil {
			*uFile.Data = *file.Data
		}
		if get := file.Get; get != nil {
			uFile.Get = &entity.HTTPClientParams{
				URL:         get.URL,
				Headers:     get.Headers,
				QueryParams: get.QueryParams,
			}
		}
		if file.Local != nil {
			uFile.Local = file.Local
		}
		return uFile
	})
}

func (c Config) DirActions() []entity.Action[[]string] {
	return toActionsSlice(c.Dirs, func(dir string) string {
		return dir
	})
}

func (c Config) RmActions() []entity.Action[[]string] {
	return toActionsSlice(c.Rm, func(rms string) string {
		return rms
	})
}

func (c Config) FsActions() []entity.Action[[]string] {
	return toActionsSlice(c.FS, func(fss string) string {
		return fss
	})
}

func toActionsSlice[S any, T any](sections []Section[[]S], mapFn func(s S) T) []entity.Action[[]T] {
	actions := make([]entity.Action[[]T], len(sections))
	for i, section := range sections {
		action := entity.Action[[]T]{
			Priority: section.Line,
			Name:     section.Tag,
			Val:      make([]T, len(section.Val)),
		}
		for j, val := range section.Val {
			action.Val[j] = mapFn(val)
		}
		actions[i] = action
	}
	return actions
}

type Settings struct {
	HTTP   *HTTPClient `yaml:"http"`
	Groups Groups      `yaml:"groups"`
}

type HTTPClient struct {
	HTTPClientParams `yaml:",inline"`
	BaseURL          AddrURL `yaml:"base_url"`
	Debug            bool    `yaml:"debug"`
}

type Groups []Group

func (g Groups) ManualActions() map[string]struct{} {
	manual := make(map[string]struct{})
	for _, group := range g {
		if !group.Manual {
			continue
		}
		for _, action := range group.Actions {
			manual[action] = struct{}{}
		}
	}
	return manual
}

func (g Groups) GroupByAction() map[string]map[string]struct{} {
	manual := make(map[string]map[string]struct{})
	for _, group := range g {
		for _, action := range group.Actions {
			groups, ok := manual[action]
			if !ok {
				groups = make(map[string]struct{})
			}
			groups[group.Name] = struct{}{}
			manual[action] = groups
		}
	}
	return manual
}

type Group struct {
	Name    string   `yaml:"name"`
	Actions []string `yaml:"actions"`
	Manual  bool     `yaml:"manual"`
}

type Section[T any] struct {
	Line int
	Tag  string
	Val  T
}

type File struct {
	Path  string  `yaml:"path"`
	Data  *Bytes  `yaml:"data"`
	Get   *Get    `yaml:"get"`
	Local *string `yaml:"local"`
}

type Bytes []byte

func (fd *Bytes) UnmarshalYAML(node *yaml.Node) error {
	*fd = Bytes(node.Value)
	return nil
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

func (c Config) Validate() error {
	for i, files := range c.Files {
		for _, file := range files.Val {
			err := validateFile(file)
			if err != nil {
				return xerrors.Errorf("files: %d [%s]: %w", i, file.Path, err)
			}
		}
	}

	if err := validateGroups(c.Settings.Groups); err != nil {
		return xerrors.Errorf("groups: %w", err)
	}

	if err := validateConfigSections(c); err != nil {
		return xerrors.Errorf("sections: %w", err)
	}
	return nil
}

func validateFile(file File) error {
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

func validateGroups(groups Groups) error {
	var (
		groupNameSet = make(map[string]int, len(groups))
		groupNames   = make([]string, 0, len(groups))
	)
	for _, group := range groups {
		var (
			name         = group.Name
			quantity, ok = groupNameSet[name]
		)
		if ok && quantity == 1 {
			groupNames = append(groupNames, name)
		}
		groupNameSet[name] = groupNameSet[name] + 1
	}

	if len(groupNames) > 0 {
		return xerrors.Errorf("duplicate names [%s]", strings.Join(groupNames, entity.LogSliceSep))
	}
	return nil
}

func validateConfigSections(conf Config) error {
	var (
		files = len(conf.Files)
		dirs  = len(conf.Dirs)
		rm    = len(conf.Rm)
		cmd   = len(conf.Cmd)
		fs    = len(conf.FS)
	)
	if files == 0 && dirs == 0 && rm == 0 && cmd == 0 && fs == 0 {
		return xerrors.Errorf(
			"config not contains executable actions [dirs: %d, rm: %d, files: %d, cms: %d, fs: %d]",
			dirs, rm, files, cmd, fs,
		)
	}
	return nil
}

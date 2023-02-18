package config

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/kozmod/progen/internal/entity"
)

func Test_NewRawPreprocessor_Process(t *testing.T) {
	t.Parallel()

	const (
		name = "conf"
	)

	t.Run("success_preprocess_raw_config_data", func(t *testing.T) {
		const (
			in = `
matrix:
 version: 1.19
steps:
 name: Setup Go {{ .matrix.version }}
`
			expected = `
matrix:
 version: 1.19
steps:
 name: Setup Go 1.19
`
		)

		rawConf, mapConf, err := NewRawPreprocessor(name, nil, nil, nil).Process([]byte(in))
		assert.NoError(t, err)
		assert.Equal(t, expected, string(rawConf))
		assert.NotEmpty(t, mapConf)
	})
	t.Run("success_process_with_custom_fn", func(t *testing.T) {
		const (
			in = `
steps:
 name:{{ random.Alpha 8 }}
`
			exp = `
steps:
 name:[a-z, A-z]{8}
`
		)

		res, _, err := NewRawPreprocessor(name, nil, entity.TemplateFnsMap, nil).Process([]byte(in))
		assert.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile(exp), string(res))
	})
	t.Run("success_process_with_custom_vars_map", func(t *testing.T) {
		const (
			in = `
steps:
 name:"{{.vars.service_name}}"
`
			exp = `
steps:
 name:"SOME"
`
		)

		res, _, err := NewRawPreprocessor(
			name,
			map[string]any{"vars": map[string]any{"service_name": "SOME"}},
			nil,
			nil).
			Process([]byte(in))
		assert.NoError(t, err)
		assert.Equal(t, exp, string(res))
	})
	t.Run("success_preprocess_raw_config_data_when_template_key_not_set", func(t *testing.T) {
		const (
			in = `
steps:
 name: Setup Go {{ .matrix.version }}
`
			expected = `
steps:
 name: Setup Go <no value>
`
		)

		res, _, err := NewRawPreprocessor(name, nil, nil, nil).Process([]byte(in))
		assert.NoError(t, err)
		assert.Equal(t, expected, string(res))
	})
	t.Run("error_preprocess_raw_config_data_when_template_missingkey_option_is_error", func(t *testing.T) {
		const (
			in = `
steps:
 name: Setup Go {{ .matrix.version }}
`
		)

		options := []string{fmt.Sprintf("%v=%v", entity.TemplateOptionsMissingKey, entity.MissingKeyError)}
		_, _, err := NewRawPreprocessor(name, nil, nil, options).Process([]byte(in))
		assert.Error(t, err)
	})
}

func Test_ValidateFile(t *testing.T) {
	t.Parallel()

	const (
		path = "some_path"
	)
	t.Run("not_error_when_get_is_not_nil", func(t *testing.T) {
		in := File{
			Path: path,
			Data: nil,
			Get:  &Get{},
		}
		err := ValidateFile(in)
		assert.NoError(t, err)
	})
	t.Run("not_error_when_data_is_not_nil", func(t *testing.T) {
		in := File{
			Path: path,
			Data: func(s string) *string { return &s }("some data"),
			Get:  nil,
		}
		err := ValidateFile(in)
		assert.NoError(t, err)
	})
	t.Run("error_when_data_and_get_are_nil", func(t *testing.T) {
		in := File{
			Path: path,
			Data: nil,
			Get:  nil,
		}
		err := ValidateFile(in)
		assert.Error(t, err)
	})
	t.Run("error_when_data_and_get_are_not_nil_both", func(t *testing.T) {
		in := File{
			Path: path,
			Data: func(s string) *string { return &s }("some data"),
			Get:  &Get{},
		}
		err := ValidateFile(in)
		assert.Error(t, err)
	})
	t.Run("error_when_path_is_empty", func(t *testing.T) {
		in := File{
			Path: "",
			Get:  &Get{},
		}
		err := ValidateFile(in)
		assert.Error(t, err)

		in.Path = "     "
		err = ValidateFile(in)
		assert.Error(t, err)
	})
}

func Test_Read(t *testing.T) {
	t.Parallel()

	t.Run("success_read_config_data", func(t *testing.T) {
		const (
			in = `
dirs1:
  - api/{{.vars.pn}}/v1
cmd1:
  - exec: [chmod -R 777 api]
    dir: .
dirs2:
  - api/{{.vars2.pn}}/v1
cmd2:
  - exec: [chmod -R 777 api]
    dir: .
`
		)

		var stdin bytes.Buffer
		stdin.Write([]byte(in))

		reader := Reader{reader: &stdin}
		b, err := reader.Read()
		assert.NoError(t, err)
		assert.Equal(t, in, string(b))
	})
	t.Run("success_read_empty_config", func(t *testing.T) {
		const (
			in = ``
		)

		var stdin bytes.Buffer
		stdin.Write([]byte(in))

		reader := Reader{reader: &stdin}
		b, err := reader.Read()
		assert.NoError(t, err)
		assert.Equal(t, in, string(b))
	})
}

func Test_YamlUnmarshaler_Unmarshal(t *testing.T) {
	t.Parallel()

	t.Run("success_unmarshal", func(t *testing.T) {
		const (
			in = `
cmd:
  - exec: pwd
    dir: ..
  - exec: ls
    args: [-l]
    dir: ..

dirs:
  - x/api/{{.vars.service_name}}/v1
  - s

dirs2:
  - y/api

cmd1:
  - ls -a

files:
  - path: x/DDDDDD
    tmpl_skip: true
    data: |
      ENV GOPROXY "{{.vars.GOPROXY}} ,proxy.golang.org,direct"

cmd2:
  - exec: ls
    args: [ -lS ]
    dir: /
  - exec: whoami
`
		)

		conf, err := NewYamlConfigUnmarshaler(nil, nil).Unmarshal([]byte(in))
		a := assert.New(t)
		a.NoError(err)
		a.Len(conf.Files, 1)
		a.Len(conf.Dirs, 2)
		a.Len(conf.Cmd, 3)

		assertSection(a, conf.Cmd, "cmd", 3,
			Command{Dir: "..", Exec: "pwd"},
			Command{Dir: "..", Exec: "ls", Args: []string{"-l"}},
		)
		assertSection(a, conf.Cmd, "cmd1", 17,
			Command{Dir: entity.Empty, Exec: "ls", Args: []string{"-a"}},
		)
		assertSection(a, conf.Cmd, "cmd2", 26,
			Command{Dir: "/", Exec: "ls", Args: []string{"-lS"}},
			Command{Exec: "whoami"})

		assertSection(a, conf.Dirs, "dirs", 10, "x/api/{{.vars.service_name}}/v1", "s")
		assertSection(a, conf.Dirs, "dirs2", 14, "y/api")

		files := conf.Files[0]
		a.Equal(int32(20), files.Line)
		a.Equal(TagFiles, files.Tag)
		a.Len(files.Val, 1)

		file := files.Val[0]
		a.Equal("x/DDDDDD", file.Path)
		a.NotNil(file.Data)
		a.Equal("ENV GOPROXY \"{{.vars.GOPROXY}} ,proxy.golang.org,direct\"\n", *file.Data)
	})

	t.Run("success_unmarshal_skip_all_cmd_and_dirs2_sections", func(t *testing.T) {
		const (
			in = `
cmd:
  - exec: [pwd]
    dir: /

dirs:
  - x/api/{{.vars.service_name}}/v1
  - s

dirs2:
  - y/api

cmd1:
  - exec: [ls -a]
    dir: .

files:
  - path: x/DDDDDD
    tmpl_skip: true
    data: |
      ENV GOPROXY "{{.vars.GOPROXY}} ,proxy.golang.org,direct"

cmd2:
  - exec: [ls -lS, whoami]
    dir:
`
		)

		var (
			a      = assert.New(t)
			filter = entity.NewRegexpChain("cmd.?", "dirs2")
			logger = MockLogger{
				infof: func(format string, args ...any) {
					a.Equal("action tag will be skipped: %s", format)
					for i, arg := range args {
						a.Containsf([]any{"cmd", "cmd1", "cmd2", "dirs2"}, arg, "contains_%d", i)
					}
				},
			}
		)
		conf, err := NewYamlConfigUnmarshaler(filter, logger).Unmarshal([]byte(in))
		a.NoError(err)
		a.Len(conf.Files, 1)
		a.Len(conf.Dirs, 1)
		a.Len(conf.Cmd, 0)

		assertSection(a, conf.Dirs, "dirs", 7, "x/api/{{.vars.service_name}}/v1", "s")

		files := conf.Files[0]
		a.Equal(int32(18), files.Line)
		a.Equal(TagFiles, files.Tag)
		a.Len(files.Val, 1)

		file := files.Val[0]
		a.Equal("x/DDDDDD", file.Path)
		a.NotNil(file.Data)
		a.Equal("ENV GOPROXY \"{{.vars.GOPROXY}} ,proxy.golang.org,direct\"\n", *file.Data)
	})
	t.Run("error_when-config_not_contains_executable_actions", func(t *testing.T) {
		const (
			in = `
var:
  some_1: val_1
`
		)
		_, err := NewYamlConfigUnmarshaler(nil, nil).Unmarshal([]byte(in))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config not contains executable actions")
	})
}

func Test_Command_UnmarshalYAML(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const (
			in = `
cmd:
 - exec: pwd
   dir: /
 - ls -a
`
		)
		var (
			a = assert.New(t)
		)
		conf, err := NewYamlConfigUnmarshaler(nil, nil).Unmarshal([]byte(in))
		a.NoError(err)
		assertSection(a, conf.Cmd, "cmd", 3,
			Command{Dir: "/", Exec: "pwd"},
			Command{Dir: entity.Empty, Exec: "ls", Args: []string{"-a"}},
		)
	})
	t.Run("success_v2", func(t *testing.T) {
		const (
			in = `
cmd:
 - ls -a
 - dh -f
 - exec: pwd
   dir: /
`
		)

		var (
			a       = assert.New(t)
			TestCmd struct {
				Commands []Command `yaml:"cmd,flow"`
			}
		)
		err := yaml.Unmarshal([]byte(in), &TestCmd)
		a.NoError(err)
		a.Equal([]Command{
			{Dir: entity.Empty, Exec: "ls", Args: []string{"-a"}},
			{Dir: entity.Empty, Exec: "dh", Args: []string{"-f"}},
			{Dir: "/", Exec: "pwd"},
		}, TestCmd.Commands)

	})
}

func Test_AddrURL_UnmarshalYAML(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		const (
			expUrl = "https://gitlab.sberlabs.com/api/v4/projects/5/repository/files/"
		)
		var (
			in         = fmt.Sprintf(`url: %s`, expUrl)
			a          = assert.New(t)
			TestAddUrl struct {
				URL AddrURL `yaml:"url"`
			}
		)
		err := yaml.Unmarshal([]byte(in), &TestAddUrl)
		a.NoError(err)
		a.Equal(expUrl, TestAddUrl.URL.String())
	})
	t.Run("error_when_url_is_invalid", func(t *testing.T) {
		const (
			expUrl = ":::gitlab.sberlabs.com/api/"
		)
		var (
			in         = fmt.Sprintf(`url: %s`, expUrl)
			a          = assert.New(t)
			TestAddUrl struct {
				URL AddrURL `yaml:"url"`
			}
		)
		err := yaml.Unmarshal([]byte(in), &TestAddUrl)
		a.Error(err)
	})
}

func assertSection[T any](a *assert.Assertions, sections []Section[[]T], tag string, line int32, expected ...T) {
	var (
		section Section[[]T]
		found   bool
	)

	for _, s := range sections {
		if s.Tag == tag {
			section = s
			found = true
			break
		}
	}
	a.Truef(found, "section not contains tag [%s]", tag)
	a.NotNil(section)
	a.Equal(line, section.Line)
	d := section.Val
	a.Equal(expected, d)
}

func Test_commandFromString(t *testing.T) {
	t.Parallel()

	const (
		command = "ls"
		arg     = "-a"
	)

	t.Run("success", func(t *testing.T) {
		var (
			commandStr = strings.Join(append([]string{command}, arg), entity.Space)
		)

		cmd, err := commandFromString(commandStr)
		assert.NoError(t, err)
		assert.Equal(t, Command{Dir: entity.Empty, Exec: command, Args: []string{arg}}, cmd)
	})
	t.Run("error_when_command_is_empty", func(t *testing.T) {
		cmd, err := commandFromString("")
		assert.Error(t, err)
		assert.NotNil(t, cmd)
		assert.ErrorIs(t, err, ErrCommandEmpty)
	})
}

type MockLogger struct {
	entity.Logger
	infof func(format string, args ...any)
}

func (m MockLogger) Infof(format string, args ...any) {
	m.infof(format, args...)
}

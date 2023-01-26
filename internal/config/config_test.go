package config

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_preprocessRawConfigData(t *testing.T) {
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

		rawConf, mapConf, err := PreprocessRawConfigData(name, []byte(in), nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(rawConf))
		assert.NotEmpty(t, mapConf)
	})

	t.Run("success_preprocess_raw_config_data", func(t *testing.T) {
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

		rawConf, mapConf, err := PreprocessRawConfigData(name, []byte(in), nil)
		assert.NoError(t, err)
		assert.Equal(t, expected, string(rawConf))
		assert.NotEmpty(t, mapConf)
	})
}

func Test_ValidateFile(t *testing.T) {
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

func Test_UnmarshalYamlConfig(t *testing.T) {
	t.Run("success_unmarshal", func(t *testing.T) {
		const (
			in = `
cmd:
  - pwd

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
  - ls -lS
  - whoami
`
		)

		conf, err := UnmarshalYamlConfig([]byte(in))
		a := assert.New(t)
		a.NoError(err)
		a.Len(conf.Files, 1)
		a.Len(conf.Dirs, 2)
		a.Len(conf.Cmd, 3)

		assertCmdFn := func(tag string, line int32, expected ...string) {
			var (
				section Section[[]string]
				found   bool
			)

			for _, cmd := range conf.Cmd {
				if cmd.Tag == tag {
					section = cmd
					found = true
					break
				}
			}
			a.True(found)
			a.NotNil(section)
			a.Equal(line, section.Line)
			a.Contains(section.Tag, TagCmd)
			commands := section.Val
			a.Equal(expected, commands)

		}

		assertCmdFn("cmd", 3, "pwd")
		assertCmdFn("cmd1", 13, "ls -a")
		assertCmdFn("cmd2", 22, "ls -lS", "whoami")

		assertDirsFn := func(tag string, line int32, expected ...string) {
			var (
				section Section[[]string]
				found   bool
			)

			for _, d := range conf.Dirs {
				if d.Tag == tag {
					section = d
					found = true
					break
				}
			}
			a.True(found)
			a.NotNil(section)
			a.Equal(line, section.Line)
			a.Contains(section.Tag, TagDirs)
			d := section.Val
			a.Equal(expected, d)

		}

		assertDirsFn("dirs", 6, "x/api/{{.vars.service_name}}/v1", "s")
		assertDirsFn("dirs2", 10, "y/api")

		files := conf.Files[0]
		a.Equal(int32(16), files.Line)
		a.Equal(TagFiles, files.Tag)
		a.Len(files.Val, 1)

		file := files.Val[0]
		a.Equal("x/DDDDDD", file.Path)
		a.NotNil(file.Data)
		a.Equal("ENV GOPROXY \"{{.vars.GOPROXY}} ,proxy.golang.org,direct\"\n", *file.Data)
		a.True(file.ExecTmplSkip)
	})
	t.Run("error_when-config_not_contains_executable_actions", func(t *testing.T) {
		const (
			in = `
var:
  some_1: val_1
`
		)
		_, err := UnmarshalYamlConfig([]byte(in))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config not contains executable actions")
	})
}

func Test_Read(t *testing.T) {
	t.Run("success_read_config_data", func(t *testing.T) {
		const (
			in = `
dirs1:
  - api/{{.vars.pn}}/v1
cmd1:
  - chmod -R 777 api
dirs2:
  - api/{{.vars2.pn}}/v1
cmd2:
  - chmod -R 777 api
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

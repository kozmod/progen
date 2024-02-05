package config

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"

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

func Test_validateFile(t *testing.T) {
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
		err := validateFile(in)
		assert.NoError(t, err)
	})
	t.Run("not_error_when_data_is_not_nil", func(t *testing.T) {
		in := File{
			Path: path,
			Data: func(d Bytes) *Bytes { return &d }(Bytes("some data")),
			Get:  nil,
		}
		err := validateFile(in)
		assert.NoError(t, err)
	})
	t.Run("error_when_data_and_get_are_nil", func(t *testing.T) {
		in := File{
			Path: path,
			Data: nil,
			Get:  nil,
		}
		err := validateFile(in)
		assert.Error(t, err)
	})
	t.Run("error_when_data_and_get_are_not_nil_both", func(t *testing.T) {
		in := File{
			Path: path,
			Data: func(d Bytes) *Bytes { return &d }(Bytes("some data")),
			Get:  &Get{},
		}
		err := validateFile(in)
		assert.Error(t, err)
	})
	t.Run("error_when_path_is_empty", func(t *testing.T) {
		in := File{
			Path: "",
			Get:  &Get{},
		}
		err := validateFile(in)
		assert.Error(t, err)

		in.Path = "     "
		err = validateFile(in)
		assert.Error(t, err)
	})
}

func Test_validateGroups(t *testing.T) {
	t.Parallel()

	const (
		groupA = "grpA"
		groupB = "grpB"
	)

	t.Run("not_error", func(t *testing.T) {
		in := Groups{
			{Name: groupA},
			{Name: groupB},
		}
		err := validateGroups(in)
		assert.NoError(t, err)
	})
	t.Run("error_when_duplicate_name", func(t *testing.T) {
		in := Groups{
			{
				Name: groupB,
			},
			{
				Name: groupB,
			},
		}
		err := validateGroups(in)
		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("duplicate names [%s]", groupB), err.Error())

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

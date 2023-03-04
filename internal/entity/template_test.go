package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	tmplName = "test_template"
)

func Test_TemplateProc(t *testing.T) {
	t.Parallel()

	const (
		in = `{{ .var }}`
	)
	t.Run("success", func(t *testing.T) {
		const (
			exp = `VAR_1`
		)
		proc := TmplProc{
			templateData: map[string]any{
				"var": "VAR_1",
			},
		}
		res, err := proc.Process(tmplName, in)
		assert.NoError(t, err)
		assert.Equal(t, exp, res)
	})
	t.Run("error", func(t *testing.T) {
		proc := TmplProc{
			templateOptions: []string{
				fmt.Sprintf("%v=%v", TemplateOptionsMissingKey, MissingKeyError),
			},
		}
		_, err := proc.Process(tmplName, in)
		assert.Error(t, err)
	})
}

func Test_TemplateFnsMap(t *testing.T) {
	t.Run("slice_func", func(t *testing.T) {
		t.Run("new", func(t *testing.T) {
			t.Run("var", func(t *testing.T) {
				const (
					in = `{{ $element := slice.New "a" 1 "b" }}{{ $element }}`
				)

				proc := TmplProc{templateFns: TemplateFnsMap}
				res, err := proc.Process(tmplName, in)
				assert.NoError(t, err)
				assert.Equal(t, "[a 1 b]", res)
			})

			t.Run("range", func(t *testing.T) {
				const (
					in = `{{range $element := slice.New "a" 1 "b" }} {{$element}} {{end}}`
				)

				proc := TmplProc{templateFns: TemplateFnsMap}
				res, err := proc.Process(tmplName, in)
				assert.NoError(t, err)
				assert.Equal(t, " a  1  b ", res)
			})

			t.Run("range_v2", func(t *testing.T) {
				const (
					in = `
#{{- $cp_config := "\n	- cp configs/config.yaml configs/config-%s.yaml"}}
cmd:
	# {{range $element := slice.New "dev1" "dev2" -}}{{ printf $cp_config  $element }}{{end}}
`
					exp = `
#
cmd:
	# 
	- cp configs/config.yaml configs/config-dev1.yaml
	- cp configs/config.yaml configs/config-dev2.yaml
`
				)

				proc := TmplProc{templateFns: TemplateFnsMap}
				res, err := proc.Process(tmplName, in)
				assert.NoError(t, err)
				assert.Equal(t, exp, res)
			})
		})
		t.Run("append", func(t *testing.T) {
			t.Run("single", func(t *testing.T) {
				const (
					in = `{{ $element := slice.New "a" }}{{ $element := slice.Append $element "b"}}{{ $element }}`
				)
				proc := TmplProc{templateFns: TemplateFnsMap}
				res, err := proc.Process(tmplName, in)
				assert.NoError(t, err)
				assert.Equal(t, "[a b]", res)
			})
			t.Run("multiple", func(t *testing.T) {
				const (
					in = `{{ $element := slice.New "a" }}{{ $element := slice.Append $element "b" "c"}}{{ $element }}`
				)
				proc := TmplProc{templateFns: TemplateFnsMap}
				res, err := proc.Process(tmplName, in)
				assert.NoError(t, err)
				assert.Equal(t, "[a b c]", res)

			})
		})
	})
}

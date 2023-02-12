package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TemplateProc(t *testing.T) {
	t.Parallel()

	const (
		tmplName = "test_template"
		in       = `{{ .var }}`
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

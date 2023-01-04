package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_YamlRootNodesOrder(t *testing.T) {
	t.Parallel()

	const (
		dirs  = "dirs"
		files = "files"
		cmd   = "cmd"
	)

	t.Run("success_get_order", func(t *testing.T) {
		var (
			in = []byte(fmt.Sprintf(`
%s:
  - "a/a"
  - "ccc"
%s:
  - path: a/x
    template: >
      Hello {{.}}

      You are doing great. Keep learning.
      Do not stop {{.}}
  - path: a/b
    template: xxx
%s:
  - pwd
  - ls -al
`, dirs, files, cmd))
		)
		order, err := YamlRootNodesOrder(in)
		assert.NoError(t, err)
		assert.Len(t, order, 3)
		assert.Equal(t, order[dirs], 0)
		assert.Equal(t, order[files], 1)
		assert.Equal(t, order[cmd], 2)
	})
	t.Run("success_with_empty_node", func(t *testing.T) {
		var (
			in = []byte(fmt.Sprintf(`%s:`, dirs))
		)
		order, err := YamlRootNodesOrder(in)
		assert.NoError(t, err)
		assert.Len(t, order, 1)
		assert.Equal(t, order[dirs], 0)
	})
	t.Run("success_and_empty_order_when_node_without_colon_sign", func(t *testing.T) {
		var (
			in = []byte(fmt.Sprintf(`%s`, dirs))
		)
		order, err := YamlRootNodesOrder(in)
		assert.NoError(t, err)
		assert.Len(t, order, 0)
	})
}

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

		res, err := preprocessRawConfigData(name, []byte(in))
		assert.NoError(t, err)
		assert.Equal(t, expected, string(res))
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

		res, err := preprocessRawConfigData(name, []byte(in))
		assert.NoError(t, err)
		assert.Equal(t, expected, string(res))
	})
}

package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseVars(t *testing.T) {
	const (
		emptyString = ""
	)

	t.Run("empty_vars", func(t *testing.T) {
		vars := parseVars([]string{})
		assert.Len(t, vars, 0)
	})
	t.Run("one_empty_var", func(t *testing.T) {
		const (
			a = "A"
		)
		vars := parseVars([]string{
			a + "=",
		})
		assert.Len(t, vars, 1)
		assert.Contains(t, vars, a)
		assert.Equal(t, emptyString, vars[a])
	})
	t.Run("one_not_separated_value", func(t *testing.T) {
		const (
			a    = "A"
			aVal = "A=A=X;Mm"
		)
		vars := parseVars([]string{
			a + "=" + aVal,
		})
		assert.Len(t, vars, 1)
		assert.Contains(t, vars, a)
		assert.Equal(t, aVal, vars[a])
	})
	t.Run("override_value", func(t *testing.T) {
		const (
			a            = "A"
			aVal         = "A=A=X;Mm"
			aOverrideVal = "override"
		)
		vars := parseVars([]string{
			a + "=" + aVal,
			a + "=" + aOverrideVal,
		})
		assert.Len(t, vars, 1)
		assert.Contains(t, vars, a)
		assert.Equal(t, aOverrideVal, vars[a])
	})
}

func Test_YamlRootNodesOrder(t *testing.T) {
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

func Test_tryFindAllVars(t *testing.T) {
	const (
		a = "a"
		b = "b"
		c = "c"
		d = "d"
		e = "e"
	)
	t.Run("success_get_vars", func(t *testing.T) {
		var (
			in = []byte(fmt.Sprintf(`
${%s}:
	${%s}
	var:${%s}
	var:github.com/${%s}/some_path
	
${%s}
`, a, b, c, d, e))
		)
		vars, err := tryFindAllVars(in)
		assert.NoError(t, err)
		assert.Len(t, vars, 5)
		assert.Contains(t, vars, a)
		assert.Contains(t, vars, b)
		assert.Contains(t, vars, c)
		assert.Contains(t, vars, d)
		assert.Contains(t, vars, d)
	})
	t.Run("success_not_get_when_vars_not exists", func(t *testing.T) {
		var (
			in = []byte(``)
		)
		vars, err := tryFindAllVars(in)
		assert.NoError(t, err)
		assert.Len(t, vars, 0)
	})
	t.Run("error_when_var_is_invalid", func(t *testing.T) {
		var (
			in = []byte(`${a ${b}`)
		)
		_, err := tryFindAllVars(in)
		assert.Error(t, err)
	})
}

func Test_replaceVars(t *testing.T) {
	t.Run("replace_vars", func(t *testing.T) {
		out := replaceVars([]byte("${a}${b}"), Vars{
			"a": "A",
			"b": "B",
		})
		assert.Equal(t, []byte("AB"), out)
	})
	t.Run("not_replace_vars_if_input_not_contains_vars", func(t *testing.T) {
		out := replaceVars([]byte("ab"), Vars{
			"a": "A",
			"b": "B",
		})
		assert.Equal(t, []byte("ab"), out)
	})
	t.Run("not_replace_vars_if_var_to_replace_is_empty", func(t *testing.T) {
		out := replaceVars([]byte("${a}${b}"), Vars{})
		assert.Equal(t, []byte("${a}${b}"), out)
	})
}

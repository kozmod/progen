package flag

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TemplateVarsFlag(t *testing.T) {
	const (
		usage    = "parse vars"
		setName  = "template vars test flag set"
		flagName = "tvar"
	)

	var (
		flagKey = fmt.Sprintf("-%s", flagName)
	)

	t.Run("error_when_flag_set_without_value", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ContinueOnError)
			vars TemplateVarsFlag
		)
		fs.Var(&vars, flagName, usage)

		err := fs.Parse([]string{
			flagKey, ".var",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), ErrVariableNotSet.Error())
		assert.Len(t, vars.Vars, 0)
	})
	t.Run("success_when_flag_not_set", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ContinueOnError)
			vars TemplateVarsFlag
		)
		fs.Var(&vars, flagName, usage)

		err := fs.Parse([]string{
			flagKey, "",
		})
		assert.NoError(t, err)
		assert.Len(t, vars.Vars, 0)
	})
	t.Run("success_parse_single_variable_L1", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ExitOnError)
			vars TemplateVarsFlag
		)
		fs.Var(&vars, flagName, usage)

		err := fs.Parse([]string{
			flagKey, ".var=1",
		})
		assert.NoError(t, err)
		assert.Len(t, vars.Vars, 1)
		assert.True(t, reflect.DeepEqual(
			map[string]any{"var": "1"},
			vars.Vars))
	})
	t.Run("success_parse_single_variable_L2", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ExitOnError)
			vars TemplateVarsFlag
		)
		fs.Var(&vars, flagName, usage)

		err := fs.Parse([]string{
			flagKey, ".var.X=1",
		})
		assert.NoError(t, err)
		assert.Len(t, vars.Vars, 1)
		assert.True(t, reflect.DeepEqual(
			map[string]any{"var": map[string]any{"X": "1"}},
			vars.Vars))
	})
	t.Run("success_parse_single_variable_L3", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ExitOnError)
			vars TemplateVarsFlag
		)
		fs.Var(&vars, flagName, usage)

		err := fs.Parse([]string{
			flagKey, ".var.Y.X=1",
		})
		assert.NoError(t, err)
		assert.Len(t, vars.Vars, 1)
		assert.True(t, reflect.DeepEqual(
			map[string]any{"var": map[string]any{"Y": map[string]any{"X": "1"}}},
			vars.Vars))
	})
	t.Run("success_parse_multi_different_variables_L2_L3", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ExitOnError)
			vars TemplateVarsFlag
		)
		fs.Var(&vars, flagName, usage)

		err := fs.Parse([]string{
			flagKey, ".var.X=1",
			flagKey, ".var.B=SOME_B",
			flagKey, ".var.Y.Y=SOME_Y",
		})
		assert.NoError(t, err)
		assert.Len(t, vars.Vars, 1)
		assert.True(t, reflect.DeepEqual(
			map[string]any{"var": map[string]any{"X": "1", "B": "SOME_B", "Y": map[string]any{"Y": "SOME_Y"}}},
			vars.Vars))
	})
	t.Run("success_parse_and_replace_multi_same_variables", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ExitOnError)
			vars TemplateVarsFlag
		)
		fs.Var(&vars, flagName, usage)

		err := fs.Parse([]string{
			flagKey, ".var.X=1",
			flagKey, ".var.X=SOME",
		})
		assert.NoError(t, err)
		assert.Len(t, vars.Vars, 1)
		assert.True(t, reflect.DeepEqual(
			map[string]any{"var": map[string]any{"X": "SOME"}},
			vars.Vars))
	})
}

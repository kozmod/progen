package flag

import (
	"errors"
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kozmod/progen/internal/entity"
)

func Test_TemplateVarsFlag(t *testing.T) {
	t.Parallel()

	const (
		usage    = "parse vars"
		setName  = "template vars test flag set"
		flagName = flagKeyTemplateVariables
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
	t.Run("success_parse_with_equal_sign", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ExitOnError)
			vars TemplateVarsFlag
		)
		fs.Var(&vars, flagName, usage)

		err := fs.Parse([]string{
			flagKey, ".var.X=SOME=x=y",
		})
		assert.NoError(t, err)
		assert.Len(t, vars.Vars, 1)
		assert.True(t, reflect.DeepEqual(
			map[string]any{"var": map[string]any{"X": "SOME=x=y"}},
			vars.Vars))
	})
}

func Test_SkipFlag(t *testing.T) {
	t.Parallel()

	const (
		usage    = "skip_flag_test_usage"
		setName  = "skip_fs"
		flagName = setName
	)

	var (
		flagKey = fmt.Sprintf("-%s", flagName)
	)

	t.Run("success_when_flag_not_set", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ContinueOnError)
			skip SkipFlag
		)
		fs.Var(&skip, flagName, usage)

		err := fs.Parse([]string{
			flagKey, "",
			flagKey, "",
		})
		assert.NoError(t, err)
		assert.Len(t, skip, 0)
	})
	t.Run("success_when_flag_not_set_v2", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ContinueOnError)
			skip SkipFlag
		)
		fs.Var(&skip, flagName, usage)

		err := fs.Parse(nil)
		assert.NoError(t, err)
		assert.Len(t, skip, 0)
	})
	t.Run("success_when_flag_set", func(t *testing.T) {
		var (
			fs   = flag.NewFlagSet(setName, flag.ContinueOnError)
			skip SkipFlag

			flagA = "cmd"
			flagB = "dirs"
		)
		fs.Var(&skip, flagName, usage)

		err := fs.Parse([]string{
			flagKey, flagA,
			flagKey, flagB,
		})
		assert.NoError(t, err)
		assert.Len(t, skip, 2)
		assert.ElementsMatch(t, []string{flagA, flagB}, skip)
	})
}

func Test_group(t *testing.T) {
	t.Parallel()

	const (
		usage    = "group_flag_test_usage"
		setName  = "group"
		flagName = setName
	)

	var (
		flagKey = fmt.Sprintf("-%s", flagName)
	)

	t.Run("success_when_flag_not_set", func(t *testing.T) {
		var (
			fs    = flag.NewFlagSet(setName, flag.ContinueOnError)
			group GroupFlag
		)
		fs.Var(&group, flagName, usage)

		err := fs.Parse([]string{
			flagKey, "",
			flagKey, "",
		})
		assert.NoError(t, err)
		assert.Len(t, group, 0)
	})
	t.Run("success_when_flag_not_set_v2", func(t *testing.T) {
		var (
			fs    = flag.NewFlagSet(setName, flag.ContinueOnError)
			group GroupFlag
		)
		fs.Var(&group, flagName, usage)

		err := fs.Parse(nil)
		assert.NoError(t, err)
		assert.Len(t, group, 0)
	})
	t.Run("success_when_flag_set", func(t *testing.T) {
		var (
			fs    = flag.NewFlagSet(setName, flag.ContinueOnError)
			group GroupFlag

			flagA = "group1"
			flagB = "group2"
		)
		fs.Var(&group, flagName, usage)

		err := fs.Parse([]string{
			flagKey, flagA,
			flagKey, flagB,
		})
		assert.NoError(t, err)
		assert.Len(t, group, 2)
		assert.ElementsMatch(t, []string{flagA, flagB}, group)
	})
}

func Test_parseFlags(t *testing.T) {
	t.Parallel()

	//goland:noinspection SpellCheckingInspection
	const (
		fsName = "testFlagSet"

		configPath = "progen.yml"
		dot        = entity.Dot
		dash       = entity.Dash
		lessThan   = entity.LessThan
	)

	dashPrefixFn := func(s string) string {
		t.Helper()
		return dash + s
	}

	//goland:noinspection SpellCheckingInspection
	var (
		v           = dashPrefixFn(flagKeyVarbose)
		errstack    = dashPrefixFn(flagKeyErrorStackTrace)
		printconfig = dashPrefixFn(flagKeyPrintConfig)
		dr          = dashPrefixFn(flagKeyDryRun)
		f           = dashPrefixFn(flagKeyConfigFile)
		pf          = dashPrefixFn(flagKeyPreprocessingAllFiles)
		missingkey  = dashPrefixFn(flagKeyMissingKey)
	)

	t.Run("success", func(t *testing.T) {
		testFs := flag.NewFlagSet(fsName, flag.ContinueOnError)
		flags, err := parseFlags(testFs, []string{v, printconfig, errstack, dr, f, configPath})
		assert.NoError(t, err)
		assert.Equal(t,
			Flags{
				Verbose:              true,
				PrintErrorStackTrace: true,
				PrintProcessedConfig: true,
				PreprocessFiles:      true,
				DryRun:               true,
				ConfigPath:           configPath,
				AWD:                  dot},
			flags)
	})
	t.Run("success_when_dash_last", func(t *testing.T) {
		testFs := flag.NewFlagSet(fsName, flag.ContinueOnError)
		flags, err := parseFlags(testFs, []string{v, dr, dash})
		assert.NoError(t, err)
		assert.Equal(t,
			Flags{Verbose: true, PreprocessFiles: true, DryRun: true, ConfigPath: configPath, AWD: dot, ReadStdin: true},
			flags)
	})
	t.Run("success_when_dash_last_and_before_less_than", func(t *testing.T) {
		testFs := flag.NewFlagSet(fsName, flag.ContinueOnError)
		flags, err := parseFlags(testFs, []string{v, dr, dash, lessThan, configPath})
		assert.NoError(t, err)
		assert.Equal(t,
			Flags{Verbose: true, PreprocessFiles: true, DryRun: true, ConfigPath: configPath, AWD: dot, ReadStdin: true},
			flags)
	})
	t.Run("error_when_flag_not_specified", func(t *testing.T) {
		testFs := flag.NewFlagSet(fsName, flag.ContinueOnError)
		testFs.SetOutput(MockWriter{
			assertWriteFn: func(p []byte) {
				assert.NotEmpty(t, p)
			},
		})
		_, err := parseFlags(testFs, []string{v, dr, f})
		assert.Error(t, err)
	})
	t.Run("error_when_dash_flag_not_last", func(t *testing.T) {
		testFs := flag.NewFlagSet(fsName, flag.ContinueOnError)
		_, err := parseFlags(testFs, []string{v, dash, dr, f, configPath})
		assert.Error(t, err)
		assert.True(t, errors.Is(err, ErrDashFlagNotLast))
	})
	t.Run("success_with_preprocess_flag_false", func(t *testing.T) {
		testFs := flag.NewFlagSet(fsName, flag.ContinueOnError)
		flags, err := parseFlags(testFs, []string{fmt.Sprintf("%s=%v", pf, false)})
		assert.NoError(t, err)
		assert.Equal(t,
			Flags{Verbose: false, PreprocessFiles: false, DryRun: false, ConfigPath: configPath, AWD: dot, ReadStdin: false},
			flags)
	})
	t.Run("missingkey", func(t *testing.T) {
		t.Run("success_with_set_default", func(t *testing.T) {
			var (
				missingKeyValueDefault = entity.MissingKeyDefault
				missingKeyValue        = fmt.Sprintf("%s=%v", missingkey, missingKeyValueDefault)
				missingKeyValueExp     = missingKeyValue[1:]
			)

			testFs := flag.NewFlagSet(fsName, flag.ContinueOnError)
			flags, err := parseFlags(testFs, []string{missingKeyValue})
			assert.NoError(t, err)
			assert.Equal(t,
				Flags{
					Verbose:         false,
					PreprocessFiles: true,
					DryRun:          false,
					ConfigPath:      configPath,
					AWD:             dot,
					ReadStdin:       false,
					MissingKey:      MissingKeyFlag(missingKeyValueDefault)},
				flags)
			assert.Equal(t, missingKeyValueExp, flags.MissingKey.String())
		})
		t.Run("error_with_set_unexpected", func(t *testing.T) {
			var (
				missingKeyValueDefault = "xxx_unexpected_xxx"
				missingKeyValue        = fmt.Sprintf("%s=%v", missingkey, missingKeyValueDefault)
			)

			testFs := flag.NewFlagSet(fsName, flag.ContinueOnError)
			testFs.SetOutput(MockWriter{
				assertWriteFn: func(p []byte) {
					assert.NotEmpty(t, p)
				},
			})
			_, err := parseFlags(testFs, []string{missingKeyValue})
			assert.Error(t, err)
		})
	})
}

type MockWriter struct {
	assertWriteFn func(p []byte)
}

func (m MockWriter) Write(p []byte) (n int, err error) {
	m.assertWriteFn(p)
	return 0, nil
}

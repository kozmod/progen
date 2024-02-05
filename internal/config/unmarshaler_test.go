package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kozmod/progen/internal/entity"
)

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

func Test_settings_tag(t *testing.T) {
	t.Run("parse_http", func(t *testing.T) {
		const (
			baseUrl   = `https://gitlab.repo_2.com/api/v4/projects/5/repository/files/`
			headerKey = `PRIVATE-TOKEN`
			headerVal = `glpat-SOME_TOKEN`
			queryKey  = `Q_PARAM_1`
			queryVal  = `Q_Val_1`
		)
		in := fmt.Sprintf(`
settings:
  http:
    debug: true
    base_url: %s
    headers:
      %s: %s
    query_params:
      %s: %s
`,
			baseUrl,
			headerKey,
			headerVal,
			queryKey,
			queryVal,
		)

		conf, err := NewYamlConfigUnmarshaler().Unmarshal([]byte(in))
		assert.NoError(t, err)
		assert.NotNil(t, conf.Settings.HTTP)

		httpConf := conf.Settings.HTTP
		assert.True(t, httpConf.Debug)
		assert.Equal(t, baseUrl, httpConf.BaseURL.String())
		assert.Len(t, httpConf.Headers, 1)

		header, ok := httpConf.Headers[headerKey]
		assert.True(t, ok)
		assert.Equal(t, headerVal, header)

		assert.Len(t, httpConf.QueryParams, 1)

		query, ok := httpConf.QueryParams[queryKey]
		assert.True(t, ok)
		assert.Equal(t, queryVal, query)
	})

	t.Run("groups_http", func(t *testing.T) {
		const (
			groupA  = "a"
			groupB  = "b"
			actionA = "cmdA"
			actionB = "cmdB"
		)
		in := fmt.Sprintf(`
settings:
   groups:
      - name: %s
        actions: [%s]
      - name: %s
        actions:
         - %s
         - %s
        manual: true
`,
			groupA,
			actionA,
			groupB,
			actionA,
			actionB,
		)

		conf, err := NewYamlConfigUnmarshaler().Unmarshal([]byte(in))
		assert.NoError(t, err)

		assert.NotEmpty(t, conf.Settings.Groups)
		assert.Len(t, conf.Settings.Groups, 2)

		assert.Equal(t, conf.Settings.Groups[0], Group{Name: groupA, Actions: []string{actionA}, Manual: false})
		assert.Equal(t, conf.Settings.Groups[1], Group{Name: groupB, Actions: []string{actionA, actionB}, Manual: true})

	})
}

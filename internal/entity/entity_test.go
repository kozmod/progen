package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RegexpChain(t *testing.T) {
	const (
		regex1 = `cmd.?`
		regex2 = "^dirs$"
	)
	t.Run("single_check", func(t *testing.T) {
		var (
			a  = assert.New(t)
			rc = NewRegexpChain(regex1)
		)

		a.Len(rc.matchers, 1)
		a.Equal(regex1, rc.matchers[0].String())
		a.False(rc.MatchString("adam[23]"))
		a.False(rc.MatchString("eve[7]"))
		a.False(rc.MatchString("Job[48]"))
		a.False(rc.MatchString("snakey"))
		a.True(rc.MatchString("cmd"))
		a.True(rc.MatchString("cmd___"))
	})
	t.Run("multiple_check", func(t *testing.T) {
		var (
			a  = assert.New(t)
			rc = NewRegexpChain(regex1, regex2)
		)

		a.Len(rc.matchers, 2)
		a.Equal(regex1, rc.matchers[0].String())
		a.Equal(regex2, rc.matchers[1].String())
		a.True(rc.MatchString("cmd"))
		a.True(rc.MatchString("cmd___"))
		a.True(rc.MatchString("dirs"))
		a.False(rc.MatchString("dirs__"))
		a.False(rc.MatchString("c"))
	})

}

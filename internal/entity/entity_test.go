package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RegexpChain(t *testing.T) {
	t.Parallel()

	const (
		regex1 = `cmd.?`
		regex2 = "^dirs$"
	)

	if testing.Short() {
		t.Skip("skipping slow test")
	}

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

func Test_RandomFn(t *testing.T) {
	t.Parallel()

	//goland:noinspection SpellCheckingInspection
	const (
		alnum = "^[[:alnum:]]*$"
		print = "^[[:print:]]*$"
	)

	var (
		f = RandomFn{}
	)

	if testing.Short() {
		t.Skip("skipping slow test")
	}

	t.Run("AlphaNum", func(t *testing.T) {
		t.Run("generate_zero", func(t *testing.T) {

			s := f.AlphaNum(0)
			assert.Empty(t, s)
		})
		t.Run("generate_50", func(t *testing.T) {
			s := f.AlphaNum(50)
			assert.Len(t, s, 50)
			assert.Regexp(t, alnum, s)
		})
	})
	t.Run("Alpha", func(t *testing.T) {
		t.Run("generate_zero", func(t *testing.T) {
			s := f.Alpha(0)
			assert.Empty(t, s)
		})
		t.Run("generate_50", func(t *testing.T) {
			s := f.Alpha(50)
			assert.Len(t, s, 50)
			assert.Regexp(t, alnum, s)
		})
	})
	t.Run("ASCII", func(t *testing.T) {
		t.Run("generate_zero", func(t *testing.T) {
			s := f.ASCII(0)
			assert.Empty(t, s)
		})
		t.Run("generate_50", func(t *testing.T) {
			s := f.ASCII(50)
			assert.Len(t, s, 50)
			assert.Regexp(t, print, s)
		})
	})
}

func Test_MissingKye(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		err := MissingKeyValue("default").Valid()
		assert.NoError(t, err)
	})
	t.Run("error", func(t *testing.T) {
		err := MissingKeyValue("some").Valid()
		assert.Error(t, err)
	})
}

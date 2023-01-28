package entity

import "regexp"

const (
	Space      = " "
	Empty      = ""
	Dash       = "-"
	Dot        = "."
	EqualsSign = "="

	NewLine = '\n'
)

type FileProducer interface {
	Get() (*DataFile, error)
}

//goland:noinspection SpellCheckingInspection
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, any ...any)
	Warnf(format string, any ...any)
	Debugf(format string, any ...any)
}

type DataFile struct {
	Template
	Data []byte
}

type LocalFile struct {
	Template
	LocalPath string
}

type RemoteFile struct {
	Template
	HTTPClientParams
}

type Template struct {
	Path     string
	Name     string
	ExecTmpl bool
}

type HTTPClientParams struct {
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
}

type Command struct {
	Cmd  string
	Args []string
}

type RegexpChain struct {
	matchers []*regexp.Regexp
}

func NewRegexpChain(regexps ...string) *RegexpChain {
	matchers := make([]*regexp.Regexp, 0, len(regexps))
	for _, regex := range regexps {
		matcher := regexp.MustCompile(regex)
		if matcher != nil {
			matchers = append(matchers, matcher)
		}
	}
	return &RegexpChain{
		matchers: matchers,
	}
}

func (c *RegexpChain) MatchString(s string) bool {
	for _, matcher := range c.matchers {
		if matcher.MatchString(s) {
			return true
		}
	}
	return false
}

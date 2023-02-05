package entity

import (
	"fmt"
	"path/filepath"
	"regexp"
)

type (
	TemplateOptionsKey string
	MissingKeyValue    string
)

func (v MissingKeyValue) Valid() error {
	switch v {
	case MissingKeyDefault,
		MissingKeyInvalid,
		MissingKeyZero,
		MissingKeyError:
		return nil
	default:
		return fmt.Errorf("templte option [%v] is not valid: %v", TemplateOptionsMissingKey, v)
	}
}

//goland:noinspection SpellCheckingInspection
const (
	TemplateOptionsMissingKey TemplateOptionsKey = "missingkey"

	MissingKeyDefault MissingKeyValue = "default"
	MissingKeyInvalid MissingKeyValue = "invalid"
	MissingKeyZero    MissingKeyValue = "zero"
	MissingKeyError   MissingKeyValue = "error"

	Space      = " "
	Empty      = ""
	Dash       = "-"
	Dot        = "."
	EqualsSign = "="
	LessThan   = "<"

	NewLine = '\n'
)

type FileProducer interface {
	Get() (DataFile, error)
}

type FileProc interface {
	Process(file DataFile) (DataFile, error)
}

type Executor interface {
	Exec() error
}

type Preprocessor interface {
	Process() error
}

//goland:noinspection SpellCheckingInspection
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, any ...any)
	Warnf(format string, any ...any)
	Debugf(format string, any ...any)
}

type DataFile struct {
	FileInfo
	Data []byte
}

type LocalFile struct {
	FileInfo
	LocalPath string
}

type RemoteFile struct {
	FileInfo
	HTTPClientParams
}

type FileInfo struct {
	Dir  string
	Name string
	path *string
}

func (f *FileInfo) Path() string {
	if f.path == nil {
		path := filepath.Join(f.Dir, f.Name)
		f.path = &path
	}
	return *f.path
}

type HTTPClientParams struct {
	URL         string
	Headers     map[string]string
	QueryParams map[string]string
}

type Command struct {
	Cmd  string
	Args []string
	Dir  string
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

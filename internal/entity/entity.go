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
	SpacedPipe = " | "

	NewLine = '\n'
)

type (
	FileProducer interface {
		Get() (DataFile, error)
	}

	FileProc interface {
		Process(file DataFile) (DataFile, error)
	}

	DirProc interface {
		Process(path string) (string, error)
	}

	TemplateProc interface {
		Process(name, text string) (string, error)
	}

	CommandProc interface {
		Process(commands []Command) error
	}

	Executor interface {
		Exec() error
	}

	Preprocessor interface {
		Process() error
	}

	//goland:noinspection SpellCheckingInspection
	Logger interface {
		Infof(format string, args ...any)
		Errorf(format string, any ...any)
		Warnf(format string, any ...any)
		Debugf(format string, any ...any)
	}
)

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
	dir  string
	name string
	path *string
}

func (f *FileInfo) Name() string {
	return f.name
}

func (f *FileInfo) Dir() string {
	return f.dir
}

func NewFileInfo(path string) FileInfo {
	return FileInfo{
		name: filepath.Base(path),
		dir:  filepath.Dir(path),
	}
}

func (f *FileInfo) Path() string {
	if f.path == nil {
		path := filepath.Join(f.dir, f.name)
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

package entity

import "strings"

const (
	Space = " "
	Empty = ""
	Dash  = "-"

	NewLine = '\n'
)

type Flags struct {
	ConfigPath string
	Verbose    bool
	DryRun     bool
	Version    bool
	ReadStdin  bool
}

func (f *Flags) FileLocationMessage() string {
	switch {
	case f == nil:
		return Empty
	case f.ReadStdin:
		return "stdin"
	case strings.TrimSpace(f.ConfigPath) != Empty:
		return f.ConfigPath
	default:
		return Empty
	}
}

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
	Path     string
	Name     string
	Data     []byte
	ExecTmpl bool
}

type LocalFile struct {
	Path      string
	Name      string
	LocalPath string
	ExecTmpl  bool
}

type RemoteFile struct {
	Path     string
	Name     string
	URL      string
	Headers  map[string]string
	ExecTmpl bool
}

type Command struct {
	Cmd  string
	Args []string
}

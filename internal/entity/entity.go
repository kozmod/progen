package entity

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

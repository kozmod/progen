package entity

import "os"

const (
	Space = " "
	Empty = ""
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

type Dir struct {
	Path string
	Perm os.FileMode
}

type File struct {
	Path string
	Name string
	Perm os.FileMode
}

type DataFile struct {
	File
	Data     []byte
	ExecTmpl bool
}

type LocalFile struct {
	File
	LocalPath string
	ExecTmpl  bool
}

type RemoteFile struct {
	File
	URL      string
	Headers  map[string]string
	ExecTmpl bool
}

type Command struct {
	Cmd  string
	Args []string
}

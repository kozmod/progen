package entity

const (
	Space = " "
	Empty = ""
)

type FileProducer interface {
	Get() (*File, error)
}

//goland:noinspection SpellCheckingInspection
type Logger interface {
	Infof(template string, args ...any)
}

type File struct {
	Path     string
	Name     string
	Data     []byte
	Template bool
}

type LocalFile struct {
	Path      string
	Name      string
	LocalPath string
	Template  bool
}

type RemoteFile struct {
	Path     string
	Name     string
	URL      string
	Headers  map[string]string
	Template bool
}

type Command struct {
	Cmd  string
	Args []string
}

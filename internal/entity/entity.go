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
	Path string
	Name string
	Data []byte
}

type RemoteFile struct {
	Path    string
	Name    string
	URL     string
	Headers map[string]string
}

type Command struct {
	Cmd  string
	Args []string
}

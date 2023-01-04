package entity

const (
	Space = " "
	Empty = ""
)

//goland:noinspection SpellCheckingInspection
type Logger interface {
	Infof(template string, args ...any)
}

type File struct {
	Path string
	Name string
	Data []byte
}

type Command struct {
	Cmd  string
	Args []string
}

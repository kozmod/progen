package entity

const (
	Space = " "
	Empty = ""
)

//goland:noinspection SpellCheckingInspection
type Logger interface {
	Infof(template string, args ...any)
}

type Template struct {
	Path string
	Name string
	Text []byte
	Data map[string]any
}

type Command struct {
	Cmd  string
	Args []string
}

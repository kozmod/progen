package entity

import (
	"fmt"
	"strings"
)

const (
	Space = " "
	Empty = ""
)

type File struct {
	Path string
	Name string
	Data []byte
}

type Command struct {
	Cmd  string
	Args []string
}

func NewCommand(cmd string) (Command, error) {
	cmd = strings.TrimSpace(cmd)
	if len(cmd) == 0 {
		return Command{}, fmt.Errorf("executor.NewCommand: command is empty")
	}
	splitCmd := make([]string, 0, len(cmd))
	for _, val := range strings.Split(cmd, Space) {
		if trimmed := strings.TrimSpace(val); val != Empty {
			splitCmd = append(splitCmd, trimmed)
		}
	}

	return Command{
		Cmd:  splitCmd[0],
		Args: splitCmd[1:],
	}, nil
}

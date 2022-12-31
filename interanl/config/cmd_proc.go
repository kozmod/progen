package config

import (
	"fmt"

	"github.com/kozmod/progen/interanl/entity"
	"github.com/kozmod/progen/interanl/proc"
)

func MustConfigureRunCommandProc(conf Config) (proc.Proc, error) {
	if len(conf.Cmd) == 0 {
		return nil, nil
	}
	commands, err := stringsToCommands(conf.Cmd)
	if err != nil {
		return nil, fmt.Errorf("config.MustConfigureRunCommandProc: %w", err)
	}

	return proc.NewRunCommandProc(commands), nil
}

func stringsToCommands(in []string) ([]entity.Command, error) {
	commands := make([]entity.Command, 0, len(in))
	for _, cmd := range in {
		command, err := entity.NewCommand(cmd)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		commands = append(commands, command)
	}
	return commands, nil
}

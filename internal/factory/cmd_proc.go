package factory

import (
	"fmt"
	"strings"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewRunCommandProc(conf config.Config, logger entity.Logger, dryRun bool) (proc.Proc, error) {
	if len(conf.Cmd) == 0 {
		return nil, nil
	}

	commands := make([]entity.Command, 0, len(conf.Cmd))
	for _, cmd := range conf.Cmd {
		cmd = strings.TrimSpace(cmd)
		if len(cmd) == 0 {
			return nil, fmt.Errorf("command is empty")
		}

		command := commandFromString(cmd)
		commands = append(commands, command)
	}

	if dryRun {
		return proc.NewDryRunCommandProc(commands, logger), nil
	}

	return proc.NewCommandProc(commands, logger), nil
}

func commandFromString(cmd string) entity.Command {
	splitCmd := make([]string, 0, len(cmd))
	for _, val := range strings.Split(cmd, entity.Space) {
		if trimmed := strings.TrimSpace(val); val != entity.Empty {
			splitCmd = append(splitCmd, trimmed)
		}
	}

	return entity.Command{
		Cmd:  splitCmd[0],
		Args: splitCmd[1:],
	}
}

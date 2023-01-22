package factory

import (
	"fmt"
	"strings"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

//goland:noinspection SpellCheckingInspection
func NewRunCommandProc(cmds []string, logger entity.Logger, dryRun bool) (proc.Proc, error) {
	if len(cmds) == 0 {
		logger.Infof("`cmd` section is empty")
		return nil, nil
	}

	commands := make([]entity.Command, 0, len(cmds))
	for _, cmd := range cmds {
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

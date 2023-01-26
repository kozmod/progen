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
	var (
		splitCmd = strings.Split(cmd, entity.Space)
		res      = make([]string, 0, len(splitCmd))
	)

	for _, val := range splitCmd {
		if trimmed := strings.TrimSpace(val); val != entity.Empty {
			res = append(res, trimmed)
		}
	}

	return entity.Command{
		Cmd:  res[0],
		Args: res[1:],
	}
}

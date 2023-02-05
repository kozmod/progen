package factory

import (
	"fmt"
	"strings"

	"github.com/kozmod/progen/internal/config"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewRunCommandExecutor(cmds []config.Command, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(cmds) == 0 {
		logger.Infof("`cmd` section is empty")
		return nil, nil
	}

	commands := make([]entity.Command, 0, len(cmds))
	for i, cmd := range cmds {

		dir := strings.TrimSpace(cmd.Dir)
		if dir == entity.Empty {
			dir = entity.Dot
		}

		for j, exec := range cmd.Exec {
			exec = strings.TrimSpace(exec)
			if len(exec) == 0 {
				return nil, fmt.Errorf("command is empty [section: %d, exec: %d, dir: %s]", i, j, dir)
			}
			command := commandFromString(exec, dir)
			commands = append(commands, command)
		}
	}

	if dryRun {
		return proc.NewDryRunCommandProc(commands, logger), nil
	}

	return proc.NewCommandProc(commands, logger), nil
}

func commandFromString(cmd, dir string) entity.Command {
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
		Dir:  dir,
	}
}

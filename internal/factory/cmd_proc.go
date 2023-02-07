package factory

import (
	"fmt"
	"strings"

	"github.com/kozmod/progen/internal/config"

	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

//goland:noinspection SpellCheckingInspection
func NewRunCommandExecutor(cmds []config.Command, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(cmds) == 0 {
		logger.Infof("`cmd` section is empty")
		return nil, nil
	}

	executors := make([]entity.Executor, 0, len(cmds))
	for i, cmd := range cmds {
		dir := strings.TrimSpace(cmd.Dir)
		if dir == entity.Empty {
			dir = entity.Dot
		}
		commands := make([]entity.Command, 0, len(cmd.Exec))
		for j, exec := range cmd.Exec {
			exec = strings.TrimSpace(exec)
			if len(exec) == 0 {
				return nil, fmt.Errorf("command is empty [section: %d, exec: %d, dir: %s]", i, j, dir)
			}
			command := commandFromString(exec)
			commands = append(commands, command)
		}

		switch {
		case dryRun:
			executors = append(executors, proc.NewDryRunCommandExecutor(commands, dir, logger))
		case cmd.Pipe:
			executors = append(executors, proc.NewPipeCommandExecutor(commands, dir, logger))
		default:
			executors = append(executors, proc.NewCommandExecutor(commands, dir, logger))
		}
	}

	return proc.NewCommandListExecutor(executors), nil
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

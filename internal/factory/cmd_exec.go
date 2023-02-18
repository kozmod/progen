package factory

import (
	"strings"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

//goland:noinspection SpellCheckingInspection
func NewRunCommandExecutor(cmds []config.Command, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(cmds) == 0 {
		logger.Infof("`cmd` section is empty")
		return nil, nil
	}

	commands := make([]entity.Command, 0, len(cmds))
	for _, cmd := range cmds {
		dir := strings.TrimSpace(cmd.Dir)
		if dir == entity.Empty {
			dir = entity.Dot
		}

		args := make([]string, 0, len(cmd.Args)+1)
		args = append(args, cmd.Args...)

		commands = append(commands, entity.Command{
			Cmd:  cmd.Exec,
			Args: args,
			Dir:  dir,
		})
	}

	switch {
	case dryRun:
		return exec.NewDryRunCommandExecutor(commands, logger), nil
	default:
		return exec.NewCommandExecutor(commands, logger), nil
	}
}

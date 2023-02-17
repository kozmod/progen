package factory

import (
	"strings"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/exec"
)

func NewRunScriptsExecutor(scripts []config.Script, logger entity.Logger, dryRun bool) (entity.Executor, error) {
	if len(scripts) == 0 {
		logger.Infof("`scripts` section is empty")
		return nil, nil
	}

	executors := make([]entity.Executor, 0, len(scripts))
	for _, script := range scripts {
		dir := strings.TrimSpace(script.Dir)
		if dir == entity.Empty {
			dir = entity.Dot
		}

		args := make([]string, 0, len(script.Args)+1)
		args = append(args, script.Args...)
		args = append(args, script.Script)

		commands := []entity.Command{
			{
				Cmd:  script.Exec,
				Args: args,
			},
		}

		switch {
		case dryRun:
			executors = append(executors, exec.NewDryRunCommandExecutor(commands, dir, logger))
		default:
			executors = append(executors, exec.NewCommandExecutor(commands, dir, logger))
		}
	}

	return exec.NewCommandListExecutor(executors), nil
}

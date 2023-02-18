package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kozmod/progen/internal/entity"
)

type CommandExecutor struct {
	commands []entity.Command
	logger   entity.Logger
}

func NewCommandExecutor(commands []entity.Command, logger entity.Logger) *CommandExecutor {
	return &CommandExecutor{
		commands: commands,
		logger:   logger,
	}
}

func (p *CommandExecutor) Exec() error {
	for _, command := range p.commands {
		dir := command.Dir
		cmd := exec.Command(command.Cmd, command.Args...)
		var (
			stdout bytes.Buffer
			stderr bytes.Buffer
		)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		cmd.Dir = dir

		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("execute command [dir: %s] %s\nerror: %w", dir, prepareCmdMessage(&stderr, command), err)
		}

		p.logger.Infof("execute [dir: %s]: %s", dir, prepareCmdMessage(&stdout, command))
	}
	return nil
}

func prepareCmdMessage(out fmt.Stringer, command entity.Command) string {
	const (
		outMsg = "out:"
	)
	message := strings.Join(append([]string{command.Cmd}, command.Args...), entity.Space)
	if strings.Contains(message, entity.NewLine) {
		message = entity.NewLine + message
	}

	if out != nil {
		if output := out.String(); strings.TrimSpace(output) != entity.Empty {
			message += fmt.Sprintf("\n%s\n%s", outMsg, output)
		}
	}

	return message
}

type DryRunCommandExecutor struct {
	commands []entity.Command
	logger   entity.Logger
}

func NewDryRunCommandExecutor(commands []entity.Command, logger entity.Logger) *DryRunCommandExecutor {
	return &DryRunCommandExecutor{
		commands: commands,
		logger:   logger,
	}
}

func (p *DryRunCommandExecutor) Exec() error {
	for _, command := range p.commands {
		p.logger.Infof("execute [dir: %s]: %s", command.Dir, prepareCmdMessage(nil, command))
	}
	return nil
}

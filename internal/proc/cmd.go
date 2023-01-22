package proc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kozmod/progen/internal/entity"
)

type CommandProc struct {
	commands []entity.Command
	logger   entity.Logger
}

func NewCommandProc(commands []entity.Command, logger entity.Logger) *CommandProc {
	return &CommandProc{
		commands: commands,
		logger:   logger,
	}
}

func (p *CommandProc) Exec() error {
	for _, command := range p.commands {
		cmd := exec.Command(command.Cmd, command.Args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("run command: %w", err)
		}

		p.logger.Infof(func() string {
			info := fmt.Sprintf("execute: %s", strings.Join(append([]string{command.Cmd}, command.Args...), entity.Space))
			if output := out.String(); strings.TrimSpace(output) != entity.Empty {
				info += fmt.Sprintf("\nout:\n%s", output)
			}
			return info
		}())
	}
	return nil
}

type DryRunCommandProc struct {
	commands []entity.Command
	logger   entity.Logger
}

func NewDryRunCommandProc(commands []entity.Command, logger entity.Logger) *DryRunCommandProc {
	return &DryRunCommandProc{
		commands: commands,
		logger:   logger,
	}
}

func (p *DryRunCommandProc) Exec() error {
	for _, command := range p.commands {
		p.logger.Infof("execute cmd: %s",
			strings.Join(append([]string{command.Cmd}, command.Args...), entity.Space))
	}
	return nil
}

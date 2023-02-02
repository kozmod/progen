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
			return fmt.Errorf("run command: [%s] : %w", p.prepareMessage(command, &out), err)
		}

		p.logger.Infof("execute: %s", p.prepareMessage(command, &out))
	}
	return nil
}

func (p *CommandProc) prepareMessage(command entity.Command, out fmt.Stringer) string {
	message := strings.Join(append([]string{command.Cmd}, command.Args...), entity.Space)
	if out == nil {
		return message
	}
	if output := out.String(); strings.TrimSpace(output) != entity.Empty {
		message += fmt.Sprintf("\nout:\n%s", output)
	}
	return message
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

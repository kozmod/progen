package proc

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/kozmod/progen/internal/entity"
)

type RunCommandProc struct {
	commands []entity.Command
	logger   entity.Logger
}

func NewRunCommandProc(commands []entity.Command, logger entity.Logger) *RunCommandProc {
	return &RunCommandProc{
		commands: commands,
		logger:   logger,
	}
}

func (p *RunCommandProc) Exec() error {
	for _, command := range p.commands {
		cmd := exec.Command(command.Cmd, command.Args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("run command: %w", err)
		}
		p.logger.Infof("execute:\ncmd: %s\nout: %s",
			strings.Join(append([]string{command.Cmd}, command.Args...), entity.Space),
			out.String())
	}
	return nil
}

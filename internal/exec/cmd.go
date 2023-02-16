package exec

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/kozmod/progen/internal/entity"
)

type CommandListExecutor struct {
	executors []entity.Executor
}

func NewCommandListExecutor(executors []entity.Executor) *CommandListExecutor {
	return &CommandListExecutor{
		executors: executors,
	}
}

func (e *CommandListExecutor) Exec() error {
	for i, executor := range e.executors {
		err := executor.Exec()
		if err != nil {
			return fmt.Errorf("run command executor [%d]: %w", i, err)
		}
	}
	return nil
}

type CommandExecutor struct {
	commands []entity.Command
	logger   entity.Logger
	dir      string
}

func NewCommandExecutor(commands []entity.Command, dir string, logger entity.Logger) *CommandExecutor {
	return &CommandExecutor{
		commands: commands,
		logger:   logger,
		dir:      dir,
	}
}

func (p *CommandExecutor) Exec() error {
	dir := p.dir
	for _, command := range p.commands {
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

type PipeCommandExecutor struct {
	commands []entity.Command
	logger   entity.Logger
	dir      string
}

func NewPipeCommandExecutor(commands []entity.Command, dir string, logger entity.Logger) *PipeCommandExecutor {
	return &PipeCommandExecutor{
		commands: commands,
		logger:   logger,
		dir:      dir,
	}
}

func (p *PipeCommandExecutor) Exec() error {
	if len(p.commands) == 0 {
		return nil
	}

	dir := p.dir
	stack := make([]*exec.Cmd, 0, len(p.commands))
	for _, cmd := range p.commands {
		command := exec.Command(cmd.Cmd, cmd.Args...)
		command.Dir = dir
		stack = append(stack, command)
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	err := p.execute(&stdout, &stderr, stack)
	if err != nil {
		return fmt.Errorf("execute command in pipe [dir: %s] %s\nerror: %w", dir, prepareCmdMessage(&stderr, p.commands...), err)
	}
	p.logger.Infof("execute pipe [dir: %s]: %s", dir, prepareCmdMessage(&stdout, p.commands...))
	return nil
}

func (p *PipeCommandExecutor) execute(outBuf, errBuf *bytes.Buffer, stack []*exec.Cmd) error {
	pipeStack := make([]*io.PipeWriter, len(stack)-1)
	i := 0
	for ; i < len(stack)-1; i++ {
		stdinPipe, stdoutPipe := io.Pipe()
		stack[i].Stdout = stdoutPipe
		stack[i+1].Stdin = stdinPipe
		pipeStack[i] = stdoutPipe
	}
	stack[i].Stdout = outBuf
	stack[i].Stderr = errBuf

	if err := p.call(stack, pipeStack); err != nil {
		return err
	}
	return nil
}

func (p *PipeCommandExecutor) call(stack []*exec.Cmd, pipes []*io.PipeWriter) (err error) {
	if stack[0].Process == nil {
		if err = stack[0].Start(); err != nil {
			return err
		}
	}
	if len(stack) > 1 {
		if err = stack[1].Start(); err != nil {
			return err
		}
		defer func() {
			if err == nil {
				_ = pipes[0].Close()
				err = p.call(stack[1:], pipes[1:])
			}
		}()
	}
	return stack[0].Wait()
}

func prepareCmdMessage(out fmt.Stringer, commands ...entity.Command) string {
	const (
		execMsg = "exec: "
		outMsg  = "out:"
	)
	var message string
	for _, command := range commands {
		cmd := strings.Join(append([]string{command.Cmd}, command.Args...), entity.Space)
		if strings.TrimSpace(message) == entity.Empty {
			message = cmd
			continue
		}
		message = strings.Join([]string{message, cmd}, entity.SpacedPipe)
	}
	if len(commands) > 0 {
		message = execMsg + message
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
	dir      string
}

func NewDryRunCommandExecutor(commands []entity.Command, dir string, logger entity.Logger) *DryRunCommandExecutor {
	return &DryRunCommandExecutor{
		commands: commands,
		logger:   logger,
		dir:      dir,
	}
}

func (p *DryRunCommandExecutor) Exec() error {
	for _, command := range p.commands {
		p.logger.Infof("execute [dir: %s]: %s",
			p.dir, strings.Join(append([]string{command.Cmd}, command.Args...), entity.Space))
	}
	return nil
}

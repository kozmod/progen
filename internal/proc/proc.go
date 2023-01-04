package proc

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kozmod/progen/internal/entity"
)

type MkdirAllProc struct {
	fileMode os.FileMode
	dirs     []string
}

func NewMkdirAllProc(dirs []string) *MkdirAllProc {
	return &MkdirAllProc{
		fileMode: os.ModePerm,
		dirs:     dirs,
	}
}

func (p *MkdirAllProc) Exec() error {
	for _, dir := range p.dirs {
		err := os.MkdirAll(dir, p.fileMode)
		if err != nil {
			return fmt.Errorf("create dir [%s]: %w", dir, err)
		}
	}
	return nil
}

type FileProc struct {
	fileMode os.FileMode
	files    []entity.File
}

func NewFileProc(files []entity.File) *FileProc {
	return &FileProc{
		fileMode: os.ModePerm,
		files:    files,
	}
}

func (p *FileProc) Exec() error {
	for _, file := range p.files {
		fileDir := file.Path
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			err = os.MkdirAll(fileDir, p.fileMode)
			if err != nil {
				return fmt.Errorf("create template dir [%s]: %w", fileDir, err)
			}
		}

		err := os.WriteFile(path.Join(file.Path, file.Name), file.Data, p.fileMode)
		if err != nil {
			return fmt.Errorf("create file [%s]: %w", file.Name, err)
		}
	}
	return nil
}

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
		p.logger.Infof("cmd: %s\n%s",
			strings.Join(append([]string{command.Cmd}, command.Args...), entity.Space),
			out.String())
	}
	return nil
}

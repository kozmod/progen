package proc

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

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
			return err
		}
	}
	return nil
}

type WriteFileProc struct {
	fileMode os.FileMode
	files    []entity.File
}

func NewWriteFileProc(files []entity.File) *WriteFileProc {
	return &WriteFileProc{
		fileMode: os.ModePerm,
		files:    files,
	}
}

func (p *WriteFileProc) Exec() error {
	for _, file := range p.files {
		fileDir := file.Path
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			err = os.MkdirAll(fileDir, p.fileMode)
			if err != nil {
				return fmt.Errorf("generator.CreateFiles.MkdirAll: %w", err)
			}
		}
		err := os.WriteFile(path.Join(file.Path, file.Name), file.Data, p.fileMode)
		if err != nil {
			return fmt.Errorf("generator.CreateFiles.WriteFile: %w", err)
		}
	}
	return nil
}

type RunCommandProc struct {
	commands []entity.Command
}

func NewRunCommandProc(commands []entity.Command) *RunCommandProc {
	return &RunCommandProc{
		commands: commands,
	}
}

func (p *RunCommandProc) Exec() error {
	for _, command := range p.commands {
		cmd := exec.Command(command.Cmd, command.Args...)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("executor.ExecCommands.Run: %w", err)
		}
		log.Printf("command: %s %v\nout: %s", command.Cmd, command.Args,
			out.String())
	}
	return nil
}

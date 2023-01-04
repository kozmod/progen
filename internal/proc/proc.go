package proc

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

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

type TextTemplateProc struct {
	fileMode  os.FileMode
	templates []entity.Template
}

func NewTextTemplateProc(tmpl []entity.Template) *TextTemplateProc {
	return &TextTemplateProc{
		fileMode:  os.ModePerm,
		templates: tmpl,
	}
}

func (p *TextTemplateProc) Exec() error {
	for _, tmpl := range p.templates {
		fileDir := tmpl.Path
		if _, err := os.Stat(fileDir); os.IsNotExist(err) {
			err = os.MkdirAll(fileDir, p.fileMode)
			if err != nil {
				return fmt.Errorf("create template dir [%s]: %w", fileDir, err)
			}
		}

		temp, err := template.New(tmpl.Name).Parse(string(tmpl.Text))
		if err != nil {
			return fmt.Errorf("create template: %w", err)
		}

		var buffer bytes.Buffer
		err = temp.Execute(&buffer, tmpl.Data)
		if err != nil {
			return fmt.Errorf("execure template [%s]: %w", tmpl.Name, err)
		}

		err = os.WriteFile(path.Join(tmpl.Path, tmpl.Name), buffer.Bytes(), p.fileMode)
		if err != nil {
			return fmt.Errorf("create template file [%s]: %w", tmpl.Name, err)
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

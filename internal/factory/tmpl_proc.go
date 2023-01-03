package factory

import (
	"path/filepath"

	"github.com/kozmod/progen/internal/config"
	"github.com/kozmod/progen/internal/entity"
	"github.com/kozmod/progen/internal/proc"
)

func NewTextTemplateProc(conf config.Config) (proc.Proc, error) {
	if len(conf.Templates) == 0 {
		return nil, nil
	}
	templates := toProcTemplates(conf.Templates)
	return proc.NewTextTemplateProc(templates), nil
}

func toProcTemplates(templates []config.Template) []entity.Template {
	procTemplates := make([]entity.Template, 0, len(templates))
	for _, t := range templates {
		tmpl := entity.Template{
			Name: filepath.Base(t.Path),
			Path: filepath.Dir(t.Path),
			Text: []byte(t.Text),
			Data: t.Data,
		}
		procTemplates = append(procTemplates, tmpl)
	}
	return procTemplates
}

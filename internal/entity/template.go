package entity

import (
	"bytes"
	"text/template"

	"golang.org/x/xerrors"
)

type TmplProc struct {
	templateData    map[string]any
	templateFns     map[string]any
	templateOptions []string
}

func NewTemplateProc(
	templateData,
	templateFns map[string]any,
	templateOptions []string) *TmplProc {
	return &TmplProc{
		templateData:    templateData,
		templateFns:     templateFns,
		templateOptions: templateOptions,
	}
}

func (p *TmplProc) Process(name, text string) (string, error) {
	tmpl, err := template.New(name).
		Funcs(p.templateFns).
		Option(p.templateOptions...).
		Parse(text)
	if err != nil {
		return Empty, xerrors.Errorf("process template: new template [%s]: %w", name, err)
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, p.templateData)
	if err != nil {
		return Empty, xerrors.Errorf("process template: execute [%s]: %w", name, err)
	}
	return buf.String(), err
}

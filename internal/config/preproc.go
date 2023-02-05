package config

import (
	"bytes"
	"fmt"
	"text/template"

	"gopkg.in/yaml.v3"

	"github.com/kozmod/progen/internal/entity"
)

type RawPreprocessor struct {
	templateName    string
	templateVars    map[string]any
	templateFns     map[string]any
	templateOptions []string
}

func NewRawPreprocessor(templateName string, templateVars, templateFns map[string]any, templateOptions []string) *RawPreprocessor {
	return &RawPreprocessor{
		templateName:    templateName,
		templateVars:    templateVars,
		templateFns:     templateFns,
		templateOptions: templateOptions,
	}
}

func (p *RawPreprocessor) Process(data []byte) ([]byte, map[string]any, error) {
	var (
		conf map[string]any
		name = p.templateName
	)

	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, nil, fmt.Errorf("parse config to map: %w", err)
	}

	conf = entity.MergeKeys(conf, p.templateVars)

	temp, err := template.New(name).
		Funcs(p.templateFns).
		Option(p.templateOptions...).
		Parse(string(data))
	if err != nil {
		return nil, nil, fmt.Errorf("new template [%s]: %w", name, err)
	}

	var buf bytes.Buffer
	err = temp.Execute(&buf, conf)
	if err != nil {
		return nil, nil, fmt.Errorf("execute template [%s]: %w", name, err)
	}
	return buf.Bytes(), conf, nil
}

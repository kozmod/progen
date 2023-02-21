package config

import (
	"fmt"

	yaml "gopkg.in/yaml.v3"

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

	res, err := entity.NewTemplateProc(conf, p.templateFns, p.templateOptions).Process(name, string(data))
	if err != nil {
		return nil, nil, fmt.Errorf("preprocess config: %w", err)
	}

	return []byte(res), conf, nil
}

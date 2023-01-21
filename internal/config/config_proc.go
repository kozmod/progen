package config

import (
	"bytes"
	"fmt"
	"text/template"

	"gopkg.in/yaml.v3"
)

func PreprocessRawConfigData(name string, data []byte) ([]byte, map[string]any, error) {
	var conf map[string]any
	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, nil, fmt.Errorf("parse config to map: %w", err)
	}

	temp, err := template.New(name).Parse(string(data))
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

func PrepareFiles(conf Config, files map[string][]File) (Config, error) {
	for i, fs := range conf.Files {
		f := files[fs.Tag]
		if cl, fl := len(conf.Files), len(files); cl != fl {
			return conf, fmt.Errorf("len of files is not match [%d:%d]: %s", cl, fl, fs.Tag)
		}
		for j, file := range fs.Val {
			if file.ExecTmplSkip {
				if file.Data == nil {
					continue
				}
				f := f[j]
				conf.Files[i].Val[j].Data = f.Data
			}
		}
	}
	return conf, nil
}

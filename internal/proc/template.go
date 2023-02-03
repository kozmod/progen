package proc

import (
	"bytes"
	"fmt"
	"text/template"
)

type templateExecutor struct {
	templateData map[string]any
	templateFns  map[string]any
}

func (e *templateExecutor) Exec(name string, data []byte) ([]byte, error) {
	temp, err := template.New(name).Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("new file template [%s]: %w", name, err)
	}

	var buf bytes.Buffer
	err = temp.Execute(&buf, e.templateData)
	if err != nil {
		return nil, fmt.Errorf("execute template [%s]: %w", name, err)
	}
	return buf.Bytes(), nil
}

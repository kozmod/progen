package config

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/kozmod/progen/internal/entity"
	"gopkg.in/yaml.v3"
)

func PreprocessRawConfigData(path string) ([]byte, map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read config: %w", err)
	}

	rawConf, mapConf, err := preprocessRawConfigData(path, data)
	if err != nil {
		return nil, nil, fmt.Errorf("preprocess raw config: %w", err)
	}
	return rawConf, mapConf, nil
}

func preprocessRawConfigData(name string, data []byte) ([]byte, map[string]any, error) {
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

func YamlRootNodesOrder(rowConfig []byte) (map[string]int, error) {
	root := yaml.Node{}
	err := yaml.Unmarshal(rowConfig, &root)
	if err != nil {
		return nil, fmt.Errorf("define order: unmarshall config: %w", err)
	}
	rootContent := root.Content
	if len(rootContent) != 1 {
		return nil, fmt.Errorf("define order: content is empty")
	}
	contentNode := rootContent[0]
	if contentNode == nil {
		return nil, fmt.Errorf("define order: config content is empty")
	}

	contentOrder := make(map[string]int, len(contentNode.Content))
	index := 0
	for _, node := range contentNode.Content {
		if node == nil {
			return nil, fmt.Errorf("define order: config content is empty: index [%d]", index)
		}
		value := strings.TrimSpace(node.Value)
		if value == entity.Empty {
			continue
		}
		contentOrder[node.Value] = index
		index++
	}

	return contentOrder, nil
}

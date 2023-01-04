package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/kozmod/progen/internal/entity"
	"gopkg.in/yaml.v3"
)

const (
	varStart           = "${"
	varEnd             = "}"
	SeparatorEqualSign = "="
)

type Vars map[string]string

func PreprocessRawConfigData(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	varsConf, err := UnmarshalYamlVarsConfig(data)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	existsVars, err := tryFindAllVars(data)
	if err != nil {
		return nil, fmt.Errorf("preprocess vars: %w", err)
	}
	configVars := parseVars(varsConf.Vars)
	for key := range existsVars {
		_, ok := configVars[key]
		if !ok {
			return nil, fmt.Errorf("preprocess vars: var is not set: %s", key)
		}
	}
	processedCof := replaceVars(data, configVars)
	return processedCof, nil
}

func tryFindAllVars(in []byte) (map[string]struct{}, error) {
	content := string(in)
	vars := make(map[string]struct{})
	for {
		start := strings.Index(content, varStart)
		if start == -1 {
			break
		}
		end := strings.Index(content, varEnd)
		trimStartIndex := start + 2
		trimEndIndex := end + 1

		nextStart := strings.Index(content[trimStartIndex:], varStart)
		switch {
		case end == -1:
			return nil, fmt.Errorf(
				"config with vars is invalid: config contains [%s], but not contains [%s]: position [%d]",
				varStart, varEnd, start)
		case nextStart != -1 && end >= nextStart+trimStartIndex:
			return nil, fmt.Errorf(
				"config with vars is invalid: next 'start' index [%d] lover then 'end' index[%d]",
				nextStart, end)
		}

		vars[content[trimStartIndex:end]] = struct{}{}
		content = content[trimEndIndex:]
	}
	return vars, nil
}

func parseVars(rawVars []string) Vars {
	varSet := make(Vars, len(rawVars))
	for _, v := range rawVars {
		split := strings.SplitN(v, SeparatorEqualSign, 2)
		key := strings.TrimSpace(split[0])
		if len(split) < 2 {
			varSet[key] = ""
			continue
		}
		varSet[key] = strings.TrimSpace(split[1])
	}
	return varSet
}

func replaceVars(rowConfig []byte, vars Vars) []byte {
	config := string(rowConfig)
	for key, val := range vars {
		config = strings.ReplaceAll(config, varStart+key+varEnd, val)
	}
	return []byte(config)
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

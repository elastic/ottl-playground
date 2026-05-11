/*
 * Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
 * or more contributor license agreements. See the NOTICE file distributed with
 * this work for additional information regarding copyright
 * ownership. Elasticsearch B.V. licenses this file to you under
 * the Apache License, Version 2.0 (the "License"); you may
 * not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package internal

import (
	"fmt"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/xconfmap"
	"gopkg.in/yaml.v3"
)

type parsedConfig[C any] struct {
	Key   string
	Value *C
}

func (p parsedConfig[C]) clone() (C, error) {
	var intermediate map[string]any
	err := mapstructure.Decode(p.Value, &intermediate)
	if err != nil {
		return *new(C), fmt.Errorf("failed to decode config: %w", err)
	}
	var dst C
	err = mapstructure.Decode(intermediate, &dst)
	if err != nil {
		return *new(C), fmt.Errorf("failed to decode config: %w", err)
	}
	return dst, err
}

func parseConfig[C any](id component.ID, yamlConfig string, createDefaultConfig func() *C) ([]parsedConfig[C], error) {
	deserializedYaml, err := confmap.NewRetrievedFromYAML([]byte(yamlConfig))
	if err != nil {
		return nil, err
	}

	deserializedConf, err := deserializedYaml.AsConf()
	if err != nil {
		return nil, err
	}

	if hasMultipleConfigs(id, deserializedConf) {
		return parseMultipleConfigs[C](yamlConfig, deserializedConf, createDefaultConfig)
	}

	defaultConfig := createDefaultConfig()
	err = unmarshalValidConfig(deserializedConf, defaultConfig)
	if err != nil {
		return nil, err
	}

	return []parsedConfig[C]{{"", defaultConfig}}, nil
}

func parseMultipleConfigs[C any](yamlConfig string, deserializedConf *confmap.Conf, createDefaultConfig func() *C) ([]parsedConfig[C], error) {
	var configs []parsedConfig[C]
	sortedKeys, err := sortedConfigKeys(yamlConfig)
	if err != nil {
		return nil, err
	}
	for _, configKey := range sortedKeys {
		sub, err := deserializedConf.Sub(configKey)
		if err != nil {
			return nil, err
		}

		defaultConfig := createDefaultConfig()
		err = unmarshalValidConfig(sub, defaultConfig)
		if err != nil {
			return nil, err
		}

		configs = append(configs, parsedConfig[C]{configKey, defaultConfig})
	}

	return configs, nil
}

func unmarshalValidConfig[C any](cfg *confmap.Conf, defaultConfig C) error {
	subConfigMap := map[string]any{}
	for k, v := range cfg.ToStringMap() {
		subConfigMap[k] = escapeDollarSigns(v)
	}

	err := confmap.NewFromStringMap(subConfigMap).Unmarshal(&defaultConfig)
	if err != nil {
		return err
	}

	validator, ok := any(defaultConfig).(xconfmap.Validator)
	if ok {
		return validator.Validate()
	}

	return nil
}

func sortedConfigKeys(yamlData string) ([]string, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(yamlData), &root); err != nil {
		return nil, fmt.Errorf("error unmarshalling YAML configuration: %v", err)
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, fmt.Errorf("unexpected YAML configuration structure")
	}
	mapNode := root.Content[0]
	if mapNode.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("expected configuration root to be a mapping node/configuration key")
	}
	var keys []string
	for i := 0; i < len(mapNode.Content); i += 2 {
		keys = append(keys, mapNode.Content[i].Value)
	}
	return keys, nil
}

func hasMultipleConfigs(id component.ID, conf *confmap.Conf) bool {
	allKeys := conf.AllKeys()
	if len(allKeys) == 0 {
		return false
	}
	for _, key := range allKeys {
		if !strings.HasPrefix(key, id.Type().String()) {
			return false
		}
	}
	return true
}

func escapeDollarSigns(val any) any {
	switch v := val.(type) {
	case string:
		return strings.ReplaceAll(v, "$$", "$")
	case []any:
		escapedVals := make([]any, len(v))
		for i, x := range v {
			escapedVals[i] = escapeDollarSigns(x)
		}
		return escapedVals
	case map[string]any:
		escapedMap := make(map[string]any, len(v))
		for k, x := range v {
			escapedMap[k] = escapeDollarSigns(x)
		}
		return escapedMap
	default:
		return val
	}
}

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
	"encoding/json"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"gopkg.in/yaml.v3"
)

type transformProcessorDebugger struct {
	consumer *processorConsumer[transformprocessor.Config]
}

func findYAMLPathIndex(yamlData string, configID, configKey string, configIndex int) (*yaml.Node, error) {
	var root yaml.Node
	if err := yaml.Unmarshal([]byte(yamlData), &root); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return nil, fmt.Errorf("unexpected YAML structure")
	}

	current := root.Content[0]
	parts := []string{configID, configKey}
	for _, key := range parts {
		if key == "" {
			continue
		}
		if current.Kind != yaml.MappingNode {
			return nil, fmt.Errorf("path '%s' not found (expected mapping node, got kind %d)", key, current.Kind)
		}
		found := false
		for i := 0; i < len(current.Content); i += 2 {
			k := current.Content[i]
			v := current.Content[i+1]
			if k.Value == key {
				current = v
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("'%s' not found in the configuration", key)
		}
	}

	if current.Kind == yaml.SequenceNode && current.Content[0].Kind == yaml.MappingNode {
		current = current.Content[configIndex]
	}

	if current.Kind == yaml.MappingNode {
		for i, c := range current.Content {
			if c.Value == "statements" {
				current = current.Content[i+1]
				break
			}
		}
	}

	return current, nil
}

type configWithLine struct {
	config transformprocessor.Config
	line   int64
}

func (t transformProcessorDebugger) DebugLogs(config, input string) (*Result, error) {
	configs, err := parseConfig[transformprocessor.Config](t.consumer.id, config, func() *transformprocessor.Config {
		return t.consumer.factory.CreateDefaultConfig().(*transformprocessor.Config)
	})
	if err != nil {
		return nil, err
	}

	res := &Result{}
	res.Debug = true
	var results []*Result
	for _, cfg := range configs {
		cfgs := make([]configWithLine, 0, len(configs)*len(cfg.Value.LogStatements))
		for i, contextStatements := range cfg.Value.LogStatements {
			pathIndex, err := findYAMLPathIndex(config, cfg.Key, "log_statements", i)
			if err != nil {
				return nil, fmt.Errorf("failed to find YAML path index: %w", err)
			}
			for i, sv := range contextStatements.Statements {
				cp, err := cfg.clone()
				if err != nil {
					return nil, err
				}
				contextStatementsCp := contextStatements
				contextStatementsCp.Statements = append(contextStatements.Statements[:i], []string{sv}...)
				cp.LogStatements = nil
				cp.LogStatements = append(cp.LogStatements, contextStatementsCp)
				cfgs = append(cfgs, configWithLine{cp, int64(pathIndex.Content[i].Line)})
			}
		}

		um := plog.JSONUnmarshaler{}
		inputLogs, err := um.UnmarshalLogs([]byte(input))
		if err != nil {
			return nil, err
		}

		ma := plog.JSONMarshaler{}
		for i, c := range cfgs {
			cpLogs := plog.NewLogs()
			inputLogs.CopyTo(cpLogs)

			result, err := newExecutionResult(t, ma.MarshalLogs, func() (plog.Logs, error) {
				return t.consumer.ConsumeLogs(&c.config, cpLogs)
			})
			if err != nil {
				return nil, err
			}
			result.Line = c.line
			results = append(results, result)
			if i+1 == len(cfgs) {
				input = result.Value
			}
		}
	}

	marshal, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	res.Value = string(marshal)
	return res, nil
}

func (t transformProcessorDebugger) DebugTraces(config, input string) (*Result, error) {
	configs, err := parseConfig[transformprocessor.Config](t.consumer.id, config, func() *transformprocessor.Config {
		return t.consumer.factory.CreateDefaultConfig().(*transformprocessor.Config)
	})
	if err != nil {
		return nil, err
	}

	res := &Result{}
	res.Debug = true
	var results []*Result
	for _, cfg := range configs {
		cfgs := make([]configWithLine, 0, len(configs)*len(cfg.Value.TraceStatements))
		for i, contextStatements := range cfg.Value.TraceStatements {
			pathIndex, err := findYAMLPathIndex(config, cfg.Key, "trace_statements", i)
			if err != nil {
				return nil, fmt.Errorf("failed to find YAML path index: %w", err)
			}
			for i, sv := range contextStatements.Statements {
				cp, err := cfg.clone()
				if err != nil {
					return nil, err
				}
				contextStatementsCp := contextStatements
				contextStatementsCp.Statements = append(contextStatements.Statements[:i], []string{sv}...)
				cp.TraceStatements = nil
				cp.TraceStatements = append(cp.TraceStatements, contextStatementsCp)
				cfgs = append(cfgs, configWithLine{cp, int64(pathIndex.Content[i].Line)})
			}
		}

		um := ptrace.JSONUnmarshaler{}
		inputTraces, err := um.UnmarshalTraces([]byte(input))
		if err != nil {
			return nil, err
		}

		ma := ptrace.JSONMarshaler{}
		for i, c := range cfgs {
			cpTraces := ptrace.NewTraces()
			inputTraces.CopyTo(cpTraces)

			result, err := newExecutionResult(t, ma.MarshalTraces, func() (ptrace.Traces, error) {
				return t.consumer.ConsumeTraces(&c.config, cpTraces)
			})
			if err != nil {
				return nil, err
			}
			result.Line = c.line
			results = append(results, result)
			if i+1 == len(cfgs) {
				input = result.Value
			}
		}
	}

	marshal, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	res.Value = string(marshal)
	return res, nil
}

func (t transformProcessorDebugger) DebugMetrics(config, input string) (*Result, error) {
	configs, err := parseConfig[transformprocessor.Config](t.consumer.id, config, func() *transformprocessor.Config {
		return t.consumer.factory.CreateDefaultConfig().(*transformprocessor.Config)
	})
	if err != nil {
		return nil, err
	}

	res := &Result{}
	res.Debug = true
	var results []*Result
	for _, cfg := range configs {
		cfgs := make([]configWithLine, 0, len(configs)*len(cfg.Value.MetricStatements))
		for i, contextStatements := range cfg.Value.MetricStatements {
			pathIndex, err := findYAMLPathIndex(config, cfg.Key, "metric_statements", i)
			if err != nil {
				return nil, fmt.Errorf("failed to find YAML path index: %w", err)
			}
			for i, sv := range contextStatements.Statements {
				cp, err := cfg.clone()
				if err != nil {
					return nil, err
				}
				contextStatementsCp := contextStatements
				contextStatementsCp.Statements = append(contextStatements.Statements[:i], []string{sv}...)
				cp.MetricStatements = nil
				cp.MetricStatements = append(cp.MetricStatements, contextStatementsCp)
				cfgs = append(cfgs, configWithLine{cp, int64(pathIndex.Content[i].Line)})
			}
		}

		um := pmetric.JSONUnmarshaler{}
		inputMetrics, err := um.UnmarshalMetrics([]byte(input))
		if err != nil {
			return nil, err
		}

		ma := pmetric.JSONMarshaler{}
		for i, c := range cfgs {
			cpMetrics := pmetric.NewMetrics()
			inputMetrics.CopyTo(cpMetrics)

			result, err := newExecutionResult(t, ma.MarshalMetrics, func() (pmetric.Metrics, error) {
				return t.consumer.ConsumeMetrics(&c.config, cpMetrics)
			})
			if err != nil {
				return nil, err
			}
			result.Line = c.line
			results = append(results, result)
			if i+1 == len(cfgs) {
				input = result.Value
			}
		}
	}

	marshal, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	res.Value = string(marshal)
	return res, nil
}

func (t transformProcessorDebugger) DebugProfiles(config, input string) (*Result, error) {
	configs, err := parseConfig[transformprocessor.Config](t.consumer.id, config, func() *transformprocessor.Config {
		return t.consumer.factory.CreateDefaultConfig().(*transformprocessor.Config)
	})
	if err != nil {
		return nil, err
	}

	res := &Result{}
	res.Debug = true
	var results []*Result
	for _, cfg := range configs {
		cfgs := make([]configWithLine, 0, len(configs)*len(cfg.Value.ProfileStatements))
		for i, contextStatements := range cfg.Value.ProfileStatements {
			pathIndex, err := findYAMLPathIndex(config, cfg.Key, "profile_statements", i)
			if err != nil {
				return nil, fmt.Errorf("failed to find YAML path index: %w", err)
			}
			for i, sv := range contextStatements.Statements {
				cp, err := cfg.clone()
				if err != nil {
					return nil, err
				}
				contextStatementsCp := contextStatements
				contextStatementsCp.Statements = append(contextStatements.Statements[:i], []string{sv}...)
				cp.ProfileStatements = nil
				cp.ProfileStatements = append(cp.ProfileStatements, contextStatementsCp)
				cfgs = append(cfgs, configWithLine{cp, int64(pathIndex.Content[i].Line)})
			}
		}

		um := pprofile.JSONUnmarshaler{}
		inputProfiles, err := um.UnmarshalProfiles([]byte(input))
		if err != nil {
			return nil, err
		}

		ma := pprofile.JSONMarshaler{}
		for i, c := range cfgs {
			cpProfiles := pprofile.NewProfiles()
			inputProfiles.CopyTo(cpProfiles)

			result, err := newExecutionResult(t, ma.MarshalProfiles, func() (pprofile.Profiles, error) {
				return t.consumer.ConsumeProfiles(&c.config, cpProfiles)
			})
			if err != nil {
				return nil, err
			}
			result.Line = c.line
			results = append(results, result)
			if i+1 == len(cfgs) {
				input = result.Value
			}
		}
	}

	marshal, err := json.Marshal(results)
	if err != nil {
		return nil, err
	}

	res.Value = string(marshal)
	return res, nil
}

func (t transformProcessorDebugger) ObservedLogs() *ObservedLogs {
	return t.consumer.ObservedLogs()
}

func NewTransformProcessorDebugger() Debugger {
	consumer := newProcessorConsumer[transformprocessor.Config](transformprocessor.NewFactory())
	return &transformProcessorDebugger{consumer}
}

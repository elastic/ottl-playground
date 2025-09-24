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
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const (
	transformprocessorMultipleConfig = "transformprocessor_multiple.yaml"
)

func Test_NewTransformProcessorDebugger(t *testing.T) {
	debugger := NewTransformProcessorDebugger()
	require.NotNil(t, debugger)

	// Verify it implements the Debugger interface
	_, ok := debugger.(*transformProcessorDebugger)
	require.True(t, ok)
}

func Test_transformProcessorDebugger_ObservedLogs(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)
	observedLogs := debugger.ObservedLogs()
	require.NotNil(t, observedLogs)

	// Test that it's the same instance as the consumer's observed logs
	assert.Equal(t, debugger.consumer.ObservedLogs(), observedLogs)
}

func Test_findYAMLPathIndex_SimpleConfig(t *testing.T) {
	yamlData := `trace_statements:
  - context: resource
    statements:
      - set(attributes["test"], true)
      - set(attributes["test2"], false)`

	node, err := findYAMLPathIndex(yamlData, "", "trace_statements", 0)
	require.NoError(t, err)
	require.NotNil(t, node)
	assert.Equal(t, 4, node.Line)

	// Should return the statements sequence
	assert.Equal(t, yaml.SequenceNode, node.Kind)
	assert.Len(t, node.Content, 2)
}

func Test_findYAMLPathIndex_WithConfigID(t *testing.T) {
	yamlData := `transform:
  trace_statements:
    - context: resource
      statements:
        - set(attributes["test"], true)
        - set(attributes["test2"], false)`

	node, err := findYAMLPathIndex(yamlData, "transform", "trace_statements", 0)
	require.NoError(t, err)
	require.NotNil(t, node)
	assert.Equal(t, 5, node.Line)

	// Should return the statements sequence
	assert.Equal(t, yaml.SequenceNode, node.Kind)
	assert.Len(t, node.Content, 2)
}

func Test_findYAMLPathIndex_WithFlatConfig(t *testing.T) {
	yamlData := `transform:
  trace_statements:
    - set(resource.attributes["test"], true)
    - set(resource.attributes["test2"], false)`

	node, err := findYAMLPathIndex(yamlData, "transform", "trace_statements", 0)
	require.NoError(t, err)
	require.NotNil(t, node)
	assert.Equal(t, 3, node.Line)

	// Should return the statements sequence
	assert.Equal(t, yaml.SequenceNode, node.Kind)
	assert.Len(t, node.Content, 2)
}

func Test_findYAMLPathIndex_WithMultipleConfigs(t *testing.T) {
	yamlData := `transform/first:
  trace_statements:
    - set(resource.attributes["test"], true)
    - set(resource.attributes["test2"], false)
transform/second:
  log_statements:
    - context: resource
      statements:
        - set(attributes["test2"], false)
`

	firstConfig, err := findYAMLPathIndex(yamlData, "transform/first", "trace_statements", 0)
	require.NoError(t, err)
	require.NotNil(t, firstConfig)
	assert.Equal(t, 3, firstConfig.Line)
	assert.Equal(t, yaml.SequenceNode, firstConfig.Kind)
	assert.Len(t, firstConfig.Content, 2)

	secondConfig, err := findYAMLPathIndex(yamlData, "transform/second", "log_statements", 0)
	require.NoError(t, err)
	require.NotNil(t, secondConfig)
	assert.Equal(t, 9, secondConfig.Line)
	assert.Equal(t, yaml.SequenceNode, secondConfig.Kind)
	assert.Len(t, secondConfig.Content, 1)

}

func Test_findYAMLPathIndex_InvalidYAML(t *testing.T) {
	yamlData := `invalid: yaml: [unclosed`

	_, err := findYAMLPathIndex(yamlData, "", "trace_statements", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse YAML")
}

func Test_findYAMLPathIndex_ConfigNotFound(t *testing.T) {
	yamlData := `trace_statements:
  - context: resource
    statements:
      - set(attributes["test"], true)`

	_, err := findYAMLPathIndex(yamlData, "nonexistent", "trace_statements", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'nonexistent' not found in the configuration")
}

func Test_findYAMLPathIndex_KeyNotFound(t *testing.T) {
	yamlData := `trace_statements:
  - context: resource
    statements:
      - set(attributes["test"], true)`

	_, err := findYAMLPathIndex(yamlData, "", "nonexistent_statements", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "'nonexistent_statements' not found in the configuration")
}

func Test_findYAMLPathIndex_EmptyDocument(t *testing.T) {
	yamlData := ``

	_, err := findYAMLPathIndex(yamlData, "", "trace_statements", 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected YAML structure")
}

func Test_transformProcessorDebugger_DebugLogs(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "logs.json")

	result, err := debugger.DebugLogs(config, payload)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify debug result structure
	assert.True(t, result.Debug)
	assert.NotEmpty(t, result.Value)

	// Parse the result JSON to verify it's valid
	var debugResults []*Result
	err = json.Unmarshal([]byte(result.Value), &debugResults)
	require.NoError(t, err)
	assert.NotEmpty(t, debugResults)

	// Verify each debug result has proper structure
	for _, debugResult := range debugResults {
		assert.NotEmpty(t, debugResult.Value)
		assert.Greater(t, debugResult.Line, int64(0))
	}
}

func Test_transformProcessorDebugger_DebugTraces(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "traces.json")

	result, err := debugger.DebugTraces(config, payload)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify debug result structure
	assert.True(t, result.Debug)
	assert.NotEmpty(t, result.Value)

	// Parse the result JSON to verify it's valid
	var debugResults []*Result
	err = json.Unmarshal([]byte(result.Value), &debugResults)
	require.NoError(t, err)
	assert.NotEmpty(t, debugResults)

	// Verify each debug result has proper structure
	for _, debugResult := range debugResults {
		assert.NotEmpty(t, debugResult.Value)
		assert.Greater(t, debugResult.Line, int64(0))
	}
}

func Test_transformProcessorDebugger_DebugMetrics(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "metrics.json")

	result, err := debugger.DebugMetrics(config, payload)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify debug result structure
	assert.True(t, result.Debug)
	assert.NotEmpty(t, result.Value)

	// Parse the result JSON to verify it's valid
	var debugResults []*Result
	err = json.Unmarshal([]byte(result.Value), &debugResults)
	require.NoError(t, err)
	assert.NotEmpty(t, debugResults)

	// Verify each debug result has proper structure
	for _, debugResult := range debugResults {
		assert.NotEmpty(t, debugResult.Value)
		assert.Greater(t, debugResult.Line, int64(0))
	}
}

func Test_transformProcessorDebugger_DebugLogs_MultipleConfigs(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)
	config := readTestData(t, transformprocessorMultipleConfig)
	payload := readTestData(t, "logs.json")

	result, err := debugger.DebugLogs(config, payload)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify debug result structure
	assert.True(t, result.Debug)
	assert.NotEmpty(t, result.Value)

	// Parse the result JSON to verify it's valid
	var debugResults []*Result
	err = json.Unmarshal([]byte(result.Value), &debugResults)
	require.NoError(t, err)
	assert.NotEmpty(t, debugResults)
}

func Test_transformProcessorDebugger_DebugLogs_InvalidConfig(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)
	config := "invalid: yaml: [unclosed"
	payload := readTestData(t, "logs.json")

	_, err := debugger.DebugLogs(config, payload)
	assert.Error(t, err)
}

func Test_transformProcessorDebugger_DebugLogs_InvalidPayload(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)
	config := readTestData(t, transformprocessorConfig)
	payload := "invalid json"

	_, err := debugger.DebugLogs(config, payload)
	assert.Error(t, err)
}

func Test_transformProcessorDebugger_DebugTraces_InvalidPayload(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)
	config := readTestData(t, transformprocessorConfig)
	payload := "invalid json"

	_, err := debugger.DebugTraces(config, payload)
	assert.Error(t, err)
}

func Test_transformProcessorDebugger_DebugMetrics_InvalidPayload(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)
	config := readTestData(t, transformprocessorConfig)
	payload := "invalid json"

	_, err := debugger.DebugMetrics(config, payload)
	assert.Error(t, err)
}

func Test_transformProcessorDebugger_DebugWithComplexConfig(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)

	// Create a more complex config with multiple statements
	config := `log_statements:
  - context: resource
    statements:
      - set(attributes["step1"], true)
      - set(attributes["step2"], "value")
      - set(attributes["step3"], 42)`

	payload := readTestData(t, "logs.json")

	result, err := debugger.DebugLogs(config, payload)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Parse the result JSON
	var debugResults []*Result
	err = json.Unmarshal([]byte(result.Value), &debugResults)
	require.NoError(t, err)

	// Should have 3 debug results (one for each statement)
	assert.Len(t, debugResults, 3)

	// Verify line numbers are set
	for _, debugResult := range debugResults {
		assert.Greater(t, debugResult.Line, int64(0))
	}
}

func Test_configWithLine_Structure(t *testing.T) {
	config := &transformprocessor.Config{}
	line := int64(42)

	cwl := configWithLine{
		config: *config,
		line:   line,
	}

	assert.NotNil(t, cwl.config)
	assert.Equal(t, line, cwl.line)
}

func Test_transformProcessorDebugger_DebugLogs_LineNumberAccuracy(t *testing.T) {
	debugger := NewTransformProcessorDebugger().(*transformProcessorDebugger)

	// Multi-line config with specific line structure
	config := `log_statements:
  - context: resource
    statements:
      - set(attributes["line4"], true)
      - set(attributes["line5"], false)`

	payload := readTestData(t, "logs.json")

	result, err := debugger.DebugLogs(config, payload)
	require.NoError(t, err)
	require.NotNil(t, result)

	var debugResults []*Result
	err = json.Unmarshal([]byte(result.Value), &debugResults)
	require.NoError(t, err)
	require.Len(t, debugResults, 2)

	// Verify line numbers correspond to statements (should be line 4 and 5)
	assert.Equal(t, int64(4), debugResults[0].Line)
	assert.Equal(t, int64(5), debugResults[1].Line)
}

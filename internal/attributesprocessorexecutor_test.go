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
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	attributesprocessorConfig = "attributesprocessor.yaml"
)

func Test_AttributeProcessorExecutor_ParseConfig(t *testing.T) {
	yamlConfig := readTestData(t, attributesprocessorConfig)
	cfgs, err := parseConfig[attributesprocessor.Config](
		component.NewIDWithName(attributesprocessor.NewFactory().Type(), "test_attributes_processor"),
		yamlConfig,
		func() *attributesprocessor.Config {
			return attributesprocessor.NewFactory().CreateDefaultConfig().(*attributesprocessor.Config)
		},
	)
	require.NoError(t, err)

	pc := cfgs[0].Value
	require.NotNil(t, pc)
	require.NotEmpty(t, pc.Actions)
}

func Test_AttributeProcessorExecutor_ExecuteLogs(t *testing.T) {
	executor := NewAttributesProcessorExecutor()
	config := readTestData(t, attributesprocessorConfig)
	payload := readTestData(t, "logs.json")

	output, err := executor.ExecuteLogs(config, payload)
	require.NoError(t, err)

	unmarshaler := &plog.JSONUnmarshaler{}
	outputLogs, err := unmarshaler.UnmarshalLogs([]byte(output.Value))
	require.NoError(t, err)
	require.NotNil(t, outputLogs)
}

func Test_AttributeProcessorExecutor_ExecuteTraces(t *testing.T) {
	executor := NewAttributesProcessorExecutor()
	config := readTestData(t, attributesprocessorConfig)
	payload := readTestData(t, "traces.json")

	output, err := executor.ExecuteTraces(config, payload)
	require.NoError(t, err)

	unmarshaler := &ptrace.JSONUnmarshaler{}
	outputTraces, err := unmarshaler.UnmarshalTraces([]byte(output.Value))
	require.NoError(t, err)
	require.NotNil(t, outputTraces)
}

func Test_AttributeProcessorExecutor_ExecuteMetrics(t *testing.T) {
	executor := NewAttributesProcessorExecutor()
	config := readTestData(t, attributesprocessorConfig)
	payload := readTestData(t, "metrics.json")

	output, err := executor.ExecuteMetrics(config, payload)
	require.NoError(t, err)

	unmarshaler := &pmetric.JSONUnmarshaler{}
	outputMetrics, err := unmarshaler.UnmarshalMetrics([]byte(output.Value))
	require.NoError(t, err)
	require.NotNil(t, outputMetrics)
}

func Test_AttributeProcessorExecutor_ObservedLogs(t *testing.T) {
	executor := NewAttributesProcessorExecutor().(*defaultExecutor[attributesprocessor.Config])
	executor.consumer.TelemetrySettings().Logger.Sugar().Debug("this is a log")
	logEntries := executor.ObservedLogs().TakeAll()

	assert.Len(t, logEntries, 1)
	assert.Contains(t, logEntries[0].ConsoleEncodedEntry(), "this is a log")
}
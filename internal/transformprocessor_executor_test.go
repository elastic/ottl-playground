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

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	transformprocessorConfig = "transformprocessor.yaml"
)

func Test_TransformProcessorExecutor_parseConfig(t *testing.T) {
	yamlConfig := readTestData(t, transformprocessorConfig)
	cfgs, err := parseConfig[transformprocessor.Config](
		component.NewIDWithName(transformprocessor.NewFactory().Type(), "test_transform_processor"),
		yamlConfig,
		func() *transformprocessor.Config {
			return transformprocessor.NewFactory().CreateDefaultConfig().(*transformprocessor.Config)
		},
	)

	require.NoError(t, err)
	require.Len(t, cfgs, 1)

	pc := cfgs[0].Value
	require.NotNil(t, pc)
	require.NotEmpty(t, pc.ErrorMode)
	require.NotEmpty(t, pc.TraceStatements)
	require.NotEmpty(t, pc.MetricStatements)
	require.NotEmpty(t, pc.MetricStatements)
}

func Test_TransformProcessorExecutor_ExecuteLogs(t *testing.T) {
	executor := NewTransformProcessorExecutor()
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "logs.json")

	output, err := executor.ExecuteLogs(config, payload)
	require.NoError(t, err)

	unmarshaler := &plog.JSONUnmarshaler{}
	outputLogs, err := unmarshaler.UnmarshalLogs([]byte(output.Value))
	require.NoError(t, err)
	require.NotNil(t, outputLogs)

	val, ok := outputLogs.ResourceLogs().At(0).Resource().Attributes().Get("log_statements")
	require.True(t, ok)
	require.True(t, val.Bool())
}

func Test_TransformProcessorExecutor_ExecuteTraces(t *testing.T) {
	executor := NewTransformProcessorExecutor()
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "traces.json")

	output, err := executor.ExecuteTraces(config, payload)
	require.NoError(t, err)

	unmarshaler := &ptrace.JSONUnmarshaler{}
	outputTraces, err := unmarshaler.UnmarshalTraces([]byte(output.Value))
	require.NoError(t, err)
	require.NotNil(t, outputTraces)

	val, ok := outputTraces.ResourceSpans().At(0).Resource().Attributes().Get("trace_statements")
	require.True(t, ok)
	require.True(t, val.Bool())
}

func Test_TransformProcessorExecutor_ExecuteMetrics(t *testing.T) {
	executor := NewTransformProcessorExecutor()
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "metrics.json")

	output, err := executor.ExecuteMetrics(config, payload)
	require.NoError(t, err)

	unmarshaler := &pmetric.JSONUnmarshaler{}
	outputMetrics, err := unmarshaler.UnmarshalMetrics([]byte(output.Value))
	require.NoError(t, err)
	require.NotNil(t, outputMetrics)

	val, ok := outputMetrics.ResourceMetrics().At(0).Resource().Attributes().Get("metric_statements")
	require.True(t, ok)
	require.True(t, val.Bool())
}

func Test_TransformProcessorExecutor_ObservedLogs(t *testing.T) {
	executor := NewJSONExecutor[transformprocessor.Config](
		newProcessorConsumer[transformprocessor.Config](transformprocessor.NewFactory()),
		&Metadata{},
	).(*defaultExecutor[transformprocessor.Config])

	executor.consumer.TelemetrySettings().Logger.Sugar().Debug("this is a log")
	logEntries := executor.ObservedLogs().TakeAll()
	assert.Len(t, logEntries, 1)
	assert.Contains(t, logEntries[0].ConsoleEncodedEntry(), "this is a log")
}

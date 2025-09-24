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
	"os"
	"path/filepath"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/stretchr/testify/require"
)

const (
	transformprocessorConfigMultiple = "transformprocessor_multiple.yaml"
)

func readTestData(t *testing.T, file string) string {
	content, err := os.ReadFile(filepath.Join("..", "testdata", file))
	require.NoError(t, err)
	return string(content)
}

func Test_Executor_ExecuteLogsMultiple(t *testing.T) {
	executor := NewJSONExecutor[transformprocessor.Config](
		newProcessorConsumer[transformprocessor.Config](transformprocessor.NewFactory()),
		&Metadata{},
	)

	config := readTestData(t, transformprocessorConfigMultiple)
	payload := readTestData(t, "logs.json")

	output, err := executor.ExecuteLogs(config, payload)
	require.NoError(t, err)

	unmarshaler := &plog.JSONUnmarshaler{}
	outputLogs, err := unmarshaler.UnmarshalLogs([]byte(output.Value))

	require.NoError(t, err)
	require.NotNil(t, outputLogs)

	val, ok := outputLogs.ResourceLogs().At(0).Resource().Attributes().Get("log_statements")
	assert.True(t, ok)
	assert.True(t, val.Bool())
	assert.Contains(t, output.Logs, "configuration: transform")

	val, ok = outputLogs.ResourceLogs().At(0).Resource().Attributes().Get("log_multiple")
	assert.True(t, ok)
	assert.True(t, val.Bool())
	assert.Contains(t, output.Logs, "configuration: transform/multiple")
}

func Test_Executor_ExecuteTracesMultiple(t *testing.T) {
	executor := NewJSONExecutor[transformprocessor.Config](
		newProcessorConsumer[transformprocessor.Config](transformprocessor.NewFactory()),
		&Metadata{},
	)
	config := readTestData(t, transformprocessorConfigMultiple)
	payload := readTestData(t, "traces.json")

	output, err := executor.ExecuteTraces(config, payload)
	require.NoError(t, err)

	unmarshaler := &ptrace.JSONUnmarshaler{}
	outputTraces, err := unmarshaler.UnmarshalTraces([]byte(output.Value))

	require.NoError(t, err)
	require.NotNil(t, outputTraces)

	val, ok := outputTraces.ResourceSpans().At(0).Resource().Attributes().Get("trace_statements")
	assert.True(t, ok)
	assert.True(t, val.Bool())
	assert.Contains(t, output.Logs, "configuration: transform")

	val, ok = outputTraces.ResourceSpans().At(0).Resource().Attributes().Get("trace_multiple")
	assert.True(t, ok)
	assert.True(t, val.Bool())
	assert.Contains(t, output.Logs, "configuration: transform/multiple")
}

func Test_Executor_ExecuteMetricsMultiple(t *testing.T) {
	executor := NewJSONExecutor[transformprocessor.Config](
		newProcessorConsumer[transformprocessor.Config](transformprocessor.NewFactory()),
		&Metadata{},
	)
	config := readTestData(t, transformprocessorConfigMultiple)
	payload := readTestData(t, "metrics.json")

	output, err := executor.ExecuteMetrics(config, payload)
	require.NoError(t, err)

	unmarshaler := &pmetric.JSONUnmarshaler{}
	outputMetrics, err := unmarshaler.UnmarshalMetrics([]byte(output.Value))

	require.NoError(t, err)
	require.NotNil(t, outputMetrics)

	val, ok := outputMetrics.ResourceMetrics().At(0).Resource().Attributes().Get("metric_statements")
	assert.True(t, ok)
	assert.True(t, val.Bool())
	assert.Contains(t, output.Logs, "configuration: transform")

	val, ok = outputMetrics.ResourceMetrics().At(0).Resource().Attributes().Get("metric_multiple")
	assert.True(t, ok)
	assert.True(t, val.Bool())
	assert.Contains(t, output.Logs, "configuration: transform/multiple")
}

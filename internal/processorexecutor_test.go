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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
)

func Test_ParseConfig_Success(t *testing.T) {
	executor := newProcessorExecutor[transformprocessor.Config](transformprocessor.NewFactory())
	config := readTestData(t, transformprocessorConfig)
	parsedConfig, err := executor.parseConfig(config)

	require.NoError(t, err)
	require.NotNil(t, parsedConfig)
	require.NotEmpty(t, parsedConfig.ErrorMode)
	require.NotEmpty(t, parsedConfig.TraceStatements)
	require.NotEmpty(t, parsedConfig.MetricStatements)
	require.NotEmpty(t, parsedConfig.MetricStatements)
}

func Test_ParseConfig_Error(t *testing.T) {
	executor := newProcessorExecutor[transformprocessor.Config](transformprocessor.NewFactory())
	_, err := executor.parseConfig("---invalid---")
	require.ErrorContains(t, err, "cannot be used as a Conf")
}

func Test_ExecuteLogStatements(t *testing.T) {
	executor := newProcessorExecutor[transformprocessor.Config](transformprocessor.NewFactory())
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "logs.json")

	output, err := executor.ExecuteLogStatements(config, payload)
	require.NoError(t, err)

	unmarshaler := &plog.JSONUnmarshaler{}
	outputLogs, err := unmarshaler.UnmarshalLogs(output)
	require.NoError(t, err)
	require.NotNil(t, outputLogs)
}

func Test_ExecuteTraceStatements(t *testing.T) {
	executor := newProcessorExecutor[transformprocessor.Config](transformprocessor.NewFactory())
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "traces.json")

	output, err := executor.ExecuteTraceStatements(config, payload)
	require.NoError(t, err)

	unmarshaler := &ptrace.JSONUnmarshaler{}
	outputTraces, err := unmarshaler.UnmarshalTraces(output)
	require.NoError(t, err)
	require.NotNil(t, outputTraces)
}

func Test_ExecuteMetricStatements(t *testing.T) {
	executor := newProcessorExecutor[transformprocessor.Config](transformprocessor.NewFactory())
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "metrics.json")

	output, err := executor.ExecuteMetricStatements(config, payload)
	require.NoError(t, err)

	unmarshaler := &pmetric.JSONUnmarshaler{}
	outputMetrics, err := unmarshaler.UnmarshalMetrics(output)
	require.NoError(t, err)
	require.NotNil(t, outputMetrics)
}

func Test_ObservedLogs(t *testing.T) {
	executor := newProcessorExecutor[transformprocessor.Config](transformprocessor.NewFactory())
	executor.settings.Logger.Sugar().Debug("this is a log")
	logEntries := executor.ObservedLogs().TakeAll()
	assert.Len(t, logEntries, 1)
	assert.Contains(t, logEntries[0].ConsoleEncodedEntry(), "this is a log")
}

func readTestData(t *testing.T, file string) string {
	content, err := os.ReadFile(filepath.Join("..", "testdata", file))
	require.NoError(t, err)
	return string(content)
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

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
	assert.Equal(t, "this is a log", logEntries[0].Message)
}

func readTestData(t *testing.T, file string) string {
	content, err := os.ReadFile(filepath.Join("..", "testdata", file))
	require.NoError(t, err)
	return string(content)
}

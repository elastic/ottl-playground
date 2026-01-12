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

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

const (
	filterprocessorConfig = "filterprocessor.yaml"
)

func Test_FilterProcessorExecutor_parseConfig(t *testing.T) {
	yamlConfig := readTestData(t, filterprocessorConfig)
	cfgs, err := parseConfig[filterprocessor.Config](
		component.NewIDWithName(filterprocessor.NewFactory().Type(), "test_filter_processor"),
		yamlConfig,
		func() *filterprocessor.Config {
			return filterprocessor.NewFactory().CreateDefaultConfig().(*filterprocessor.Config)
		},
	)
	require.NoError(t, err)

	pc := cfgs[0].Value
	require.NotNil(t, pc)
	require.NotEmpty(t, pc.ErrorMode)
	require.NotEmpty(t, pc.Logs)
	require.NotEmpty(t, pc.Traces)
	require.NotEmpty(t, pc.Metrics)
}

func Test_FilterProcessorExecutor_ExecuteLogs(t *testing.T) {
	executor := NewFilterProcessorExecutor()
	config := readTestData(t, filterprocessorConfig)
	payload := readTestData(t, "logs.json")

	output, err := executor.ExecuteLogs(config, payload)
	require.NoError(t, err)

	unmarshaler := &plog.JSONUnmarshaler{}
	outputLogs, err := unmarshaler.UnmarshalLogs([]byte(output.Value))

	require.NoError(t, err)
	require.NotNil(t, outputLogs)
	assert.Equal(t, outputLogs.LogRecordCount(), 1)
}

func Test_FilterProcessorExecutor_ExecuteTraces(t *testing.T) {
	executor := NewFilterProcessorExecutor()
	config := readTestData(t, filterprocessorConfig)
	payload := readTestData(t, "traces.json")

	output, err := executor.ExecuteTraces(config, payload)
	require.NoError(t, err)

	unmarshaler := &ptrace.JSONUnmarshaler{}
	outputTraces, err := unmarshaler.UnmarshalTraces([]byte(output.Value))
	require.NoError(t, err)
	require.NotNil(t, outputTraces)

	scopeSpans := outputTraces.ResourceSpans().At(0).ScopeSpans()
	assert.Equal(t, 1, scopeSpans.Len())
	assert.Equal(t, "eee19b7ec3c1b174", scopeSpans.At(0).Spans().At(0).SpanID().String())
}

func Test_FilterProcessorExecutor_ExecuteMetrics(t *testing.T) {
	executor := NewFilterProcessorExecutor()
	config := readTestData(t, filterprocessorConfig)
	payload := readTestData(t, "metrics.json")

	output, err := executor.ExecuteMetrics(config, payload)
	require.NoError(t, err)

	unmarshaler := &pmetric.JSONUnmarshaler{}
	outputMetrics, err := unmarshaler.UnmarshalMetrics([]byte(output.Value))
	require.NoError(t, err)
	require.NotNil(t, outputMetrics)

	metrics := outputMetrics.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
	for _, v := range metrics.All() {
		require.NotEqual(t, "my.counter", v.Name())
	}
}

func Test_FilterProcessorExecutor_ObservedLogs(t *testing.T) {
	executor := NewFilterProcessorExecutor().(*defaultExecutor[filterprocessor.Config])
	executor.consumer.TelemetrySettings().Logger.Sugar().Debug("this is a log")
	logEntries := executor.ObservedLogs().TakeAll()

	assert.Len(t, logEntries, 1)
	assert.Contains(t, logEntries[0].ConsoleEncodedEntry(), "this is a log")
}

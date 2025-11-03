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
	"go.opentelemetry.io/collector/pdata/pprofile"
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

	var count int
	for _, rl := range outputLogs.ResourceLogs().All() {
		for _, sl := range rl.ScopeLogs().All() {
			for _, lr := range sl.LogRecords().All() {
				val, ok := lr.Attributes().Get("log_statements")
				if assert.True(t, ok) && assert.Equal(t, "log_value", val.Str()) {
					count++
				}
			}
		}
	}

	require.Positive(t, count)
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

	var count int
	for _, rs := range outputTraces.ResourceSpans().All() {
		for _, ss := range rs.ScopeSpans().All() {
			for _, sp := range ss.Spans().All() {
				val, ok := sp.Attributes().Get("trace_statements")
				if assert.True(t, ok) && assert.Equal(t, "trace_value", val.Str()) {
					count++
				}
			}
		}
	}

	require.Positive(t, count)
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

	var count int
	for _, rm := range outputMetrics.ResourceMetrics().All() {
		for _, sm := range rm.ScopeMetrics().All() {
			for _, me := range sm.Metrics().All() {
				val, ok := me.Metadata().Get("metric_statements")
				if assert.True(t, ok) && assert.Equal(t, "metric_value", val.Str()) {
					count++
				}
			}
		}
	}

	require.Positive(t, count)
}

func Test_TransformProcessorExecutor_ExecuteProfiles(t *testing.T) {
	executor := NewTransformProcessorExecutor()
	config := readTestData(t, transformprocessorConfig)
	payload := readTestData(t, "profiles.json")

	output, err := executor.ExecuteProfiles(config, payload)
	require.NoError(t, err)

	unmarshaler := &pprofile.JSONUnmarshaler{}
	outputProfiles, err := unmarshaler.UnmarshalProfiles([]byte(output.Value))
	require.NoError(t, err)
	require.NotNil(t, outputProfiles)

	var count int
	for _, rp := range outputProfiles.ResourceProfiles().All() {
		for _, sp := range rp.ScopeProfiles().All() {
			for _, pr := range sp.Profiles().All() {
				attrs := pprofile.FromAttributeIndices(outputProfiles.Dictionary().AttributeTable(), pr, outputProfiles.Dictionary())
				val, ok := attrs.Get("profile_statements")
				if assert.True(t, ok) && assert.Equal(t, "profile_value", val.Str()) {
					count++
				}
			}
		}
	}

	require.Positive(t, count)
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

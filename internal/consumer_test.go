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
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func Test_newProcessorConsumer(t *testing.T) {
	factory := transformprocessor.NewFactory()
	consumer := newProcessorConsumer[transformprocessor.Config](factory)

	require.NotNil(t, consumer)
	assert.Equal(t, "transform/ottl_playground", consumer.ComponentID().String())
	assert.NotNil(t, consumer.TelemetrySettings())
	assert.NotNil(t, consumer.observedLogs)
	assert.Equal(t, factory, consumer.factory)
}

func Test_processorConsumer_CreateDefaultConfig(t *testing.T) {
	factory := transformprocessor.NewFactory()
	consumer := newProcessorConsumer[transformprocessor.Config](factory)

	config := consumer.CreateDefaultConfig()
	require.NotNil(t, config)

	expectedConfig := factory.CreateDefaultConfig().(*transformprocessor.Config)
	assert.Equal(t, expectedConfig.ErrorMode, config.ErrorMode)
}

func Test_processorConsumer_ConsumeLogs(t *testing.T) {
	factory := transformprocessor.NewFactory()
	consumer := newProcessorConsumer[transformprocessor.Config](factory)

	// Create test input logs
	inputLogs := plog.NewLogs()
	resourceLogs := inputLogs.ResourceLogs().AppendEmpty()
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
	logRecord := scopeLogs.LogRecords().AppendEmpty()
	logRecord.Body().SetStr("test log message")
	logRecord.SetSeverityText("INFO")

	// Create a basic config
	config := &transformprocessor.Config{}

	outputLogs, err := consumer.ConsumeLogs(config, inputLogs)
	require.NoError(t, err)
	require.NotNil(t, outputLogs)
	assert.Equal(t, inputLogs.LogRecordCount(), outputLogs.LogRecordCount())
}

func Test_processorConsumer_ConsumeMetrics(t *testing.T) {
	factory := transformprocessor.NewFactory()
	consumer := newProcessorConsumer[transformprocessor.Config](factory)

	// Create test input metrics
	inputMetrics := pmetric.NewMetrics()
	resourceMetrics := inputMetrics.ResourceMetrics().AppendEmpty()
	scopeMetrics := resourceMetrics.ScopeMetrics().AppendEmpty()
	metric := scopeMetrics.Metrics().AppendEmpty()
	metric.SetName("test.metric")

	// Create a basic config
	config := &transformprocessor.Config{}

	outputMetrics, err := consumer.ConsumeMetrics(config, inputMetrics)
	require.NoError(t, err)
	require.NotNil(t, outputMetrics)
	assert.Equal(t, inputMetrics.MetricCount(), outputMetrics.MetricCount())
}

func Test_processorConsumer_ConsumeTraces(t *testing.T) {
	factory := transformprocessor.NewFactory()
	consumer := newProcessorConsumer[transformprocessor.Config](factory)

	// Create test input traces
	inputTraces := ptrace.NewTraces()
	resourceSpans := inputTraces.ResourceSpans().AppendEmpty()
	scopeSpans := resourceSpans.ScopeSpans().AppendEmpty()
	span := scopeSpans.Spans().AppendEmpty()
	span.SetName("test-span")
	span.SetTraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	span.SetSpanID([8]byte{1, 2, 3, 4, 5, 6, 7, 8})

	// Create a basic config
	config := &transformprocessor.Config{}

	outputTraces, err := consumer.ConsumeTraces(config, inputTraces)
	require.NoError(t, err)
	require.NotNil(t, outputTraces)
	assert.Equal(t, inputTraces.SpanCount(), outputTraces.SpanCount())
}

func Test_processorConsumer_ObservedLogs(t *testing.T) {
	factory := transformprocessor.NewFactory()
	consumer := newProcessorConsumer[transformprocessor.Config](factory)

	observedLogs := consumer.ObservedLogs()
	require.NotNil(t, observedLogs)

	// Test that it's the same instance
	assert.Equal(t, consumer.observedLogs, observedLogs)
}

func Test_processorConsumer_TelemetrySettings(t *testing.T) {
	factory := transformprocessor.NewFactory()
	consumer := newProcessorConsumer[transformprocessor.Config](factory)

	telemetrySettings := consumer.TelemetrySettings()
	require.NotNil(t, telemetrySettings)
	require.NotNil(t, telemetrySettings.Logger)

	// Test that it's the same instance
	assert.Equal(t, consumer.telemetrySettings, telemetrySettings)
}

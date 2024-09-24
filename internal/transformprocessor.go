// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
)

type transformProcessorExecutor struct {
	factory           processor.Factory
	settings          processor.Settings
	telemetrySettings component.TelemetrySettings
	errorMode         ottl.ErrorMode
}

func NewTransformProcessorExecutor() Executor {
	telemetrySettings := componenttest.NewNopTelemetrySettings()
	settings := processor.Settings{
		ID:                component.MustNewID("ottl_playground"),
		TelemetrySettings: telemetrySettings,
	}

	return &transformProcessorExecutor{
		factory:           transformprocessor.NewFactory(),
		telemetrySettings: telemetrySettings,
		settings:          settings,
	}
}

func (p *transformProcessorExecutor) parseConfig(yamlConfig string) (*transformprocessor.Config, error) {
	deserializedYaml, err := confmap.NewRetrievedFromYAML([]byte(yamlConfig))
	if err != nil {
		return nil, err
	}

	yamlConfigMap, err := deserializedYaml.AsConf()
	if err != nil {
		return nil, err
	}

	defaultConfig := p.factory.CreateDefaultConfig().(*transformprocessor.Config)
	err = yamlConfigMap.Unmarshal(&defaultConfig)
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

func (p *transformProcessorExecutor) ExecuteLogStatements(yamlConfig, input string) ([]byte, error) {
	config, err := p.parseConfig(yamlConfig)
	if err != nil {
		return nil, err
	}

	var transformedLogs plog.Logs
	logsConsumer, _ := consumer.NewLogs(func(ctx context.Context, ld plog.Logs) error {
		transformedLogs = ld
		return nil
	})

	logsProcessor, err := p.factory.CreateLogsProcessor(context.Background(), p.settings, config, logsConsumer)
	if err != nil {
		return nil, err
	}

	logsUnmarshaler := &plog.JSONUnmarshaler{}
	inputLogs, err := logsUnmarshaler.UnmarshalLogs([]byte(input))
	if err != nil {
		return nil, err
	}

	err = logsProcessor.ConsumeLogs(context.Background(), inputLogs)
	if err != nil {
		return nil, err
	}

	logsMarshaler := plog.JSONMarshaler{}
	json, err := logsMarshaler.MarshalLogs(transformedLogs)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func (p *transformProcessorExecutor) ExecuteTraceStatements(yamlConfig, input string) ([]byte, error) {
	config, err := p.parseConfig(yamlConfig)
	if err != nil {
		return nil, err
	}

	var transformedTraces ptrace.Traces
	tracesConsumer, _ := consumer.NewTraces(func(ctx context.Context, ld ptrace.Traces) error {
		transformedTraces = ld
		return nil
	})

	tracesProcessor, err := p.factory.CreateTracesProcessor(context.Background(), p.settings, config, tracesConsumer)
	if err != nil {
		return nil, err
	}

	tracesUnmarshaler := &ptrace.JSONUnmarshaler{}
	inputLogs, err := tracesUnmarshaler.UnmarshalTraces([]byte(input))
	if err != nil {
		return nil, err
	}

	err = tracesProcessor.ConsumeTraces(context.Background(), inputLogs)
	if err != nil {
		return nil, err
	}

	tracesMarshaler := ptrace.JSONMarshaler{}
	json, err := tracesMarshaler.MarshalTraces(transformedTraces)
	if err != nil {
		return nil, err
	}

	return json, nil
}

func (p *transformProcessorExecutor) ExecuteMetricStatements(yamlConfig, input string) ([]byte, error) {
	config, err := p.parseConfig(yamlConfig)
	if err != nil {
		return nil, err
	}

	var transformedMetrics pmetric.Metrics
	metricsConsumer, _ := consumer.NewMetrics(func(ctx context.Context, ld pmetric.Metrics) error {
		transformedMetrics = ld
		return nil
	})

	metricsProcessor, err := p.factory.CreateMetricsProcessor(context.Background(), p.settings, config, metricsConsumer)
	if err != nil {
		return nil, err
	}

	tracesUnmarshaler := &pmetric.JSONUnmarshaler{}
	inputLogs, err := tracesUnmarshaler.UnmarshalMetrics([]byte(input))
	if err != nil {
		return nil, err
	}

	err = metricsProcessor.ConsumeMetrics(context.Background(), inputLogs)
	if err != nil {
		return nil, err
	}

	metricsMarshaler := pmetric.JSONMarshaler{}
	json, err := metricsMarshaler.MarshalMetrics(transformedMetrics)
	if err != nil {
		return nil, err
	}

	return json, nil
}

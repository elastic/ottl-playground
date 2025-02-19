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
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type processorExecutor[T any] struct {
	factory           processor.Factory
	settings          processor.Settings
	telemetrySettings component.TelemetrySettings
	observedLogs      *ObservedLogs
}

func newProcessorExecutor[C any](factory processor.Factory) *processorExecutor[C] {
	observedLogger, observedLogs := NewLogObserver(zap.DebugLevel, zap.NewDevelopmentEncoderConfig())
	logger, _ := zap.NewDevelopmentConfig().Build(zap.WrapCore(func(z zapcore.Core) zapcore.Core {
		return observedLogger
	}))

	telemetrySettings := componenttest.NewNopTelemetrySettings()
	telemetrySettings.Logger = logger
	settings := processor.Settings{
		ID:                component.MustNewID("ottl_playground"),
		TelemetrySettings: telemetrySettings,
	}

	return &processorExecutor[C]{
		factory:           factory,
		telemetrySettings: telemetrySettings,
		settings:          settings,
		observedLogs:      observedLogs,
	}
}

func (p *processorExecutor[C]) parseConfig(yamlConfig string) (*C, error) {
	deserializedYaml, err := confmap.NewRetrievedFromYAML([]byte(yamlConfig))
	if err != nil {
		return nil, err
	}

	yamlConfigMap, err := deserializedYaml.AsConf()
	if err != nil {
		return nil, err
	}

	defaultConfig := p.factory.CreateDefaultConfig().(*C)
	err = yamlConfigMap.Unmarshal(&defaultConfig)
	if err != nil {
		return nil, err
	}

	return defaultConfig, nil
}

func (p *processorExecutor[C]) ExecuteLogStatements(yamlConfig, input string) ([]byte, error) {
	config, err := p.parseConfig(yamlConfig)
	if err != nil {
		return nil, err
	}

	transformedLogs := plog.NewLogs()
	logsConsumer, _ := consumer.NewLogs(func(_ context.Context, ld plog.Logs) error {
		transformedLogs = ld
		return nil
	})

	logsProcessor, err := p.factory.CreateLogs(context.Background(), p.settings, config, logsConsumer)
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

func (p *processorExecutor[C]) ExecuteTraceStatements(yamlConfig, input string) ([]byte, error) {
	config, err := p.parseConfig(yamlConfig)
	if err != nil {
		return nil, err
	}

	transformedTraces := ptrace.NewTraces()
	tracesConsumer, _ := consumer.NewTraces(func(_ context.Context, ld ptrace.Traces) error {
		transformedTraces = ld
		return nil
	})

	tracesProcessor, err := p.factory.CreateTraces(context.Background(), p.settings, config, tracesConsumer)
	if err != nil {
		return nil, err
	}

	tracesUnmarshaler := &ptrace.JSONUnmarshaler{}
	inputTraces, err := tracesUnmarshaler.UnmarshalTraces([]byte(input))
	if err != nil {
		return nil, err
	}

	err = tracesProcessor.ConsumeTraces(context.Background(), inputTraces)
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

func (p *processorExecutor[C]) ExecuteMetricStatements(yamlConfig, input string) ([]byte, error) {
	config, err := p.parseConfig(yamlConfig)
	if err != nil {
		return nil, err
	}

	transformedMetrics := pmetric.NewMetrics()
	metricsConsumer, _ := consumer.NewMetrics(func(_ context.Context, ld pmetric.Metrics) error {
		transformedMetrics = ld
		return nil
	})

	metricsProcessor, err := p.factory.CreateMetrics(context.Background(), p.settings, config, metricsConsumer)
	if err != nil {
		return nil, err
	}

	tracesUnmarshaler := &pmetric.JSONUnmarshaler{}
	inputMetrics, err := tracesUnmarshaler.UnmarshalMetrics([]byte(input))
	if err != nil {
		return nil, err
	}

	err = metricsProcessor.ConsumeMetrics(context.Background(), inputMetrics)
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

func (p *processorExecutor[C]) ObservedLogs() *ObservedLogs {
	return p.observedLogs
}

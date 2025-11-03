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
	"context"
	"errors"

	"go.opentelemetry.io/collector/consumer/xconsumer"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.opentelemetry.io/collector/processor/xprocessor"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Consumer[C any] interface {
	Observable
	// ComponentID returns the component.ID of the component.
	ComponentID() component.ID
	// ConsumeLogs processes the input logs and returns the transformed logs or an error.
	ConsumeLogs(config *C, input plog.Logs) (plog.Logs, error)
	// ConsumeMetrics processes the input metrics and returns the transformed metrics or an error.
	ConsumeMetrics(config *C, input pmetric.Metrics) (pmetric.Metrics, error)
	// ConsumeTraces processes the input traces and returns the transformed traces or an error.
	ConsumeTraces(config *C, input ptrace.Traces) (ptrace.Traces, error)
	// ConsumeProfiles processes the input profiles and returns the transformed profiles or an error.
	ConsumeProfiles(config *C, input pprofile.Profiles) (pprofile.Profiles, error)
	// CreateDefaultConfig returns the default configuration for the given component.
	CreateDefaultConfig() *C
	// TelemetrySettings returns the telemetry settings used by the component.
	TelemetrySettings() component.TelemetrySettings
}

type processorConsumer[C any] struct {
	id                component.ID
	factory           processor.Factory
	settings          processor.Settings
	telemetrySettings component.TelemetrySettings
	observedLogs      *ObservedLogs
}

func newProcessorConsumer[C any](
	factory processor.Factory,
) *processorConsumer[C] {
	observedLogger, observedLogs := NewLogObserver(zap.DebugLevel, zap.NewDevelopmentEncoderConfig())
	logger, _ := zap.NewDevelopmentConfig().Build(zap.WrapCore(func(z zapcore.Core) zapcore.Core {
		return observedLogger
	}))

	telemetrySettings := componenttest.NewNopTelemetrySettings()
	telemetrySettings.Logger = logger

	componentID := component.MustNewIDWithName(factory.Type().String(), "ottl_playground")
	buildInfo := component.NewDefaultBuildInfo()
	buildInfo.Description = "OTTL Playground"
	buildInfo.Version = CollectorContribProcessorsVersion
	buildInfo.Command = "wasm"

	settings := processor.Settings{
		ID:                componentID,
		TelemetrySettings: telemetrySettings,
		BuildInfo:         buildInfo,
	}

	return &processorConsumer[C]{
		id:                componentID,
		factory:           factory,
		telemetrySettings: telemetrySettings,
		settings:          settings,
		observedLogs:      observedLogs,
	}
}

func (p processorConsumer[C]) ConsumeLogs(config *C, input plog.Logs) (plog.Logs, error) {
	transformedLogs := plog.NewLogs()
	logsConsumer, _ := consumer.NewLogs(func(_ context.Context, ld plog.Logs) error {
		transformedLogs = ld
		return nil
	})

	logsProcessor, err := p.factory.CreateLogs(context.Background(), p.settings, config, logsConsumer)
	if err != nil {
		return plog.Logs{}, err
	}

	err = logsProcessor.Start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		return plog.Logs{}, err
	}

	err = logsProcessor.ConsumeLogs(context.Background(), input)
	if err != nil {
		return plog.Logs{}, err
	}

	return transformedLogs, nil
}

func (p processorConsumer[C]) ConsumeMetrics(config *C, input pmetric.Metrics) (pmetric.Metrics, error) {
	transformedMetrics := pmetric.NewMetrics()
	metricsConsumer, _ := consumer.NewMetrics(func(_ context.Context, ld pmetric.Metrics) error {
		transformedMetrics = ld
		return nil
	})

	metricsProcessor, err := p.factory.CreateMetrics(context.Background(), p.settings, config, metricsConsumer)
	if err != nil {
		return pmetric.Metrics{}, err
	}

	err = metricsProcessor.Start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		return pmetric.Metrics{}, err
	}

	err = metricsProcessor.ConsumeMetrics(context.Background(), input)
	if err != nil {
		return pmetric.Metrics{}, err
	}

	return transformedMetrics, nil
}

func (p processorConsumer[C]) ConsumeTraces(config *C, input ptrace.Traces) (ptrace.Traces, error) {
	transformedTraces := ptrace.NewTraces()
	tracesConsumer, _ := consumer.NewTraces(func(_ context.Context, ld ptrace.Traces) error {
		transformedTraces = ld
		return nil
	})

	tracesProcessor, err := p.factory.CreateTraces(context.Background(), p.settings, config, tracesConsumer)
	if err != nil {
		return ptrace.Traces{}, err
	}

	err = tracesProcessor.Start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		return ptrace.Traces{}, err
	}

	err = tracesProcessor.ConsumeTraces(context.Background(), input)
	if err != nil {
		return ptrace.Traces{}, err
	}

	return transformedTraces, nil
}

func (p processorConsumer[C]) ConsumeProfiles(config *C, input pprofile.Profiles) (pprofile.Profiles, error) {
	factory, ok := p.factory.(xprocessor.Factory)
	if !ok {
		return pprofile.Profiles{}, errors.New("profiles are not supported by this OTel Collector version or component")
	}

	transformedProfiles := pprofile.NewProfiles()
	profilesConsumer, _ := xconsumer.NewProfiles(func(_ context.Context, ld pprofile.Profiles) error {
		transformedProfiles = ld
		return nil
	})

	profilesProcessor, err := factory.CreateProfiles(context.Background(), p.settings, config, profilesConsumer)
	if err != nil {
		return pprofile.Profiles{}, err
	}

	err = profilesProcessor.Start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		return pprofile.Profiles{}, err
	}

	err = profilesProcessor.ConsumeProfiles(context.Background(), input)
	if err != nil {
		return pprofile.Profiles{}, err
	}

	return transformedProfiles, nil
}

func (p processorConsumer[C]) ObservedLogs() *ObservedLogs {
	return p.observedLogs
}

func (p processorConsumer[C]) TelemetrySettings() component.TelemetrySettings {
	return p.telemetrySettings
}

func (p processorConsumer[C]) CreateDefaultConfig() *C {
	return p.factory.CreateDefaultConfig().(*C)
}

func (p processorConsumer[C]) ComponentID() component.ID {
	return p.id
}

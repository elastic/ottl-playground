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
	"fmt"
	"log"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type ResultView string

// Keep this list in sync with the views defined in the frontend.
const (
	ResultViewVisualDiff    ResultView = "visual_delta"
	ResultViewAnnotatedDiff ResultView = "annotated_delta"
	ResultViewJSON          ResultView = "json"
	ResultViewLogs          ResultView = "logs"
)

type ResultViewConfig struct {
	Enabled bool `json:"enabled"`
}

func newDefaultResultViewConfig() map[ResultView]*ResultViewConfig {
	return map[ResultView]*ResultViewConfig{
		ResultViewVisualDiff:    {Enabled: true},
		ResultViewAnnotatedDiff: {Enabled: true},
		ResultViewJSON:          {Enabled: true},
		ResultViewLogs:          {Enabled: true},
	}
}

type PayloadExample struct {
	Name   string `json:"name"`
	Signal string `json:"signal"` // "logs", "traces", or "metrics"
	Value  string `json:"value"`
}

type ConfigExample struct {
	Name    string `json:"name"`
	Signal  string `json:"signal"` // "logs", "traces", or "metrics"
	Config  string `json:"config"`
	Payload string `json:"payload"`
}

type Examples struct {
	Configs  []ConfigExample  `json:"configs"`
	Payloads []PayloadExample `json:"payloads"`
}

type ComponentType string

const (
	ComponentTypeProcessor ComponentType = "processor"
)

// Metadata contains information about the playground executor, such as its ID, name,
// path, version, documentation URL, and configuration for result views.
type Metadata struct {
	Type             ComponentType                    `json:"type"`
	ID               string                           `json:"id"`
	Name             string                           `json:"name"`
	Path             string                           `json:"path"`
	Version          string                           `json:"version"`
	DocsURL          string                           `json:"docsURL"`
	ResultViewConfig map[ResultView]*ResultViewConfig `json:"resultViewConfig"`
	Examples         Examples                         `json:"examples"`
	Debuggable       bool                             `json:"debuggable"`
}

// MedataOption is a function that modifies the Metadata configuration.
type medataOption func(*Metadata) error

// enableResultViews enables the specified result views in the metadata.
func enableResultViews(views ...ResultView) medataOption {
	return func(metadata *Metadata) error {
		for _, cfg := range metadata.ResultViewConfig {
			cfg.Enabled = false
		}
		for _, view := range views {
			cfg, ok := metadata.ResultViewConfig[view]
			if !ok {
				return fmt.Errorf("unsupported result view %q", view)
			}
			cfg.Enabled = true
		}
		return nil
	}
}

// disableResultViews disables the specified result views in the metadata.
func disableResultViews(views ...ResultView) medataOption {
	return func(metadata *Metadata) error {
		for _, v := range views {
			cfg, ok := metadata.ResultViewConfig[v]
			if !ok {
				return fmt.Errorf("unsupported result view %q", v)
			}
			cfg.Enabled = false
		}
		return nil
	}
}

// withConfigExamples adds examples to the executor metadata.
func withConfigExamples(examples ...ConfigExample) medataOption {
	return func(metadata *Metadata) error {
		metadata.Examples.Configs = append(metadata.Examples.Configs, examples...)
		return nil
	}
}

// withPayloadExamples adds examples to the executor metadata.
func withPayloadExamples(examples ...PayloadExample) medataOption {
	return func(metadata *Metadata) error {
		metadata.Examples.Payloads = append(metadata.Examples.Payloads, examples...)
		return nil
	}
}

func newMetadata(ct ComponentType, id, name, path, docsURL string, options ...medataOption) *Metadata {
	meta := Metadata{
		Type:             ct,
		ID:               id,
		Name:             name,
		Path:             path,
		DocsURL:          docsURL,
		Version:          CollectorContribProcessorsVersion,
		ResultViewConfig: newDefaultResultViewConfig(),
		Examples:         Examples{},
	}

	for _, opt := range options {
		if err := opt(&meta); err != nil {
			log.Printf("error applying metadata option: %v", err)
		}
	}

	return &meta
}

// Executor evaluates OTTL statements using specific configurations and inputs.
type Executor interface {
	Observable
	// ExecuteLogs evaluates log statements using the given configuration and JSON payload.
	// The returned value must be a valid plog.Logs JSON representing the input transformation.
	ExecuteLogs(config, input string) (*Result, error)
	// ExecuteTraces is like ExecuteLogs, but for traces.
	ExecuteTraces(config, input string) (*Result, error)
	// ExecuteMetrics is like ExecuteLogs, but for metrics.
	ExecuteMetrics(config, input string) (*Result, error)
	// ExecuteProfiles is like ExecuteLogs, but for profiles.
	ExecuteProfiles(config, input string) (*Result, error)
	// Metadata returns information about the executor
	Metadata() *Metadata
}

// DebuggableExecutor is an Executor that supports debugging.
type DebuggableExecutor interface {
	Debugger() (Debugger, error)
}

// Debugger provides debugging capabilities for OTTL statements.
type Debugger interface {
	Observable
	// DebugLogs evaluates log statements using the given configuration and JSON
	// payload with debugging enabled.
	DebugLogs(config, input string) (*Result, error)
	// DebugTraces is like DebugLogs, but for traces.
	DebugTraces(config, input string) (*Result, error)
	// DebugMetrics is like DebugLogs, but for metrics.
	DebugMetrics(config, input string) (*Result, error)
	// DebugProfiles is like DebugLogs, but for profiles.
	DebugProfiles(config, input string) (*Result, error)
}

// Observable represents an entity that can provide observed logs.
type Observable interface {
	// ObservedLogs returns the statements execution's logs
	ObservedLogs() *ObservedLogs
}

type defaultExecutor[C any] struct {
	consumer         Consumer[C]
	metadata         *Metadata
	logMarshaler     plog.Marshaler
	metricMarshaler  pmetric.Marshaler
	traceMarshaler   ptrace.Marshaler
	profileMarshaler pprofile.Marshaler
	debugger         Debugger
}

type executorOption[C any] func(*defaultExecutor[C])

func withDebugger[C any](debugger Debugger) executorOption[C] {
	return func(e *defaultExecutor[C]) {
		e.metadata.Debuggable = true
		e.debugger = debugger
	}
}

func NewJSONExecutor[C any](
	consumer Consumer[C],
	metadata *Metadata,
	options ...executorOption[C],
) Executor {
	exec := &defaultExecutor[C]{
		consumer:         consumer,
		metadata:         metadata,
		logMarshaler:     &plog.JSONMarshaler{},
		metricMarshaler:  &pmetric.JSONMarshaler{},
		traceMarshaler:   &ptrace.JSONMarshaler{},
		profileMarshaler: &pprofile.JSONMarshaler{},
	}
	for _, opt := range options {
		opt(exec)
	}
	return exec
}

func (e *defaultExecutor[C]) ExecuteLogs(config, input string) (*Result, error) {
	logsUnmarshaler := &plog.JSONUnmarshaler{}
	inputLogs, err := logsUnmarshaler.UnmarshalLogs([]byte(input))
	if err != nil {
		return nil, err
	}

	cfgs, err := parseConfig[C](e.consumer.ComponentID(), config, e.consumer.CreateDefaultConfig)
	if err != nil {
		return nil, err
	}

	return newExecutionResult(e, e.logMarshaler.MarshalLogs, func() (plog.Logs, error) {
		transformedLogs := inputLogs
		for _, cfg := range cfgs {
			if len(cfgs) > 1 {
				e.consumer.TelemetrySettings().Logger.Sugar().Debugf("[playground] Running configuration: %s", cfg.Key)
			}
			transformedLogs, err = e.consumer.ConsumeLogs(cfg.Value, transformedLogs)
			if err != nil {
				return plog.Logs{}, err
			}
		}
		return transformedLogs, nil
	})
}

func (e *defaultExecutor[C]) ExecuteTraces(config, input string) (*Result, error) {
	tracesUnmarshaler := &ptrace.JSONUnmarshaler{}
	inputTraces, err := tracesUnmarshaler.UnmarshalTraces([]byte(input))
	if err != nil {
		return nil, err
	}

	cfgs, err := parseConfig[C](e.consumer.ComponentID(), config, e.consumer.CreateDefaultConfig)
	if err != nil {
		return nil, err
	}

	return newExecutionResult(e, e.traceMarshaler.MarshalTraces, func() (ptrace.Traces, error) {
		transformedTraces := inputTraces
		for _, cfg := range cfgs {
			if len(cfgs) > 1 {
				e.consumer.TelemetrySettings().Logger.Sugar().Debugf("[playground] Running configuration: %s", cfg.Key)
			}
			transformedTraces, err = e.consumer.ConsumeTraces(cfg.Value, transformedTraces)
			if err != nil {
				return ptrace.Traces{}, err
			}
		}
		return transformedTraces, nil
	})
}

func (e *defaultExecutor[C]) ExecuteMetrics(config, input string) (*Result, error) {
	metricsUnmarshaler := &pmetric.JSONUnmarshaler{}
	inputMetrics, err := metricsUnmarshaler.UnmarshalMetrics([]byte(input))
	if err != nil {
		return nil, err
	}

	cfgs, err := parseConfig[C](e.consumer.ComponentID(), config, e.consumer.CreateDefaultConfig)
	if err != nil {
		return nil, err
	}

	return newExecutionResult(e, e.metricMarshaler.MarshalMetrics, func() (pmetric.Metrics, error) {
		transformedMetrics := inputMetrics
		for _, cfg := range cfgs {
			if len(cfgs) > 1 {
				e.consumer.TelemetrySettings().Logger.Sugar().Debugf("[playground] Running configuration: %s", cfg.Key)
			}
			transformedMetrics, err = e.consumer.ConsumeMetrics(cfg.Value, inputMetrics)
			if err != nil {
				return pmetric.Metrics{}, err
			}
		}
		return transformedMetrics, nil
	})
}

func (e *defaultExecutor[C]) ExecuteProfiles(config, input string) (*Result, error) {
	profilesUnmarshaler := &pprofile.JSONUnmarshaler{}
	inputProfiles, err := profilesUnmarshaler.UnmarshalProfiles([]byte(input))
	if err != nil {
		return nil, err
	}

	cfgs, err := parseConfig[C](e.consumer.ComponentID(), config, e.consumer.CreateDefaultConfig)
	if err != nil {
		return nil, err
	}

	return newExecutionResult(e, e.profileMarshaler.MarshalProfiles, func() (pprofile.Profiles, error) {
		transformedProfiles := inputProfiles
		for _, cfg := range cfgs {
			if len(cfgs) > 1 {
				e.consumer.TelemetrySettings().Logger.Sugar().Debugf("[playground] Running configuration: %s", cfg.Key)
			}
			transformedProfiles, err = e.consumer.ConsumeProfiles(cfg.Value, transformedProfiles)
			if err != nil {
				return pprofile.Profiles{}, err
			}
		}
		return transformedProfiles, nil
	})
}

func (e *defaultExecutor[C]) ObservedLogs() *ObservedLogs {
	return e.consumer.ObservedLogs()
}

func (e *defaultExecutor[C]) Metadata() *Metadata {
	return e.metadata
}

func (e *defaultExecutor[C]) Debugger() (Debugger, error) {
	if !e.metadata.Debuggable || e.debugger == nil {
		return nil, fmt.Errorf("executor %q does not support debugging", e.metadata.Name)
	}
	return e.debugger, nil
}

// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"go.uber.org/zap/zaptest/observer"
)

type Metadata struct {
	ID      string
	Name    string
	Path    string
	Version string
	DocsURL string
}

func newMetadata(id, name, path, docsURL string) Metadata {
	return Metadata{
		ID:      id,
		Name:    name,
		Path:    path,
		DocsURL: docsURL,
		Version: CollectorContribProcessorsVersion,
	}
}

// Executor evaluates OTTL statements using specific configurations and inputs.
type Executor interface {
	// ExecuteLogStatements evaluates log statements using the given configuration and JSON payload.
	// The returned value must be a valid plog.Logs JSON representing the input transformation.
	ExecuteLogStatements(config, input string) ([]byte, error)
	// ExecuteTraceStatements is like ExecuteLogStatements, but for traces.
	ExecuteTraceStatements(config, input string) ([]byte, error)
	// ExecuteMetricStatements is like ExecuteLogStatements, but for metrics.
	ExecuteMetricStatements(config, input string) ([]byte, error)
	// ObservedLogs returns the statements execution's logs
	ObservedLogs() *observer.ObservedLogs
	// Metadata returns information about the executor
	Metadata() Metadata
}

func Executors() []Executor {
	return []Executor{
		NewTransformProcessorExecutor(),
		NewFilterProcessorExecutor(),
	}
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"go.uber.org/zap/zaptest/observer"
)

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
}

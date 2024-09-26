// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import "go.uber.org/zap/zaptest/observer"

type Executor interface {
	ExecuteLogStatements(yamlConfig, input string) ([]byte, error)

	ExecuteTraceStatements(yamlConfig, input string) ([]byte, error)

	ExecuteMetricStatements(yamlConfig, input string) ([]byte, error)

	ObservedLogs() *observer.ObservedLogs
}

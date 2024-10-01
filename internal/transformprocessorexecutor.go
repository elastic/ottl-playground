// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
)

type transformProcessorExecutor struct {
	*processorExecutor[transformprocessor.Config]
}

// NewTransformProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [transformprocessor].
func NewTransformProcessorExecutor() Executor {
	executor := newProcessorExecutor[transformprocessor.Config](transformprocessor.NewFactory())
	return &transformProcessorExecutor{executor}
}

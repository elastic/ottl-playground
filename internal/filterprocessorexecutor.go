// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
)

type filterProcessorExecutor struct {
	*processorExecutor[filterprocessor.Config]
}

// NewFilterProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [filterprocessor].
func NewFilterProcessorExecutor() Executor {
	executor := newProcessorExecutor[filterprocessor.Config](filterprocessor.NewFactory())
	return &filterProcessorExecutor{executor}
}

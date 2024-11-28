// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
)

type transformProcessorExecutor struct {
	*processorExecutor[transformprocessor.Config]
}

func (t transformProcessorExecutor) Metadata() Metadata {
	return newMetadata(
		"transform_processor",
		"Transform processor",
		"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/transformprocessor",
		"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/transformprocessor",
	)
}

// NewTransformProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [transformprocessor].
func NewTransformProcessorExecutor() Executor {
	executor := newProcessorExecutor[transformprocessor.Config](transformprocessor.NewFactory())
	return &transformProcessorExecutor{executor}
}

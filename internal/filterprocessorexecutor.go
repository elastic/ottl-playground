// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
)

type filterProcessorExecutor struct {
	*processorExecutor[filterprocessor.Config]
}

func (f filterProcessorExecutor) Metadata() Metadata {
	return newMetadata(
		"filter_processor",
		"Filter processor",
		"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor",
		"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/filterprocessor",
	)
}

// NewFilterProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [filterprocessor].
func NewFilterProcessorExecutor() Executor {
	executor := newProcessorExecutor[filterprocessor.Config](filterprocessor.NewFactory())
	return &filterProcessorExecutor{executor}
}

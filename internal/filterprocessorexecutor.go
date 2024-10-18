// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
)

type filterProcessorExecutor struct {
	*processorExecutor[filterprocessor.Config]
}

func (f filterProcessorExecutor) Metadata() Metadata {
	return Metadata{
		Id:      "filter_processor",
		Name:    "Filter processor",
		DocsURL: "https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/filterprocessor",
		Version: "v0.110.0",
	}
}

// NewFilterProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [filterprocessor].
func NewFilterProcessorExecutor() Executor {
	executor := newProcessorExecutor[filterprocessor.Config](filterprocessor.NewFactory())
	return &filterProcessorExecutor{executor}
}

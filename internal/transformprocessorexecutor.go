// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
)

type transformProcessorExecutor struct {
	*processorExecutor[transformprocessor.Config]
}

func (t transformProcessorExecutor) Metadata() Metadata {
	return Metadata{
		Id:      "transform_processor",
		Name:    "Transform processor",
		DocsURL: "https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/transformprocessor",
		Version: "v0.110.0",
	}
}

// NewTransformProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [transformprocessor].
func NewTransformProcessorExecutor() Executor {
	executor := newProcessorExecutor[transformprocessor.Config](transformprocessor.NewFactory())
	return &transformProcessorExecutor{executor}
}

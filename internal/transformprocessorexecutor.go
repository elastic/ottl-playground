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

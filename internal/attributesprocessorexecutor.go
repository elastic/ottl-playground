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
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
)

type attributesProcessorExecutor struct {
	*processorExecutor[attributesprocessor.Config]
}

func (t attributesProcessorExecutor) Metadata() Metadata {
	return newMetadata(
		"attributes_processor",
		"Attributes processor",
		"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor",
		"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/attributesprocessor",
	)
}

// NewAttributesProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [attributesprocessor].
func NewAttributesProcessorExecutor() Executor {
	executor := newProcessorExecutor[attributesprocessor.Config](attributesprocessor.NewFactory())
	return &attributesProcessorExecutor{executor}
}

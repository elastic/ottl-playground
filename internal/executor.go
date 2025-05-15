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

type Metadata struct {
	ID      string
	Name    string
	Path    string
	Version string
	DocsURL string
}

func newMetadata(id, name, path, docsURL string) Metadata {
	return Metadata{
		ID:      id,
		Name:    name,
		Path:    path,
		DocsURL: docsURL,
		Version: CollectorContribProcessorsVersion,
	}
}

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
	ObservedLogs() *ObservedLogs
	// Metadata returns information about the executor
	Metadata() Metadata
}

func Executors() []Executor {
	return []Executor{
		NewTransformProcessorExecutor(),
		NewFilterProcessorExecutor(),
		NewGroupByAttrsProcessorExecutor(),
	}
}

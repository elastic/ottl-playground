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
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/ottlplayground/internal"
)

var (
	statementsExecutors       []internal.Executor
	statementsExecutorsLookup = map[string]internal.Executor{}
)

func init() {
	executors := internal.Executors()
	slices.SortFunc(executors, func(a, b internal.Executor) int {
		return strings.Compare(a.Metadata().Name, b.Metadata().Name)
	})

	for _, executor := range executors {
		registerStatementsExecutor(executor)
	}
}

func registerStatementsExecutor(executor internal.Executor) {
	statementsExecutors = append(statementsExecutors, executor)
	statementsExecutorsLookup[executor.Metadata().ID] = executor
}

func Execute(config, signal, ottlDataPayload, executorName string, debug bool) map[string]any {
	executor, ok := statementsExecutorsLookup[executorName]
	if !ok {
		return internal.NewErrorResult(fmt.Sprintf("unsupported executor %s", executorName), "").AsRaw()
	}

	var result *internal.Result
	var err error
	if debug {
		result, err = debugConfig(config, signal, ottlDataPayload, executorName, executor)
	} else {
		result, err = executeConfig(config, signal, ottlDataPayload, executorName, executor)
	}

	if err != nil {
		result = internal.NewErrorResult(fmt.Sprintf("unable to run %s configuration. Error: %v", signal, err), executor.ObservedLogs().TakeAllString())
	}

	return result.AsRaw()
}

func executeConfig(config, signal, ottlDataPayload, _ string, executor internal.Executor) (*internal.Result, error) {
	switch signal {
	case "logs":
		return executor.ExecuteLogs(config, ottlDataPayload)
	case "traces":
		return executor.ExecuteTraces(config, ottlDataPayload)
	case "metrics":
		return executor.ExecuteMetrics(config, ottlDataPayload)
	default:
		return internal.NewErrorResult(fmt.Sprintf("unsupported OTLP signal type %s", signal), ""), nil
	}
}

func debugConfig(config, signal, ottlDataPayload, executorName string, executor internal.Executor) (*internal.Result, error) {
	debuggableExecutor, ok := executor.(internal.DebuggableExecutor)
	if !ok {
		return internal.NewErrorResult(fmt.Sprintf("executor %q does not support debugging", executorName), ""), nil
	}

	debugger, err := debuggableExecutor.Debugger()
	if err != nil {
		return nil, err
	}

	switch signal {
	case "logs":
		return debugger.DebugLogs(config, ottlDataPayload)
	case "traces":
		return debugger.DebugTraces(config, ottlDataPayload)
	case "metrics":
		return debugger.DebugMetrics(config, ottlDataPayload)
	default:
		return internal.NewErrorResult(fmt.Sprintf("unsupported OTLP signal type %s", signal), ""), nil
	}
}

func Executors() []any {
	var res []any
	for _, executor := range statementsExecutors {
		var metadataValue map[string]any
		if metadataBytes, err := json.Marshal(executor.Metadata()); err == nil {
			_ = json.Unmarshal(metadataBytes, &metadataValue)
		}

		var examples map[string]any
		if examplesJson, err := json.Marshal(executor.Metadata().Examples); err == nil {
			_ = json.Unmarshal(examplesJson, &examples)
		}

		res = append(res, metadataValue)
	}
	return res
}

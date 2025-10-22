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
	"fmt"
	"strings"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/ottlplayground/internal"
)

var (
	statementsExecutors       []internal.Executor
	statementsExecutorsLookup = map[string]internal.Executor{}
)

func init() {
	for _, executor := range internal.Executors() {
		registerStatementsExecutor(executor)
	}
}

func registerStatementsExecutor(executor internal.Executor) {
	statementsExecutors = append(statementsExecutors, executor)
	statementsExecutorsLookup[executor.Metadata().ID] = executor
}

func newResult(json string, err string, logs string, executionTime int64) map[string]any {
	v := map[string]any{
		"value":         json,
		"logs":          logs,
		"executionTime": executionTime,
	}
	if err != "" {
		v["error"] = err
	}
	return v
}

func NewErrorResult(err string, logs string) map[string]any {
	return newResult("", err, logs, 0)
}

func takeObservedLogs(executor internal.Executor) string {
	all := executor.ObservedLogs().TakeAll()
	var s strings.Builder
	for _, entry := range all {
		s.WriteString(entry.ConsoleEncodedEntry())
	}
	return s.String()
}

func ExecuteStatements(config, ottlDataType, ottlDataPayload, executorName string) map[string]any {
	executor, ok := statementsExecutorsLookup[executorName]
	if !ok {
		return NewErrorResult(fmt.Sprintf("unsupported evaluator %s", executorName), "")
	}

	start := time.Now()
	var output []byte
	var err error
	switch ottlDataType {
	case "logs":
		output, err = executor.ExecuteLogStatements(config, ottlDataPayload)
	case "traces":
		output, err = executor.ExecuteTraceStatements(config, ottlDataPayload)
	case "metrics":
		output, err = executor.ExecuteMetricStatements(config, ottlDataPayload)
	case "profiles":
		output, err = executor.ExecuteProfileStatements(config, ottlDataPayload)
	default:
		return NewErrorResult(fmt.Sprintf("unsupported OTLP data type %s", ottlDataType), "")
	}

	if err != nil {
		return NewErrorResult(fmt.Sprintf("unable to run %s statements. Error: %v", ottlDataType, err), takeObservedLogs(executor))
	}

	executionTime := time.Since(start).Milliseconds()
	return newResult(string(output), "", takeObservedLogs(executor), executionTime)
}

func StatementsExecutors() []any {
	var res []any
	for _, executor := range statementsExecutors {
		meta := executor.Metadata()
		res = append(res, map[string]any{
			"id":      meta.ID,
			"name":    meta.Name,
			"path":    meta.Path,
			"docsURL": meta.DocsURL,
			"version": meta.Version,
		})
	}
	return res
}

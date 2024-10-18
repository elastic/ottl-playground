// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/ottlplayground/internal"
)

type statementExecutor struct {
	id       string
	name     string
	helpLink string
	internal.Executor
}

var (
	defaultLogEncoder         = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	statementsExecutorsLookup = map[string]*statementExecutor{}
	statementsExecutors       = []statementExecutor{
		{
			"transform_processor",
			"Transform processor",
			"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/transformprocessor",
			internal.NewTransformProcessorExecutor(),
		},
		{
			"filter_processor",
			"Filter processor",
			"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/filterprocessor",
			internal.NewFilterProcessorExecutor(),
		},
	}
)

func init() {
	for _, executor := range statementsExecutors {
		statementsExecutorsLookup[executor.id] = &executor
	}
}

func newResult(json string, err string, logs string) map[string]any {
	v := map[string]any{
		"value": json,
		"logs":  logs,
	}
	if err != "" {
		v["error"] = err
	}
	return v
}

func NewErrorResult(err string, logs string) map[string]any {
	return newResult("", err, logs)
}

func takeObservedLogs(executor internal.Executor) string {
	all := executor.ObservedLogs().TakeAll()
	var s strings.Builder
	for _, entry := range all {
		v, err := defaultLogEncoder.EncodeEntry(entry.Entry, entry.Context)
		if err == nil {
			s.Write(v.Bytes())
		}
	}
	return s.String()
}

func ExecuteStatements(config, ottlDataType, ottlDataPayload, executorName string) map[string]any {
	executor, ok := statementsExecutorsLookup[executorName]
	if !ok {
		return NewErrorResult(fmt.Sprintf("unsupported evaluator %s", executorName), "")
	}

	var output []byte
	var err error
	switch ottlDataType {
	case "logs":
		output, err = executor.ExecuteLogStatements(config, ottlDataPayload)
	case "traces":
		output, err = executor.ExecuteTraceStatements(config, ottlDataPayload)
	case "metrics":
		output, err = executor.ExecuteMetricStatements(config, ottlDataPayload)
	default:
		return NewErrorResult(fmt.Sprintf("unsupported OTLP data type %s", ottlDataType), "")
	}

	if err != nil {
		return NewErrorResult(fmt.Sprintf("unable to run %s statements. Error: %v", ottlDataType, err), takeObservedLogs(executor))
	}

	return newResult(string(output), "", takeObservedLogs(executor))
}

func StatementsExecutors() []any {
	var res []any
	for _, executor := range statementsExecutors {
		res = append(res, map[string]any{
			"id":       executor.id,
			"name":     executor.name,
			"helpLink": executor.helpLink,
		})
	}
	return res
}

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

var (
	defaultLogEncoder   = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	statementsExecutors = map[string]internal.Executor{
		"transform_processor": internal.NewTransformProcessorExecutor(),
		"filter_processor":    internal.NewFilterProcessorExecutor(),
	}
)

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
	executor, ok := statementsExecutors[executorName]
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

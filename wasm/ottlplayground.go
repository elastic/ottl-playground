// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build js && wasm

package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/ottlplayground/internal"
)

var (
	defaultLogEncoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
)

func newResult(json string, error string, logs string) map[string]any {
	v := map[string]any{
		"value": json,
		"logs":  logs,
	}
	if error != "" {
		v["error"] = error
	}
	return v
}

func newErrorResult(err string, logs string) map[string]any {
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

func executeStatementsWrapper() js.Func {
	executors := map[string]internal.Executor{
		"transform_processor": internal.NewTransformProcessorExecutor(),
		"filter_processor":    internal.NewFilterProcessorExecutor(),
	}

	executeStatementsFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from", r)
			}
		}()

		if len(args) != 4 {
			return newErrorResult("invalid number of arguments", "")
		}

		config := args[0].String()
		dataType := args[1].String()
		payload := args[2].String()
		executorName := args[3].String()

		executor, ok := executors[executorName]
		if !ok {
			return newErrorResult(fmt.Sprintf("unsupported evaluator %s", executorName), "")
		}

		var output []byte
		var err error
		switch dataType {
		case "logs":
			output, err = executor.ExecuteLogStatements(config, payload)
		case "traces":
			output, err = executor.ExecuteTraceStatements(config, payload)
		case "metrics":
			output, err = executor.ExecuteMetricStatements(config, payload)
		default:
			return newErrorResult(fmt.Sprintf("unsupported OTLP data type %s", dataType), "")
		}

		if err != nil {
			return newErrorResult(fmt.Sprintf("unable to run %s statements. Error %w", dataType, err), takeObservedLogs(executor))
		}

		return newResult(string(output), "", takeObservedLogs(executor))
	})

	return executeStatementsFunc
}

func main() {
	js.Global().Set("executeStatements", executeStatementsWrapper())
	<-make(chan struct{})
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/ottlplayground/internal"
)

func newError(message string, args ...interface{}) map[string]any {
	return map[string]any{
		"error": fmt.Sprintf(message, args...),
	}
}

func executeStatementsWrapper() js.Func {
	executor := internal.NewTransformProcessorExecutor()
	executeStatementsFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return newError("Invalid number of arguments passed")
		}

		config := args[0].String()
		dataType := args[1].String()
		payload := args[2].String()

		var output []byte
		var err error
		switch dataType {
		case "logs":
			output, err = executor.ExecuteLogStatements(config, payload)
		case "trace":
			output, err = executor.ExecuteTraceStatements(config, payload)
		case "metrics":
			output, err = executor.ExecuteMetricStatements(config, payload)
		default:
			return newError("unsupported OTLP data type %s", dataType)
		}

		if err != nil {
			return newError("Unable to run %s statements. Error %w", dataType, err)
		}

		return string(output)
	})

	return executeStatementsFunc
}

func main() {
	js.Global().Set("executeStatements", executeStatementsWrapper())
	<-make(chan struct{})
}

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build js && wasm

package main

import (
	"fmt"
	"syscall/js"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/ottlplayground/wasm/internal"
)

func executeStatementsWrapper() js.Func {
	executeStatementsFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from", r)
			}
		}()

		if len(args) != 4 {
			return internal.NewErrorResult("invalid number of arguments", "")
		}

		config := args[0].String()
		ottlDataType := args[1].String()
		ottlDataPayload := args[2].String()
		executorName := args[3].String()
		return internal.ExecuteStatements(config, ottlDataType, ottlDataPayload, executorName)
	})

	return executeStatementsFunc
}

func main() {
	js.Global().Set("executeStatements", executeStatementsWrapper())
	<-make(chan struct{})
}

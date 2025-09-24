//go:build js && wasm

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

package main

import (
	"fmt"
	"syscall/js"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/ottlplayground/wasm/internal"
)

func handlePanic() {
	defer func() {
		recover()
	}()
	if r := recover(); r != nil {
		js.Global().Call("wasmPanicHandler", fmt.Sprintf("An error occurred in the WASM module: %v", r))
		js.Global().Get("console").Call("error", "stack trace:", js.Global().Get("Error").New().Get("stack").String())
	}
}

func executeWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer handlePanic()
		if len(args) != 5 {
			return map[string]any{"error": "invalid number of arguments"}
		}

		config := args[0].String()
		ottlDataType := args[1].String()
		ottlDataPayload := args[2].String()
		executorName := args[3].String()
		debug := args[4].Bool()
		return js.ValueOf(internal.Execute(config, ottlDataType, ottlDataPayload, executorName, debug))
	})
}

func getExecutorsWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, _ []js.Value) any {
		defer handlePanic()
		return js.ValueOf(internal.Executors())
	})
}

func main() {
	js.Global().Set("execute", executeWrapper())
	js.Global().Set("getExecutors", getExecutorsWrapper())
	<-make(chan struct{})
}

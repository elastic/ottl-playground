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

	"github.com/elastic/ottl-playground/wasm/internal"
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

func getOTTLFunctionsWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, _ []js.Value) any {
		defer handlePanic()
		return js.ValueOf(internal.GetOTTLFunctions())
	})
}

func getContextPathsWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer handlePanic()
		if len(args) != 1 {
			return nil
		}
		context := args[0].String()
		return js.ValueOf(internal.GetContextPaths(context))
	})
}

func getContextEnumsWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer handlePanic()
		if len(args) != 1 {
			return nil
		}
		context := args[0].String()
		return js.ValueOf(internal.GetContextEnums(context))
	})
}

func getOTTLMetadataWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, _ []js.Value) any {
		defer handlePanic()
		return js.ValueOf(internal.GetOTTLMetadata())
	})
}

func validateStatementsWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer handlePanic()
		if len(args) != 4 {
			return js.ValueOf([]any{
				map[string]any{
					"message":   "invalid number of arguments (expected: config, dataType, payload, executor)",
					"severity":  "error",
					"line":      1,
					"column":    1,
					"endLine":   1,
					"endColumn": 1,
				},
			})
		}
		config := args[0].String()
		dataType := args[1].String()
		payload := args[2].String()
		executorName := args[3].String()
		return js.ValueOf(internal.ValidateStatements(config, dataType, payload, executorName))
	})
}

func getCompletionContextWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer handlePanic()
		if len(args) != 1 {
			return nil
		}
		statement := args[0].String()
		return js.ValueOf(internal.GetCompletionContext(statement))
	})
}

func main() {
	js.Global().Set("execute", executeWrapper())
	js.Global().Set("getExecutors", getExecutorsWrapper())

	// Autocomplete metadata exports
	js.Global().Set("getOTTLFunctions", getOTTLFunctionsWrapper())
	js.Global().Set("getContextPaths", getContextPathsWrapper())
	js.Global().Set("getContextEnums", getContextEnumsWrapper())
	js.Global().Set("getOTTLMetadata", getOTTLMetadataWrapper())

	// Validation export
	js.Global().Set("validateStatements", validateStatementsWrapper())

	// Completion context export (lexer-based)
	js.Global().Set("getCompletionContext", getCompletionContextWrapper())

	<-make(chan struct{})
}

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
	if r := recover(); r != nil {
		fmt.Println("recovered from", r)
	}
}

func executeStatementsWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer handlePanic()
		if len(args) != 4 {
			return internal.NewErrorResult("invalid number of arguments", "")
		}

		config := args[0].String()
		ottlDataType := args[1].String()
		ottlDataPayload := args[2].String()
		executorName := args[3].String()
		return js.ValueOf(internal.ExecuteStatements(config, ottlDataType, ottlDataPayload, executorName))
	})
}

func getEvaluatorsWrapper() js.Func {
	return js.FuncOf(func(_ js.Value, _ []js.Value) any {
		defer handlePanic()
		return js.ValueOf(internal.StatementsExecutors())
	})
}

func main() {
	js.Global().Set("executeStatements", executeStatementsWrapper())
	js.Global().Set("statementsExecutors", getEvaluatorsWrapper())
	<-make(chan struct{})
}

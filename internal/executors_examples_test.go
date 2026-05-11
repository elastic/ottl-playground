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
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testingSignalPayloads = map[string]string{
		"logs":    mustReadTestData("logs.json"),
		"traces":  mustReadTestData("traces.json"),
		"metrics": mustReadTestData("metrics.json"),
	}
)

func mustReadTestData(file string) string {
	content, err := os.ReadFile(filepath.Join("..", "testdata", file))
	if err != nil {
		panic(fmt.Sprintf("failed to read test data file %s: %v", file, err))
	}
	return string(content)
}

func Test_Executors_ConfigExamples(t *testing.T) {
	for _, executor := range Executors() {
		examples := executor.Metadata().Examples.Configs
		if len(examples) == 0 {
			t.Skip("No config examples available for this executor")
		}

		t.Run(fmt.Sprintf("%s", executor.Metadata().ID), func(t *testing.T) {
			for _, example := range examples {
				t.Run(example.Name, func(t *testing.T) {
					var result *Result
					var err error
					examplePayload := example.Payload
					if examplePayload == "" {
						examplePayload = testingSignalPayloads[example.Signal]
					}

					switch example.Signal {
					case "logs":
						result, err = executor.ExecuteLogs(example.Config, examplePayload)
					case "traces":
						result, err = executor.ExecuteTraces(example.Config, examplePayload)
					case "metrics":
						result, err = executor.ExecuteMetrics(example.Config, examplePayload)
					default:
						t.Fatalf("Unknown signal type: %s", example.Signal)
					}

					require.NoError(t, err)
					require.NotNil(t, result)
					require.NotEmpty(t, result.Value)
					require.Empty(t, result.Error)

				})

				debuggableExecutor, debuggable := executor.(DebuggableExecutor)
				if debuggable {
					debugger, err := debuggableExecutor.Debugger()
					if err == nil {
						t.Run(fmt.Sprintf("%s_debugger", example.Name), func(t *testing.T) {
							var result *Result
							var err error
							examplePayload := example.Payload
							if examplePayload == "" {
								examplePayload = testingSignalPayloads[example.Signal]
							}

							switch example.Signal {
							case "logs":
								result, err = debugger.DebugLogs(example.Config, examplePayload)
							case "traces":
								result, err = debugger.DebugTraces(example.Config, examplePayload)
							case "metrics":
								result, err = debugger.DebugMetrics(example.Config, examplePayload)
							default:
								t.Fatalf("Unknown signal type: %s", example.Signal)
							}

							require.NoError(t, err)
							require.NotNil(t, result)
							require.NotEmpty(t, result.Value)
							require.Empty(t, result.Error)

						})
					}
				}
			}
		})
	}
}

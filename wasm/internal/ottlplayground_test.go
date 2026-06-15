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
	"errors"
	"fmt"
	"testing"

	"github.com/elastic/ottl-playground/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func Test_NewErrorResult(t *testing.T) {
	logs := "execution logs"
	err := "error"
	expected := map[string]any{
		"debug":         false,
		"line":          float64(0),
		"value":         "",
		"logs":          logs,
		"error":         err,
		"executionTime": float64(0),
	}

	result := internal.NewErrorResult(err, logs).AsRaw()
	assert.Equal(t, expected, result)
}

func Test_NewErrorResult_Debug(t *testing.T) {
	logs := "debug execution logs"
	err := "debug error"

	result := internal.NewErrorResult(err, logs)
	result.Debug = true

	expected := map[string]any{
		"debug":         true,
		"line":          float64(0),
		"value":         "",
		"logs":          logs,
		"error":         err,
		"executionTime": float64(0),
	}

	assert.Equal(t, expected, result.AsRaw())
}

func Test_ExecuteStatements_UnsupportedExecutor(t *testing.T) {
	config := "empty"
	otlpDataType := "logs"
	otlpDataPayload := "{}"
	executorName := "unsupported_processor"

	expectedError := fmt.Sprintf("unsupported executor %s", executorName)
	result := Execute(config, otlpDataType, otlpDataPayload, executorName, false)
	assert.Equal(t, "", result["value"])
	assert.Equal(t, "", result["logs"])
	assert.Equal(t, expectedError, result["error"])
	assert.GreaterOrEqual(t, result["executionTime"], float64(0))
}

func Test_ExecuteStatements_UnsupportedOTLPType(t *testing.T) {
	config := "empty"
	otlpDataType := "unsupported_datatype"
	otlpDataPayload := "{}"
	executorName := "transform_processor"

	expectedError := fmt.Sprintf("unsupported OTLP signal type %s", otlpDataType)

	result := Execute(config, otlpDataType, otlpDataPayload, executorName, false)
	assert.Equal(t, "", result["value"])
	assert.Equal(t, "", result["logs"])
	assert.Equal(t, expectedError, result["error"])
	assert.GreaterOrEqual(t, result["executionTime"], float64(0))
}

func Test_ExecuteStatements(t *testing.T) {
	tests := []struct {
		name           string
		otlpDataType   string
		executorFunc   string
		expectedOutput string
		expectedError  error
		debug          bool
	}{
		{
			name:           "Logs Success",
			otlpDataType:   "logs",
			executorFunc:   "ExecuteLogs",
			expectedOutput: "log output",
			debug:          false,
		},
		{
			name:          "Logs Error",
			otlpDataType:  "logs",
			executorFunc:  "ExecuteLogs",
			expectedError: errors.New("ExecuteLogs execution error"),
			debug:         false,
		},
		{
			name:           "Traces Success",
			otlpDataType:   "traces",
			executorFunc:   "ExecuteTraces",
			expectedOutput: "trace output",
			debug:          false,
		},
		{
			name:          "Traces Error",
			otlpDataType:  "traces",
			executorFunc:  "ExecuteTraces",
			expectedError: errors.New("ExecuteTraces execution error"),
			debug:         false,
		},
		{
			name:           "Metrics Success",
			otlpDataType:   "metrics",
			executorFunc:   "ExecuteMetrics",
			expectedOutput: "metric output",
			debug:          false,
		},
		{
			name:          "Metrics Error",
			otlpDataType:  "metrics",
			executorFunc:  "ExecuteMetrics",
			expectedError: errors.New("ExecuteMetrics execution error"),
			debug:         false,
		},
		{
			name:           "Profiles Success",
			otlpDataType:   "profiles",
			executorFunc:   "ExecuteProfiles",
			expectedOutput: "profile output",
		},
		{
			name:          "Profiles Error",
			otlpDataType:  "profiles",
			executorFunc:  "ExecuteProfiles",
			expectedError: errors.New("ExecuteProfileStatements execution error"),
		},
	}

	var (
		testConfig      = "empty"
		ottlDataPayload = "{}"
	)

	_, observedLogs := internal.NewLogObserver(zap.NewNop().Core(), zap.NewDevelopmentEncoderConfig())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executorName := tt.name
			mockExecutor := &MockExecutor{
				metadata: internal.Metadata{
					ID:   executorName,
					Name: executorName,
				},
				debuggable: false, // Explicitly not debuggable for these tests
				debugger:   nil,
			}

			registerStatementsExecutor(mockExecutor)
			mockExecutor.On(tt.executorFunc, testConfig, ottlDataPayload).Return(&internal.Result{Value: tt.expectedOutput, Debug: tt.debug}, tt.expectedError)
			if tt.expectedError != nil {
				mockExecutor.On("ObservedLogs").Return(observedLogs)
			}

			result := Execute(testConfig, tt.otlpDataType, ottlDataPayload, executorName, tt.debug)

			if tt.expectedError != nil {
				assert.Empty(t, result["value"])
				expectedErrorMsg := fmt.Sprintf("unable to run %s configuration. Error: %v", tt.otlpDataType, tt.expectedError)
				assert.Contains(t, result["error"], expectedErrorMsg)
				assert.Equal(t, result["executionTime"], float64(0))
			} else {
				assert.Equal(t, tt.expectedOutput, result["value"])
				assert.Equal(t, tt.debug, result["debug"])
				assert.NotContains(t, result, "error")
				assert.GreaterOrEqual(t, result["executionTime"], float64(0))
			}

			mockExecutor.AssertExpectations(t)
		})
	}
}

func Test_TakeObserved_Logs(t *testing.T) {
	mockExecutor := new(MockExecutor)
	core, observedLogs := internal.NewLogObserver(zap.DebugLevel, zap.NewDevelopmentEncoderConfig())
	mockExecutor.On("ObservedLogs").Return(observedLogs)

	logger := zap.New(core)
	logger.Debug("debug logs")

	logs := mockExecutor.ObservedLogs().TakeAllString()

	assert.Contains(t, logs, "debug logs")
	mockExecutor.AssertExpectations(t)
}

func Test_ExecuteStatements_Debug_UnsupportedExecutor(t *testing.T) {
	config := "empty"
	otlpDataType := "logs"
	otlpDataPayload := "{}"
	executorName := "non_debuggable_processor"

	mockExecutor := &MockExecutor{
		metadata: internal.Metadata{
			ID:   executorName,
			Name: executorName,
		},
		debuggable: false,
		debugger:   nil,
	}

	registerStatementsExecutor(mockExecutor)
	_, observedLogs := internal.NewLogObserver(zap.NewNop().Core(), zap.NewDevelopmentEncoderConfig())
	mockExecutor.On("ObservedLogs").Return(observedLogs)

	result := Execute(config, otlpDataType, otlpDataPayload, executorName, true)

	expectedError := fmt.Sprintf("unable to run %s configuration. Error: executor %q does not support debugging", otlpDataType, executorName)
	assert.Equal(t, "", result["value"])
	assert.Contains(t, result["error"], expectedError)
	assert.GreaterOrEqual(t, result["executionTime"], float64(0))
}

func Test_ExecuteStatements_Debug(t *testing.T) {
	tests := []struct {
		name           string
		otlpDataType   string
		debugFunc      string
		expectedOutput string
		expectedError  error
	}{
		{
			name:           "Debug Logs Success",
			otlpDataType:   "logs",
			debugFunc:      "DebugLogs",
			expectedOutput: "[{\"value\":\"debug log output\",\"debug\":true,\"line\":10}]",
		},
		{
			name:          "Debug Logs Error",
			otlpDataType:  "logs",
			debugFunc:     "DebugLogs",
			expectedError: errors.New("DebugLogs execution error"),
		},
		{
			name:           "Debug Traces Success",
			otlpDataType:   "traces",
			debugFunc:      "DebugTraces",
			expectedOutput: "[{\"value\":\"debug trace output\",\"debug\":true,\"line\":20}]",
		},
		{
			name:          "Debug Traces Error",
			otlpDataType:  "traces",
			debugFunc:     "DebugTraces",
			expectedError: errors.New("DebugTraces execution error"),
		},
		{
			name:           "Debug Metrics Success",
			otlpDataType:   "metrics",
			debugFunc:      "DebugMetrics",
			expectedOutput: "[{\"value\":\"debug metric output\",\"debug\":true,\"line\":30}]",
		},
		{
			name:          "Debug Metrics Error",
			otlpDataType:  "metrics",
			debugFunc:     "DebugMetrics",
			expectedError: errors.New("DebugMetrics execution error"),
		},
		{
			name:           "Debug Profiles Success",
			otlpDataType:   "profiles",
			debugFunc:      "DebugProfiles",
			expectedOutput: "[{\"value\":\"debug profile output\",\"debug\":true,\"line\":30}]",
		},
		{
			name:          "Debug Profiles Error",
			otlpDataType:  "profiles",
			debugFunc:     "DebugProfiles",
			expectedError: errors.New("DebugProfiles execution error"),
		},
	}

	var (
		testConfig      = "empty"
		ottlDataPayload = "{}"
	)

	_, observedLogs := internal.NewLogObserver(zap.NewNop().Core(), zap.NewDevelopmentEncoderConfig())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executorName := tt.name
			mockDebugger := &MockDebugger{}
			mockExecutor := &MockExecutor{
				metadata: internal.Metadata{
					ID:   executorName,
					Name: executorName,
				},
				debuggable: true,
				debugger:   mockDebugger,
			}

			registerStatementsExecutor(mockExecutor)

			if tt.expectedError != nil {
				mockDebugger.On(tt.debugFunc, testConfig, ottlDataPayload).Return(&internal.Result{Debug: true}, tt.expectedError)
				mockExecutor.On("ObservedLogs").Return(observedLogs)
			} else {
				result := &internal.Result{
					Value: tt.expectedOutput,
					Debug: true,
				}
				mockDebugger.On(tt.debugFunc, testConfig, ottlDataPayload).Return(result, nil)
			}

			result := Execute(testConfig, tt.otlpDataType, ottlDataPayload, executorName, true)

			if tt.expectedError != nil {
				assert.Empty(t, result["value"])
				expectedErrorMsg := fmt.Sprintf("unable to run %s configuration. Error: %v", tt.otlpDataType, tt.expectedError)
				assert.Contains(t, result["error"], expectedErrorMsg)
				assert.Equal(t, result["executionTime"], float64(0))
			} else {
				assert.Equal(t, tt.expectedOutput, result["value"])
				assert.Equal(t, true, result["debug"])
				assert.NotContains(t, result, "error")
				assert.GreaterOrEqual(t, result["executionTime"], float64(0))
			}

			mockDebugger.AssertExpectations(t)
		})
	}
}

func Test_ExecuteStatements_Debug_UnsupportedOTLPType(t *testing.T) {
	config := "empty"
	otlpDataType := "unsupported_datatype"
	otlpDataPayload := "{}"
	executorName := "debug_processor"

	mockDebugger := &MockDebugger{}
	mockExecutor := &MockExecutor{
		metadata: internal.Metadata{
			ID:   executorName,
			Name: executorName,
		},
		debuggable: true,
		debugger:   mockDebugger,
	}

	registerStatementsExecutor(mockExecutor)

	result := Execute(config, otlpDataType, otlpDataPayload, executorName, true)

	expectedError := fmt.Sprintf("unsupported OTLP signal type %s", otlpDataType)
	assert.Equal(t, "", result["value"])
	assert.Equal(t, "", result["logs"])
	assert.Equal(t, expectedError, result["error"])
	assert.GreaterOrEqual(t, result["executionTime"], float64(0))
}

func Test_ExecuteStatements_Debug_DebuggerCreationError(t *testing.T) {
	config := "empty"
	otlpDataType := "logs"
	otlpDataPayload := "{}"
	executorName := "failing_debugger_processor"

	// Create a mock executor that returns an error when Debugger() is called
	mockExecutor := &MockExecutor{
		metadata: internal.Metadata{
			ID:   executorName,
			Name: executorName,
		},
		debuggable: true,
		debugger:   nil, // This will cause Debugger() to return an error
	}

	registerStatementsExecutor(mockExecutor)
	_, observedLogs := internal.NewLogObserver(zap.NewNop().Core(), zap.NewDevelopmentEncoderConfig())
	mockExecutor.On("ObservedLogs").Return(observedLogs)

	result := Execute(config, otlpDataType, otlpDataPayload, executorName, true)

	expectedErrorMsg := fmt.Sprintf("unable to run %s configuration. Error: executor %q does not support debugging", otlpDataType, executorName)
	assert.Equal(t, "", result["value"])
	assert.Contains(t, result["error"], expectedErrorMsg)
	assert.Equal(t, result["executionTime"], float64(0))
}

type MockExecutor struct {
	mock.Mock
	metadata   internal.Metadata
	debuggable bool
	debugger   *MockDebugger
}

func (m *MockExecutor) ExecuteLogs(config, payload string) (*internal.Result, error) {
	args := m.Called(config, payload)
	return args.Get(0).(*internal.Result), args.Error(1)
}

func (m *MockExecutor) ExecuteTraces(config, payload string) (*internal.Result, error) {
	args := m.Called(config, payload)
	return args.Get(0).(*internal.Result), args.Error(1)
}

func (m *MockExecutor) ExecuteMetrics(config, payload string) (*internal.Result, error) {
	args := m.Called(config, payload)
	return args.Get(0).(*internal.Result), args.Error(1)
}

func (m *MockExecutor) ExecuteProfiles(config, payload string) (*internal.Result, error) {
	args := m.Called(config, payload)
	return args.Get(0).(*internal.Result), args.Error(1)
}

func (m *MockExecutor) ObservedLogs() *internal.ObservedLogs {
	args := m.Called()
	return args.Get(0).(*internal.ObservedLogs)
}

func (m *MockExecutor) Metadata() *internal.Metadata {
	return &m.metadata
}

func (m *MockExecutor) Debugger() (internal.Debugger, error) {
	if !m.debuggable || m.debugger == nil {
		return nil, fmt.Errorf("executor %q does not support debugging", m.metadata.Name)
	}
	return m.debugger, nil
}

type MockDebugger struct {
	mock.Mock
}

func (m *MockDebugger) DebugLogs(config, payload string) (*internal.Result, error) {
	args := m.Called(config, payload)
	return args.Get(0).(*internal.Result), args.Error(1)
}

func (m *MockDebugger) DebugTraces(config, payload string) (*internal.Result, error) {
	args := m.Called(config, payload)
	return args.Get(0).(*internal.Result), args.Error(1)
}

func (m *MockDebugger) DebugMetrics(config, payload string) (*internal.Result, error) {
	args := m.Called(config, payload)
	return args.Get(0).(*internal.Result), args.Error(1)
}

func (m *MockDebugger) DebugProfiles(config, payload string) (*internal.Result, error) {
	args := m.Called(config, payload)
	return args.Get(0).(*internal.Result), args.Error(1)
}

func (m *MockDebugger) ObservedLogs() *internal.ObservedLogs {
	args := m.Called()
	return args.Get(0).(*internal.ObservedLogs)
}

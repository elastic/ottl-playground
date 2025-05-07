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

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/ottlplayground/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func Test_NewErrorResult(t *testing.T) {
	logs := "execution logs"
	err := "error"
	expected := map[string]any{
		"value":         "",
		"logs":          logs,
		"error":         err,
		"executionTime": int64(0),
	}

	result := NewErrorResult(err, logs)
	assert.Equal(t, expected, result)
}

func Test_ExecuteStatements_UnsupportedExecutor(t *testing.T) {
	config := "empty"
	otlpDataType := "logs"
	otlpDataPayload := "{}"
	executorName := "unsupported_processor"

	expectedError := fmt.Sprintf("unsupported evaluator %s", executorName)
	result := ExecuteStatements(config, otlpDataType, otlpDataPayload, executorName)
	assert.Equal(t, "", result["value"])
	assert.Equal(t, "", result["logs"])
	assert.Equal(t, expectedError, result["error"])
	assert.GreaterOrEqual(t, result["executionTime"], int64(0))
}

func Test_ExecuteStatements_UnsupportedOTLPType(t *testing.T) {
	config := "empty"
	otlpDataType := "unsupported_datatype"
	otlpDataPayload := "{}"
	executorName := "transform_processor"

	expectedError := fmt.Sprintf("unsupported OTLP data type %s", otlpDataType)

	result := ExecuteStatements(config, otlpDataType, otlpDataPayload, executorName)
	assert.Equal(t, "", result["value"])
	assert.Equal(t, "", result["logs"])
	assert.Equal(t, expectedError, result["error"])
	assert.GreaterOrEqual(t, result["executionTime"], int64(0))
}

func Test_ExecuteStatements(t *testing.T) {
	tests := []struct {
		name           string
		otlpDataType   string
		executorFunc   string
		expectedOutput string
		expectedError  error
	}{
		{
			name:           "Logs Success",
			otlpDataType:   "logs",
			executorFunc:   "ExecuteLogStatements",
			expectedOutput: "log output",
		},
		{
			name:          "Logs Error",
			otlpDataType:  "logs",
			executorFunc:  "ExecuteLogStatements",
			expectedError: errors.New("ExecuteLogStatements execution error"),
		},
		{
			name:           "Traces Success",
			otlpDataType:   "traces",
			executorFunc:   "ExecuteTraceStatements",
			expectedOutput: "trace output",
		},
		{
			name:          "Traces Error",
			otlpDataType:  "traces",
			executorFunc:  "ExecuteTraceStatements",
			expectedError: errors.New("ExecuteTraceStatements execution error"),
		},
		{
			name:           "Metrics Success",
			otlpDataType:   "metrics",
			executorFunc:   "ExecuteMetricStatements",
			expectedOutput: "metric output",
		},
		{
			name:          "Metrics Error",
			otlpDataType:  "metrics",
			executorFunc:  "ExecuteMetricStatements",
			expectedError: errors.New("ExecuteMetricStatements execution error"),
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
			}

			registerStatementsExecutor(mockExecutor)
			mockExecutor.On(tt.executorFunc, testConfig, ottlDataPayload).Return(tt.expectedOutput, tt.expectedError)
			mockExecutor.On("ObservedLogs").Return(observedLogs)

			result := ExecuteStatements(testConfig, tt.otlpDataType, ottlDataPayload, executorName)

			if tt.expectedError != nil {
				assert.Empty(t, result["value"])
				expectedErrorMsg := fmt.Sprintf("unable to run %s statements. Error: %v", tt.otlpDataType, tt.expectedError)
				assert.Contains(t, result["error"], expectedErrorMsg)
				assert.Equal(t, result["executionTime"], int64(0))
			} else {
				assert.Equal(t, tt.expectedOutput, result["value"])
				assert.NotContains(t, result, "error")
				assert.GreaterOrEqual(t, result["executionTime"], int64(0))
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

	logs := takeObservedLogs(mockExecutor)

	assert.Contains(t, logs, "debug logs")
	mockExecutor.AssertExpectations(t)
}

type MockExecutor struct {
	mock.Mock
	metadata internal.Metadata
}

func (m *MockExecutor) ExecuteLogStatements(config, payload string) ([]byte, error) {
	args := m.Called(config, payload)
	return []byte(args.String(0)), args.Error(1)
}

func (m *MockExecutor) ExecuteTraceStatements(config, payload string) ([]byte, error) {
	args := m.Called(config, payload)
	return []byte(args.String(0)), args.Error(1)
}

func (m *MockExecutor) ExecuteMetricStatements(config, payload string) ([]byte, error) {
	args := m.Called(config, payload)
	return []byte(args.String(0)), args.Error(1)
}

func (m *MockExecutor) ObservedLogs() *internal.ObservedLogs {
	args := m.Called()
	return args.Get(0).(*internal.ObservedLogs)
}

func (m *MockExecutor) Metadata() internal.Metadata {
	return m.metadata
}

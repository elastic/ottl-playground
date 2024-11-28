// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/open-telemetry/opentelemetry-collector-contrib/cmd/ottlplayground/internal"
)

func Test_NewErrorResult(t *testing.T) {
	logs := "execution logs"
	err := "error"
	expected := map[string]any{
		"value": "",
		"logs":  logs,
		"error": err,
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
	expected := map[string]any{
		"value": "",
		"logs":  "",
		"error": expectedError,
	}

	result := ExecuteStatements(config, otlpDataType, otlpDataPayload, executorName)
	assert.Equal(t, expected, result)
}

func Test_ExecuteStatements_UnsupportedOTLPType(t *testing.T) {
	config := "empty"
	otlpDataType := "unsupported_datatype"
	otlpDataPayload := "{}"
	executorName := "transform_processor"

	expectedError := fmt.Sprintf("unsupported OTLP data type %s", otlpDataType)
	expected := map[string]any{
		"value": "",
		"logs":  "",
		"error": expectedError,
	}

	result := ExecuteStatements(config, otlpDataType, otlpDataPayload, executorName)
	assert.Equal(t, expected, result)
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

	_, observedLogs := observer.New(zap.NewNop().Core())
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
			} else {
				assert.Equal(t, tt.expectedOutput, result["value"])
				assert.NotContains(t, result, "error")
			}

			mockExecutor.AssertExpectations(t)
		})
	}
}

func Test_TakeObserved_Logs(t *testing.T) {
	mockExecutor := new(MockExecutor)
	core, observedLogs := observer.New(zap.DebugLevel)
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

func (m *MockExecutor) ObservedLogs() *observer.ObservedLogs {
	args := m.Called()
	return args.Get(0).(*observer.ObservedLogs)
}

func (m *MockExecutor) Metadata() internal.Metadata {
	return m.metadata
}

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
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewErrorResult(t *testing.T) {
	errorMsg := "test error message"
	logs := "test logs"

	result := NewErrorResult(errorMsg, logs)

	require.NotNil(t, result)
	assert.Equal(t, errorMsg, result.Error)
	assert.Equal(t, logs, result.Logs)
	assert.Empty(t, result.Value)
	assert.Equal(t, int64(0), result.ExecutionTime)
	assert.False(t, result.Debug)
	assert.Equal(t, int64(0), result.Line)
}

func Test_Result_executeWithTimer(t *testing.T) {
	result := &Result{}

	executed := false
	err := result.executeWithTimer(func() error {
		executed = true
		time.Sleep(10 * time.Millisecond) // Add some execution time
		return nil
	})

	require.NoError(t, err)
	assert.True(t, executed)
	assert.Greater(t, result.ExecutionTime, int64(0))
	assert.Less(t, result.ExecutionTime, int64(1000)) // Should be less than 1 second
}

func Test_Result_executeWithTimer_WithError(t *testing.T) {
	result := &Result{}
	expectedError := errors.New("test error")

	err := result.executeWithTimer(func() error {
		time.Sleep(5 * time.Millisecond)
		return expectedError
	})

	assert.Equal(t, expectedError, err)
	assert.Greater(t, result.ExecutionTime, int64(0))
}

func Test_Result_startTimer_and_stopTimer(t *testing.T) {
	result := &Result{}

	// Test initial state
	assert.Equal(t, int64(0), result.ExecutionTime)

	result.startTimer()
	time.Sleep(10 * time.Millisecond)
	result.stopTimer()

	assert.Greater(t, result.ExecutionTime, int64(0))
	assert.Less(t, result.ExecutionTime, int64(1000))
}

func Test_Result_AsRaw(t *testing.T) {
	result := &Result{
		Value:         "test value",
		ExecutionTime: 123,
		Error:         "test error",
		Logs:          "test logs",
		Debug:         true,
		Line:          456,
	}

	raw := result.AsRaw()

	require.NotNil(t, raw)
	assert.Equal(t, "test value", raw["value"])
	assert.Equal(t, float64(123), raw["executionTime"]) // JSON unmarshalling converts numbers to float64
	assert.Equal(t, "test error", raw["error"])
	assert.Equal(t, "test logs", raw["logs"])
	assert.Equal(t, true, raw["debug"])
	assert.Equal(t, float64(456), raw["line"])
}

func Test_Result_AsRaw_WithJSON(t *testing.T) {
	jsonStr := `{"key": "value"}`
	result := &Result{
		Value: "test value",
		JSON:  &jsonStr,
	}

	raw := result.AsRaw()

	require.NotNil(t, raw)
	assert.Equal(t, "test value", raw["value"])
	assert.Equal(t, jsonStr, raw["json"])
}

func Test_Result_AsRaw_WithNilJSON(t *testing.T) {
	result := &Result{
		Value: "test value",
		JSON:  nil,
	}

	raw := result.AsRaw()

	require.NotNil(t, raw)
	assert.Equal(t, "test value", raw["value"])
	_, jsonExists := raw["json"]
	assert.False(t, jsonExists, "json field should not exist when nil")
}

// Mock Observable for testing
type mockObservable struct {
	observedLogs *ObservedLogs
}

func (m *mockObservable) ObservedLogs() *ObservedLogs {
	return m.observedLogs
}

func Test_newExecutionResult_Success(t *testing.T) {
	observedLogs := &ObservedLogs{}
	observedLogs.add(LoggedEntry{consoleEncodedEntry: "test log entry"})

	mock := &mockObservable{observedLogs: observedLogs}

	valueMarshaller := func(input string) ([]byte, error) {
		return []byte(input), nil
	}

	command := func() (string, error) {
		time.Sleep(5 * time.Millisecond) // Simulate some work
		return "test result", nil
	}

	result, err := newExecutionResult(mock, valueMarshaller, command)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "test result", result.Value)
	assert.Greater(t, result.ExecutionTime, int64(0))
	assert.Equal(t, "test log entry", result.Logs)
}

func Test_newExecutionResult_CommandError(t *testing.T) {
	observedLogs := &ObservedLogs{}
	mock := &mockObservable{observedLogs: observedLogs}

	valueMarshaller := func(input string) ([]byte, error) {
		return []byte(input), nil
	}

	expectedError := errors.New("command failed")
	command := func() (string, error) {
		return "", expectedError
	}

	result, err := newExecutionResult(mock, valueMarshaller, command)

	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

func Test_newExecutionResult_MarshallerError(t *testing.T) {
	observedLogs := &ObservedLogs{}
	mock := &mockObservable{observedLogs: observedLogs}

	expectedError := errors.New("marshaller failed")
	valueMarshaller := func(input string) ([]byte, error) {
		return nil, expectedError
	}

	command := func() (string, error) {
		return "test result", nil
	}

	result, err := newExecutionResult(mock, valueMarshaller, command)

	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

func Test_newExecutionResult_WithComplexType(t *testing.T) {
	observedLogs := &ObservedLogs{}
	mock := &mockObservable{observedLogs: observedLogs}

	type complexType struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	valueMarshaller := func(input complexType) ([]byte, error) {
		return json.Marshal(input)
	}

	command := func() (complexType, error) {
		return complexType{Name: "test", Value: 42}, nil
	}

	result, err := newExecutionResult(mock, valueMarshaller, command)

	require.NoError(t, err)
	require.NotNil(t, result)

	var unmarshalledValue complexType
	err = json.Unmarshal([]byte(result.Value), &unmarshalledValue)
	require.NoError(t, err)
	assert.Equal(t, "test", unmarshalledValue.Name)
	assert.Equal(t, 42, unmarshalledValue.Value)
}

func Test_Result_TimingAccuracy(t *testing.T) {
	result := &Result{}

	// Test that timing is reasonably accurate
	expectedSleep := 50 * time.Millisecond

	err := result.executeWithTimer(func() error {
		time.Sleep(expectedSleep)
		return nil
	})

	require.NoError(t, err)

	// Allow some tolerance for timing variations
	assert.GreaterOrEqual(t, result.ExecutionTime, int64(40)) // At least 40ms
	assert.LessOrEqual(t, result.ExecutionTime, int64(100))   // At most 100ms
}

func Test_Result_MultipleTimerCalls(t *testing.T) {
	result := &Result{}

	// First timing
	result.startTimer()
	time.Sleep(10 * time.Millisecond)
	result.stopTimer()
	firstTime := result.ExecutionTime

	// Second timing should reset
	result.startTimer()
	time.Sleep(20 * time.Millisecond)
	result.stopTimer()
	secondTime := result.ExecutionTime

	assert.Greater(t, firstTime, int64(0))
	assert.Greater(t, secondTime, int64(0))
	assert.Greater(t, secondTime, firstTime) // Second timing should be longer
}

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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Test_ObservedLogs_Len(t *testing.T) {
	observedLogs := &ObservedLogs{}

	assert.Equal(t, 0, observedLogs.Len())

	observedLogs.add(LoggedEntry{consoleEncodedEntry: "entry1"})
	observedLogs.add(LoggedEntry{consoleEncodedEntry: "entry2"})

	assert.Equal(t, 2, observedLogs.Len())
}

func Test_ObservedLogs_All(t *testing.T) {
	observedLogs := &ObservedLogs{}

	result := observedLogs.All()
	assert.Len(t, result, 0)

	entry1 := LoggedEntry{consoleEncodedEntry: "entry1"}
	entry2 := LoggedEntry{consoleEncodedEntry: "entry2"}
	observedLogs.add(entry1)
	observedLogs.add(entry2)

	result = observedLogs.All()
	require.Len(t, result, 2)
	assert.Equal(t, entry1.consoleEncodedEntry, result[0].consoleEncodedEntry)
	assert.Equal(t, entry2.consoleEncodedEntry, result[1].consoleEncodedEntry)

	assert.Equal(t, 2, observedLogs.Len())
}

func Test_ObservedLogs_TakeAll(t *testing.T) {
	observedLogs := &ObservedLogs{}

	entry1 := LoggedEntry{consoleEncodedEntry: "entry1"}
	entry2 := LoggedEntry{consoleEncodedEntry: "entry2"}
	observedLogs.add(entry1)
	observedLogs.add(entry2)

	result := observedLogs.TakeAll()
	require.Len(t, result, 2)
	assert.Equal(t, entry1.consoleEncodedEntry, result[0].consoleEncodedEntry)
	assert.Equal(t, entry2.consoleEncodedEntry, result[1].consoleEncodedEntry)

	assert.Equal(t, 0, observedLogs.Len())

	result = observedLogs.TakeAll()
	assert.Len(t, result, 0)
}

func Test_ObservedLogs_TakeAllString(t *testing.T) {
	observedLogs := &ObservedLogs{}

	result := observedLogs.TakeAllString()
	assert.Equal(t, "", result)

	observedLogs.add(LoggedEntry{consoleEncodedEntry: "entry1"})
	observedLogs.add(LoggedEntry{consoleEncodedEntry: "entry2"})

	result = observedLogs.TakeAllString()
	assert.Equal(t, "entry1entry2", result)

	assert.Equal(t, 0, observedLogs.Len())
}

func Test_ObservedLogs_add(t *testing.T) {
	observedLogs := &ObservedLogs{}

	entry := LoggedEntry{consoleEncodedEntry: "test entry"}
	observedLogs.add(entry)

	assert.Equal(t, 1, observedLogs.Len())
	result := observedLogs.All()
	require.Len(t, result, 1)
	assert.Equal(t, entry.consoleEncodedEntry, result[0].consoleEncodedEntry)
}

func Test_NewLogObserver(t *testing.T) {
	level := zap.DebugLevel
	config := zap.NewDevelopmentEncoderConfig()
	core, observedLogs := NewLogObserver(level, config)

	require.NotNil(t, core)
	require.NotNil(t, observedLogs)

	observer, ok := core.(*contextObserver)
	require.True(t, ok)
	assert.Equal(t, observedLogs, observer.logs)
	assert.NotNil(t, observer.config)
}

func Test_contextObserver_Level(t *testing.T) {
	level := zap.DebugLevel
	config := zap.NewDevelopmentEncoderConfig()
	core, _ := NewLogObserver(level, config)

	observer, ok := core.(*contextObserver)
	require.True(t, ok)
	assert.Equal(t, zap.DebugLevel, observer.Level())
}

func Test_contextObserver_Enabled(t *testing.T) {
	level := zap.InfoLevel
	config := zap.NewDevelopmentEncoderConfig()
	core, _ := NewLogObserver(level, config)

	contextObserver, ok := core.(*contextObserver)
	require.True(t, ok)

	assert.True(t, contextObserver.Enabled(zap.InfoLevel))
	assert.True(t, contextObserver.Enabled(zap.WarnLevel))
	assert.True(t, contextObserver.Enabled(zap.ErrorLevel))
	assert.False(t, contextObserver.Enabled(zap.DebugLevel))
}

func Test_contextObserver_Check(t *testing.T) {
	level := zap.DebugLevel
	config := zap.NewDevelopmentEncoderConfig()
	core, _ := NewLogObserver(level, config)

	contextObserver, ok := core.(*contextObserver)
	require.True(t, ok)
	entry := zapcore.Entry{Level: zap.InfoLevel, Message: "test message"}
	checkedEntry := &zapcore.CheckedEntry{}

	result := contextObserver.Check(entry, checkedEntry)
	require.NotNil(t, result)
}

func Test_contextObserver_With(t *testing.T) {
	level := zap.DebugLevel
	config := zap.NewDevelopmentEncoderConfig()
	core, observedLogs := NewLogObserver(level, config)

	observer, ok := core.(*contextObserver)
	require.True(t, ok)
	fields := []zapcore.Field{zap.String("key", "value")}

	newCore := observer.With(fields)
	require.NotNil(t, newCore)

	newContextObserver, ok := newCore.(*contextObserver)
	require.True(t, ok)
	assert.Equal(t, observedLogs, newContextObserver.logs)
	assert.Len(t, newContextObserver.context, 1)
	assert.Equal(t, "key", newContextObserver.context[0].Key)
}

func Test_contextObserver_Write(t *testing.T) {
	observer, logs := NewLogObserver(
		zapcore.NewNopCore(),
		zap.NewDevelopmentEncoderConfig())

	expectedConsoleEntry := "INFO\ttest message\t{\"key\": \"value\"}\n"
	entry := zapcore.Entry{
		Level:   zap.InfoLevel,
		Message: "test message",
	}

	err := observer.Write(entry, []zap.Field{zap.String("key", "value")})
	require.NoError(t, err)
	require.Len(t, logs.All(), 1)
	assert.Equal(t, expectedConsoleEntry, logs.All()[0].ConsoleEncodedEntry())
}

func Test_contextObserver_Sync(t *testing.T) {
	level := zap.DebugLevel
	config := zap.NewDevelopmentEncoderConfig()
	core, _ := NewLogObserver(level, config)

	contextObserver, ok := core.(*contextObserver)
	require.True(t, ok)
	err := contextObserver.Sync()
	assert.NoError(t, err)
}

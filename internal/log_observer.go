// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"go.uber.org/zap/zapcore"
	"sync"
)

type LoggedEntry struct {
	entry               zapcore.Entry
	consoleEncodedEntry string
}

func (e *LoggedEntry) ConsoleEncodedEntry() string {
	return e.consoleEncodedEntry
}

type ObservedLogs struct {
	mu   sync.RWMutex
	logs []LoggedEntry
}

func (o *ObservedLogs) Len() int {
	o.mu.RLock()
	n := len(o.logs)
	o.mu.RUnlock()
	return n
}

func (o *ObservedLogs) All() []LoggedEntry {
	o.mu.RLock()
	ret := make([]LoggedEntry, len(o.logs))
	copy(ret, o.logs)
	o.mu.RUnlock()
	return ret
}

func (o *ObservedLogs) TakeAll() []LoggedEntry {
	o.mu.Lock()
	ret := o.logs
	o.logs = nil
	o.mu.Unlock()
	return ret
}

func (o *ObservedLogs) add(log LoggedEntry) {
	o.mu.Lock()
	o.logs = append(o.logs, log)
	o.mu.Unlock()
}

func NewLogObserver(level zapcore.LevelEnabler, config zapcore.EncoderConfig) (zapcore.Core, *ObservedLogs) {
	ol := &ObservedLogs{}
	return &contextObserver{
		config:       config,
		LevelEnabler: level,
		logs:         ol,
	}, ol
}

type contextObserver struct {
	zapcore.LevelEnabler
	config  zapcore.EncoderConfig
	logs    *ObservedLogs
	context []zapcore.Field
}

var _ zapcore.Core = (*contextObserver)(nil)

func (co *contextObserver) Level() zapcore.Level {
	return zapcore.LevelOf(co.LevelEnabler)
}

func (co *contextObserver) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if co.Enabled(ent.Level) {
		return ce.AddCore(ent, co)
	}
	return ce
}

func (co *contextObserver) With(fields []zapcore.Field) zapcore.Core {
	return &contextObserver{
		LevelEnabler: co.LevelEnabler,
		logs:         co.logs,
		context:      append(co.context[:len(co.context):len(co.context)], fields...),
	}
}

func (co *contextObserver) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	encoder := zapcore.NewConsoleEncoder(co.config)
	encodedEntryBuffer, err := encoder.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}

	co.logs.add(LoggedEntry{entry, encodedEntryBuffer.String()})
	return nil
}

func (co *contextObserver) Sync() error {
	return nil
}

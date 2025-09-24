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
	"time"
)

type Result struct {
	Value         string  `json:"value"`
	JSON          *string `json:"json,omitempty"`
	ExecutionTime int64   `json:"executionTime"`
	Error         string  `json:"error,omitempty"`
	Logs          string  `json:"logs"`
	Debug         bool    `json:"debug"`
	Line          int64   `json:"line"`
	start         time.Time
}

func NewErrorResult(err string, logs string) *Result {
	return &Result{
		Error: err,
		Logs:  logs,
	}
}

func (r *Result) executeWithTimer(fn func() error) error {
	r.startTimer()
	defer r.stopTimer()
	return fn()
}

func (r *Result) startTimer() {
	r.start = time.Now().Add(time.Duration(-r.ExecutionTime) * time.Millisecond)
}

func (r *Result) stopTimer() {
	r.ExecutionTime = time.Since(r.start).Milliseconds()
}

func (r *Result) AsRaw() map[string]any {
	res, err := structToMap(r)
	if err != nil {
		return map[string]any{
			"error": fmt.Sprintf("failed to convert result %v into raw map: %v", r, err),
		}
	}
	return res
}

func newExecutionResult[T any](
	observable Observable,
	valueMarshaller func(T) ([]byte, error),
	command func() (T, error),

) (*Result, error) {
	res := &Result{}
	res.start = time.Now()
	b, err := command()
	if err != nil {
		return nil, err
	}
	res.ExecutionTime = time.Since(res.start).Milliseconds()
	valueBytes, err := valueMarshaller(b)
	if err != nil {
		return nil, err
	}
	res.Value = string(valueBytes)
	res.Logs = observable.ObservedLogs().TakeAllString()
	return res, nil
}

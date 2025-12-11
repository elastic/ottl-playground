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
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor"
)

var filterProcessorConfigExamples = []ConfigExample{
	{
		Name:   "Drop specific metric and value",
		Signal: "metrics",
		Config: "filter: \n" +
			"  metrics:\n" +
			"    datapoint:\n" +
			`      - metric.name == "my.histogram" and count == 2`,
	},
	{
		Name:   "Drop spans",
		Signal: "traces",
		Config: "filter: \n" +
			"  traces:\n" +
			"    span:\n" +
			"      - kind == SPAN_KIND_INTERNAL",
	},
	{
		Name:   "Drop data by resource attribute",
		Signal: "traces",
		Config: "filter: \n" +
			"  traces:\n" +
			"    span:\n" +
			`      - IsMatch(resource.attributes["service.name"], "my-*")`,
	},
	{
		Name:   "Drop debug and trace logs",
		Signal: "logs",
		Config: "filter: \n" +
			"  logs:\n" +
			"    log_record:\n" +
			"      - severity_number != SEVERITY_NUMBER_UNSPECIFIED and severity_number < SEVERITY_NUMBER_INFO",
		Payload: `{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeLogs":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"logRecords":[{"timeUnixNano":"1544712660300000000","observedTimeUnixNano":"1544712660300000000","severityNumber":10,"severityText":"Information","traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","body":{"stringValue":"I'm an INFO log record"}},{"timeUnixNano":"1544712660300000000","observedTimeUnixNano":"1544712660300000000","severityNumber":5,"severityText":"Debug","traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","body":{"stringValue":"I'm a DEBUG log record"}}]}]}]}`,
	},
}

// NewFilterProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [filterprocessor].
func NewFilterProcessorExecutor() Executor {
	return NewJSONExecutor[filterprocessor.Config](
		newProcessorConsumer[filterprocessor.Config](filterprocessor.NewFactory()),
		newMetadata(
			ComponentTypeProcessor,
			"filter_processor",
			"Filter",
			"github.com/open-telemetry/opentelemetry-collector-contrib/processor/filterprocessor",
			"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/filterprocessor",
			withConfigExamples(filterProcessorConfigExamples...),
		),
	)
}

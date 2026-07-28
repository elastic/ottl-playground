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
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor"
)

var attributesProcessorConfigExamples = []ConfigExample{
	{
		Name:   "Insert Attribute for All Except Certain Metric Name",
		Signal: "metrics",
		Config: "attributes: \n" +
			"  exclude:\n" +
			"    match_type: strict\n" +
			"    metric_names: [\"my.gauge\"]\n" +
			"  actions:\n" +
			"    - key: my.insert.attr\n" +
			"      action: insert\n" +
			"      value: \"some inserted value\"",
	},
	{
		Name:   "Upsert Attribute for All Logs Above Certain Severity",
		Signal: "logs",
		Config: "attributes: \n" +
			"  include:\n" +
			"    log_severity_number:\n" +
			"      min: 10\n" +
			"      match_undefined: false\n" +
			"  actions:\n" +
			"    - key: my.upsert.attr\n" +
			"      action: upsert\n" +
			"      value: true",
	},
	{
		Name:   "Delete Attributes Matching Regex",
		Signal: "metrics",
		Config: "attributes: \n" +
			"  actions:\n" +
			"    - pattern: '^my\\.[a-z]+\\.attr$'\n" +
			"      action: delete",
	},
	{
		Name:   "Convert Attribute Type",
		Signal: "traces",
		Config: "attributes: \n" +
			"  actions:\n" +
			"    - key: http.response.status_code\n" +
			"      action: convert\n" +
			"      converted_type: string",
		Payload: `{"resourceSpans":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeSpans":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"spans":[{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","parentSpanId":"eee19b7ec3c1b173","name":"I'm a server span","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","kind":2,"attributes":[{"key":"http.response.status_code","value":{"intValue":"500"}}],"status":{}},{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b173","parentSpanId":"eee19b7ec3c1b173","name":"Me too","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","kind":1,"attributes":[{"key":"http.response.status_code","value":{"intValue":"500"}}],"status":{}}]}]}]}`,
	},
}

// NewAttributesProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [attributesprocessor].
func NewAttributesProcessorExecutor() Executor {
	return NewJSONExecutor[attributesprocessor.Config](
		newProcessorConsumer[attributesprocessor.Config](attributesprocessor.NewFactory()),
		newMetadata(
			ComponentTypeProcessor,
			"attributes_processor",
			"Attributes",
			"github.com/open-telemetry/opentelemetry-collector-contrib/processor/attributesprocessor",
			"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/attributesprocessor",
			withConfigExamples(attributesProcessorConfigExamples...),
		),
	)
}
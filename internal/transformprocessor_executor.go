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
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
)

var transformProcessorConfigExamples = []ConfigExample{
	{
		Name:   "Rename an attribute",
		Signal: "traces",
		Config: "transform: \n" +
			"  trace_statements:\n" +
			`    - set(resource.attributes["service.new_name"], resource.attributes["service.name"])` + "\n" +
			`    - delete_key(resource.attributes, "service.name")`,
	},
	{
		Name:   "Copy field to attributes",
		Signal: "logs",
		Config: "transform: \n" +
			"  log_statements:\n" +
			`    - set(log.attributes["body"], log.body)`,
	},
	{
		Name:   "Combine two attributes",
		Signal: "logs",
		Config: "transform: \n" +
			"  log_statements:\n" +
			`    - set(log.attributes["combined"], Concat([log.attributes["string.attribute"], log.attributes["boolean.attribute"]], " "))`,
	},
	{
		Name:   "Set a field",
		Signal: "logs",
		Config: "transform: \n" +
			"  log_statements:\n" +
			"    - set(log.severity_number, SEVERITY_NUMBER_INFO)\n" +
			`    - set(log.severity_text, "INFO")`,
	},
	{
		Name:   "Parse unstructured log",
		Signal: "logs",
		Config: "transform: \n" +
			"  log_statements:\n" +
			`    - 'merge_maps(log.attributes, ExtractPatterns(log.body, "Example (?P<example_type>[a-z\\.]+)"), "upsert")'`,
	},
	{
		Name:   "Conditionally set a field",
		Signal: "traces",
		Config: "transform: \n" +
			"  trace_statements:\n" +
			`   - set(span.attributes["server"], true) where span.kind == SPAN_KIND_SERVER`,
	},
	{
		Name:   "Update a resource attribute",
		Signal: "logs",
		Config: "transform: \n" +
			"  log_statements:\n" +
			`    - set(resource.attributes["service.name"], "mycompany-application") `,
	},
	{
		Name:   "Parse and manipulate JSON",
		Signal: "logs",
		Config: "transform: \n" +
			"  log_statements:\n" +
			`    - merge_maps(log.cache, ParseJSON(log.body), "upsert") where IsMatch(log.body, "^\\{")` + "\n" +
			`    - set(log.time, Time(log.cache["timestamp"], "%Y-%m-%dT%H:%M:%SZ"))` + "\n" +
			`    - set(log.severity_text, log.cache["level"])` + "\n" +
			`    - set(log.body, log.cache["message"])`,
		Payload: `{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeLogs":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"logRecords":[{"timeUnixNano":"1544712660300000000","observedTimeUnixNano":"1544712660300000000","severityNumber":10,"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","body":{"stringValue":"{\"timestamp\": \"2025-03-01T12:12:14Z\", \"level\":\"INFO\",\"message\":\"Elapsed time: 10ms\"}"}}]}]}]}`,
	},
	{
		Name:   "Parse and manipulate Timestamps",
		Signal: "metrics",
		Config: "transform: \n" +
			"  metric_statements:\n" +
			`    - set(resource.attributes["date"], String(TruncateTime(Time(resource.attributes["timestamp"], "%Y-%m-%dT%H:%M:%SZ"), Duration("24h"))))`,
	},
	{
		Name:   "Manipulate strings",
		Signal: "metrics",
		Config: "transform: \n" +
			"  metric_statements:\n" +
			`    - set(resource.attributes["service.name"], ConvertCase(Concat([resource.attributes["service.name"], scope.version], ".v"), "upper"))`,
	},
	{
		Name:   "Scale a metric",
		Signal: "metrics",
		Config: "transform: \n" +
			"  metric_statements:\n" +
			`    - scale_metric(10.0, "kWh") where metric.name == "my.gauge"`,
	},
	{
		Name:   "Dynamically rename a metric",
		Signal: "metrics",
		Config: "transform: \n" +
			"  metric_statements:\n" +
			`    - replace_pattern(metric.name, "my.(.+)", "metrics.$$1")`,
	},
	{
		Name:   "Aggregate a metric",
		Signal: "metrics",
		Config: "transform: \n" +
			"  metric_statements:\n" +
			`    - copy_metric(name="my.second.histogram") where metric.name == "my.histogram"` + "\n" +
			`    - aggregate_on_attributes("sum", []) where metric.name == "my.second.histogram"`,
	},
}

// NewTransformProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [transformprocessor].
func NewTransformProcessorExecutor() Executor {
	debugger := NewTransformProcessorDebugger()
	return NewJSONExecutor[transformprocessor.Config](
		newProcessorConsumer[transformprocessor.Config](transformprocessor.NewFactory()),
		newMetadata(
			ComponentTypeProcessor,
			"transform_processor",
			"Transform",
			"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/transformprocessor",
			"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/transformprocessor",
			withConfigExamples(transformProcessorConfigExamples...),
		),
		withDebugger[transformprocessor.Config](debugger),
	)
}

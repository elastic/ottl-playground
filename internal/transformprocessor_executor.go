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
	"go.opentelemetry.io/collector/featuregate"
)

var transformProcessorConfigExamples = []ConfigExample{
	{
		Name:   "Rename an attribute",
		Signal: "traces",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  trace_statements:\n" +
			`    - set(resource.attributes["service.new_name"], resource.attributes["service.name"])` + "\n" +
			`    - delete_key(resource.attributes, "service.name")`,
	},
	{
		Name:   "Copy field to attributes",
		Signal: "logs",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  log_statements:\n" +
			`    - set(log.attributes["body"], log.body)`,
	},
	{
		Name:   "Combine two attributes",
		Signal: "logs",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  log_statements:\n" +
			`    - set(log.attributes["combined"], Concat([log.attributes["string.attribute"], log.attributes["boolean.attribute"]], " "))`,
	},
	{
		Name:   "Set a field",
		Signal: "logs",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  log_statements:\n" +
			"    - set(log.severity_number, SEVERITY_NUMBER_INFO)\n" +
			`    - set(log.severity_text, "INFO")`,
	},
	{
		Name:   "Parse unstructured log",
		Signal: "logs",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  log_statements:\n" +
			`    - 'merge_maps(log.attributes, ExtractPatterns(log.body, "Example (?P<example_type>[a-z\\.]+)"), "upsert")'`,
	},
	{
		Name:   "Conditionally set a field",
		Signal: "traces",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  trace_statements:\n" +
			`   - set(span.attributes["server"], true) where span.kind == SPAN_KIND_SERVER`,
	},
	{
		Name:   "Update a resource attribute",
		Signal: "logs",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  log_statements:\n" +
			`    - set(resource.attributes["service.name"], "mycompany-application") `,
	},
	{
		Name:   "Parse and manipulate JSON",
		Signal: "logs",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
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
			"  error_mode: propagate\n" +
			"  metric_statements:\n" +
			`    - set(resource.attributes["date"], String(TruncateTime(Time(resource.attributes["timestamp"], "%Y-%m-%dT%H:%M:%SZ"), Duration("24h"))))`,
	},
	{
		Name:   "Manipulate strings",
		Signal: "metrics",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  metric_statements:\n" +
			`    - set(resource.attributes["service.name"], ConvertCase(Concat([resource.attributes["service.name"], scope.version], ".v"), "upper"))`,
	},
	{
		Name:   "Scale a metric",
		Signal: "metrics",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  metric_statements:\n" +
			`    - scale_metric(10.0, "kWh") where metric.name == "my.gauge"`,
	},
	{
		Name:   "Dynamically rename a metric",
		Signal: "metrics",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  metric_statements:\n" +
			`    - replace_pattern(metric.name, "my.(.+)", "metrics.$$1")`,
	},
	{
		Name:   "Aggregate a metric",
		Signal: "metrics",
		Config: "transform: \n" +
			"  error_mode: propagate\n" +
			"  metric_statements:\n" +
			`    - copy_metric(name="my.second.histogram") where metric.name == "my.histogram"` + "\n" +
			`    - aggregate_on_attributes("sum", []) where metric.name == "my.second.histogram"`,
	},
	{
		Name:   "Filter and rename attribute keys",
		Signal: "traces",
		Config: "transform:\n" +
			"  error_mode: propagate\n" +
			"  trace_statements:\n" +
			`    - set(span.attributes, Filter(span.attributes, (k, _) => HasPrefix(k, "http.")))` + "\n" +
			`    - set(resource.attributes, MapKeys(resource.attributes, (k, _) => Format("app.%s", [k])))`,
		Payload: `{"resourceSpans":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}},{"key":"host.name","value":{"stringValue":"web-01"}}]},"scopeSpans":[{"scope":{"name":"my.library","version":"1.0.0"},"spans":[{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","parentSpanId":"eee19b7ec3c1b173","name":"HTTP GET /api/users","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","kind":2,"attributes":[{"key":"http.method","value":{"stringValue":"GET"}},{"key":"http.url","value":{"stringValue":"/api/users"}},{"key":"http.status_code","value":{"intValue":"200"}},{"key":"db.system","value":{"stringValue":"postgresql"}},{"key":"internal.retry","value":{"boolValue":false}}],"status":{}}]}]}]}`,
	},
	{
		Name:   "Normalize attribute keys to snake_case",
		Signal: "logs",
		Config: "transform:\n" +
			"  error_mode: propagate\n" +
			"  log_statements:\n" +
			"    - set(log.attributes, MapKeys(log.attributes, (k, _) => ToSnakeCase(k)))",
		Payload: `{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"checkout-service"}}]},"scopeLogs":[{"scope":{"name":"checkout.handler"},"logRecords":[{"timeUnixNano":"1792764000000000000","severityNumber":9,"severityText":"INFO","body":{"stringValue":"Request completed"},"attributes":[{"key":"userId","value":{"stringValue":"u-48213"}},{"key":"requestPath","value":{"stringValue":"/cart/checkout"}},{"key":"responseTimeMs","value":{"intValue":"182"}},{"key":"OrderStatus","value":{"stringValue":"confirmed"}},{"key":"already_snake","value":{"stringValue":"unchanged"}}]}]}]}]}`,
	},
	{
		Name:   "Normalize attribute values to lowercase",
		Signal: "metrics",
		Config: "transform:\n" +
			"  error_mode: propagate\n" +
			"  metric_statements:\n" +
			"    - set(datapoint.attributes, MapEach(datapoint.attributes, (_, v) => When(() => IsString(v), ToLowerCase(String(v)), v)))",
		Payload: `{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"order-worker"}}]},"scopeMetrics":[{"scope":{"name":"worker.queue"},"metrics":[{"name":"queue.depth","unit":"{message}","gauge":{"dataPoints":[{"attributes":[{"key":"env","value":{"stringValue":"Prod"}},{"key":"region","value":{"stringValue":"EU-West-1"}},{"key":"queue","value":{"stringValue":"Orders"}}],"timeUnixNano":"1792764000000000000","asInt":"42"},{"attributes":[{"key":"env","value":{"stringValue":"PROD"}},{"key":"region","value":{"stringValue":"eu-west-1"}},{"key":"queue","value":{"stringValue":"orders"}}],"timeUnixNano":"1792764000000000000","asInt":"17"}]}}]}]}]}`,
	},
}

var transformProcessorSyntaxHighlightPatterns = []metadataOption{
	withOTTLSyntaxHighlightPatterns(OTTLGrammarStatement, []string{`^(transform(/.*)?\.)?(metric|log|trace|profile)_statements\.\*(\.statements\.\*)?$`}),
	withOTTLSyntaxHighlightPatterns(OTTLGrammarCondition, []string{`^(transform(/.*)?\.)?(metric|log|trace|profile)_statements\.\*\.conditions\.\*$`}),
}

// NewTransformProcessorExecutor creates an internal.Executor that runs OTTL statements using
// the [transformprocessor].
func NewTransformProcessorExecutor() Executor {
	debugger := NewTransformProcessorDebugger()
	factory := transformprocessor.NewFactory()
	enableTransformProcessorFeatureGates()

	return NewJSONExecutor[transformprocessor.Config](
		newProcessorConsumer[transformprocessor.Config](factory),
		newMetadata(
			ComponentTypeProcessor,
			"transform_processor",
			"Transform",
			"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/transformprocessor",
			"https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/processor/transformprocessor",
			append([]metadataOption{
				withConfigExamples(transformProcessorConfigExamples...)},
				transformProcessorSyntaxHighlightPatterns...,
			)...,
		),
		withDebugger[transformprocessor.Config](debugger),
	)
}

func enableTransformProcessorFeatureGates() {
	// Enable flatten for logs; ignore error so older/different contrib versions still run.
	_ = featuregate.GlobalRegistry().Set("transform.flatten.logs", true)
}

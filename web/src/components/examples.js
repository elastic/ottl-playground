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

export const PAYLOAD_EXAMPLES = {
  logs: '{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeLogs":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"logRecords":[{"timeUnixNano":"1544712660300000000","observedTimeUnixNano":"1544712660300000000","severityNumber":18,"severityText":"Information","traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","body":{"stringValue":"Example log record"},"attributes":[{"key":"error.type","value":{"stringValue":"Error"}}]}]}]}]}',
  traces: '{"resourceSpans":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeSpans":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"spans":[{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","parentSpanId":"eee19b7ec3c1b173","name":"Im a server span","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","kind":2,"attributes":[{"key":"my.span.attr","value":{"stringValue":"some value"}},{"key":"error.type","value":{"stringValue":"hola"}}],"status":{"code":2}},{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b173","parentSpanId":"eee19b7ec3c1b173","name":"Me too","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","kind":1,"attributes":[{"key":"my.span.attr","value":{"stringValue":"some value"}},{"key":"http.response.status_code","value":{"intValue":"500"}},{"key":"error.type","value":{"stringValue":"mundo"}}],"status":{}}]}]}]}',
  metrics: '{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}},{"key":"timestamp","value":{"stringValue":"2018-12-01T16:17:18Z"}}]},"scopeMetrics":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"metrics":[{"name":"my.counter","unit":"1","description":"I am a Counter","sum":{"aggregationTemporality":1,"isMonotonic":true,"dataPoints":[{"asDouble":5,"startTimeUnixNano":"1544712660300000000","timeUnixNano":"1544712660300000000","attributes":[{"key":"my.counter.attr","value":{"stringValue":"some value"}}]},{"asDouble":2,"startTimeUnixNano":"1544712660300000000","timeUnixNano":"1544712660300000000","attributes":[{"key":"another.counter.attr","value":{"stringValue":"another value"}}]}]}},{"name":"my.gauge","unit":"1","description":"I am a Gauge","gauge":{"dataPoints":[{"asDouble":10,"timeUnixNano":"1544712660300000000","attributes":[{"key":"my.gauge.attr","value":{"stringValue":"some value"}},{"key":"error.type","value":{"stringValue":"hello"}}]}]}},{"name":"my.histogram","unit":"1","description":"I am a Histogram","histogram":{"aggregationTemporality":1,"dataPoints":[{"startTimeUnixNano":"1544712660300000000","timeUnixNano":"1544712660300000000","count":"2","sum":2,"bucketCounts":["1","1"],"explicitBounds":[1],"min":0,"max":2,"attributes":[{"key":"my.histogram.attr","value":{"stringValue":"some value"}},{"key":"error.type","value":{"stringValue":"world"}}]}]}}]}]}]}'
 };

const TRANSFORM_PROCESSOR_CONFIG_EXAMPLES = [
  {
    name: 'Rename an attribute',
    otlp_type: 'traces',
    config:
      'error_mode: ignore\n' +
      'trace_statements:\n' +
      ' - context: resource\n' +
      '   statements:\n' +
      '    - set(attributes["service.new_name"], attributes["service.name"])\n' +
      '    - delete_key(attributes, "service.name")',
  },
  {
    name: 'Copy field to attributes',
    otlp_type: 'logs',
    config:
      'error_mode: ignore\n' +
      'log_statements:\n' +
      ' - context: log\n' +
      '   statements:\n' +
      '    - set(attributes["body"], body)',
  },
  {
    name: 'Combine two attributes',
    otlp_type: 'logs',
    config:
      'error_mode: ignore\n' +
      'log_statements:\n' +
      ' - context: log\n' +
      '   statements:\n' +
      '    - set(attributes["combined"], Concat([attributes["string.attribute"], attributes["boolean.attribute"]], " "))',
  },
  {
    name: 'Set a field',
    otlp_type: 'logs',
    config:
      'log_statements:\n' +
      ' - context: log\n' +
      '   statements:\n' +
      '    - set(severity_number, SEVERITY_NUMBER_INFO)\n' +
      '    - set(severity_text, "INFO")',
  },
  {
    name: 'Parse unstructured log',
    otlp_type: 'logs',
    config:
      'log_statements:\n' +
      ' - context: log\n' +
      '   statements:\n' +
      '    - \'merge_maps(attributes, ExtractPatterns(body, "Example (?P<example_type>[a-z\\\\.]+)"), "upsert")\'',
  },
  {
    name: 'Conditionally set a field',
    otlp_type: 'traces',
    config:
      'trace_statements:\n' +
      ' - context: span\n' +
      '   statements:\n' +
      '    - set(attributes["server"], true) where kind == SPAN_KIND_SERVER',
  },
  {
    name: 'Update a resource attribute',
    otlp_type: 'logs',
    config:
      'log_statements:\n' +
      ' - context: resource\n' +
      '   statements:\n' +
      '    - set(attributes["service.name"], "mycompany-application")',
  },
  {
    name: 'Parse and manipulate JSON',
    otlp_type: 'logs',
    config:
      'log_statements:\n' +
      ' - context: log\n' +
      '   statements:\n' +
      '    - merge_maps(cache, ParseJSON(body), "upsert") where IsMatch(body, "^\\\\{")\n' +
      '    - set(time, Time(cache["timestamp"], "%Y-%m-%dT%H:%M:%SZ"))\n' +
      '    - set(severity_text, cache["level"])\n' +
      '    - set(body, cache["message"])',
    payload:
      '{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeLogs":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"logRecords":[{"timeUnixNano":"1544712660300000000","observedTimeUnixNano":"1544712660300000000","severityNumber":10,"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","body":{"stringValue":"{\\"timestamp\\": \\"2025-03-01T12:12:14Z\\", \\"level\\":\\"INFO\\",\\"message\\":\\"Elapsed time: 10ms\\"}"}}]}]}]}',
  },
  {
    name: 'Parse and manipulate Timestamps',
    otlp_type: 'metrics',
    config:
      'metric_statements:\n' +
      ' - context: resource\n' +
      '   statements:\n' +
      '    - set(attributes["date"], String(TruncateTime(Time(attributes["timestamp"], "%Y-%m-%dT%H:%M:%SZ"), Duration("24h"))))',
  },
  {
    name: 'Manipulate strings',
    otlp_type: 'metrics',
    config:
      'metric_statements:\n' +
      ' - context: scope\n' +
      '   statements:\n' +
      '    - set(resource.attributes["service.name"], ConvertCase(Concat([resource.attributes["service.name"], version], ".v"), "upper"))',
  },
  {
    name: 'Scale a metric',
    otlp_type: 'metrics',
    config:
      'metric_statements:\n' +
      ' - context: metric\n' +
      '   statements:\n' +
      '    - scale_metric(10.0, "kWh") where name == "my.gauge"',
  },
  {
    name: 'Dynamically rename a metric',
    otlp_type: 'metrics',
    config:
      'metric_statements:\n' +
      ' - context: metric\n' +
      '   statements:\n' +
      '    - replace_pattern(name, "my.(.+)", "metrics.$$1")',
  },
  {
    name: 'Aggregate a metric',
    otlp_type: 'metrics',
    config:
      'metric_statements:\n' +
      ' - context: metric\n' +
      '   statements:\n' +
      '    - copy_metric(name="my.second.histogram") where name == "my.histogram"\n' +
      '    - aggregate_on_attributes("sum", []) where name == "my.second.histogram"',
  },
];

const FILTER_PROCESSOR_CONFIG_EXAMPLES = [
  {
    name: 'Drop specific metric and value',
    otlp_type: 'metrics',
    config:
      'metrics:\n' +
      '  datapoint:\n' +
      '    - metric.name == "my.histogram" and count == 2',
  },
  {
    name: 'Drop spans',
    otlp_type: 'traces',
    config:
      'error_mode: ignore\n' +
      'traces:\n' +
      '  span:\n' +
      '    - kind == SPAN_KIND_INTERNAL',
  },
  {
    name: 'Drop data by resource attribute',
    otlp_type: 'traces',
    config:
      'traces:\n' +
      '  span:\n' +
      '    - IsMatch(resource.attributes["service.name"], "my-*")',
  },
  {
    name: 'Drop debug and trace logs',
    otlp_type: 'logs',
    config:
      'logs:\n' +
      '  log_record:\n' +
      '    - severity_number != SEVERITY_NUMBER_UNSPECIFIED and severity_number < SEVERITY_NUMBER_INFO',
    payload:
      '{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeLogs":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"logRecords":[{"timeUnixNano":"1544712660300000000","observedTimeUnixNano":"1544712660300000000","severityNumber":10,"severityText":"Information","traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","body":{"stringValue":"I\'m an INFO log record"}},{"timeUnixNano":"1544712660300000000","observedTimeUnixNano":"1544712660300000000","severityNumber":5,"severityText":"Debug","traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","body":{"stringValue":"I\'m a DEBUG log record"}}]}]}]}',
  },
];

const ATTRIBUTES_PROCESSOR_CONFIG_EXAMPLES = [
  {
    name: 'Insert Attribute for All Except Certain Metric Name',
    otlp_type: 'metrics',
    config:
      'exclude:\n' +
      ' match_type: strict\n' +
      ' metric_names: ["my.gauge"]\n' +
      'actions:\n' +
      ' - key: my.insert.attr\n' +
      '   action: insert\n' +
      '   value: "some inserted value"',
  },
  {
    name: 'Upsert Attribute for All Logs Above Certain Severity',
    otlp_type: 'logs',
    config:
      'include:\n' +
      ' log_severity_number:\n' +
      '   min: 10\n' +
      '   match_undefined: false\n' +
      'actions:\n' +
      ' - key: my.upsert.attr\n' +
      '   action: upsert\n' +
      '   value: true',
  },
  {
    name: 'Delete Attributes Matching Regex',
    otlp_type: 'metrics',
    config:
      'actions:\n' +
      ' - pattern: \'^my\\.[a-z]+\\.attr$\'\n' +
      '   action: delete',
  },
  {
    name: 'Convert Attribute Type',
    otlp_type: 'traces',
    config:
      'actions:\n' +
      ' - key: http.response.status_code\n' +
      '   action: convert\n' +
      '   converted_type: string',
  },
]

const sortBy = (property) => {
  return function (a, b) {
    return a[property] < b[property] ? -1 : a[property] > b[property] ? 1 : 0;
  };
};

export const CONFIG_EXAMPLES = {
  transform_processor: TRANSFORM_PROCESSOR_CONFIG_EXAMPLES.sort(sortBy('name')),
  filter_processor: FILTER_PROCESSOR_CONFIG_EXAMPLES.sort(sortBy('name')),
  attributes_processor: ATTRIBUTES_PROCESSOR_CONFIG_EXAMPLES.sort(sortBy('name'))
};

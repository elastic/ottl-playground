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

export const DEFAULT_PAYLOADS = {
  logs: '{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeLogs":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"logRecords":[{"timeUnixNano":"1544712660300000000","observedTimeUnixNano":"1544712660300000000","severityNumber":10,"severityText":"Information","traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","body":{"stringValue":"Example log record"},"attributes":[{"key":"string.attribute","value":{"stringValue":"some string"}},{"key":"boolean.attribute","value":{"boolValue":true}},{"key":"int.attribute","value":{"intValue":"10"}},{"key":"double.attribute","value":{"doubleValue":637.704}},{"key":"array.attribute","value":{"arrayValue":{"values":[{"stringValue":"many"},{"stringValue":"values"}]}}},{"key":"map.attribute","value":{"kvlistValue":{"values":[{"key":"some.map.key","value":{"stringValue":"some value"}}]}}}]}]}]}]}',
  traces:
    '{"resourceSpans":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeSpans":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"spans":[{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","parentSpanId":"eee19b7ec3c1b173","name":"I\'m a server span","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","kind":2,"attributes":[{"key":"my.span.attr","value":{"stringValue":"some value"}}],"status":{}},{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b173","parentSpanId":"eee19b7ec3c1b173","name":"Me too","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","kind":1,"attributes":[{"key":"my.span.attr","value":{"stringValue":"some value"}}],"status":{}}]}]}]}',
  metrics:
    '{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}},{"key":"timestamp","value":{"stringValue":"2018-12-01T16:17:18Z"}}]},"scopeMetrics":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"metrics":[{"name":"my.counter","unit":"1","description":"I am a Counter","sum":{"aggregationTemporality":1,"isMonotonic":true,"dataPoints":[{"asDouble":5,"startTimeUnixNano":"1544712660300000000","timeUnixNano":"1544712660300000000","attributes":[{"key":"my.counter.attr","value":{"stringValue":"some value"}}]},{"asDouble":2,"startTimeUnixNano":"1544712660300000000","timeUnixNano":"1544712660300000000","attributes":[{"key":"another.counter.attr","value":{"stringValue":"another value"}}]}]}},{"name":"my.gauge","unit":"1","description":"I am a Gauge","gauge":{"dataPoints":[{"asDouble":10,"timeUnixNano":"1544712660300000000","attributes":[{"key":"my.gauge.attr","value":{"stringValue":"some value"}}]}]}},{"name":"my.histogram","unit":"1","description":"I am a Histogram","histogram":{"aggregationTemporality":1,"dataPoints":[{"startTimeUnixNano":"1544712660300000000","timeUnixNano":"1544712660300000000","count":"2","sum":2,"bucketCounts":["1","1"],"explicitBounds":[1],"min":0,"max":2,"attributes":[{"key":"my.histogram.attr","value":{"stringValue":"some value"}}]}]}}]}]}]}',
};

export const DEFAULT_PAYLOAD_EXAMPLES = [
  {
    name: 'Logs',
    signal: 'logs',
    value: DEFAULT_PAYLOADS.logs,
  },
  {
    name: 'Traces',
    signal: 'traces',
    value: DEFAULT_PAYLOADS.traces,
  },
  {
    name: 'Metrics',
    signal: 'metrics',
    value: DEFAULT_PAYLOADS.metrics,
  },
];

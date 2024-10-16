export const PAYLOAD_EXAMPLES = {
  logs: '{"resourceLogs":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeLogs":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"logRecords":[{"timeUnixNano":"1544712660300000000","observedTimeUnixNano":"1544712660300000000","severityNumber":10,"severityText":"Information","traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","body":{"stringValue":"Example log record"},"attributes":[{"key":"string.attribute","value":{"stringValue":"some string"}},{"key":"boolean.attribute","value":{"boolValue":true}},{"key":"int.attribute","value":{"intValue":"10"}},{"key":"double.attribute","value":{"doubleValue":637.704}},{"key":"array.attribute","value":{"arrayValue":{"values":[{"stringValue":"many"},{"stringValue":"values"}]}}},{"key":"map.attribute","value":{"kvlistValue":{"values":[{"key":"some.map.key","value":{"stringValue":"some value"}}]}}}]}]}]}]}',
  traces:
    '{"resourceSpans":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeSpans":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"spans":[{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b174","parentSpanId":"eee19b7ec3c1b173","name":"I\'m a server span","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","kind":2,"attributes":[{"key":"my.span.attr","value":{"stringValue":"some value"}}]},{"traceId":"5b8efff798038103d269b633813fc60c","spanId":"eee19b7ec3c1b173","parentSpanId":"eee19b7ec3c1b173","name":"Me too","startTimeUnixNano":"1544712660000000000","endTimeUnixNano":"1544712661000000000","kind":1,"attributes":[{"key":"my.span.attr","value":{"stringValue":"some value"}}]}]}]}]}',
  metrics:
    '{"resourceMetrics":[{"resource":{"attributes":[{"key":"service.name","value":{"stringValue":"my.service"}}]},"scopeMetrics":[{"scope":{"name":"my.library","version":"1.0.0","attributes":[{"key":"my.scope.attribute","value":{"stringValue":"some scope attribute"}}]},"metrics":[{"name":"my.counter","unit":"1","description":"I am a Counter","sum":{"aggregationTemporality":1,"isMonotonic":true,"dataPoints":[{"asDouble":5,"startTimeUnixNano":"1544712660300000000","timeUnixNano":"1544712660300000000","attributes":[{"key":"my.counter.attr","value":{"stringValue":"some value"}}]}]}},{"name":"my.gauge","unit":"1","description":"I am a Gauge","gauge":{"dataPoints":[{"asDouble":10,"timeUnixNano":"1544712660300000000","attributes":[{"key":"my.gauge.attr","value":{"stringValue":"some value"}}]}]}},{"name":"my.histogram","unit":"1","description":"I am a Histogram","histogram":{"aggregationTemporality":1,"dataPoints":[{"startTimeUnixNano":"1544712660300000000","timeUnixNano":"1544712660300000000","count":2,"sum":2,"bucketCounts":[1,1],"explicitBounds":[1],"min":0,"max":2,"attributes":[{"key":"my.histogram.attr","value":{"stringValue":"some value"}}]}]}}]}]}]}',
};

const TRANSFORM_PROCESSOR_CONFIG_EXAMPLES = [
  {
    name: 'Rename attribute',
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
    name: 'Move field to attribute',
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
      '    # Use Concat function to combine any number of string, separated by a delimiter. \n' +
      '    - set(attributes["combined"], Concat([attributes["string.attribute"], attributes["boolean.attribute"]], " "))',
  },
];

const FILTER_PROCESSOR_CONFIG_EXAMPLES = [
  {
    name: 'Dropping specific metric and value',
    otlp_type: 'metrics',
    config:
      'error_mode: ignore\n' +
      'metrics:\n' +
      '  datapoint:\n' +
      '    - metric.name == "my.histogram" and count == 2',
  },
  {
    name: 'Dropping spans',
    otlp_type: 'traces',
    config:
      'error_mode: ignore\n' + 'traces:\n' + '  span:\n' + '    - kind == 1',
  },
  {
    name: 'Dropping data by resource attribute',
    otlp_type: 'traces',
    config:
      'error_mode: ignore\n' +
      'traces:\n' +
      '  span:\n' +
      '    - IsMatch(resource.attributes["service.name"], "my-*")',
  },
];

export const CONFIG_EXAMPLES = {
  transform_processor: TRANSFORM_PROCESSOR_CONFIG_EXAMPLES,
  filter_processor: FILTER_PROCESSOR_CONFIG_EXAMPLES,
};

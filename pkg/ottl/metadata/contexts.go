// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package metadata

// GetContextEnums returns enums available for a specific context.
// These are derived from the OTTL symbol tables.
func GetContextEnums(ctx ContextType) []EnumInfo {
	switch ctx {
	case ContextLog:
		return logEnums()
	case ContextSpan, ContextSpanEvent:
		return spanEnums()
	case ContextMetric:
		return metricEnums()
	case ContextDataPoint:
		return dataPointEnums()
	default:
		// Resource, Scope, Profile don't have context-specific enums
		return nil
	}
}

// logEnums returns severity enums for the log context.
func logEnums() []EnumInfo {
	return []EnumInfo{
		{Name: "SEVERITY_NUMBER_UNSPECIFIED", Value: 0},
		{Name: "SEVERITY_NUMBER_TRACE", Value: 1},
		{Name: "SEVERITY_NUMBER_TRACE2", Value: 2},
		{Name: "SEVERITY_NUMBER_TRACE3", Value: 3},
		{Name: "SEVERITY_NUMBER_TRACE4", Value: 4},
		{Name: "SEVERITY_NUMBER_DEBUG", Value: 5},
		{Name: "SEVERITY_NUMBER_DEBUG2", Value: 6},
		{Name: "SEVERITY_NUMBER_DEBUG3", Value: 7},
		{Name: "SEVERITY_NUMBER_DEBUG4", Value: 8},
		{Name: "SEVERITY_NUMBER_INFO", Value: 9},
		{Name: "SEVERITY_NUMBER_INFO2", Value: 10},
		{Name: "SEVERITY_NUMBER_INFO3", Value: 11},
		{Name: "SEVERITY_NUMBER_INFO4", Value: 12},
		{Name: "SEVERITY_NUMBER_WARN", Value: 13},
		{Name: "SEVERITY_NUMBER_WARN2", Value: 14},
		{Name: "SEVERITY_NUMBER_WARN3", Value: 15},
		{Name: "SEVERITY_NUMBER_WARN4", Value: 16},
		{Name: "SEVERITY_NUMBER_ERROR", Value: 17},
		{Name: "SEVERITY_NUMBER_ERROR2", Value: 18},
		{Name: "SEVERITY_NUMBER_ERROR3", Value: 19},
		{Name: "SEVERITY_NUMBER_ERROR4", Value: 20},
		{Name: "SEVERITY_NUMBER_FATAL", Value: 21},
		{Name: "SEVERITY_NUMBER_FATAL2", Value: 22},
		{Name: "SEVERITY_NUMBER_FATAL3", Value: 23},
		{Name: "SEVERITY_NUMBER_FATAL4", Value: 24},
	}
}

// spanEnums returns span kind and status code enums for span contexts.
func spanEnums() []EnumInfo {
	return []EnumInfo{
		{Name: "SPAN_KIND_UNSPECIFIED", Value: 0},
		{Name: "SPAN_KIND_INTERNAL", Value: 1},
		{Name: "SPAN_KIND_SERVER", Value: 2},
		{Name: "SPAN_KIND_CLIENT", Value: 3},
		{Name: "SPAN_KIND_PRODUCER", Value: 4},
		{Name: "SPAN_KIND_CONSUMER", Value: 5},
		{Name: "STATUS_CODE_UNSET", Value: 0},
		{Name: "STATUS_CODE_OK", Value: 1},
		{Name: "STATUS_CODE_ERROR", Value: 2},
	}
}

// metricEnums returns metric type and aggregation temporality enums.
func metricEnums() []EnumInfo {
	return []EnumInfo{
		// Metric types
		{Name: "METRIC_DATA_TYPE_NONE", Value: 0},
		{Name: "METRIC_DATA_TYPE_GAUGE", Value: 1},
		{Name: "METRIC_DATA_TYPE_SUM", Value: 2},
		{Name: "METRIC_DATA_TYPE_HISTOGRAM", Value: 3},
		{Name: "METRIC_DATA_TYPE_EXPONENTIAL_HISTOGRAM", Value: 4},
		{Name: "METRIC_DATA_TYPE_SUMMARY", Value: 5},
		// Aggregation temporality
		{Name: "AGGREGATION_TEMPORALITY_UNSPECIFIED", Value: 0},
		{Name: "AGGREGATION_TEMPORALITY_DELTA", Value: 1},
		{Name: "AGGREGATION_TEMPORALITY_CUMULATIVE", Value: 2},
	}
}

// dataPointEnums returns enums for the data point context.
func dataPointEnums() []EnumInfo {
	return []EnumInfo{
		// Data point flags
		{Name: "FLAG_NONE", Value: 0},
		{Name: "FLAG_NO_RECORDED_VALUE", Value: 1},
	}
}

// GetContextPaths returns paths available for a specific context.
// These are derived from the OTTL documentation and context implementations.
func GetContextPaths(ctx ContextType) []PathInfo {
	switch ctx {
	case ContextLog:
		return logPaths()
	case ContextSpan:
		return spanPaths()
	case ContextSpanEvent:
		return spanEventPaths()
	case ContextMetric:
		return metricPaths()
	case ContextDataPoint:
		return dataPointPaths()
	case ContextResource:
		return resourcePaths()
	case ContextScope:
		return scopePaths()
	case ContextProfile:
		return profilePaths()
	default:
		return nil
	}
}

// logPaths returns paths available in the log context.
func logPaths() []PathInfo {
	return []PathInfo{
		{Path: "time_unix_nano", Type: "int64", Description: "Timestamp as Unix nanoseconds"},
		{Path: "observed_time_unix_nano", Type: "int64", Description: "Observed timestamp as Unix nanoseconds"},
		{Path: "time", Type: "time.Time", Description: "Timestamp as time.Time"},
		{Path: "observed_time", Type: "time.Time", Description: "Observed timestamp as time.Time"},
		{Path: "severity_number", Type: "int64", Description: "Severity number (use SEVERITY_NUMBER_* enums)"},
		{Path: "severity_text", Type: "string", Description: "Severity text"},
		{Path: "body", Type: "any", Description: "Log body (can be string, map, or slice)", SupportsKeys: true},
		{Path: "body.string", Type: "string", Description: "Log body as string"},
		{Path: "attributes", Type: "pcommon.Map", Description: "Log attributes", SupportsKeys: true},
		{Path: "dropped_attributes_count", Type: "int64", Description: "Number of dropped attributes"},
		{Path: "flags", Type: "int64", Description: "Log flags"},
		{Path: "trace_id", Type: "pcommon.TraceID", Description: "Trace ID"},
		{Path: "trace_id.string", Type: "string", Description: "Trace ID as hex string"},
		{Path: "span_id", Type: "pcommon.SpanID", Description: "Span ID"},
		{Path: "span_id.string", Type: "string", Description: "Span ID as hex string"},
		{Path: "event_name", Type: "string", Description: "Event name"},
		// Resource and scope paths (inherited)
		{Path: "resource", Type: "pcommon.Resource", Description: "Resource associated with the log"},
		{Path: "resource.attributes", Type: "pcommon.Map", Description: "Resource attributes", SupportsKeys: true},
		{Path: "resource.dropped_attributes_count", Type: "int64", Description: "Resource dropped attributes count"},
		{Path: "resource.schema_url", Type: "string", Description: "Resource schema URL"},
		{Path: "instrumentation_scope", Type: "pcommon.InstrumentationScope", Description: "Instrumentation scope"},
		{Path: "instrumentation_scope.name", Type: "string", Description: "Scope name"},
		{Path: "instrumentation_scope.version", Type: "string", Description: "Scope version"},
		{Path: "instrumentation_scope.attributes", Type: "pcommon.Map", Description: "Scope attributes", SupportsKeys: true},
		{Path: "instrumentation_scope.dropped_attributes_count", Type: "int64", Description: "Scope dropped attributes count"},
		{Path: "instrumentation_scope.schema_url", Type: "string", Description: "Scope schema URL"},
		{Path: "cache", Type: "pcommon.Map", Description: "Temporary cache for the log record", SupportsKeys: true},
	}
}

// spanPaths returns paths available in the span context.
func spanPaths() []PathInfo {
	return []PathInfo{
		{Path: "trace_id", Type: "pcommon.TraceID", Description: "Trace ID"},
		{Path: "trace_id.string", Type: "string", Description: "Trace ID as hex string"},
		{Path: "span_id", Type: "pcommon.SpanID", Description: "Span ID"},
		{Path: "span_id.string", Type: "string", Description: "Span ID as hex string"},
		{Path: "parent_span_id", Type: "pcommon.SpanID", Description: "Parent span ID"},
		{Path: "parent_span_id.string", Type: "string", Description: "Parent span ID as hex string"},
		{Path: "trace_state", Type: "string", Description: "Trace state"},
		{Path: "name", Type: "string", Description: "Span name"},
		{Path: "kind", Type: "int64", Description: "Span kind (use SPAN_KIND_* enums)"},
		{Path: "start_time_unix_nano", Type: "int64", Description: "Start time as Unix nanoseconds"},
		{Path: "end_time_unix_nano", Type: "int64", Description: "End time as Unix nanoseconds"},
		{Path: "start_time", Type: "time.Time", Description: "Start time as time.Time"},
		{Path: "end_time", Type: "time.Time", Description: "End time as time.Time"},
		{Path: "attributes", Type: "pcommon.Map", Description: "Span attributes", SupportsKeys: true},
		{Path: "dropped_attributes_count", Type: "int64", Description: "Number of dropped attributes"},
		{Path: "events", Type: "ptrace.SpanEventSlice", Description: "Span events"},
		{Path: "dropped_events_count", Type: "int64", Description: "Number of dropped events"},
		{Path: "links", Type: "ptrace.SpanLinkSlice", Description: "Span links"},
		{Path: "dropped_links_count", Type: "int64", Description: "Number of dropped links"},
		{Path: "status", Type: "ptrace.Status", Description: "Span status"},
		{Path: "status.code", Type: "int64", Description: "Status code (use STATUS_CODE_* enums)"},
		{Path: "status.message", Type: "string", Description: "Status message"},
		{Path: "flags", Type: "int64", Description: "Span flags"},
		// Resource and scope paths (inherited)
		{Path: "resource", Type: "pcommon.Resource", Description: "Resource associated with the span"},
		{Path: "resource.attributes", Type: "pcommon.Map", Description: "Resource attributes", SupportsKeys: true},
		{Path: "resource.dropped_attributes_count", Type: "int64", Description: "Resource dropped attributes count"},
		{Path: "resource.schema_url", Type: "string", Description: "Resource schema URL"},
		{Path: "instrumentation_scope", Type: "pcommon.InstrumentationScope", Description: "Instrumentation scope"},
		{Path: "instrumentation_scope.name", Type: "string", Description: "Scope name"},
		{Path: "instrumentation_scope.version", Type: "string", Description: "Scope version"},
		{Path: "instrumentation_scope.attributes", Type: "pcommon.Map", Description: "Scope attributes", SupportsKeys: true},
		{Path: "instrumentation_scope.dropped_attributes_count", Type: "int64", Description: "Scope dropped attributes count"},
		{Path: "instrumentation_scope.schema_url", Type: "string", Description: "Scope schema URL"},
		{Path: "cache", Type: "pcommon.Map", Description: "Temporary cache for the span", SupportsKeys: true},
	}
}

// spanEventPaths returns paths available in the span event context.
func spanEventPaths() []PathInfo {
	return []PathInfo{
		{Path: "time_unix_nano", Type: "int64", Description: "Event time as Unix nanoseconds"},
		{Path: "time", Type: "time.Time", Description: "Event time as time.Time"},
		{Path: "name", Type: "string", Description: "Event name"},
		{Path: "attributes", Type: "pcommon.Map", Description: "Event attributes", SupportsKeys: true},
		{Path: "dropped_attributes_count", Type: "int64", Description: "Number of dropped attributes"},
		// Span fields accessible from span event context
		{Path: "span.trace_id", Type: "pcommon.TraceID", Description: "Parent span trace ID"},
		{Path: "span.span_id", Type: "pcommon.SpanID", Description: "Parent span ID"},
		{Path: "span.name", Type: "string", Description: "Parent span name"},
		// Resource and scope paths (inherited)
		{Path: "resource", Type: "pcommon.Resource", Description: "Resource associated with the span event"},
		{Path: "resource.attributes", Type: "pcommon.Map", Description: "Resource attributes", SupportsKeys: true},
		{Path: "instrumentation_scope", Type: "pcommon.InstrumentationScope", Description: "Instrumentation scope"},
		{Path: "instrumentation_scope.name", Type: "string", Description: "Scope name"},
		{Path: "instrumentation_scope.version", Type: "string", Description: "Scope version"},
		{Path: "cache", Type: "pcommon.Map", Description: "Temporary cache", SupportsKeys: true},
	}
}

// metricPaths returns paths available in the metric context.
func metricPaths() []PathInfo {
	return []PathInfo{
		{Path: "name", Type: "string", Description: "Metric name"},
		{Path: "description", Type: "string", Description: "Metric description"},
		{Path: "unit", Type: "string", Description: "Metric unit"},
		{Path: "type", Type: "int64", Description: "Metric type"},
		{Path: "aggregation_temporality", Type: "int64", Description: "Aggregation temporality"},
		{Path: "is_monotonic", Type: "bool", Description: "Whether the metric is monotonic"},
		// Resource and scope paths (inherited)
		{Path: "resource", Type: "pcommon.Resource", Description: "Resource associated with the metric"},
		{Path: "resource.attributes", Type: "pcommon.Map", Description: "Resource attributes", SupportsKeys: true},
		{Path: "resource.dropped_attributes_count", Type: "int64", Description: "Resource dropped attributes count"},
		{Path: "resource.schema_url", Type: "string", Description: "Resource schema URL"},
		{Path: "instrumentation_scope", Type: "pcommon.InstrumentationScope", Description: "Instrumentation scope"},
		{Path: "instrumentation_scope.name", Type: "string", Description: "Scope name"},
		{Path: "instrumentation_scope.version", Type: "string", Description: "Scope version"},
		{Path: "instrumentation_scope.attributes", Type: "pcommon.Map", Description: "Scope attributes", SupportsKeys: true},
		{Path: "cache", Type: "pcommon.Map", Description: "Temporary cache", SupportsKeys: true},
	}
}

// dataPointPaths returns paths available in the data point context.
func dataPointPaths() []PathInfo {
	return []PathInfo{
		{Path: "attributes", Type: "pcommon.Map", Description: "Data point attributes", SupportsKeys: true},
		{Path: "start_time_unix_nano", Type: "int64", Description: "Start time as Unix nanoseconds"},
		{Path: "time_unix_nano", Type: "int64", Description: "Time as Unix nanoseconds"},
		{Path: "start_time", Type: "time.Time", Description: "Start time as time.Time"},
		{Path: "time", Type: "time.Time", Description: "Time as time.Time"},
		{Path: "value_double", Type: "float64", Description: "Value as double (for gauge/sum)"},
		{Path: "value_int", Type: "int64", Description: "Value as int (for gauge/sum)"},
		{Path: "exemplars", Type: "pmetric.ExemplarSlice", Description: "Exemplars"},
		{Path: "flags", Type: "int64", Description: "Data point flags"},
		{Path: "count", Type: "int64", Description: "Count (for histogram/summary)"},
		{Path: "sum", Type: "float64", Description: "Sum (for histogram/summary)"},
		{Path: "min", Type: "float64", Description: "Min (for histogram)"},
		{Path: "max", Type: "float64", Description: "Max (for histogram)"},
		{Path: "bucket_counts", Type: "[]uint64", Description: "Bucket counts (for histogram)"},
		{Path: "explicit_bounds", Type: "[]float64", Description: "Explicit bounds (for histogram)"},
		{Path: "scale", Type: "int64", Description: "Scale (for exponential histogram)"},
		{Path: "zero_count", Type: "int64", Description: "Zero count (for exponential histogram)"},
		{Path: "positive", Type: "pmetric.Buckets", Description: "Positive buckets (for exponential histogram)"},
		{Path: "negative", Type: "pmetric.Buckets", Description: "Negative buckets (for exponential histogram)"},
		{Path: "quantile_values", Type: "pmetric.ValueAtQuantileSlice", Description: "Quantile values (for summary)"},
		// Metric paths (inherited from metric context)
		{Path: "metric", Type: "pmetric.Metric", Description: "Parent metric"},
		{Path: "metric.name", Type: "string", Description: "Metric name"},
		{Path: "metric.description", Type: "string", Description: "Metric description"},
		{Path: "metric.unit", Type: "string", Description: "Metric unit"},
		// Resource and scope paths (inherited)
		{Path: "resource", Type: "pcommon.Resource", Description: "Resource associated with the data point"},
		{Path: "resource.attributes", Type: "pcommon.Map", Description: "Resource attributes", SupportsKeys: true},
		{Path: "instrumentation_scope", Type: "pcommon.InstrumentationScope", Description: "Instrumentation scope"},
		{Path: "instrumentation_scope.name", Type: "string", Description: "Scope name"},
		{Path: "instrumentation_scope.version", Type: "string", Description: "Scope version"},
		{Path: "cache", Type: "pcommon.Map", Description: "Temporary cache", SupportsKeys: true},
	}
}

// resourcePaths returns paths available in the resource context.
func resourcePaths() []PathInfo {
	return []PathInfo{
		{Path: "attributes", Type: "pcommon.Map", Description: "Resource attributes", SupportsKeys: true},
		{Path: "dropped_attributes_count", Type: "int64", Description: "Number of dropped attributes"},
		{Path: "schema_url", Type: "string", Description: "Schema URL"},
		{Path: "cache", Type: "pcommon.Map", Description: "Temporary cache", SupportsKeys: true},
	}
}

// scopePaths returns paths available in the scope context.
func scopePaths() []PathInfo {
	return []PathInfo{
		{Path: "name", Type: "string", Description: "Scope name"},
		{Path: "version", Type: "string", Description: "Scope version"},
		{Path: "attributes", Type: "pcommon.Map", Description: "Scope attributes", SupportsKeys: true},
		{Path: "dropped_attributes_count", Type: "int64", Description: "Number of dropped attributes"},
		{Path: "schema_url", Type: "string", Description: "Schema URL"},
		{Path: "cache", Type: "pcommon.Map", Description: "Temporary cache", SupportsKeys: true},
	}
}

// profilePaths returns paths available in the profile context.
func profilePaths() []PathInfo {
	return []PathInfo{
		{Path: "profile_id", Type: "pprofile.ProfileID", Description: "Profile ID"},
		{Path: "profile_id.string", Type: "string", Description: "Profile ID as hex string"},
		{Path: "start_time_unix_nano", Type: "int64", Description: "Start time as Unix nanoseconds"},
		{Path: "end_time_unix_nano", Type: "int64", Description: "End time as Unix nanoseconds"},
		{Path: "start_time", Type: "time.Time", Description: "Start time as time.Time"},
		{Path: "end_time", Type: "time.Time", Description: "End time as time.Time"},
		{Path: "attributes", Type: "pcommon.Map", Description: "Profile attributes", SupportsKeys: true},
		{Path: "dropped_attributes_count", Type: "int64", Description: "Number of dropped attributes"},
		{Path: "duration", Type: "int64", Description: "Profile duration"},
		{Path: "period", Type: "int64", Description: "Profile period"},
		{Path: "default_sample_type", Type: "int64", Description: "Default sample type"},
		// Resource and scope paths (inherited)
		{Path: "resource", Type: "pcommon.Resource", Description: "Resource associated with the profile"},
		{Path: "resource.attributes", Type: "pcommon.Map", Description: "Resource attributes", SupportsKeys: true},
		{Path: "instrumentation_scope", Type: "pcommon.InstrumentationScope", Description: "Instrumentation scope"},
		{Path: "instrumentation_scope.name", Type: "string", Description: "Scope name"},
		{Path: "instrumentation_scope.version", Type: "string", Description: "Scope version"},
		{Path: "cache", Type: "pcommon.Map", Description: "Temporary cache", SupportsKeys: true},
	}
}

// GetAllContextMetadata returns metadata for all contexts.
func GetAllContextMetadata(version string) *OTTLMetadata {
	metadata := NewOTTLMetadata(version)

	for _, ctx := range AllContextTypes() {
		metadata.Contexts[ctx] = &ContextMetadata{
			Name:  ctx,
			Paths: GetContextPaths(ctx),
			Enums: GetContextEnums(ctx),
		}
	}

	return metadata
}

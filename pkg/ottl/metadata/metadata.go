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

// Package metadata provides introspection capabilities for OTTL functions,
// contexts, paths, and enums. This package is designed to be upstream-ready
// with no playground-specific dependencies.
package metadata

// ParameterKind represents the semantic type of an OTTL function parameter.
type ParameterKind string

const (
	// Getter types - read values
	KindGetter           ParameterKind = "Getter"
	KindStringGetter     ParameterKind = "StringGetter"
	KindStringLikeGetter ParameterKind = "StringLikeGetter"
	KindIntGetter        ParameterKind = "IntGetter"
	KindIntLikeGetter    ParameterKind = "IntLikeGetter"
	KindFloatGetter      ParameterKind = "FloatGetter"
	KindFloatLikeGetter  ParameterKind = "FloatLikeGetter"
	KindBoolGetter       ParameterKind = "BoolGetter"
	KindBoolLikeGetter   ParameterKind = "BoolLikeGetter"
	KindPMapGetter       ParameterKind = "PMapGetter"
	KindPSliceGetter     ParameterKind = "PSliceGetter"
	KindDurationGetter   ParameterKind = "DurationGetter"
	KindTimeGetter       ParameterKind = "TimeGetter"
	KindByteSliceGetter  ParameterKind = "ByteSliceLikeGetter"
	KindFunctionGetter   ParameterKind = "FunctionGetter"

	// Setter types - write values
	KindSetter ParameterKind = "Setter"

	// GetSetter types - read and write values
	KindGetSetter     ParameterKind = "GetSetter"
	KindPMapGetSetter ParameterKind = "PMapGetSetter"

	// Primitive types
	KindString  ParameterKind = "string"
	KindFloat64 ParameterKind = "float64"
	KindInt64   ParameterKind = "int64"
	KindBool    ParameterKind = "bool"
	KindBytes   ParameterKind = "[]byte"

	// Special types
	KindEnum    ParameterKind = "Enum"
	KindUnknown ParameterKind = "unknown"
)

// ParameterInfo describes a function parameter.
type ParameterInfo struct {
	// Name is the parameter name in snake_case (as used in OTTL syntax).
	Name string `json:"name"`

	// Kind indicates the semantic type of the parameter.
	Kind ParameterKind `json:"kind"`

	// IsOptional indicates whether this parameter can be omitted.
	IsOptional bool `json:"optional,omitempty"`

	// IsSlice indicates whether this parameter accepts multiple values.
	IsSlice bool `json:"isSlice,omitempty"`

	// Description provides documentation for the parameter.
	Description string `json:"description,omitempty"`
}

// FunctionInfo describes an OTTL function.
type FunctionInfo struct {
	// Name is the function name as used in OTTL statements.
	Name string `json:"name"`

	// Parameters describes the function's parameters in order.
	Parameters []ParameterInfo `json:"parameters"`

	// IsEditor indicates whether this is an editor function (modifies data)
	// as opposed to a converter function (returns a value).
	IsEditor bool `json:"isEditor"`

	// Description provides documentation for the function.
	Description string `json:"description,omitempty"`

	// ReturnType describes what the function returns (for converters).
	ReturnType string `json:"returnType,omitempty"`
}

// PathInfo describes a context path (like "body", "attributes", etc.).
type PathInfo struct {
	// Path is the full path expression (e.g., "attributes", "body.string").
	Path string `json:"path"`

	// Type describes the Go type of this path.
	Type string `json:"type"`

	// Description provides documentation for the path.
	Description string `json:"description,omitempty"`

	// SupportsKeys indicates whether this path supports key access (e.g., attributes["key"]).
	SupportsKeys bool `json:"supportsKeys,omitempty"`

	// IsReadOnly indicates whether this path can only be read, not written.
	IsReadOnly bool `json:"readOnly,omitempty"`
}

// EnumInfo describes an OTTL enum value.
type EnumInfo struct {
	// Name is the enum constant name (e.g., "SEVERITY_NUMBER_INFO").
	Name string `json:"name"`

	// Value is the numeric value of the enum.
	Value int64 `json:"value"`

	// Description provides documentation for the enum.
	Description string `json:"description,omitempty"`
}

// ContextType represents an OTTL context type.
type ContextType string

const (
	ContextLog       ContextType = "log"
	ContextSpan      ContextType = "span"
	ContextSpanEvent ContextType = "spanevent"
	ContextMetric    ContextType = "metric"
	ContextDataPoint ContextType = "datapoint"
	ContextResource  ContextType = "resource"
	ContextScope     ContextType = "scope"
	ContextProfile   ContextType = "profile"
)

// AllContextTypes returns all supported context types.
func AllContextTypes() []ContextType {
	return []ContextType{
		ContextLog,
		ContextSpan,
		ContextSpanEvent,
		ContextMetric,
		ContextDataPoint,
		ContextResource,
		ContextScope,
		ContextProfile,
	}
}

// ContextMetadata represents all metadata for a specific OTTL context.
type ContextMetadata struct {
	// Name is the context type name.
	Name ContextType `json:"name"`

	// Functions available in this context.
	Functions []FunctionInfo `json:"functions"`

	// Paths available in this context.
	Paths []PathInfo `json:"paths"`

	// Enums available in this context.
	Enums []EnumInfo `json:"enums"`
}

// OTTLMetadata represents complete OTTL metadata for a version.
type OTTLMetadata struct {
	// Version is the collector-contrib version this metadata was extracted from.
	Version string `json:"version"`

	// Contexts contains metadata for each context type.
	Contexts map[ContextType]*ContextMetadata `json:"contexts"`
}

// NewOTTLMetadata creates a new OTTLMetadata instance.
func NewOTTLMetadata(version string) *OTTLMetadata {
	return &OTTLMetadata{
		Version:  version,
		Contexts: make(map[ContextType]*ContextMetadata),
	}
}

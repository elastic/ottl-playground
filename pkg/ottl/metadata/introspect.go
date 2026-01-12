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

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
)

// ExtractFunctionInfo extracts metadata from an OTTL factory using reflection.
func ExtractFunctionInfo[K any](factory ottl.Factory[K]) FunctionInfo {
	info := FunctionInfo{
		Name:       factory.Name(),
		Parameters: []ParameterInfo{},
	}

	args := factory.CreateDefaultArguments()
	if args == nil {
		return info
	}

	val := reflect.ValueOf(args)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return info
	}

	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		param := extractParameterInfo(field)
		info.Parameters = append(info.Parameters, param)
	}

	// Determine if editor (has side effects) vs converter (pure function)
	// Editors typically have Setter/GetSetter as first param
	info.IsEditor = determineIsEditor(info.Parameters)

	return info
}

// extractParameterInfo extracts parameter metadata from a struct field.
func extractParameterInfo(field reflect.StructField) ParameterInfo {
	typeName := field.Type.Name()
	typeStr := field.Type.String()

	// Check if it's a slice first
	isSlice := field.Type.Kind() == reflect.Slice
	if isSlice {
		// Get the element type for slices
		elemType := field.Type.Elem()
		typeName = elemType.Name()
		typeStr = elemType.String()
	}

	// Check if it's optional
	isOptional := strings.Contains(typeStr, "Optional")

	// Determine the parameter kind
	kind := determineKind(typeName, typeStr)

	return ParameterInfo{
		Name:       toSnakeCase(field.Name),
		Kind:       kind,
		IsOptional: isOptional,
		IsSlice:    isSlice,
	}
}

// determineKind determines the ParameterKind from type information.
func determineKind(typeName, typeStr string) ParameterKind {
	// Handle Optional wrapper - extract inner type
	if strings.Contains(typeStr, "Optional") {
		// Extract inner type from Optional[T]
		inner := extractOptionalInner(typeStr)
		return determineKindFromString(inner)
	}

	return determineKindFromString(typeStr)
}

// determineKindFromString maps a type string to a ParameterKind.
func determineKindFromString(typeStr string) ParameterKind {
	// Order matters - check more specific types first
	switch {
	case strings.Contains(typeStr, "PMapGetSetter"):
		return KindPMapGetSetter
	case strings.Contains(typeStr, "GetSetter"):
		return KindGetSetter
	case strings.Contains(typeStr, "Setter"):
		return KindSetter
	case strings.Contains(typeStr, "StringLikeGetter"):
		return KindStringLikeGetter
	case strings.Contains(typeStr, "StringGetter"):
		return KindStringGetter
	case strings.Contains(typeStr, "IntLikeGetter"):
		return KindIntLikeGetter
	case strings.Contains(typeStr, "IntGetter"):
		return KindIntGetter
	case strings.Contains(typeStr, "FloatLikeGetter"):
		return KindFloatLikeGetter
	case strings.Contains(typeStr, "FloatGetter"):
		return KindFloatGetter
	case strings.Contains(typeStr, "BoolLikeGetter"):
		return KindBoolLikeGetter
	case strings.Contains(typeStr, "BoolGetter"):
		return KindBoolGetter
	case strings.Contains(typeStr, "PMapGetter"):
		return KindPMapGetter
	case strings.Contains(typeStr, "PSliceGetter"):
		return KindPSliceGetter
	case strings.Contains(typeStr, "DurationGetter"):
		return KindDurationGetter
	case strings.Contains(typeStr, "TimeGetter"):
		return KindTimeGetter
	case strings.Contains(typeStr, "ByteSliceLikeGetter"):
		return KindByteSliceGetter
	case strings.Contains(typeStr, "FunctionGetter"):
		return KindFunctionGetter
	case strings.Contains(typeStr, "Getter"):
		return KindGetter
	case strings.Contains(typeStr, "Enum"):
		return KindEnum
	case typeStr == "string":
		return KindString
	case typeStr == "float64":
		return KindFloat64
	case typeStr == "int64":
		return KindInt64
	case typeStr == "int":
		return KindInt64
	case typeStr == "bool":
		return KindBool
	case strings.Contains(typeStr, "[]byte"):
		return KindBytes
	default:
		return KindUnknown
	}
}

// extractOptionalInner extracts the inner type from an Optional[T] type string.
func extractOptionalInner(typeStr string) string {
	// Handle patterns like "ottl.Optional[ottl.Getter[K]]"
	start := strings.Index(typeStr, "Optional[")
	if start == -1 {
		return typeStr
	}

	// Find the matching closing bracket
	start += len("Optional[")
	depth := 1
	end := start
	for i := start; i < len(typeStr) && depth > 0; i++ {
		switch typeStr[i] {
		case '[':
			depth++
		case ']':
			depth--
		}
		if depth > 0 {
			end = i + 1
		}
	}

	if end > start {
		return typeStr[start:end]
	}
	return typeStr
}

// determineIsEditor determines if a function is an editor based on its parameters.
// Editors typically have Setter or GetSetter as their first parameter.
func determineIsEditor(params []ParameterInfo) bool {
	if len(params) == 0 {
		return false
	}

	firstKind := params[0].Kind
	switch firstKind {
	case KindSetter, KindGetSetter, KindPMapGetSetter:
		return true
	default:
		return false
	}
}

// toSnakeCase converts a CamelCase string to snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ExtractFunctionsFromMap extracts metadata from a factory map.
func ExtractFunctionsFromMap[K any](factories map[string]ottl.Factory[K]) []FunctionInfo {
	functions := make([]FunctionInfo, 0, len(factories))
	for _, factory := range factories {
		info := ExtractFunctionInfo(factory)
		functions = append(functions, info)
	}
	return functions
}

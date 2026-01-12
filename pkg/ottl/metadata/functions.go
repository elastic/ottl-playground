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
	"sort"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
)

// ExtractStandardFunctions extracts metadata for all standard OTTL functions.
// The type parameter K is the transform context type.
func ExtractStandardFunctions[K any]() []FunctionInfo {
	factories := ottlfuncs.StandardFuncs[K]()
	return ExtractAndSortFunctions(factories)
}

// ExtractStandardConverters extracts metadata for all standard OTTL converters.
func ExtractStandardConverters[K any]() []FunctionInfo {
	factories := ottlfuncs.StandardConverters[K]()
	return ExtractAndSortFunctions(factories)
}

// ExtractAndSortFunctions extracts and sorts function info from a factory map.
func ExtractAndSortFunctions[K any](factories map[string]ottl.Factory[K]) []FunctionInfo {
	functions := ExtractFunctionsFromMap(factories)
	SortFunctions(functions)
	return functions
}

// SortFunctions sorts a slice of FunctionInfo alphabetically by name.
func SortFunctions(functions []FunctionInfo) {
	sort.Slice(functions, func(i, j int) bool {
		return functions[i].Name < functions[j].Name
	})
}

// FilterEditors returns only editor functions from a slice.
func FilterEditors(functions []FunctionInfo) []FunctionInfo {
	var editors []FunctionInfo
	for _, f := range functions {
		if f.IsEditor {
			editors = append(editors, f)
		}
	}
	return editors
}

// FilterConverters returns only converter functions from a slice.
func FilterConverters(functions []FunctionInfo) []FunctionInfo {
	var converters []FunctionInfo
	for _, f := range functions {
		if !f.IsEditor {
			converters = append(converters, f)
		}
	}
	return converters
}

// FunctionsByName creates a map of functions keyed by name for quick lookup.
func FunctionsByName(functions []FunctionInfo) map[string]FunctionInfo {
	m := make(map[string]FunctionInfo, len(functions))
	for _, f := range functions {
		m[f.Name] = f
	}
	return m
}

// GetFunctionSignature returns a human-readable signature for a function.
func GetFunctionSignature(f FunctionInfo) string {
	if len(f.Parameters) == 0 {
		return f.Name + "()"
	}

	var params []string
	for _, p := range f.Parameters {
		param := p.Name
		if p.IsOptional {
			param += "?"
		}
		params = append(params, param)
	}

	result := f.Name + "("
	for i, p := range params {
		if i > 0 {
			result += ", "
		}
		result += p
	}
	result += ")"

	return result
}

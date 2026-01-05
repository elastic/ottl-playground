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
	"github.com/elastic/ottl-playground/internal"
	"github.com/elastic/ottl-playground/pkg/ottl/metadata"
)

// functionsCache caches the extracted function metadata.
var functionsCache []metadata.FunctionInfo

// GetOTTLFunctions returns all available OTTL functions with their metadata.
// Functions are extracted via reflection from the OTTL StandardFuncs registry.
func GetOTTLFunctions() []any {
	if functionsCache == nil {
		functionsCache = metadata.ExtractStandardFunctions[any]()
	}

	result := make([]any, len(functionsCache))
	for i, f := range functionsCache {
		params := make([]any, len(f.Parameters))
		for j, p := range f.Parameters {
			params[j] = map[string]any{
				"name":     p.Name,
				"kind":     string(p.Kind),
				"optional": p.IsOptional,
				"isSlice":  p.IsSlice,
			}
		}
		result[i] = map[string]any{
			"name":       f.Name,
			"parameters": params,
			"isEditor":   f.IsEditor,
		}
	}
	return result
}

// GetContextPaths returns paths available for a specific OTTL context.
func GetContextPaths(context string) []any {
	ctx := metadata.ContextType(context)
	paths := metadata.GetContextPaths(ctx)

	result := make([]any, len(paths))
	for i, p := range paths {
		result[i] = map[string]any{
			"path":         p.Path,
			"type":         p.Type,
			"description":  p.Description,
			"supportsKeys": p.SupportsKeys,
		}
	}
	return result
}

// GetContextEnums returns enums available for a specific OTTL context.
func GetContextEnums(context string) []any {
	ctx := metadata.ContextType(context)
	enums := metadata.GetContextEnums(ctx)

	result := make([]any, len(enums))
	for i, e := range enums {
		result[i] = map[string]any{
			"name":  e.Name,
			"value": e.Value,
		}
	}
	return result
}

// GetOTTLMetadata returns complete OTTL metadata including version information.
func GetOTTLMetadata() map[string]any {
	version := internal.CollectorContribProcessorsVersion

	return map[string]any{
		"version":   version,
		"functions": GetOTTLFunctions(),
		"contexts":  getAllContextsMetadata(),
	}
}

// getAllContextsMetadata returns metadata for all supported contexts.
func getAllContextsMetadata() map[string]any {
	contexts := make(map[string]any)

	for _, ctx := range metadata.AllContextTypes() {
		contexts[string(ctx)] = map[string]any{
			"paths": GetContextPaths(string(ctx)),
			"enums": GetContextEnums(string(ctx)),
		}
	}

	return contexts
}

// ValidateStatements validates OTTL statements using native OTTL parsers.
// This provides accurate syntax and semantic validation with proper position information.
// The dataType, payload, and executorName parameters are kept for API compatibility
// but the native parser doesn't need them for syntax validation.
func ValidateStatements(config, dataType, payload, executorName string) []any {
	// Use native OTTL parser validation for accurate error positions
	return NativeValidateStatements(config)
}

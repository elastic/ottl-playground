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
	"regexp"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottldatapoint"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottllog"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlmetric"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlresource"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlscope"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspan"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/contexts/ottlspanevent"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/ottl/ottlfuncs"
)

// StatementLocation tracks the location of a statement in the YAML config.
type StatementLocation struct {
	Statement string
	Line      int
	Column    int
}

// NativeValidateStatements validates OTTL statements using native parsers.
// This provides accurate syntax and semantic validation with proper position information.
func NativeValidateStatements(config string) []any {
	diagnostics := []any{}

	// Extract statements with their locations from the YAML config
	statements := extractStatementsWithLocations(config)

	if len(statements) == 0 {
		return diagnostics
	}

	// Create a telemetry settings with a no-op logger
	telemetrySettings := component.TelemetrySettings{
		Logger: zap.NewNop(),
	}

	// Detect the context from the config
	context := detectContext(config)

	// Validate each statement
	for _, stmt := range statements {
		err := validateStatement(stmt.Statement, context, telemetrySettings)
		if err != nil {
			// Parse the error to extract position information within the statement
			col, endCol := extractErrorColumnFromMessage(err.Error(), stmt.Statement)

			diagnostics = append(diagnostics, map[string]any{
				"message":   cleanValidationError(err.Error()),
				"severity":  "error",
				"line":      stmt.Line,
				"column":    stmt.Column + col,
				"endLine":   stmt.Line,
				"endColumn": stmt.Column + endCol,
			})
		}
	}

	return diagnostics
}

// extractStatementsWithLocations parses YAML config and extracts statements with their line numbers.
func extractStatementsWithLocations(config string) []StatementLocation {
	var statements []StatementLocation
	lines := strings.Split(config, "\n")

	inStatements := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if we're entering a statements section
		if strings.HasSuffix(trimmed, "statements:") {
			inStatements = true
			continue
		}

		// Check if we're in a new section (non-indented key or context block)
		if inStatements && len(trimmed) > 0 {
			// Check for unindented content or new section
			if !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t") && !strings.HasPrefix(trimmed, "-") {
				inStatements = false
				continue
			}
			// Check for context: declaration (nested section)
			if strings.HasPrefix(trimmed, "context:") || strings.HasPrefix(trimmed, "- context:") {
				inStatements = false
				continue
			}
		}

		// Extract statement from list item
		if inStatements && strings.HasPrefix(trimmed, "- ") {
			statement := strings.TrimPrefix(trimmed, "- ")
			// Remove quotes if present
			statement = strings.Trim(statement, "\"'")

			// Calculate column where statement starts
			col := strings.Index(line, "- ") + 3 // After "- "
			if strings.HasPrefix(strings.TrimPrefix(trimmed, "- "), "\"") ||
				strings.HasPrefix(strings.TrimPrefix(trimmed, "- "), "'") {
				col++ // Account for opening quote
			}

			statements = append(statements, StatementLocation{
				Statement: statement,
				Line:      i + 1, // 1-indexed
				Column:    col,
			})
		}
	}

	return statements
}

// detectContext determines the OTTL context from the YAML config.
func detectContext(config string) string {
	lower := strings.ToLower(config)

	// Check for explicit context declarations
	contextPattern := regexp.MustCompile(`context:\s*(\w+)`)
	if matches := contextPattern.FindStringSubmatch(lower); len(matches) > 1 {
		return normalizeContext(matches[1])
	}

	// Check for statement type prefixes
	if strings.Contains(lower, "log_statements:") || strings.Contains(lower, "log:") {
		return "log"
	}
	if strings.Contains(lower, "trace_statements:") || strings.Contains(lower, "span:") {
		return "span"
	}
	if strings.Contains(lower, "metric_statements:") || strings.Contains(lower, "metric:") {
		return "metric"
	}

	// Default to log
	return "log"
}

// normalizeContext normalizes context names.
func normalizeContext(ctx string) string {
	switch strings.ToLower(ctx) {
	case "log":
		return "log"
	case "span", "trace":
		return "span"
	case "spanevent":
		return "spanevent"
	case "metric":
		return "metric"
	case "datapoint":
		return "datapoint"
	case "resource":
		return "resource"
	case "scope", "instrumentationscope":
		return "scope"
	default:
		return "log"
	}
}

// validateStatement validates a single OTTL statement using the native parser.
func validateStatement(statement, context string, settings component.TelemetrySettings) error {
	switch context {
	case "log":
		return validateLogStatement(statement, settings)
	case "span":
		return validateSpanStatement(statement, settings)
	case "spanevent":
		return validateSpanEventStatement(statement, settings)
	case "metric":
		return validateMetricStatement(statement, settings)
	case "datapoint":
		return validateDataPointStatement(statement, settings)
	case "resource":
		return validateResourceStatement(statement, settings)
	case "scope":
		return validateScopeStatement(statement, settings)
	default:
		return validateLogStatement(statement, settings)
	}
}

func validateLogStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottllog.TransformContext]()
	parser, err := ottllog.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateSpanStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlspan.TransformContext]()
	parser, err := ottlspan.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateSpanEventStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlspanevent.TransformContext]()
	parser, err := ottlspanevent.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateMetricStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlmetric.TransformContext]()
	parser, err := ottlmetric.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateDataPointStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottldatapoint.TransformContext]()
	parser, err := ottldatapoint.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateResourceStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlresource.TransformContext]()
	parser, err := ottlresource.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

func validateScopeStatement(statement string, settings component.TelemetrySettings) error {
	funcs := ottlfuncs.StandardFuncs[ottlscope.TransformContext]()
	parser, err := ottlscope.NewParser(funcs, settings)
	if err != nil {
		return err
	}
	_, err = parser.ParseStatement(statement)
	return err
}

// extractErrorColumnFromMessage tries to extract column position from the error message.
// Returns (start column offset, end column offset) within the statement.
func extractErrorColumnFromMessage(errMsg, statement string) (int, int) {
	// Try to find the problematic segment in the error
	// e.g., 'segment "boy" from path'
	segmentPattern := regexp.MustCompile(`segment "(\w+)"`)
	if matches := segmentPattern.FindStringSubmatch(errMsg); len(matches) > 1 {
		segment := matches[1]
		if idx := strings.Index(statement, segment); idx >= 0 {
			return idx, idx + len(segment)
		}
	}

	// Try to find quoted path in error
	pathPattern := regexp.MustCompile(`path "([^"]+)"`)
	if matches := pathPattern.FindStringSubmatch(errMsg); len(matches) > 1 {
		path := matches[1]
		// Strip context prefix for matching
		pathParts := strings.Split(path, ".")
		if len(pathParts) > 1 {
			shortPath := strings.Join(pathParts[1:], ".")
			if idx := strings.Index(statement, shortPath); idx >= 0 {
				return idx, idx + len(shortPath)
			}
		}
	}

	// Try to find function name in error
	funcPattern := regexp.MustCompile(`function "(\w+)"`)
	if matches := funcPattern.FindStringSubmatch(errMsg); len(matches) > 1 {
		funcName := matches[1]
		if idx := strings.Index(statement, funcName); idx >= 0 {
			return idx, idx + len(funcName)
		}
	}

	// Default: highlight the whole statement
	return 0, len(statement)
}

// cleanValidationError cleans up the error message for display.
func cleanValidationError(errMsg string) string {
	// Remove common prefixes
	prefixes := []string{
		"statement has invalid syntax: ",
		"error while parsing arguments for call to ",
	}
	for _, prefix := range prefixes {
		errMsg = strings.TrimPrefix(errMsg, prefix)
	}

	// Shorten long messages
	if len(errMsg) > 200 {
		errMsg = errMsg[:200] + "..."
	}

	return errMsg
}

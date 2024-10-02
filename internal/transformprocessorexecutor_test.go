// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	transformprocessorConfig = "transformprocessor.yaml"
)

func Test_TransformProcessorExecutor_ParseConfig(t *testing.T) {
	yamlConfig := readTestData(t, transformprocessorConfig)
	executor := NewTransformProcessorExecutor().(*transformProcessorExecutor)
	assert.NotNil(t, executor)
	parsedConfig, err := executor.parseConfig(yamlConfig)
	require.NoError(t, err)

	require.NotNil(t, parsedConfig)
	require.NotEmpty(t, parsedConfig.ErrorMode)
	require.NotEmpty(t, parsedConfig.TraceStatements)
	require.NotEmpty(t, parsedConfig.MetricStatements)
	require.NotEmpty(t, parsedConfig.MetricStatements)
}

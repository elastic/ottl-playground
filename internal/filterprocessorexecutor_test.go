// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	filterprocessorConfig = "filterprocessor.yaml"
)

func Test_FilterProcessorExecutor_ParseConfig(t *testing.T) {
	yamlConfig, err := os.ReadFile(filepath.Join("..", "testdata", filterprocessorConfig))
	require.NoError(t, err)

	executor := NewFilterProcessorExecutor().(*filterProcessorExecutor)
	assert.NotNil(t, executor)

	parsedConfig, err := executor.parseConfig(string(yamlConfig))
	require.NoError(t, err)
	require.NotNil(t, parsedConfig)
	require.NotEmpty(t, parsedConfig.ErrorMode)
	require.NotEmpty(t, parsedConfig.Logs)
	require.NotEmpty(t, parsedConfig.Traces)
	require.NotEmpty(t, parsedConfig.Metrics)
}

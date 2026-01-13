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
	"errors"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/transformprocessor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
)

const (
	configMultiple = "config_multiple.yaml"
)

func Test_parseConfig_singleConfig(t *testing.T) {
	config := readTestData(t, transformprocessorConfig)
	allConfigs, err := parseConfig[transformprocessor.Config](
		component.NewID(transformprocessor.NewFactory().Type()),
		config,
		func() *transformprocessor.Config {
			return transformprocessor.NewFactory().CreateDefaultConfig().(*transformprocessor.Config)
		},
	)

	assert.Len(t, allConfigs, 1, "Expected exactly one config to be parsed")
	pc := allConfigs[0].Value
	require.NoError(t, err)
	require.NotNil(t, pc)
	require.NotEmpty(t, pc.ErrorMode)
	require.NotEmpty(t, pc.TraceStatements)
	require.NotEmpty(t, pc.MetricStatements)
	require.NotEmpty(t, pc.MetricStatements)
}

func Test_parseConfig_configValidator(t *testing.T) {
	id := component.MustNewID("test")
	tests := []struct {
		name      string
		config    string
		expectErr bool
	}{
		{
			name:      "valid single config",
			config:    "valid: true",
			expectErr: false,
		},
		{
			name:      "invalid single config",
			config:    "valid: false",
			expectErr: true,
		},
		{
			name:      "valid multiple config",
			config:    "test:\n valid: true\ntest/b:\n valid: true",
			expectErr: false,
		},
		{
			name:      "invalid multiple config",
			config:    "test:\n valid: true\ntest/b:\n valid: false",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfig[mockConfig](
				id,
				tt.config,
				func() *mockConfig {
					return &mockConfig{Valid: !tt.expectErr}
				},
			)

			if tt.expectErr {
				require.ErrorIs(t, err, errInvalidMockConfig)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_parseConfig_multipleConfig(t *testing.T) {
	config := readTestData(t, configMultiple)
	allConfigs, err := parseConfig[transformprocessor.Config](
		component.NewID(transformprocessor.NewFactory().Type()),
		config,
		func() *transformprocessor.Config {
			return transformprocessor.NewFactory().CreateDefaultConfig().(*transformprocessor.Config)
		},
	)

	assert.Len(t, allConfigs, 3, "Expected exactly 3 configs to be parsed")
	// Validate that the keys and values and order are as expected
	for i, key := range []string{"transform", "transform/b", "transform/a"} {
		c := allConfigs[i]
		require.NotNil(t, c)
		assert.Equal(t, key, c.Key)
		assert.NotEmpty(t, c.Value)
		require.NoError(t, err)
		require.NotEmpty(t, c.Value.ErrorMode)
		require.NotEmpty(t, c.Value.TraceStatements)
		require.NotEmpty(t, c.Value.MetricStatements)
		require.NotEmpty(t, c.Value.MetricStatements)
	}
}

func Test_parseConfig_invalidConfig(t *testing.T) {
	_, err := parseConfig[transformprocessor.Config](
		component.NewID(transformprocessor.NewFactory().Type()),
		"---invalid---",
		func() *transformprocessor.Config {
			return transformprocessor.NewFactory().CreateDefaultConfig().(*transformprocessor.Config)
		},
	)
	require.ErrorContains(t, err, "cannot be used as a Conf")
}

var (
	errInvalidMockConfig = errors.New("config validation failed")
)

type mockConfig struct {
	Valid bool `mapstructure:"valid"`
}

func (m *mockConfig) Validate() error {
	if m.Valid {
		return nil
	}
	return errInvalidMockConfig
}

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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	StringField string  `json:"stringField"`
	IntField    int     `json:"intField"`
	BoolField   bool    `json:"boolField"`
	FloatField  float64 `json:"floatField"`
}

type nestedStruct struct {
	Nested testStruct `json:"nested"`
	Array  []int      `json:"array"`
}

func Test_structToMap_SimpleStruct(t *testing.T) {
	input := testStruct{
		StringField: "test string",
		IntField:    42,
		BoolField:   true,
		FloatField:  3.14,
	}

	result, err := structToMap(input)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "test string", result["stringField"])
	assert.Equal(t, float64(42), result["intField"]) // JSON unmarshalling converts numbers to float64
	assert.Equal(t, true, result["boolField"])
	assert.Equal(t, 3.14, result["floatField"])
}

func Test_structToMap_NestedStruct(t *testing.T) {
	input := nestedStruct{
		Nested: testStruct{
			StringField: "nested string",
			IntField:    10,
			BoolField:   false,
			FloatField:  2.71,
		},
		Array: []int{1, 2, 3},
	}

	result, err := structToMap(input)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Check nested struct
	nested, ok := result["nested"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "nested string", nested["stringField"])
	assert.Equal(t, float64(10), nested["intField"])
	assert.Equal(t, false, nested["boolField"])
	assert.Equal(t, 2.71, nested["floatField"])

	// Check array
	array, ok := result["array"].([]any)
	require.True(t, ok)
	require.Len(t, array, 3)
	assert.Equal(t, float64(1), array[0])
	assert.Equal(t, float64(2), array[1])
	assert.Equal(t, float64(3), array[2])
}

func Test_structToMap_EmptyStruct(t *testing.T) {
	input := testStruct{}

	result, err := structToMap(input)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "", result["stringField"])
	assert.Equal(t, float64(0), result["intField"])
	assert.Equal(t, false, result["boolField"])
	assert.Equal(t, float64(0), result["floatField"])
}

func Test_structToMap_NilValue(t *testing.T) {
	var input *testStruct

	result, err := structToMap(input)
	require.NoError(t, err)
	assert.Nil(t, result)
}

func Test_structToMap_MapInput(t *testing.T) {
	input := map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	result, err := structToMap(input)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, float64(42), result["key2"])
	assert.Equal(t, true, result["key3"])
}

func Test_structToMap_SliceInput(t *testing.T) {
	input := []string{"a", "b", "c"}

	result, err := structToMap(input)
	require.ErrorContains(t, err, "failed to unmarshal value")
	assert.Nil(t, result)
}

func Test_structToMap_WithPointers(t *testing.T) {
	input := &testStruct{
		StringField: "pointer test",
		IntField:    100,
		BoolField:   true,
		FloatField:  1.23,
	}

	result, err := structToMap(input)
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, "pointer test", result["stringField"])
	assert.Equal(t, float64(100), result["intField"])
	assert.Equal(t, true, result["boolField"])
	assert.Equal(t, 1.23, result["floatField"])
}

package converter

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStructToJson(t *testing.T) {
	s := struct {
		TestValue string `json:"test_value"`
	}{
		TestValue: "abc",
	}

	res, err := StructToJson(s)
	require.NoError(t, err)
	assert.Equal(t, `{"test_value":"abc"}`, res)

	res, err = StructToJson(math.NaN())
	require.Error(t, err)
	assert.Empty(t, res)
}

func TestMapToStruct(t *testing.T) {
	i := map[string]interface{}{
		"testValue": "abc",
	}
	s := struct {
		Test_value string `mapstructure:"testValue"`
	}{}

	err := MapToStruct(i, &s)
	require.NoError(t, err)
	assert.Equal(t, i["testValue"], s.Test_value)

	var n int
	err = MapToStruct(i, &n)
	require.Error(t, err)
}

func TestStructToMap(t *testing.T) {
	s := struct {
		TestValue string  `json:"test_value"`
		Count     float64 `json:"count"`
	}{
		TestValue: "abc",
		Count:     2,
	}
	m, err := StructToMap(s)
	require.NoError(t, err)
	assert.Len(t, m, 2)
	assert.Equal(t, s.TestValue, m["test_value"])
	assert.InEpsilon(t, s.Count, m["count"], 0)

	m, err = StructToMap(math.NaN())
	require.Error(t, err)
	assert.Nil(t, m)

	m, err = StructToMap(123)
	require.Error(t, err)
	assert.Nil(t, m)
}

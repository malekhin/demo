package common

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dateTimeTestCases = []string{"2025-07-12", `"2000-01-01"`, "2025-07-12T02:29:05+03:00"}

func Test_DateTime(t *testing.T) {
	for _, date := range dateTimeTestCases {
		dt := DateTime{}
		err := dt.UnmarshalJSON([]byte(date))
		require.NoError(t, err, "")

		value, err := dt.Value()
		require.NoError(t, err, "")
		assert.True(t, driver.IsScanValue(value))
	}

	test := "200-07-12"
	dt := DateTime{}
	err := dt.UnmarshalJSON([]byte(test))
	require.Error(t, err)

	value, err := dt.Value()
	require.NoError(t, err)
	assert.Nil(t, value)

	test = "null"
	dt = DateTime{}
	err = dt.UnmarshalJSON([]byte(test))
	require.NoError(t, err)

	json, err := dt.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, "null", string(json))

	test = "2000-01-01"
	dt = DateTime{}
	err = dt.UnmarshalJSON([]byte(test))
	require.NoError(t, err)

	json, err = dt.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `"2000-01-01T00:00:00+03:00"`, string(json))

	dt = DateTime{}
	err = dt.Scan(time.Now())
	require.NoError(t, err)
	assert.True(t, dt.Valid)
}

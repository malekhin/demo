package database

import (
	"database/sql/driver"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonWrapper(t *testing.T) {
	jw := JsonWrapper{}
	var err error

	json := `{"test":1,"test2":"stringValue"}`
	err = jw.Scan(json)
	require.NoError(t, err)

	j, err := jw.GetJson()
	require.NoError(t, err)
	assert.Equal(t, json, j)

	value, err := jw.Value()
	require.NoError(t, err)
	assert.True(t, driver.IsScanValue(value))

	var res struct {
		Test2 string `json:"test2"`
	}
	err = jw.Get(&res)
	require.NoError(t, err)
	assert.NotEmpty(t, res.Test2)

	n, err := jw.GetInt("test")
	require.NoError(t, err)
	assert.Equal(t, 1, n)

	_, err = jw.GetInt("test2")
	require.Error(t, err)

	n, err = jw.GetInt("not_exists")
	require.Error(t, err)
	assert.Empty(t, n)

	jw = JsonWrapper{}
	err = jw.Scan(`true`)
	require.NoError(t, err)
	n, err = jw.GetInt("test")
	require.NoError(t, err)
	assert.Empty(t, n)

	jw = JsonWrapper{Data: math.NaN()}
	var res2 struct{}
	err = jw.Get(&res2)
	require.Error(t, err)
}

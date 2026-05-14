package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIn(t *testing.T) {
	tests := []struct {
		query         string
		args          []interface{}
		expectedQuery string
		expectArgs    []interface{}
	}{
		{
			query:         `select * from tariff where product_type in ($2) and is_deleted = $1`,
			args:          []interface{}{false, []string{"osago", "property"}},
			expectedQuery: `select * from tariff where product_type in ($2, $3) and is_deleted = $1`,
			expectArgs:    []interface{}{false, "osago", "property"},
		},
		{
			query:         `select * from tariff where product_type in(?) and is_deleted = ?`,
			args:          []interface{}{[]string{"osago", "property"}, false},
			expectedQuery: `select * from tariff where product_type in(?, ?) and is_deleted = ?`,
			expectArgs:    []interface{}{"osago", "property", false},
		},
		{
			query:         `select * from tariff where product_type in( ?)`,
			args:          []interface{}{[]string{"osago", "property"}},
			expectedQuery: `select * from tariff where product_type in( ?, ?)`,
			expectArgs:    []interface{}{"osago", "property"},
		},
	}

	for _, test := range tests {
		query, args, err := in(test.query, test.args)
		require.NoError(t, err)

		assert.Equal(t, test.expectedQuery, query)
		assert.Equal(t, len(args), len(test.expectArgs))
		for i, expectArg := range test.expectArgs {
			assert.Equal(t, expectArg, args[i])
		}
	}
}

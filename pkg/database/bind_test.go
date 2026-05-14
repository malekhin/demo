package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Нужно убедиться что в разных вариантах написания
// (внутри where, в конце строки, со знаком запятой в конце, без знака запятой, другие варианты)
// происходит конвертирование именованных параметров в bindVar
func TestDbBind(t *testing.T) {
	tests := []struct {
		query   string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			query: `
						select array_agg(name)[:number], product_type
						from tariff
						where product_type = :product_type and kv_from_sk_percent > :kv
						group by product_type, is_portion
		having is_portion = :isPortion`, // Это не ошибка форматирования, проверяется именованый параметр в конце строки
			args: map[string]interface{}{
				"number":       0,
				"product_type": "osago",
				"kv":           0,
				"isPortion":    true,
			},
			wantErr: false,
		},
		{
			query: `
							insert into tariff
								(name, kv_from_sk_percent, product_type, is_deleted)
							values
								(:name, :kv_from_sk_percent, :product_type , :is_deleted)
					`, // Это не ошибка форматирования, проверяется знак вопроса не после именовоного параметра
			args: map[string]interface{}{
				"name":               "Тест",
				"kv_from_sk_percent": 1.99,
				"product_type":       "osago",
				"is_deleted":         true,
			},
			wantErr: false,
		},
		{
			query: `
							update tariff
							set name = :name, kv_from_sk_percent = :kvPercent
							where id = :id;
					`,
			args: map[string]interface{}{
				"name":      "Тест",
				"kvPercent": 1.99,
				"id":        1,
			},
			wantErr: false,
		},
		{
			query: `
					update tariff
					set name = :name
					where id = :id;
			`, // Не верно переданные параметры (нужно передать id, передано kvPercent)
			args: map[string]interface{}{
				"name":      "Тест",
				"kvPercent": 1.99,
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		query, retArgs, err := bindMap(test.query, test.args)
		if !test.wantErr {
			require.NoError(t, err)
			assert.NotEmpty(t, query)
			assert.Len(t, retArgs, len(test.args))
		} else {
			require.Error(t, err)
		}
	}
}

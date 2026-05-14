package database

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const codeDuplicate = "23505"

type Error error

var ErrDuplicate Error = errors.New("duplicate")

// Обрабатываем ошибки с кодом
func getError(err error) error {
	var e *pgconn.PgError
	if errors.As(err, &e) {
		switch e.Code {
		case codeDuplicate:
			return ErrDuplicate
		}
	}

	return err
}

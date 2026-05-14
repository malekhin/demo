package database

import (
	"context"
	"demo/internal/common"
	"fmt"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type db_ struct {
	db     Sql
	driver Driver
}

func NewDb(db *sqlx.DB, driver Driver) *db_ {
	return &db_{db: db, driver: driver}
}

func (d *db_) NamedRows(ctx context.Context, result interface{}, query string, arg interface{}) error {
	query, args, err := bindNamed(query, arg)
	if err != nil {
		return common.Wrap(err, "bindNamed")
	}

	err = d.Rows(ctx, result, query, args...)
	if err != nil {
		return common.Wrap(err, "Rows")
	}

	return nil
}

func (d *db_) NamedRow(ctx context.Context, result interface{}, query string, arg interface{}) error {
	query, args, err := bindNamed(query, arg)
	if err != nil {
		return common.Wrap(err, "bindNamed")
	}

	err = d.Row(ctx, result, query, args...)
	if err != nil {
		return common.Wrap(err, "Row")
	}

	return nil
}

func (d *db_) NamedExec(ctx context.Context, query string, arg interface{}) error {
	query, args, err := bindNamed(query, arg)
	if err != nil {
		return common.Wrap(err, "bindNamed")
	}

	err = d.Exec(ctx, query, args...)
	if err != nil {
		return common.Wrap(err, "Exec")
	}

	return nil
}

func (d *db_) NamedBatch(ctx context.Context, query string, arg interface{}) error {
	_, err := d.sqlx(ctx).NamedExecContext(ctx, d.bind(query), arg)
	if err != nil {
		return getError(err)
	}

	return nil
}

func (d *db_) Rows(ctx context.Context, result interface{}, query string, args ...interface{}) error {
	query, args, err := in(query, args)
	if err != nil {
		return common.Wrap(err, "in")
	}

	err = d.sqlx(ctx).SelectContext(ctx, result, d.bind(query), args...)
	if err != nil {
		return getError(err)
	}

	return nil
}

func (d *db_) Row(ctx context.Context, result interface{}, query string, args ...interface{}) error {
	query, args, err := in(query, args)
	if err != nil {
		return common.Wrap(err, "in")
	}

	err = d.sqlx(ctx).GetContext(ctx, result, d.bind(query), args...)
	if err != nil {
		return getError(err)
	}

	return nil
}

func (d *db_) Exec(ctx context.Context, query string, args ...interface{}) error {
	query, args, err := in(query, args)
	if err != nil {
		return common.Wrap(err, "in")
	}

	_, err = d.sqlx(ctx).ExecContext(ctx, d.bind(query), args...)
	if err != nil {
		return getError(err)
	}

	return nil
}

// Определяет нужно ли использовать транзакцию
func (d *db_) sqlx(ctx context.Context) Sql {
	if tx, ok := ctx.Value(txKey{}).(Sql); ok {
		return tx
	}

	return d.db
}

func connectSqlX(c Config) (*sqlx.DB, error) {
	var dbUrl string

	switch c.Driver {
	case Pgx:
		dbUrl = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s default_query_exec_mode=%s",
			c.Host, c.Port, c.User, c.Password, c.Name, c.SslMode, c.QueryExecMode,
		)
	case Mysql:
		dbUrl = fmt.Sprintf(
			"%s:%s@(%s:%d)/%s",
			c.User, c.Password, c.Host, c.Port, c.Name,
		)
	}

	db, err := sqlx.Connect(string(c.Driver), dbUrl)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}

	return db, nil
}

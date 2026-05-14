package database

import (
	"context"
	"database/sql"
	"demo/internal/common"
	"errors"

	"github.com/jmoiron/sqlx"
)

type txKey struct{}

func NewTx(database Database) Tx {
	return &tx_{dbConnect: database.ConnectList()}
}

type Tx interface {
	Db() TxManager
	Slave() TxManager
}

type tx_ struct {
	dbConnect []*sqlx.DB
}

func (tx tx_) Db() TxManager {
	if len(tx.dbConnect) >= 1 {
		return &txManager_{db: tx.dbConnect[0], driver: Driver(tx.dbConnect[1].DriverName())}
	}

	panic("master database is not connected")
}

func (tx tx_) Slave() TxManager {
	if len(tx.dbConnect) >= 2 {
		return &txManager_{db: tx.dbConnect[1], driver: Driver(tx.dbConnect[1].DriverName())}
	}

	panic("slave database is not connected")
}

//go:generate mockery --name=TxManager
type TxManager interface {
	BeginTx(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) (context.Context, error)
	Rollback(ctx context.Context) (context.Context, error)
}

type txManager_ struct {
	db     *sqlx.DB
	driver Driver
}

func (t *txManager_) BeginTx(ctx context.Context) (context.Context, error) {
	var err error
	tx, err := t.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return nil, common.Wrap(err, "BeginTxx")
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	return ctx, nil
}

func (t *txManager_) Rollback(ctx context.Context) (context.Context, error) {
	tx, ok := ctx.Value(txKey{}).(*sqlx.Tx)
	if !ok {
		return ctx, errors.New("transaction in not started")
	}

	err := tx.Rollback()
	if err != nil {
		return ctx, common.Wrap(err, "Rollback")
	}

	return context.WithValue(ctx, txKey{}, nil), nil
}

func (t *txManager_) Commit(ctx context.Context) (context.Context, error) {
	tx, ok := ctx.Value(txKey{}).(*sqlx.Tx)
	if !ok {
		return ctx, errors.New("transaction in not started")
	}

	err := tx.Commit()
	if err != nil {
		return ctx, common.Wrap(err, "Commit")
	}

	return context.WithValue(ctx, txKey{}, nil), nil
}

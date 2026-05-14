package database

import (
	"context"
	"database/sql"
	"demo/internal/common/config"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	Db() Db
	Slave() Db
	ConnectList() []*sqlx.DB
}

type database_ struct {
	dbList      []*db_
	connectList []*sqlx.DB
}

type Db interface {
	NamedRows(ctx context.Context, result interface{}, query string, arg interface{}) error
	NamedRow(ctx context.Context, result interface{}, query string, arg interface{}) error
	NamedExec(ctx context.Context, query string, arg interface{}) error
	NamedBatch(ctx context.Context, query string, arg interface{}) error
	Rows(ctx context.Context, result interface{}, query string, args ...interface{}) error
	Row(ctx context.Context, result interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) error
}

type Sql interface {
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

func New(config *config.Config) Database {
	configs := []Config{
		{
			Driver:        Driver(config.Db.Driver),
			Host:          config.Db.Host,
			Port:          config.Db.Port,
			User:          config.Db.User,
			Password:      config.Db.Password,
			Name:          config.Db.Name,
			QueryExecMode: config.Db.QueryExecMode,
			SslMode:       config.Db.SslMode,
		},
		{
			Driver:        Driver(config.MonolithDb.Driver),
			Host:          config.MonolithDb.Host,
			Port:          config.MonolithDb.Port,
			User:          config.MonolithDb.User,
			Password:      config.MonolithDb.Password,
			Name:          config.MonolithDb.Name,
			QueryExecMode: config.MonolithDb.QueryExecMode,
			SslMode:       config.MonolithDb.SslMode,
		},
	}

	database := database_{
		dbList:      make([]*db_, len(configs)),
		connectList: make([]*sqlx.DB, len(configs)),
	}

	for i, dbConfig := range configs {
		connectDb, err := connectSqlX(dbConfig)
		if err != nil {
			panic(err)
		}
		database.dbList[i] = NewDb(connectDb, dbConfig.Driver)
		database.connectList[i] = connectDb
	}

	return &database
}

func (d database_) Db() Db {
	if len(d.dbList) >= 1 {
		return d.dbList[0]
	}

	panic("master database is not connected")
}

func (d database_) Slave() Db {
	if len(d.dbList) >= 2 {
		return d.dbList[1]
	}

	panic("slave database is not connected")
}

func (d database_) ConnectList() []*sqlx.DB {
	return d.connectList
}

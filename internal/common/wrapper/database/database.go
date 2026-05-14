package database

import (
	"demo/internal/common/metrics"
	"demo/pkg/database"

	"github.com/jmoiron/sqlx"
)

const (
	Master = "master"
	Slave  = "slave"
)

type Database interface {
	Db() Db
	Slave() Db
	ConnectList() []*sqlx.DB
}

type database_ struct {
	database database.Database
	metrics  metrics.IMetrics
}

func New(database database.Database, metrics metrics.IMetrics) Database {
	return &database_{database: database, metrics: metrics}
}

func (d database_) Db() Db {
	return NewDb(d.database.Db(), d.metrics, Master)
}

func (d database_) Slave() Db {
	return NewDb(d.database.Slave(), d.metrics, Slave)
}

func (d database_) ConnectList() []*sqlx.DB {
	return d.database.ConnectList()
}

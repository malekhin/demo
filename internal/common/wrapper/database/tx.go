package database

import (
	"demo/internal/common/metrics"
	"demo/pkg/database"
)

type Tx interface {
	Db() TxManager
	Slave() TxManager
}

type tx_ struct {
	tx      database.Tx
	metrics metrics.IMetrics
}

func NewTx(tx database.Tx, metrics metrics.IMetrics) Tx {
	return &tx_{tx: tx, metrics: metrics}
}

func (tx tx_) Db() TxManager {
	return NewTxManager(tx.tx.Db(), tx.metrics, Master)
}

func (tx tx_) Slave() TxManager {
	return NewTxManager(tx.tx.Slave(), tx.metrics, Slave)
}

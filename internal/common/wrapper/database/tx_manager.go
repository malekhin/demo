package database

import (
	"context"
	"demo/internal/common"
	"demo/internal/common/metrics"
	"demo/pkg/database"
	"fmt"
	"time"
)

type TxManager interface {
	BeginTx(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) (context.Context, error)
	Rollback(ctx context.Context) (context.Context, error)
}

type txManager_ struct {
	txManager   database.TxManager
	metrics     metrics.IMetrics
	connectName string
}

func NewTxManager(txManager database.TxManager, metrics metrics.IMetrics, connectName string) TxManager {
	return &txManager_{txManager: txManager, metrics: metrics, connectName: connectName}
}

func (t txManager_) BeginTx(ctx context.Context) (context.Context, error) {
	start := time.Now()
	ctx, err := t.txManager.BeginTx(ctx)
	t.metrics.RecordExternalCall(t.getDependency(), "begin", t.getStatus(err), time.Since(start))
	if err != nil {
		return ctx, common.Wrap(err, "Wrapper: BeginTx")
	}

	return ctx, nil
}

func (t txManager_) Commit(ctx context.Context) (context.Context, error) {
	start := time.Now()
	ctx, err := t.txManager.Commit(ctx)
	t.metrics.RecordExternalCall(t.getDependency(), "commit", t.getStatus(err), time.Since(start))
	if err != nil {
		return ctx, common.Wrap(err, "Wrapper: Commit")
	}

	return ctx, nil
}

func (t txManager_) Rollback(ctx context.Context) (context.Context, error) {
	start := time.Now()
	ctx, err := t.txManager.Rollback(ctx)
	t.metrics.RecordExternalCall(t.getDependency(), "rollback", t.getStatus(err), time.Since(start))
	if err != nil {
		return ctx, common.Wrap(err, "Wrapper: Rollback")
	}

	return ctx, nil
}

func (t txManager_) getDependency() string {
	return fmt.Sprintf("%s_%s", metrics.DependencyDatabase, t.connectName)
}

func (t txManager_) getStatus(err error) string {
	if err != nil {
		return metrics.StatusError
	}

	return metrics.StatusSuccess
}

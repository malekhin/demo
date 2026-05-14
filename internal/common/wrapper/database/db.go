package database

import (
	"context"
	"demo/internal/common"
	"demo/internal/common/metrics"
	"demo/pkg/database"
	"fmt"
	"strings"
	"time"
)

type Db interface {
	NamedRows(ctx context.Context, result interface{}, query string, arg interface{}) error
	NamedRow(ctx context.Context, result interface{}, query string, arg interface{}) error
	NamedExec(ctx context.Context, query string, arg interface{}) error
	NamedBatch(ctx context.Context, query string, arg interface{}) error
	Rows(ctx context.Context, result interface{}, query string, args ...interface{}) error
	Row(ctx context.Context, result interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) error
}

type db_ struct {
	db          database.Db
	metrics     metrics.IMetrics
	connectName string
}

func NewDb(db database.Db, metrics metrics.IMetrics, connectName string) Db {
	return &db_{db: db, metrics: metrics, connectName: connectName}
}

func (d *db_) NamedRows(ctx context.Context, result interface{}, query string, arg interface{}) error {
	start := time.Now()
	err := d.db.NamedRows(ctx, result, query, arg)
	d.metrics.RecordExternalCall(d.getDependency(), d.getOperation(query), d.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: NamedRows")
	}

	return nil
}

func (d *db_) NamedRow(ctx context.Context, result interface{}, query string, arg interface{}) error {
	start := time.Now()
	err := d.db.NamedRow(ctx, result, query, arg)
	d.metrics.RecordExternalCall(d.getDependency(), d.getOperation(query), d.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: NamedRow")
	}

	return nil
}

func (d *db_) NamedExec(ctx context.Context, query string, arg interface{}) error {
	start := time.Now()
	err := d.db.NamedExec(ctx, query, arg)
	d.metrics.RecordExternalCall(d.getDependency(), d.getOperation(query), d.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: NamedExec")
	}

	return nil
}

func (d *db_) NamedBatch(ctx context.Context, query string, arg interface{}) error {
	start := time.Now()
	err := d.db.NamedBatch(ctx, query, arg)
	d.metrics.RecordExternalCall(d.getDependency(), d.getOperation(query), d.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: NamedBatch")
	}

	return nil
}

func (d *db_) Rows(ctx context.Context, result interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := d.db.Rows(ctx, result, query, args...)
	d.metrics.RecordExternalCall(d.getDependency(), d.getOperation(query), d.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: Rows")
	}

	return nil
}

func (d *db_) Row(ctx context.Context, result interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := d.db.Row(ctx, result, query, args...)
	d.metrics.RecordExternalCall(d.getDependency(), d.getOperation(query), d.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: Row")
	}

	return nil
}

func (d *db_) Exec(ctx context.Context, query string, args ...interface{}) error {
	start := time.Now()
	err := d.db.Exec(ctx, query, args...)
	d.metrics.RecordExternalCall(d.getDependency(), d.getOperation(query), d.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: Exec")
	}

	return nil
}

func (d *db_) getDependency() string {
	return fmt.Sprintf("%s_%s", metrics.DependencyDatabase, d.connectName)
}

func (d *db_) getOperation(query string) string {
	for _, key := range []string{"select", "update", "insert"} {
		if strings.Contains(query, key) {
			return key
		}
	}

	return "unknown"
}

func (d *db_) getStatus(err error) string {
	if err != nil {
		return metrics.StatusError
	}

	return metrics.StatusSuccess
}

package clickhouse

import (
	"context"
	"demo/internal/common"
	"demo/internal/common/metrics"
	"demo/pkg/clickhouse"
	"strings"
	"time"
)

type IClickhouse interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Rows(ctx context.Context, result interface{}, query string, args ...interface{}) error
	Row(ctx context.Context, result interface{}, query string, args ...interface{}) error
	Insert(ctx context.Context, query string, arg interface{}) error
	AsyncInsert(ctx context.Context, query string) error
}

type Clickhouse struct {
	clickhouse clickhouse.IClickhouse
	metrics    metrics.IMetrics
}

func New(clickhouse clickhouse.IClickhouse, metrics metrics.IMetrics) IClickhouse {
	return &Clickhouse{clickhouse: clickhouse, metrics: metrics}
}

func (c Clickhouse) Exec(ctx context.Context, query string, args ...interface{}) error {
	start := time.Now()
	err := c.clickhouse.Exec(ctx, query, args...)
	c.metrics.RecordExternalCall(c.getDependency(), c.getOperation(query), c.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: Exec")
	}

	return nil
}

func (c Clickhouse) Rows(ctx context.Context, result interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := c.clickhouse.Rows(ctx, result, query, args...)
	c.metrics.RecordExternalCall(c.getDependency(), c.getOperation(query), c.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: Rows")
	}

	return nil
}

func (c Clickhouse) Row(ctx context.Context, result interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := c.clickhouse.Row(ctx, result, query, args...)
	c.metrics.RecordExternalCall(c.getDependency(), c.getOperation(query), c.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: Row")
	}

	return nil
}

func (c Clickhouse) Insert(ctx context.Context, query string, arg interface{}) error {
	start := time.Now()
	err := c.clickhouse.Insert(ctx, query, arg)
	c.metrics.RecordExternalCall(c.getDependency(), c.getOperation(query), c.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: Insert")
	}

	return nil
}

func (c Clickhouse) AsyncInsert(ctx context.Context, query string) error {
	start := time.Now()
	err := c.clickhouse.AsyncInsert(ctx, query)
	c.metrics.RecordExternalCall(c.getDependency(), c.getOperation(query), c.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper: AsyncInsert")
	}

	return nil
}

func (c Clickhouse) getDependency() string {
	return metrics.DependencyClickhouse
}

func (c Clickhouse) getOperation(query string) string {
	for _, key := range []string{"select", "update", "insert"} {
		if strings.Contains(query, key) {
			return key
		}
	}

	return "unknown"
}

func (c Clickhouse) getStatus(err error) string {
	if err != nil {
		return metrics.StatusError
	}

	return metrics.StatusSuccess
}

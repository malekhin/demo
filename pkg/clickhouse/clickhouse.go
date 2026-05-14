package clickhouse

import (
	"context"
	"demo/internal/common/config"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/zap"
)

type Clickhouse struct {
	conn  driver.Conn
	batch *batch
}

//go:generate mockery --name=IClickhouse
type IClickhouse interface {
	Exec(ctx context.Context, query string, args ...interface{}) error
	Rows(ctx context.Context, result interface{}, query string, args ...interface{}) error
	Row(ctx context.Context, result interface{}, query string, args ...interface{}) error
	Insert(ctx context.Context, query string, arg interface{}) error
	AsyncInsert(ctx context.Context, query string) error
}

func New(config *config.Config, logger *zap.Logger) IClickhouse {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", config.Clickhouse.Host, config.Clickhouse.Port)},
		Auth: clickhouse.Auth{
			Database: config.Clickhouse.Name,
			Username: config.Clickhouse.User,
			Password: config.Clickhouse.Password,
		},
	})
	if err != nil {
		panic(err)
	}

	return &Clickhouse{
		conn:  conn,
		batch: newBatch(conn, config.Clickhouse, logger),
	}
}

func (c *Clickhouse) Insert(ctx context.Context, query string, arg interface{}) error {
	err := c.batch.Insert(ctx, query, arg)
	if err != nil {
		return fmt.Errorf("%w: QueryRow", err)
	}

	return nil
}

func (c *Clickhouse) AsyncInsert(ctx context.Context, query string) error {
	err := c.conn.AsyncInsert(ctx, query, true)
	if err != nil {
		return fmt.Errorf("%w: QueryRow", err)
	}

	return nil
}

func (c *Clickhouse) Row(ctx context.Context, result interface{}, query string, args ...interface{}) error {
	row := c.conn.QueryRow(ctx, query, args...)
	err := row.Scan(result)
	if err != nil {
		return fmt.Errorf("%w: Scan", err)
	}

	return nil
}

func (c *Clickhouse) Rows(ctx context.Context, result interface{}, query string, args ...interface{}) error {
	err := c.conn.Select(ctx, result, query, args...)
	if err != nil {
		return fmt.Errorf("%w: Select", err)
	}

	return nil
}

func (c *Clickhouse) Exec(ctx context.Context, query string, args ...interface{}) error {
	err := c.conn.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%w: Select", err)
	}

	return nil
}

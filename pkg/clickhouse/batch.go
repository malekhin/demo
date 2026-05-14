package clickhouse

import (
	"context"
	"crypto/sha256"
	"demo/internal/common/config"
	"fmt"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"go.uber.org/zap"
)

type batch struct {
	conn     driver.Conn
	elements map[string]*elem
	config   config.ClickhouseConfig
	logger   *zap.Logger
	mu       sync.Mutex
}

type elem struct {
	driver.Batch
	runAt time.Time
}

func newBatch(conn driver.Conn, config config.ClickhouseConfig, logger *zap.Logger) *batch {
	b := &batch{
		conn:     conn,
		config:   config,
		logger:   logger,
		elements: make(map[string]*elem),
	}

	go func() {
		b.executor()
	}()

	return b
}

func (b *batch) Insert(ctx context.Context, query string, arg interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	bat, err := b.getBatch(ctx, query)
	if err != nil {
		return err
	}

	err = bat.AppendStruct(arg)
	if err != nil {
		return fmt.Errorf("%w: AppendStruct", err)
	}

	return nil
}

func (b *batch) getBatch(ctx context.Context, query string) (*elem, error) {
	hash := b.getHash(query)
	if bat, ok := b.elements[hash]; ok {
		return bat, nil
	}

	bat, err := b.conn.PrepareBatch(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: PrepareBatch", err)
	}

	b.elements[hash] = &elem{
		Batch: bat,
		runAt: time.Now().Add(b.config.BatchWaiting),
	}

	return b.elements[hash], nil
}

func (b *batch) getHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return string(bs)
}

func (b *batch) executor() {
	for {
		for hash, elem := range b.elements {
			if time.Now().After(elem.runAt) {
				b.mu.Lock()
				if err := elem.Send(); err != nil {
					b.logger.Error(fmt.Sprintf("Clickhouse: %v", err))
				}

				delete(b.elements, hash)
				b.mu.Unlock()
			}
		}

		time.Sleep(time.Second)
	}
}

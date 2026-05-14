package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/redis/go-redis/v9"
)

const prefixLock = "lock_"

func (r Redis) Lock(ctx context.Context, key string, exp time.Duration) (bool, error) {
	key = r.prefix(prefixLock + key)
	val := uuid.Must(uuid.NewV4()).String()

	ok, err := r.client.SetNX(ctx, key, val, exp).Result()
	if err != nil {
		return false, fmt.Errorf("can't redis.setnx: %w", err)
	}

	return ok, nil
}

func (r Redis) Unlock(ctx context.Context, key string) error {
	key = r.prefix(prefixLock + key)

	err := r.client.Watch(ctx, func(tx *redis.Tx) error {
		_, err := tx.Get(ctx, key).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			return fmt.Errorf("can't redis.get: %w", err)
		}

		return tx.Del(ctx, key).Err()
	}, key)

	if err != nil && !errors.Is(err, redis.TxFailedErr) {
		return fmt.Errorf("can't redis.watch: %w", err)
	}

	return nil
}

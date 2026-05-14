package redis

import (
	"context"
	"demo/internal/common"
	"demo/internal/common/metrics"
	"demo/pkg/redis"
	"time"
)

const operationSet = "set"
const operationGet = "get"
const operationSetGet = "setget"
const operationLock = "lock"
const operationUnlock = "unlock"

type IRedis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetInt(ctx context.Context, key string) (int, bool, error)
	GetFloat(ctx context.Context, key string) (float64, bool, error)
	GetString(ctx context.Context, key string) (string, bool, error)
	GetStruct(ctx context.Context, key string, output interface{}) (bool, error)
	SetAndGet(ctx context.Context, result interface{}, f func() (interface{}, error), key string, expires time.Duration) error
	Lock(ctx context.Context, key string, exp time.Duration) (bool, error)
	Unlock(ctx context.Context, key string) error
}

type Redis struct {
	redis   redis.IRedis
	metrics metrics.IMetrics
}

func NewRedis(redis redis.IRedis, metrics metrics.IMetrics) IRedis {
	return &Redis{redis: redis, metrics: metrics}
}

func (r Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	start := time.Now()
	err := r.redis.Set(ctx, key, value, expiration)
	r.metrics.RecordExternalCall(r.getDependency(), operationSet, r.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper Set")
	}

	return nil
}

func (r Redis) GetInt(ctx context.Context, key string) (int, bool, error) {
	start := time.Now()
	val, exists, err := r.redis.GetInt(ctx, key)
	r.metrics.RecordExternalCall(r.getDependency(), operationGet, r.getStatus(err), time.Since(start))
	if err != nil {
		return 0, false, common.Wrap(err, "Wrapper GetInt")
	}

	return val, exists, nil
}

func (r Redis) GetFloat(ctx context.Context, key string) (float64, bool, error) {
	start := time.Now()
	val, exists, err := r.redis.GetFloat(ctx, key)
	r.metrics.RecordExternalCall(r.getDependency(), operationGet, r.getStatus(err), time.Since(start))
	if err != nil {
		return 0, false, common.Wrap(err, "Wrapper GetFloat")
	}

	return val, exists, nil
}

func (r Redis) GetString(ctx context.Context, key string) (string, bool, error) {
	start := time.Now()
	val, exists, err := r.redis.GetString(ctx, key)
	r.metrics.RecordExternalCall(r.getDependency(), operationGet, r.getStatus(err), time.Since(start))
	if err != nil {
		return "", false, common.Wrap(err, "Wrapper GetString")
	}

	return val, exists, nil
}

func (r Redis) GetStruct(ctx context.Context, key string, output interface{}) (bool, error) {
	start := time.Now()
	exists, err := r.redis.GetStruct(ctx, key, output)
	r.metrics.RecordExternalCall(r.getDependency(), operationGet, r.getStatus(err), time.Since(start))
	if err != nil {
		return false, common.Wrap(err, "Wrapper GetStruct")
	}

	return exists, nil
}

func (r Redis) SetAndGet(ctx context.Context, result interface{}, f func() (interface{}, error), key string, expires time.Duration) error {
	start := time.Now()
	err := r.redis.SetAndGet(ctx, result, f, key, expires)
	r.metrics.RecordExternalCall(r.getDependency(), operationSetGet, r.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper SetAndGet")
	}

	return nil
}

func (r Redis) Lock(ctx context.Context, key string, exp time.Duration) (bool, error) {
	start := time.Now()
	exists, err := r.redis.Lock(ctx, key, exp)
	r.metrics.RecordExternalCall(r.getDependency(), operationLock, r.getStatus(err), time.Since(start))
	if err != nil {
		return false, common.Wrap(err, "Wrapper Lock")
	}

	return exists, nil
}

func (r Redis) Unlock(ctx context.Context, key string) error {
	start := time.Now()
	err := r.redis.Unlock(ctx, key)
	r.metrics.RecordExternalCall(r.getDependency(), operationUnlock, r.getStatus(err), time.Since(start))
	if err != nil {
		return common.Wrap(err, "Wrapper Unlock")
	}

	return nil
}

func (r Redis) getDependency() string {
	return metrics.DependencyRedis
}

func (r Redis) getStatus(err error) string {
	if err != nil {
		return metrics.StatusError
	}

	return metrics.StatusSuccess
}

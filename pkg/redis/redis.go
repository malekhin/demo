package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/net/context"
)

//go:generate mockery --name=IRedis
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
	client *redis.Client
	config Config
}

func NewRedis(config Config) IRedis {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Username: config.Username,
		Password: config.Password,
	})

	return &Redis{client: client, config: config}
}

func (r Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%w: json.Marshal", err)
	}

	err = r.client.Set(ctx, r.prefix(key), val, expiration).Err()
	if err != nil {
		return fmt.Errorf("%w: Set", err)
	}

	return nil
}

func (r Redis) GetInt(ctx context.Context, key string) (int, bool, error) {
	val, exists, err := r.get(ctx, key)
	if !exists || err != nil {
		return 0, exists, err
	}

	v, ok := val.(float64)
	if !ok {
		return 0, true, fmt.Errorf("cant convert %s to int", val)
	}

	return int(v), true, nil
}

func (r Redis) GetFloat(ctx context.Context, key string) (float64, bool, error) {
	val, exists, err := r.get(ctx, key)
	if !exists || err != nil {
		return 0, exists, err
	}

	v, ok := val.(float64)
	if !ok {
		return 0, true, fmt.Errorf("cant convert %s to int", val)
	}

	return v, true, nil
}

func (r Redis) GetString(ctx context.Context, key string) (string, bool, error) {
	val, exists, err := r.get(ctx, key)
	if !exists || err != nil {
		return "", exists, err
	}

	v, ok := val.(string)
	if !ok {
		return "", true, fmt.Errorf("cant convert %s to int", val)
	}

	return v, true, nil
}

func (r Redis) GetStruct(ctx context.Context, key string, output interface{}) (bool, error) {
	val, exists, err := r.get(ctx, key)
	if !exists || err != nil {
		return exists, err
	}

	err = mapstructure.WeakDecode(val, &output)
	if err != nil {
		return true, fmt.Errorf("%w: mapstructure.WeakDecode", err)
	}

	return true, nil
}

func (r Redis) SetAndGet(ctx context.Context, result interface{}, f func() (interface{}, error), key string, expires time.Duration) error {
	value, err := r.client.Get(ctx, r.prefix(key)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("%w: Get", err)
	}

	if value != "" {
		err = json.Unmarshal([]byte(value), result)
		if err != nil {
			return fmt.Errorf("%w: Unmarshal", err)
		}
		return nil
	}

	result2, err := f()
	if err != nil {
		return err
	}

	raw, err := json.Marshal(result2)
	if err != nil {
		return fmt.Errorf("%w: Marshal", err)
	}

	err = r.client.Set(ctx, r.prefix(key), string(raw), expires).Err()
	if err != nil {
		return fmt.Errorf("%w: Set", err)
	}

	err = json.Unmarshal(raw, result)
	if err != nil {
		return fmt.Errorf("%w: Marshal", err)
	}

	return nil
}

func (r Redis) get(ctx context.Context, key string) (interface{}, bool, error) {
	val, err := r.client.Get(ctx, r.prefix(key)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("%w: Get", err)
	}

	var value interface{}
	err = json.Unmarshal([]byte(val), &value)
	if err != nil {
		return nil, true, fmt.Errorf("%w: json.Unmarshal", err)
	}

	return value, true, nil
}

func (r Redis) prefix(key string) string {
	return r.config.Prefix + key
}

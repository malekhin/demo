package scheduler

import (
	"context"
	"demo/pkg/redis"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

const redisLockPrefix = "WorkerJobExecutionLock_"
const redisLockDuration = 3 * time.Second

type RunFunc func(ctx context.Context) error

type Cron struct {
	cron   *cron.Cron
	logger *zap.Logger
	redis  redis.IRedis
}

func NewCron(logger *zap.Logger, redis redis.IRedis, cron *cron.Cron) *Cron {
	return &Cron{logger: logger, redis: redis, cron: cron}
}

func (c *Cron) Start() {
	c.cron.Start()
}

func (c *Cron) Stop() {
	c.cron.Stop()
}

func (c *Cron) AddJob(ctx context.Context, spec string, f RunFunc) {
	_, err := c.run(ctx, spec, f, &sync.Mutex{})
	if err != nil {
		panic(err)
	}
}

func (c *Cron) run(ctx context.Context, spec string, f RunFunc, mu *sync.Mutex) (cron.EntryID, error) {
	reflexFuncName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	entryId, err := c.cron.AddFunc(
		spec, func() {
			mu.Lock()
			defer mu.Unlock()

			ok, err := c.redis.Lock(ctx, redisLockPrefix+reflexFuncName, redisLockDuration)
			if err != nil {
				c.logger.Error(fmt.Sprintf("Job is failed: %s", reflexFuncName), zap.Error(err))
				return
			}
			if !ok {
				return
			}

			c.logger.Info(fmt.Sprintf("Job is started: %s", reflexFuncName))
			err = f(ctx)
			if err != nil {
				c.logger.Error(fmt.Sprintf("Job is failed: %s", reflexFuncName), zap.Error(err))
			} else {
				c.logger.Info(fmt.Sprintf("Job is success: %s", reflexFuncName))
			}
		},
	)
	if err != nil {
		return 0, fmt.Errorf("cron.AddFunc: %w", err)
	}

	return entryId, nil
}

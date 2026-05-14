package app

import (
	"context"
	"demo/internal/common/config"
	"demo/internal/common/wrapper/redis"
	"demo/pkg/scheduler"

	"github.com/robfig/cron/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type WorkerParams struct {
	fx.In
	Lc     fx.Lifecycle
	Config *config.Config
	Logger *zap.Logger
	Redis  redis.IRedis
}

var WorkerModule = fx.Module("worker",
	fx.Provide(NewWorker),
)

type Worker struct{}

type Run interface {
	Run(ctx context.Context) error
}

func NewWorker(p WorkerParams) *Worker {
	c := scheduler.NewCron(p.Logger, p.Redis, cron.New())

	p.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				//c.AddJob(ctx, p.Config.Cron.ImportSubagent.StartTime, p.ImportSubagentJob.Run)
				c.Start()
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			c.Stop()

			return nil
		},
	})

	return &Worker{}
}

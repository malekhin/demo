package app

import (
	"context"
	"demo/internal/common/auth"
	"demo/internal/common/config"
	"demo/internal/domain/sk"
	"demo/pkg/clickhouse"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var Module = fx.Module("app",
	fx.Provide(NewErrorMiddleware),
)

var ConfigModule = fx.Module("config",
	fx.Provide(config.New),
)

var ClickhouseModule = fx.Module("clickhouse",
	fx.Provide(clickhouse.New),
)

func Modules() []fx.Option {
	return []fx.Option{
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger.Named("FX")}
		}),

		Module,

		ConfigModule,
		WrapperModule,
		DatabaseModule,
		ClickhouseModule,

		LoggerModule,
		RedisModule,
		MetricsModule,

		HttpServerModule,
		WorkerModule,

		auth.Module,

		sk.Module,
	}
}

func NewApp() *fx.App {
	modules := Modules()
	modules = append(modules, fx.Invoke(func(*http.Server) {}))
	modules = append(modules, fx.Invoke(func(worker *Worker) {}))
	return fx.New(
		modules...,
	)
}

func Populate(targets ...interface{}) {
	PopulateWith(nil, targets...)
}

func PopulateWith(option fx.Option, targets ...interface{}) {
	modules := Modules()
	modules = append(modules, fx.Populate(targets...))
	if option != nil {
		modules = append(modules, option)
	}
	app := fx.New(
		modules...,
	)
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	defer func(app *fx.App, ctx context.Context) {
		_ = app.Stop(ctx)
	}(app, context.Background())
}

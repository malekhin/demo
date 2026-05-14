package app

import (
	"demo/internal/common/wrapper/clickhouse"
	"demo/internal/common/wrapper/database"
	"demo/internal/common/wrapper/redis"

	"go.uber.org/fx"
)

var WrapperModule = fx.Module("wrapper",
	fx.Provide(database.New),
	fx.Provide(database.NewTx),
	fx.Provide(redis.NewRedis),
	fx.Provide(clickhouse.New),
)

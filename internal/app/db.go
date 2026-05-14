package app

import (
	"demo/pkg/database"

	"go.uber.org/fx"
)

var DatabaseModule = fx.Module("database",
	fx.Provide(database.New),
	fx.Provide(database.NewTx),
)

package sk

import (
	"demo/internal/domain/sk/http/admin"
	"demo/internal/domain/sk/http/public"
	"demo/internal/domain/sk/service"
	"demo/internal/domain/sk/storage"

	"go.uber.org/fx"
)

var Module = fx.Module("sk",
	fx.Provide(admin.NewSkHandlers),
	fx.Provide(public.NewSkHandlers),
	fx.Provide(service.NewSk),
	fx.Provide(storage.NewSk),
)

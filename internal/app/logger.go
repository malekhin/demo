package app

import (
	"demo/internal/common/config"
	"fmt"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	if config.GetEnvironment() == config.EnvironmentTest {
		return zap.NewNop(), nil
	}

	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Encoding = cfg.Logger.Encoding

	loggerLevel, err := zap.ParseAtomicLevel(cfg.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("%w: zap.ParseAtomicLevel", err)
	}
	loggerCfg.Level = loggerLevel

	l, err := loggerCfg.Build()
	if err != nil {
		return nil, fmt.Errorf("%w: Build", err)
	}

	return l, nil
}

var LoggerModule = fx.Module("Logger",
	fx.Provide(NewLogger),
)

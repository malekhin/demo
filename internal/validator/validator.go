package validator

import (
	"demo/internal/common/wrapper/database"
	"demo/internal/validator/custom"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Validators struct {
	custom custom.Validators
}

func RegisterCustomValidators(db database.Database, logger *zap.Logger) Validators {
	validators := Validators{
		custom: custom.NewCustomValidators(db, logger),
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("sk", validators.custom.IsRowExists(custom.SkTable, custom.IdColumn), true); err != nil {
			logger.Error("CustomValidator Error", zap.Error(err))
		}
		if err := v.RegisterValidation("requiredTime", validators.custom.RequiredTime()); err != nil {
			logger.Error("CustomValidator Error", zap.Error(err))
		}
		if err := v.RegisterValidation("tariffSort", validators.custom.IsRowNotExists(custom.TariffTable, custom.SortColumn)); err != nil {
			logger.Error("CustomValidator Error", zap.Error(err))
		}
		if err := v.RegisterValidation("tomorrow", validators.custom.IsTomorrow()); err != nil {
			logger.Error("CustomValidator Error", zap.Error(err))
		}
		if err := v.RegisterValidation("dateOnly", validators.custom.DateOnly()); err != nil {
			logger.Error("CustomValidator Error", zap.Error(err))
		}
	}

	return validators
}

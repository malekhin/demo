package custom

import (
	"context"
	"demo/internal/common"
	"demo/internal/common/wrapper/database"
	"time"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type table string
type column string

var (
	SkTable             table = "sk"
	TariffTable         table = "tariff"
	TariffSubagentTable table = "tariff_subagent"
	ProductTable        table = "product"
	TagTable            table = "tag"
)
var (
	IdColumn   column = "id"
	TypeColumn column = "type"
	SortColumn column = "sort"
)

type Validators struct {
	db     database.Database
	logger *zap.Logger
}

func NewCustomValidators(db database.Database, logger *zap.Logger) Validators {
	validators := Validators{db: db, logger: logger}
	return validators
}

// IsRowExists Проверяет на существование записи в таблице
func (v Validators) IsRowExists(table table, column column) validator.Func {
	return func(fl validator.FieldLevel) bool {
		var intNilPtr *int
		if fl.Field().Interface() == intNilPtr {
			return true
		}

		if fl.Field().CanInt() && fl.Field().Int() == 0 {
			return false
		}

		query := `
			select count(*) from ` + string(table) + ` where ` + string(column) + ` = ?
		`

		var count int
		err := v.db.Db().Row(context.Background(), &count, query, fl.Field().Interface())
		if err != nil {
			v.logger.Error("CustomValidator Error", zap.Error(err))
		}

		return count >= 1
	}
}

// IsRowNotExists Проверяет на отсутствие записи в таблице
func (v Validators) IsRowNotExists(table table, column column) validator.Func {
	return func(fl validator.FieldLevel) bool {
		var count int
		err := v.db.Db().Row(context.Background(), &count, `
			select count(*) from `+string(table)+` where `+string(column)+` = ?
		`, fl.Field().Interface())
		if err != nil {
			v.logger.Error("CustomValidator Error", zap.Error(err))
		}

		return count == 0
	}
}

// RequiredTime Условие обязательности для типа time.Time
func (v Validators) RequiredTime() validator.Func {
	return func(fl validator.FieldLevel) bool {
		t := fl.Field().Interface().(time.Time)
		return !t.IsZero()
	}
}

// IsTomorrow Проверяет что дата не раньше следующего дня
func (v Validators) IsTomorrow() validator.Func {
	return func(fl validator.FieldLevel) bool {
		t := fl.Field().Interface().(common.DateTime)

		now := time.Now()
		tomorrow := time.Date(
			now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local,
		).Add(24 * time.Hour).Add(-time.Second)

		return t.Valid && t.Time.After(tomorrow)
	}
}

func (v Validators) DateOnly() validator.Func {
	return func(fl validator.FieldLevel) bool {
		t := fl.Field().Interface().(string)
		_, err := time.Parse(time.DateOnly, t)

		return err == nil
	}
}

package app

import (
	"demo/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ErrorMiddleware gin.HandlerFunc

type ErrorMiddlewareParams struct {
	fx.In
	Logger *zap.Logger
}

func NewErrorMiddleware(params ErrorMiddlewareParams) ErrorMiddleware {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		var status int
		var errList []util.Error
		var sliceValidationErrors binding.SliceValidationError

		for _, err := range c.Errors {
			if errors.As(err, &sliceValidationErrors) {
				for _, sliceValidationError := range sliceValidationErrors {
					status, errList = validate(params, sliceValidationError)
				}
			} else {
				status, errList = validate(params, err)
			}
		}

		c.JSON(status, util.ResponseError{
			Status: util.StatusError,
			Errors: errList,
		})
	}
}

func validate(params ErrorMiddlewareParams, err error) (int, []util.Error) {
	var status int
	var errList []util.Error

	var validationErrors validator.ValidationErrors
	var numError *strconv.NumError
	var unmarshalTypeError *json.UnmarshalTypeError
	var timeError *time.ParseError
	var ginError *gin.Error
	var jsonSyntaxError *json.SyntaxError
	var badRequest util.BadRequestError

	if errors.As(err, &validationErrors) {
		for _, validationError := range validationErrors {
			errList = append(errList, util.Error{Message: getErrorForTag(validationError)})
		}
		status = http.StatusBadRequest
	} else if errors.As(err, &badRequest) {
		errList = append(errList, util.Error{Message: err.Error()})
		status = http.StatusBadRequest
	} else if errors.As(err, &numError) {
		errList = append(errList, util.Error{Message: fmt.Sprintf("invalid input format: number expected (got '%s')", numError.Num)})
		status = http.StatusBadRequest
	} else if errors.As(err, &unmarshalTypeError) {
		errList = append(errList, util.Error{Message: fmt.Sprintf("invalid input format: %s", unmarshalTypeError.Field)})
		status = http.StatusBadRequest
	} else if errors.As(err, &timeError) {
		errList = append(errList, util.Error{Message: fmt.Sprintf("invalid time format: %s (expected: %s)", timeError.Value, timeError.Layout)})
		status = http.StatusBadRequest
	} else if errors.As(err, &ginError) && strings.HasPrefix(err.Error(), "uuid:") {
		errList = append(errList, util.Error{Message: ginError.Error()})
		status = http.StatusBadRequest
	} else if errors.As(err, &jsonSyntaxError) {
		errList = append(errList, util.Error{Message: jsonSyntaxError.Error()})
		status = http.StatusBadRequest
	} else if errors.Is(err, io.EOF) {
		errList = append(errList, util.Error{Message: "request body is empty"})
		status = http.StatusBadRequest
	} else {
		params.Logger.Error(util.InternalError, zap.Error(err))
		errList = append(errList, util.Error{Message: util.InternalError})
		status = http.StatusInternalServerError
	}

	return status, errList
}

// Кастомные ошибки валидатора
func getErrorForTag(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "sk":
		return fmt.Sprintf("sk id %d is not found", fieldError.Value())
	case "tariff":
		return fmt.Sprintf("tariff id %d is not found", fieldError.Value())
	case "product":
		return fmt.Sprintf("product type '%s' is not found", fieldError.Value())
	case "tag":
		return fmt.Sprintf("tag id %d is not found", fieldError.Value())
	case "customError":
		return fieldError.Param()

	default:
		return fieldError.Error()
	}
}

func RegisterValidatorTagNameFunc() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "" {
				name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
			}
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

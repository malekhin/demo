package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var StatusOk = "OK"
var StatusError = "ERROR"

const InternalError = "internal server error"

type Buffer interface {
	Get() []byte
}

type BadRequestError struct {
	Err error
}

func (e BadRequestError) Error() string {
	return e.Err.Error()
}

func NewBadRequest(err error) BadRequestError {
	return BadRequestError{Err: err}
}

type ResponseOkPaginate[T any] struct {
	Status     string `json:"status"`
	TotalCount int    `json:"totalCount"`
	Data       T      `json:"data,omitempty"`
}

type ResponseOk[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data,omitempty"`
}

type ResponseError struct {
	Status string  `json:"status"`
	Errors []Error `json:"errors,omitempty"`
}

type Error struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

// SetStatusOk Успешный ответ
func SetStatusOk[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, ResponseOk[T]{
		Status: StatusOk,
		Data:   data,
	})
}

// SetStatusOkPaginate Успешный ответ c пагинацией
func SetStatusOkPaginate[T any](c *gin.Context, data T, count int) {
	c.JSON(http.StatusOK, ResponseOkPaginate[T]{
		Status:     StatusOk,
		TotalCount: count,
		Data:       data,
	})
}

// SetStatusError В зависимости от типа ошибки middleware сформирует ответ и код ответа
func SetStatusError(c *gin.Context, err error) {
	_ = c.Error(err)
}

func DownloadFile(c *gin.Context, filename string, buffer Buffer) {
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Get())
}

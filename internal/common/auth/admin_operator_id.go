package auth

import (
	"context"
	"demo/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

const (
	OperatorIdHeaderContextKey = "X-Operator-Id"
)

// AdminOperatorAuthMiddleware получает идентификатор агента и устанавливает его в context
type AdminOperatorAuthMiddleware gin.HandlerFunc

func NewAdminOperatorAuthMiddleware() AdminOperatorAuthMiddleware {
	return func(c *gin.Context) {
		operatorIdHeader := c.GetHeader(OperatorIdHeaderContextKey)

		operatorUuid, err := uuid.FromString(operatorIdHeader)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseError{
				Status: util.StatusError,
				Errors: []util.Error{{Message: "X-Operator-Id is uuid expected"}},
			})
			return
		}

		c.Set(OperatorIdHeaderContextKey, operatorUuid)
	}
}

func GetOperatorId(c context.Context) uuid.UUID {
	return c.Value(OperatorIdHeaderContextKey).(uuid.UUID)
}

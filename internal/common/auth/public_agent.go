package auth

import (
	"context"
	"demo/internal/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	agentIdHeader       = "X-Agent-Id"
	agentIdContextKey   = "agent_id"
	agentTypeContextKey = "agent_type"
)

type AgentType string

const (
	Agent    AgentType = "AGENT"
	SubAgent AgentType = "SUB_AGENT"
)

// PublicAgentAuthMiddleware получает идентификатор агента и устанавливает его в context
type PublicAgentAuthMiddleware gin.HandlerFunc

func NewPublicAgentAuthMiddleware(logger *zap.Logger) PublicAgentAuthMiddleware {
	logger = logger.Named("AGENT_AUTH_MIDDLEWARE")

	return func(c *gin.Context) {
		agentIdHeader := c.GetHeader(agentIdHeader)

		if agentIdHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.ResponseError{
				Status: util.StatusError,
				Errors: []util.Error{{Message: "X-Agent-Id is empty"}},
			})
			return
		}

		agentId, err := strconv.Atoi(agentIdHeader)
		if err != nil || agentId <= 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, util.ResponseError{
				Status: util.StatusError,
				Errors: []util.Error{{Message: "X-Agent-Id is not int"}},
			})
			return
		}

		// OSAGO-7689 X-Agent-Id должен быть был всегда указан для публичных эндпоинтов.
		// При этом на данном этапе существование агента проверять не будем.

		hasAgent := true
		hasSubAgent := false

		if hasAgent {
			c.Set(agentTypeContextKey, Agent)
		} else if hasSubAgent {
			c.Set(agentTypeContextKey, SubAgent)
		} else {
			logger.Warn("not found (sub-)agent",
				zap.Int("x-agent-id", agentId),
			)
			return
		}

		c.Set(agentIdContextKey, agentId)
	}
}

func GetAgentId(c context.Context) int {
	return c.Value(agentIdContextKey).(int)
}

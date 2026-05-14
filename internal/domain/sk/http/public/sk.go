package public

import (
	"demo/internal/domain/sk/model"
	"demo/internal/domain/sk/service"
	"demo/internal/util"

	"github.com/gin-gonic/gin"
)

type SkHandlers struct {
	skService *service.Sk
}

func NewSkHandlers(skService *service.Sk) *SkHandlers {
	return &SkHandlers{skService: skService}
}

// SkList
// @Tags public,sk
// @Summary Список СК
// @Router /public/v1/sk [get]
// @Param q query model.SkFilter true "Фильтр"
// @Success 200 {object} util.ResponseOkPaginate[[]model.SkItem]
func (h *SkHandlers) SkList(c *gin.Context) {
	var filter model.SkFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		util.SetStatusError(c, err)
		return
	}

	list, count, err := h.skService.List(c, filter)
	if err != nil {
		util.SetStatusError(c, err)
		return
	}

	util.SetStatusOkPaginate(c, list, count)
}

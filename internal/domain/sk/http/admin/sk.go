package admin

import (
	"demo/internal/domain/sk/model"
	"demo/internal/domain/sk/service"
	"demo/internal/util"
	"go/types"

	"github.com/gin-gonic/gin"
)

type SkHandlers struct {
	skService *service.Sk
}

func NewSkHandlers(skService *service.Sk) *SkHandlers {
	return &SkHandlers{skService: skService}
}

// SkList
// @Tags admin,sk
// @Summary Список СК
// @Router /admin/v1/sk [get]
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

// SkAdd
// @Tags admin,sk
// @Summary Добавление СК
// @Router /admin/v1/sk [post]
// @Param request body model.SkAdd true "body"
// @Success 200 {object} util.ResponseOk[types.Nil]
func (h *SkHandlers) SkAdd(c *gin.Context) {
	var request model.SkAdd
	if err := c.ShouldBindJSON(&request); err != nil {
		util.SetStatusError(c, err)
		return
	}

	err := h.skService.Add(c, request)
	if err != nil {
		util.SetStatusError(c, err)
		return
	}

	util.SetStatusOk(c, types.Nil{})
}

// SkEdit
// @Tags admin,sk
// @Summary Редактирование СК
// @Router /admin/v1/sk/{id} [post]
// @Param id path integer true "ID СК"
// @Param request body model.SkEdit true "body"
// @Success 200 {object} util.ResponseOk[types.Nil]
func (h *SkHandlers) SkEdit(c *gin.Context) {
	var uri model.SkUri
	if err := c.ShouldBindUri(&uri); err != nil {
		util.SetStatusError(c, err)
		return
	}

	var request model.SkEdit
	if err := c.ShouldBindJSON(&request); err != nil {
		util.SetStatusError(c, err)
		return
	}

	err := h.skService.Edit(c, uri.Id, request)
	if err != nil {
		util.SetStatusError(c, err)
		return
	}

	util.SetStatusOk(c, types.Nil{})
}

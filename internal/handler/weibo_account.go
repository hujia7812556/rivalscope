package handler

import (
	"github.com/gin-gonic/gin"

	"rivalscope/internal/dto"
	"rivalscope/internal/service"
)

// WeiboAccountHandler 微博账号管理接口。
type WeiboAccountHandler struct {
	accountService *service.WeiboAccountService
}

// NewWeiboAccountHandler 创建微博账号 handler。
func NewWeiboAccountHandler(accountService *service.WeiboAccountService) *WeiboAccountHandler {
	return &WeiboAccountHandler{accountService: accountService}
}

// GetList 处理 GET /api/v1/weibo/accounts,返回全部有效账号。
func (h *WeiboAccountHandler) GetList(c *gin.Context) {
	items, err := h.accountService.List()
	if err != nil {
		dto.Fail(c, err)
		return
	}
	if items == nil {
		items = []dto.WeiboAccountItem{}
	}
	dto.Success(c, items)
}

// Create 处理 POST /api/v1/weibo/accounts。
func (h *WeiboAccountHandler) Create(c *gin.Context) {
	var req dto.WeiboAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Fail(c, dto.ErrBadRequest)
		return
	}
	if err := h.accountService.Create(req); err != nil {
		dto.Fail(c, err)
		return
	}
	dto.Success(c, nil)
}

// Update 处理 PATCH /api/v1/weibo/accounts/:id。
func (h *WeiboAccountHandler) Update(c *gin.Context) {
	id, ok := parseId(c)
	if !ok {
		dto.Fail(c, dto.ErrBadRequest)
		return
	}
	var req dto.WeiboAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Fail(c, dto.ErrBadRequest)
		return
	}
	if err := h.accountService.Update(id, req); err != nil {
		dto.Fail(c, err)
		return
	}
	dto.Success(c, nil)
}

// Delete 处理 DELETE /api/v1/weibo/accounts/:id。
func (h *WeiboAccountHandler) Delete(c *gin.Context) {
	id, ok := parseId(c)
	if !ok {
		dto.Fail(c, dto.ErrBadRequest)
		return
	}
	if err := h.accountService.Delete(id); err != nil {
		dto.Fail(c, err)
		return
	}
	dto.Success(c, nil)
}

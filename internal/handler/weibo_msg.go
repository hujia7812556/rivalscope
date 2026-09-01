package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"rivalscope/internal/dto"
	"rivalscope/internal/repo"
	"rivalscope/internal/service"
)

// WeiboMsgHandler 微博文列表接口(服务端分页 + 排序)。
type WeiboMsgHandler struct {
	msgService *service.WeiboMsgService
}

// NewWeiboMsgHandler 创建微博文 handler。
func NewWeiboMsgHandler(msgService *service.WeiboMsgService) *WeiboMsgHandler {
	return &WeiboMsgHandler{msgService: msgService}
}

// GetList 处理 GET /api/v1/weibo/messages。
// 查询参数:account_id、start_date、end_date、sort_field、sort_order、page、page_size。
func (h *WeiboMsgHandler) GetList(c *gin.Context) {
	accountId, _ := strconv.Atoi(c.DefaultQuery("account_id", "0"))
	start, end, err := parseDateRange(c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		dto.Fail(c, dto.ErrBadRequest)
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}

	data, err := h.msgService.ListPage(repo.WeiboMsgQuery{
		AccountId: accountId,
		Start:     start,
		End:       end,
		SortField: c.DefaultQuery("sort_field", "id"),
		SortOrder: c.DefaultQuery("sort_order", "asc"),
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		dto.Fail(c, err)
		return
	}
	dto.Success(c, data)
}

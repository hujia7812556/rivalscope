package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"rivalscope/internal/dto"
	"rivalscope/internal/repo"
	"rivalscope/internal/service"
)

// WeiboStatisticsHandler 微博统计接口。
type WeiboStatisticsHandler struct {
	statsService *service.WeiboStatisticsService
}

// NewWeiboStatisticsHandler 创建统计 handler。
func NewWeiboStatisticsHandler(statsService *service.WeiboStatisticsService) *WeiboStatisticsHandler {
	return &WeiboStatisticsHandler{statsService: statsService}
}

// GetList 处理 GET /api/v1/weibo/statistics。
// 查询参数:account_id、start_date、end_date;按账号聚合返回,不分页。
func (h *WeiboStatisticsHandler) GetList(c *gin.Context) {
	accountId, _ := strconv.Atoi(c.DefaultQuery("account_id", "0"))
	start, end, err := parseDateRange(c.Query("start_date"), c.Query("end_date"))
	if err != nil {
		dto.Fail(c, dto.ErrBadRequest)
		return
	}

	items, err := h.statsService.List(repo.WeiboStatisticsQuery{
		AccountId: accountId,
		Start:     start,
		End:       end,
	})
	if err != nil {
		dto.Fail(c, err)
		return
	}
	if items == nil {
		items = []dto.WeiboStatisticsItem{}
	}
	dto.Success(c, items)
}

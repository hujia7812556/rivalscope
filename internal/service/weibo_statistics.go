package service

import (
	"fmt"

	"rivalscope/internal/dto"
	"rivalscope/internal/repo"
)

// WeiboStatisticsService 微博统计业务。
type WeiboStatisticsService struct {
	statsRepo *repo.WeiboStatisticsRepo
}

// NewWeiboStatisticsService 创建统计 service。
func NewWeiboStatisticsService(statsRepo *repo.WeiboStatisticsRepo) *WeiboStatisticsService {
	return &WeiboStatisticsService{statsRepo: statsRepo}
}

// List 按条件聚合统计并转为响应 DTO。
// Sum 输出整数串、Avg 保留 4 位小数(与旧接口格式保持一致,便于前端直接导出)。
func (s *WeiboStatisticsService) List(q repo.WeiboStatisticsQuery) ([]dto.WeiboStatisticsItem, error) {
	list, err := s.statsRepo.List(q)
	if err != nil {
		return nil, err
	}

	items := make([]dto.WeiboStatisticsItem, 0, len(list))
	for _, st := range list {
		items = append(items, dto.WeiboStatisticsItem{
			Name:       st.Name,
			ID:         st.ID,
			Attention:  st.Attention,
			Fans:       st.Fans,
			Feed:       st.Feed,
			CrawlTime:  formatTime(st.UpdateTime, "2006-01-02 15:04:05"),
			Count:      st.Count,
			ForwardSum: fmt.Sprintf("%d", st.ForwardSum),
			CommentSum: fmt.Sprintf("%d", st.CommentSum),
			LikeSum:    fmt.Sprintf("%d", st.LikeSum),
			ForwardAvg: fmt.Sprintf("%.4f", st.ForwardAvg),
			CommentAvg: fmt.Sprintf("%.4f", st.CommentAvg),
			LikeAvg:    fmt.Sprintf("%.4f", st.LikeAvg),
			ForwardMax: st.ForwardMax,
			CommentMax: st.CommentMax,
			LikeMax:    st.LikeMax,
			ForwardMin: st.ForwardMin,
			CommentMin: st.CommentMin,
			LikeMin:    st.LikeMin,
		})
	}
	return items, nil
}

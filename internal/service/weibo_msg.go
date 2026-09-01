package service

import (
	"fmt"

	"rivalscope/internal/dto"
	"rivalscope/internal/repo"
)

// WeiboMsgService 微博文列表业务。
type WeiboMsgService struct {
	msgRepo *repo.WeiboMsgRepo
}

// NewWeiboMsgService 创建微博文 service。
func NewWeiboMsgService(msgRepo *repo.WeiboMsgRepo) *WeiboMsgService {
	return &WeiboMsgService{msgRepo: msgRepo}
}

// ListPage 分页查询微博文并转为响应 DTO。
func (s *WeiboMsgService) ListPage(q repo.WeiboMsgQuery) (*dto.WeiboMsgListData, error) {
	q.SortField, q.SortOrder = repo.NormalizeWeiboMsgSort(q.SortField, q.SortOrder)
	list, total, err := s.msgRepo.ListPage(q)
	if err != nil {
		return nil, err
	}

	items := make([]dto.WeiboMsgItem, 0, len(list))
	for _, m := range list {
		items = append(items, dto.WeiboMsgItem{
			ID:          m.ID,
			AccountName: m.Account.Name,
			Mid:         m.Mid,
			AccountId:   m.AccountID,
			Forward:     m.Forward,
			Comment:     m.Comment,
			Like:        m.Like,
			Pubtime:     formatTime(m.Pubtime, "2006-01-02"),
			CrawlTime:   formatTime(m.UpdateTime, "2006-01-02 15:04:05"),
			Url:         fmt.Sprintf("https://m.weibo.cn/detail/%s", m.Mid),
		})
	}
	return &dto.WeiboMsgListData{
		Total:       total,
		CurrentPage: q.Page,
		PageSize:    q.PageSize,
		MsgList:     items,
	}, nil
}

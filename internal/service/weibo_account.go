package service

import (
	"fmt"
	"time"

	"rivalscope/internal/dto"
	"rivalscope/internal/model"
	"rivalscope/internal/repo"
)

// formatTime 格式化时间;零值(数据库 NULL 扫描结果)返回空串。
// 线上 weibo_account/weibo_feed 的 update_time 列可为 NULL。
func formatTime(t time.Time, layout string) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(layout)
}

// WeiboAccountService 微博账号管理业务。
type WeiboAccountService struct {
	accountRepo *repo.WeiboAccountRepo
}

// NewWeiboAccountService 创建微博账号 service。
func NewWeiboAccountService(accountRepo *repo.WeiboAccountRepo) *WeiboAccountService {
	return &WeiboAccountService{accountRepo: accountRepo}
}

// toItem 模型转响应 DTO:格式化时间并拼出主页链接。
func toWeiboAccountItem(a model.WeiboAccount) dto.WeiboAccountItem {
	return dto.WeiboAccountItem{
		ID:         a.ID,
		Name:       a.Name,
		Uid:        a.Uid,
		Status:     a.Status,
		Attention:  a.Attention,
		Fans:       a.Fans,
		Feed:       a.Feed,
		CreateTime: formatTime(a.CreateTime, "2006-01-02 15:04:05"),
		UpdateTime: formatTime(a.UpdateTime, "2006-01-02 15:04:05"),
		Url:        fmt.Sprintf("https://m.weibo.cn/u/%s", a.Uid),
		CrawlTime:  formatTime(a.UpdateTime, "2006-01-02"),
	}
}

// List 返回全部有效账号。
func (s *WeiboAccountService) List() ([]dto.WeiboAccountItem, error) {
	list, err := s.accountRepo.List()
	if err != nil {
		return nil, err
	}
	items := make([]dto.WeiboAccountItem, 0, len(list))
	for _, a := range list {
		items = append(items, toWeiboAccountItem(a))
	}
	return items, nil
}

// Create 新增账号;uid 重复时返回 2001 业务错误。
func (s *WeiboAccountService) Create(req dto.WeiboAccountRequest) error {
	exist, err := s.accountRepo.GetByUid(req.Uid)
	if err != nil {
		return dto.ErrInternalServerError
	}
	if exist != nil {
		return dto.ErrWeiboUidExists
	}
	return s.accountRepo.Create(&model.WeiboAccount{
		Name: req.Name,
		Uid:  req.Uid,
	})
}

// Update 更新账号名称与 uid。
func (s *WeiboAccountService) Update(id int, req dto.WeiboAccountRequest) error {
	a, err := s.accountRepo.GetByID(id)
	if err != nil {
		return err
	}
	// uid 改动时校验是否与其他账号冲突
	if a.Uid != req.Uid {
		other, err := s.accountRepo.GetByUid(req.Uid)
		if err != nil {
			return dto.ErrInternalServerError
		}
		if other != nil && other.ID != id {
			return dto.ErrWeiboUidExists
		}
	}
	a.Name = req.Name
	a.Uid = req.Uid
	return s.accountRepo.Update(a)
}

// Delete 删除账号(物理删除)。
func (s *WeiboAccountService) Delete(id int) error {
	if _, err := s.accountRepo.GetByID(id); err != nil {
		return err
	}
	return s.accountRepo.Delete(id)
}

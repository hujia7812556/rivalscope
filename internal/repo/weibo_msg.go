package repo

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"rivalscope/internal/model"
)

// WeiboMsgQuery 微博文列表查询条件(由 handler 解析后传入)。
type WeiboMsgQuery struct {
	AccountId int       // >0 时按账号过滤
	Start     time.Time // 非零时过滤 pubtime >= Start
	End       time.Time // 非零时过滤 pubtime < End(已含 +1 天处理)
	SortField string    // 排序字段(传入 repo 前已经过白名单校验)
	SortOrder string    // asc / desc
	Page      int
	PageSize  int
}

// weiboMsgSortFields 允许的服务端排序字段白名单,防止 SQL 列名注入。
var weiboMsgSortFields = map[string]bool{
	"id": true, "mid": true, "forward": true, "comment": true, "like": true, "pubtime": true,
}

// NormalizeWeiboMsgSort 校验排序字段与方向,非法值回落为 id asc。
func NormalizeWeiboMsgSort(field, order string) (string, string) {
	if !weiboMsgSortFields[field] {
		field = "id"
	}
	if order != "desc" {
		order = "asc"
	}
	return field, order
}

// WeiboMsgRepo 微博文数据访问。
type WeiboMsgRepo struct {
	db *gorm.DB
}

// NewWeiboMsgRepo 创建微博文 repo。
func NewWeiboMsgRepo(db *gorm.DB) *WeiboMsgRepo { return &WeiboMsgRepo{db: db} }

// ListPage 分页查询微博文;Preload Account 用于展示账号名称。
// 返回当前页数据与总条数。
func (r *WeiboMsgRepo) ListPage(q WeiboMsgQuery) ([]model.WeiboMsg, int64, error) {
	db := r.db.Model(&model.WeiboMsg{}).Preload("Account")
	if q.AccountId > 0 {
		db = db.Where("account_id = ?", q.AccountId)
	}
	if !q.Start.IsZero() {
		db = db.Where("pubtime >= ?", q.Start.Format("2006-01-02 15:04:05"))
	}
	if !q.End.IsZero() {
		db = db.Where("pubtime < ?", q.End.Format("2006-01-02 15:04:05"))
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order(clause.OrderByColumn{
		Column: clause.Column{Name: q.SortField},
		Desc:   q.SortOrder == "desc",
	})
	if q.Page > 0 && q.PageSize > 0 {
		db = db.Offset((q.Page - 1) * q.PageSize).Limit(q.PageSize)
	}

	var list []model.WeiboMsg
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

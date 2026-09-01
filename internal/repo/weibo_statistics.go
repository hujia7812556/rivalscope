package repo

import (
	"time"

	"gorm.io/gorm"

	"rivalscope/internal/model"
)

// WeiboStatisticsQuery 统计查询条件。
type WeiboStatisticsQuery struct {
	AccountId int       // >0 时只统计该账号
	Start     time.Time // 非零时过滤 pubtime >= Start
	End       time.Time // 非零时过滤 pubtime < End(已含 +1 天处理)
}

// WeiboStatisticsRepo 微博统计聚合查询。
type WeiboStatisticsRepo struct {
	db *gorm.DB
}

// NewWeiboStatisticsRepo 创建统计 repo。
func NewWeiboStatisticsRepo(db *gorm.DB) *WeiboStatisticsRepo {
	return &WeiboStatisticsRepo{db: db}
}

// likeColumn 返回 like 列在当前方言下的安全引用。
// like 是保留字:PostgreSQL/SQLite 用双引号,MySQL 用反引号(双引号默认是字符串字面量)。
func likeColumn(db *gorm.DB) string {
	if db.Dialector.Name() == "mysql" {
		return "`like`"
	}
	return `"like"`
}

// List 按账号聚合统计(每个账号一行):
// count 为时间范围内的微博条数,转发/评论/点赞各自做 SUM/AVG/MAX/MIN。
// 口径说明:
//   - INNER JOIN,时间范围内没有微博的账号不会出现在结果中;
//   - attention/fans/feed/name 取自 weibo_account 当前值(非时点快照);
//   - 只统计有效账号(status=1)。
func (r *WeiboStatisticsRepo) List(q WeiboStatisticsQuery) ([]model.WeiboStatistics, error) {
	like := likeColumn(r.db)
	db := r.db.Table("weibo_feed as wf").
		Select("wa.name, wa.id, wa.attention, wa.fans, wa.feed, wa.update_time, "+
			"COUNT(wf.id) AS count, "+
			"SUM(wf.forward) AS forward_sum, SUM(wf.comment) AS comment_sum, SUM(wf."+like+") AS like_sum, "+
			"AVG(wf.forward) AS forward_avg, AVG(wf.comment) AS comment_avg, AVG(wf."+like+") AS like_avg, "+
			"MAX(wf.forward) AS forward_max, MAX(wf.comment) AS comment_max, MAX(wf."+like+") AS like_max, "+
			"MIN(wf.forward) AS forward_min, MIN(wf.comment) AS comment_min, MIN(wf."+like+") AS like_min").
		Joins("JOIN weibo_account as wa ON wa.id = wf.account_id").
		Where("wa.status = ?", model.WeiboAccountStatusValid)

	if q.AccountId > 0 {
		db = db.Where("wf.account_id = ?", q.AccountId)
	}
	if !q.Start.IsZero() {
		db = db.Where("wf.pubtime >= ?", q.Start.Format("2006-01-02 15:04:05"))
	}
	if !q.End.IsZero() {
		db = db.Where("wf.pubtime < ?", q.End.Format("2006-01-02 15:04:05"))
	}

	var list []model.WeiboStatistics
	err := db.Group("wa.id, wa.name, wa.attention, wa.fans, wa.feed, wa.update_time").
		Order("wa.id asc").Scan(&list).Error
	return list, err
}

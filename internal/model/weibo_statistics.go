package model

import "time"

// WeiboStatistics 微博统计聚合结果。
// 说明:这不是真实表,数据库中不存在 weibo_statistics,
// 它只作为 weibo_feed JOIN weibo_account 聚合查询的 Scan 承接结构。
type WeiboStatistics struct {
	ID         int       `gorm:"primaryKey"`      // 映射 weibo_account.id
	Name       string    `gorm:"column:name"`
	Attention  int       `gorm:"column:attention"`
	Fans       int       `gorm:"column:fans"`
	Feed       int       `gorm:"column:feed"`
	UpdateTime time.Time `gorm:"column:update_time"` // 账号最近抓取时间
	Count      int       `gorm:"column:count"`       // 时间范围内的微博条数
	ForwardSum int64     `gorm:"column:forward_sum"`
	CommentSum int64     `gorm:"column:comment_sum"`
	LikeSum    int64     `gorm:"column:like_sum"`
	ForwardAvg float64   `gorm:"column:forward_avg"`
	CommentAvg float64   `gorm:"column:comment_avg"`
	LikeAvg    float64   `gorm:"column:like_avg"`
	ForwardMax int       `gorm:"column:forward_max"`
	CommentMax int       `gorm:"column:comment_max"`
	LikeMax    int       `gorm:"column:like_max"`
	ForwardMin int       `gorm:"column:forward_min"`
	CommentMin int       `gorm:"column:comment_min"`
	LikeMin    int       `gorm:"column:like_min"`
}

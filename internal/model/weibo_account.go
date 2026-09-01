package model

import "time"

// 微博账号状态:1 有效 / 2 无效(无效账号不参与列表与统计)。
const (
	WeiboAccountStatusValid   = 1
	WeiboAccountStatusInvalid = 2
)

// WeiboAccount 微博账号(表 weibo_account,由爬虫写入维护)。
type WeiboAccount struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	Name      string `gorm:"type:varchar(255);not null"`
	Uid       string `gorm:"type:varchar(255);not null"` // 微博 uid
	Status    int    `gorm:"type:int unsigned;default:1;not null;comment:'1有效，2无效'"`
	Attention int    `gorm:"type:int unsigned;default:0;not null;comment:'关注数'"`
	Fans      int    `gorm:"type:int unsigned;default:0;not null;comment:'粉丝数'"`
	Feed      int    `gorm:"type:int unsigned;default:0;not null;comment:'微博数'"`
	CreateTime time.Time `gorm:"type:timestamp;default:current_timestamp();not null"`
	UpdateTime time.Time `gorm:"type:timestamp;default:current_timestamp();not null;autoUpdateTime"`
}

// TableName 指定表名(无复数)。
func (WeiboAccount) TableName() string { return "weibo_account" }

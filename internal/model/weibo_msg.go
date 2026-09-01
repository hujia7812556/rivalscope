package model

import "time"

// WeiboMsg 微博文(表 weibo_feed,由爬虫写入维护)。
type WeiboMsg struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	Mid        string    `gorm:"size:255;not null"`        // 微博消息 ID
	AccountID  int       `gorm:"not null;index"`           // 所属账号 weibo_account.id
	Forward    int       `gorm:"not null;default:0;index"` // 转发数
	Comment    int       `gorm:"not null;default:0;index"` // 评论数
	Like       int       `gorm:"not null;default:0;index"` // 点赞数
	Pubtime    time.Time `gorm:"not null"`                 // 发布时间
	CreateTime time.Time `gorm:"default:current_timestamp();not null"`
	UpdateTime time.Time `gorm:"default:current_timestamp();not null;autoUpdateTime"`

	// 关联账号(列表查询时 Preload,用于展示账号名称)
	Account WeiboAccount `gorm:"foreignKey:AccountID"`
}

// TableName 指定表名:模型名与表名不对称,历史表名是 weibo_feed。
func (WeiboMsg) TableName() string { return "weibo_feed" }

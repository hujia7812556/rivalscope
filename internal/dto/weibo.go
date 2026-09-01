package dto

// WeiboAccountRequest 新增/更新微博账号的请求体。
type WeiboAccountRequest struct {
	Name string `json:"name" binding:"required" example:"竞品A官方微博"`
	Uid  string `json:"uid" binding:"required" example:"123456"`
}

// WeiboAccountItem 微博账号列表项。
type WeiboAccountItem struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Uid        string `json:"uid"`
	Status     int    `json:"status"`
	Attention  int    `json:"attention"`
	Fans       int    `json:"fans"`
	Feed       int    `json:"feed"`
	CreateTime string `json:"create_time"` // 2006-01-02 15:04:05
	UpdateTime string `json:"update_time"` // 2006-01-02 15:04:05
	Url        string `json:"url"`         // https://m.weibo.cn/u/{uid}
	CrawlTime  string `json:"crawl_time"`  // 2006-01-02,取自账号 UpdateTime
}

// WeiboMsgItem 微博文列表项。
type WeiboMsgItem struct {
	ID          int    `json:"id"`
	AccountName string `json:"account_name"`
	Mid         string `json:"mid"`
	AccountId   int    `json:"account_id"`
	Forward     int    `json:"forward"`
	Comment     int    `json:"comment"`
	Like        int    `json:"like"`
	Pubtime     string `json:"pubtime"`    // 2006-01-02
	CrawlTime   string `json:"crawl_time"` // 2006-01-02 15:04:05
	Url         string `json:"url"`        // https://m.weibo.cn/detail/{mid}
}

// WeiboMsgListData 微博文分页列表。
type WeiboMsgListData struct {
	Total       int64          `json:"total"`
	CurrentPage int            `json:"current_page"`
	PageSize    int            `json:"page_size"`
	MsgList     []WeiboMsgItem `json:"msg_list"`
}

// WeiboStatisticsItem 微博统计聚合行(每账号一行)。
// 说明:Sum/Avg 沿用旧接口的字符串格式(整数/保留 4 位小数),前端导出直接可用。
type WeiboStatisticsItem struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	Attention   int    `json:"attention"`
	Fans        int    `json:"fans"`
	Feed        int    `json:"feed"`
	CrawlTime   string `json:"crawl_time"` // 2006-01-02 15:04:05
	Count       int    `json:"count"`
	ForwardSum  string `json:"forward_sum"`
	CommentSum  string `json:"comment_sum"`
	LikeSum     string `json:"like_sum"`
	ForwardAvg  string `json:"forward_avg"`
	CommentAvg  string `json:"comment_avg"`
	LikeAvg     string `json:"like_avg"`
	ForwardMax  int    `json:"forward_max"`
	CommentMax  int    `json:"comment_max"`
	LikeMax     int    `json:"like_max"`
	ForwardMin  int    `json:"forward_min"`
	CommentMin  int    `json:"comment_min"`
	LikeMin     int    `json:"like_min"`
}

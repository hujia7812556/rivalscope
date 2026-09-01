/**
 * 全局接口类型定义(与后端 internal/dto 对应)。
 */

/** 统一响应结构 */
export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

/** 微博账号 */
export interface WeiboAccount {
  id: number
  name: string
  uid: string
  status: number
  attention: number
  fans: number
  feed: number
  create_time: string
  update_time: string
  url: string
  crawl_time: string
}

/** 账号新增/编辑表单 */
export interface WeiboAccountForm {
  name: string
  uid: string
}

/** 微博文 */
export interface WeiboMsg {
  id: number
  account_name: string
  mid: string
  account_id: number
  forward: number
  comment: number
  like: number
  pubtime: string
  crawl_time: string
  url: string
}

/** 微博文分页数据 */
export interface WeiboMsgListData {
  total: number
  current_page: number
  page_size: number
  msg_list: WeiboMsg[]
}

/** 微博统计(按账号聚合,Sum/Avg 为字符串格式) */
export interface WeiboStatistics {
  name: string
  id: number
  attention: number
  fans: number
  feed: number
  crawl_time: string
  count: number
  forward_sum: string
  comment_sum: string
  like_sum: string
  forward_avg: string
  comment_avg: string
  like_avg: string
  forward_max: number
  comment_max: number
  like_max: number
  forward_min: number
  comment_min: number
  like_min: number
}

/** 当前用户信息(来自 JWT claims) */
export interface Profile {
  username: string
  nickname: string
}

import client, { unwrap } from './client'
import type { Profile, WeiboAccount, WeiboAccountForm, WeiboMsgListData, WeiboStatistics } from '../types'

/** 时间范围查询参数(start_date/end_date 格式 YYYY-MM-DD) */
export interface DateRangeParams {
  account_id?: number
  start_date?: string
  end_date?: string
}

/** 微博文列表查询参数 */
export interface MsgQueryParams extends DateRangeParams {
  sort_field?: string
  sort_order?: 'asc' | 'desc'
  page?: number
  page_size?: number
}

/** 登录,成功返回 accessToken */
export async function login(username: string, password: string): Promise<string> {
  const token = await unwrap<{ accessToken: string }>(
    await client.post('/login', { username, password }),
  )
  return token.accessToken
}

/** 当前用户信息 */
export async function fetchProfile(): Promise<Profile> {
  return unwrap<Profile>(await client.get('/user'))
}

/** 微博账号列表 */
export async function fetchAccounts(): Promise<WeiboAccount[]> {
  return unwrap<WeiboAccount[]>(await client.get('/weibo/accounts'))
}

/** 新增微博账号 */
export async function createAccount(form: WeiboAccountForm): Promise<void> {
  await unwrap(await client.post('/weibo/accounts', form))
}

/** 更新微博账号 */
export async function updateAccount(id: number, form: WeiboAccountForm): Promise<void> {
  await unwrap(await client.patch(`/weibo/accounts/${id}`, form))
}

/** 删除微博账号 */
export async function deleteAccount(id: number): Promise<void> {
  await unwrap(await client.delete(`/weibo/accounts/${id}`))
}

/** 微博文分页列表(服务端排序) */
export async function fetchMessages(params: MsgQueryParams): Promise<WeiboMsgListData> {
  return unwrap<WeiboMsgListData>(await client.get('/weibo/messages', { params }))
}

/** 微博统计(按账号聚合,不分页) */
export async function fetchStatistics(params: DateRangeParams): Promise<WeiboStatistics[]> {
  return unwrap<WeiboStatistics[]>(await client.get('/weibo/statistics', { params }))
}

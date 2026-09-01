import axios from 'axios'
import { message } from 'antd'
import { clearToken, getToken } from '../auth'
import type { ApiResponse } from '../types'

/**
 * axios 实例:统一 baseURL、token 注入、错误处理。
 * - 请求拦截:自动附带 Authorization: Bearer <token>
 * - 响应拦截:401 清除 token 跳登录;业务 code !== 0 统一弹错并 reject
 */
const client = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
})

client.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

client.interceptors.response.use(
  (response) => {
    const body = response.data as ApiResponse<unknown>
    if (body.code !== 0) {
      message.error(body.message || '请求失败')
      return Promise.reject(new Error(body.message))
    }
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      clearToken()
      // 跳转登录页(保留当前路径,登录后可回跳)
      const current = window.location.pathname + window.location.search
      if (!window.location.pathname.startsWith('/login')) {
        window.location.href = `/login?redirect=${encodeURIComponent(current)}`
      }
      return Promise.reject(error)
    }
    const msg = error.response?.data?.message || '网络异常,请稍后重试'
    message.error(msg)
    return Promise.reject(error)
  },
)

/** 取响应中的 data 字段 */
export function unwrap<T>(response: { data: ApiResponse<T> }): T {
  return response.data.data
}

export default client

/**
 * 登录 token 存取工具(内部工具,直接使用 localStorage)。
 */

const TOKEN_KEY = 'rivalscope_token'

/** 读取已保存的 token */
export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) ?? ''
}

/** 保存 token */
export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

/** 清除 token(退出登录 / 401 时调用) */
export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

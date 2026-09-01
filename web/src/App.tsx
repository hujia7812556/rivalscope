import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import type { ReactElement } from 'react'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import { getToken } from './auth'
import AppLayout from './layouts/AppLayout'
import LoginPage from './pages/Login'
import StatisticsPage from './pages/Statistics'
import MessagesPage from './pages/Messages'
import AccountsPage from './pages/Accounts'

const queryClient = new QueryClient()

/** 未登录时重定向到登录页 */
function RequireAuth({ children }: { children: ReactElement }) {
  if (!getToken()) {
    return <Navigate to="/login" replace />
  }
  return children
}

/** 应用根组件:主题配置 + QueryProvider + 路由表 */
export default function App() {
  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        token: {
          // 品牌色:靛蓝,区别于 AntD 默认蓝,更有辨识度
          colorPrimary: '#4f46e5',
          borderRadius: 8,
          fontSize: 14,
        },
        components: {
          Layout: {
            headerBg: '#ffffff',
            headerHeight: 60,
            siderBg: '#101828',
            bodyBg: '#f4f6fa',
          },
          Menu: {
            darkItemBg: '#101828',
            darkItemSelectedBg: '#4f46e5',
            darkItemColor: 'rgba(255,255,255,0.65)',
            darkItemSelectedColor: '#ffffff',
          },
          Table: {
            headerBg: '#f8fafc',
            headerColor: '#475569',
            rowHoverBg: '#f0f4ff',
          },
          Card: {
            boxShadowTertiary: '0 1px 2px rgba(16,24,40,0.04)',
          },
        },
      }}
    >
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route
              element={
                <RequireAuth>
                  <AppLayout />
                </RequireAuth>
              }
            >
              <Route path="/" element={<Navigate to="/weibo/statistics" replace />} />
              <Route path="/weibo/statistics" element={<StatisticsPage />} />
              <Route path="/weibo/msg" element={<MessagesPage />} />
              <Route path="/weibo/accounts" element={<AccountsPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/weibo/statistics" replace />} />
          </Routes>
        </BrowserRouter>
      </QueryClientProvider>
    </ConfigProvider>
  )
}

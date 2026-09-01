import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import type { ReactElement } from 'react'
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

/** 应用根组件:QueryProvider + 路由表 */
export default function App() {
  return (
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
  )
}

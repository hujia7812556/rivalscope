import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { Button, Dropdown, Layout, Menu, Spin, theme } from 'antd'
import type { MenuProps } from 'antd'
import { BarChartOutlined, FileTextOutlined, TeamOutlined, UserOutlined } from '@ant-design/icons'
import { fetchAccounts, fetchProfile } from '../api'
import { clearToken } from '../auth'

const { Sider, Header, Content } = Layout

/** 侧边导航:微博模块(统计结果/详细数据/账号管理) */
const menuItems: MenuProps['items'] = [
  {
    key: 'weibo',
    label: '微博',
    children: [
      { key: '/weibo/statistics', label: '统计结果', icon: <BarChartOutlined /> },
      { key: '/weibo/msg', label: '详细数据', icon: <FileTextOutlined /> },
      { key: '/weibo/accounts', label: '账号管理', icon: <TeamOutlined /> },
    ],
  },
]

/**
 * 应用主布局:左侧菜单 + 顶栏(用户信息/退出)+ 内容区。
 * 预加载账号列表,供业务页搜索下拉直接复用缓存。
 */
export default function AppLayout() {
  const navigate = useNavigate()
  const location = useLocation()
  const { token } = theme.useToken()

  useQuery({ queryKey: ['accounts'], queryFn: fetchAccounts })
  const { data: profile, isLoading } = useQuery({ queryKey: ['profile'], queryFn: fetchProfile })

  const displayName = useMemo(
    () => profile?.nickname || profile?.username || '已登录',
    [profile],
  )

  /** 退出登录:清除 token 并回到登录页 */
  const handleLogout = () => {
    clearToken()
    navigate('/login', { replace: true })
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="dark" width={200}>
        <div
          style={{
            height: 48,
            margin: 12,
            borderRadius: 6,
            color: '#fff',
            fontSize: 16,
            fontWeight: 600,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            letterSpacing: 1,
          }}
        >
          竞品观察
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          defaultOpenKeys={['weibo']}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            background: token.colorBgContainer,
            padding: '0 24px',
            display: 'flex',
            justifyContent: 'flex-end',
            alignItems: 'center',
          }}
        >
          <Dropdown
            menu={{
              items: [{ key: 'logout', label: '退出登录' }],
              onClick: ({ key }) => key === 'logout' && handleLogout(),
            }}
          >
            <Button type="text" icon={<UserOutlined />}>
              {isLoading ? <Spin size="small" /> : displayName}
            </Button>
          </Dropdown>
        </Header>
        <Content style={{ margin: 16 }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

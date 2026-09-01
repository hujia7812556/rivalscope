import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { Avatar, Button, Dropdown, Layout, Menu, Spin, theme } from 'antd'
import type { MenuProps } from 'antd'
import {
  BarChartOutlined,
  FileTextOutlined,
  LogoutOutlined,
  RadarChartOutlined,
  TeamOutlined,
  UserOutlined,
} from '@ant-design/icons'
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

/** 页面标题映射(顶栏展示) */
const pageTitle: Record<string, { title: string; desc: string }> = {
  '/weibo/statistics': { title: '统计结果', desc: '按账号聚合的互动数据概览' },
  '/weibo/msg': { title: '详细数据', desc: '微博文明细,支持排序与分页' },
  '/weibo/accounts': { title: '账号管理', desc: '维护要观察的竞品账号' },
}

/**
 * 应用主布局:左侧菜单 + 顶栏(页面标题/用户信息)+ 内容区。
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
  const meta = pageTitle[location.pathname] ?? { title: '竞品观察', desc: '' }

  /** 退出登录:清除 token 并回到登录页 */
  const handleLogout = () => {
    clearToken()
    navigate('/login', { replace: true })
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="dark" width={168} style={{ borderRight: '1px solid rgba(255,255,255,0.06)' }}>
        {/* 品牌区 */}
        <div
          style={{
            height: 60,
            display: 'flex',
            alignItems: 'center',
            gap: 8,
            padding: '0 16px',
            borderBottom: '1px solid rgba(255,255,255,0.06)',
          }}
        >
          <RadarChartOutlined style={{ fontSize: 20, color: '#818cf8' }} />
          <div style={{ lineHeight: 1.2 }}>
            <div style={{ color: '#fff', fontSize: 15, fontWeight: 600, letterSpacing: 1 }}>
              竞品观察
            </div>
            <div style={{ color: 'rgba(255,255,255,0.35)', fontSize: 11 }}>RIVALSCOPE</div>
          </div>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          defaultOpenKeys={['weibo']}
          items={menuItems}
          onClick={({ key }) => navigate(key)}
        />
        {/* 底部版本信息 */}
        <div
          style={{
            position: 'absolute',
            bottom: 12,
            left: 0,
            right: 0,
            textAlign: 'center',
            color: 'rgba(255,255,255,0.25)',
            fontSize: 11,
          }}
        >
          竞品社媒数据平台
        </div>
      </Sider>
      <Layout>
        <Header
          style={{
            background: token.colorBgContainer,
            // 左侧与侧边栏品牌文字起点对齐(品牌区 padding 16 + 图标 20 + 间距 8 = 44)
            padding: '0 24px 0 44px',
            display: 'flex',
            justifyContent: 'space-between',
            alignItems: 'center',
            borderBottom: '1px solid #eef0f4',
          }}
        >
          <div style={{ display: 'flex', alignItems: 'baseline', gap: 12 }}>
            <span style={{ fontSize: 15, fontWeight: 600, color: '#101828', lineHeight: '22px' }}>
              {meta.title}
            </span>
            {meta.desc && (
              <span style={{ fontSize: 12, color: '#98a2b3' }}>{meta.desc}</span>
            )}
          </div>
          <Dropdown
            menu={{
              items: [{ key: 'logout', label: '退出登录', icon: <LogoutOutlined /> }],
              onClick: ({ key }) => key === 'logout' && handleLogout(),
            }}
          >
            <Button type="text" style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              <Avatar size={28} style={{ background: '#4f46e5' }} icon={<UserOutlined />}>
                {displayName.slice(0, 1)}
              </Avatar>
              {isLoading ? <Spin size="small" /> : displayName}
            </Button>
          </Dropdown>
        </Header>
        <Content style={{ padding: 0, overflow: 'auto' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

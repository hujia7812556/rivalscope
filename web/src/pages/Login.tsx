import { useState } from 'react'
import { Button, Card, Form, Input, Typography, message } from 'antd'
import { LockOutlined, RadarChartOutlined, UserOutlined } from '@ant-design/icons'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { login } from '../api'
import { setToken } from '../auth'

/** 登录页:用户名 + 密码(账号配置在服务端 config.yaml 的 auth.users) */
export default function LoginPage() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const [loading, setLoading] = useState(false)

  /** 提交登录:保存 token 后跳转 redirect 指定的页面(默认首页) */
  const handleFinish = async (values: { username: string; password: string }) => {
    setLoading(true)
    try {
      const token = await login(values.username, values.password)
      setToken(token)
      message.success('登录成功')
      const redirect = searchParams.get('redirect') || '/'
      navigate(redirect, { replace: true })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      style={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        // 品牌渐变背景
        background:
          'radial-gradient(1200px 600px at 15% -10%, rgba(99,102,241,0.35), transparent 60%),' +
          'radial-gradient(1000px 500px at 110% 110%, rgba(79,70,229,0.28), transparent 55%),' +
          '#0f172a',
      }}
    >
      <Card
        style={{
          width: 392,
          border: 'none',
          borderRadius: 16,
          boxShadow: '0 20px 50px rgba(15,23,42,0.45)',
        }}
      >
        <div style={{ textAlign: 'center', marginBottom: 28, paddingTop: 8 }}>
          {/* 品牌图标 */}
          <div
            style={{
              width: 56,
              height: 56,
              margin: '0 auto 16px',
              borderRadius: 16,
              background: 'linear-gradient(135deg, #6366f1, #4f46e5)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              boxShadow: '0 8px 20px rgba(79,70,229,0.35)',
            }}
          >
            <RadarChartOutlined style={{ fontSize: 28, color: '#fff' }} />
          </div>
          <Typography.Title level={3} style={{ marginBottom: 4 }}>
            竞品观察
          </Typography.Title>
          <Typography.Text type="secondary">RivalScope · 竞品社媒数据平台</Typography.Text>
        </div>
        <Form layout="vertical" onFinish={handleFinish} size="large">
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined style={{ color: '#94a3b8' }} />} placeholder="用户名" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined style={{ color: '#94a3b8' }} />} placeholder="密码" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0, marginTop: 6 }}>
            <Button type="primary" htmlType="submit" block size="large" loading={loading}>
              登 录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

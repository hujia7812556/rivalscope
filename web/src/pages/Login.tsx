import { useState } from 'react'
import { Button, Card, Form, Input, Typography, message } from 'antd'
import { LockOutlined, UserOutlined } from '@ant-design/icons'
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
        background: '#f0f2f5',
      }}
    >
      <Card style={{ width: 380 }}>
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Typography.Title level={3} style={{ marginBottom: 4 }}>
            竞品观察
          </Typography.Title>
          <Typography.Text type="secondary">RivalScope · 竞品社媒数据平台</Typography.Text>
        </div>
        <Form layout="vertical" onFinish={handleFinish}>
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名" size="large" />
          </Form.Item>
          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
          </Form.Item>
          <Form.Item style={{ marginBottom: 0 }}>
            <Button type="primary" htmlType="submit" block size="large" loading={loading}>
              登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  )
}

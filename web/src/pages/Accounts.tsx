import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Button, Card, Form, Input, Modal, Popconfirm, Space, Table, message } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { DeleteOutlined, EditOutlined, PlusOutlined } from '@ant-design/icons'
import { createAccount, deleteAccount, fetchAccounts, updateAccount } from '../api'
import type { WeiboAccount } from '../types'

/** 弹窗表单值 */
interface AccountFormValues {
  name: string
  uid: string
}

/**
 * 账号管理页:微博账号的增删改。
 * 名称列链接到微博主页,弹窗表单带必填校验。
 */
export default function AccountsPage() {
  const queryClient = useQueryClient()
  const [form] = Form.useForm<AccountFormValues>()
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<WeiboAccount | null>(null)

  const { data: rows = [], isFetching } = useQuery({
    queryKey: ['accounts'],
    queryFn: fetchAccounts,
  })

  /** 成功后刷新列表并关闭弹窗 */
  const refresh = () => queryClient.invalidateQueries({ queryKey: ['accounts'] })

  const saveMutation = useMutation({
    mutationFn: (values: AccountFormValues) =>
      editing ? updateAccount(editing.id, values) : createAccount(values),
    onSuccess: async () => {
      message.success('提交成功!')
      setModalOpen(false)
      await refresh()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteAccount(id),
    onSuccess: async () => {
      message.success('删除成功!')
      await refresh()
    },
  })

  /** 打开新增弹窗 */
  const openCreate = () => {
    setEditing(null)
    form.resetFields()
    setModalOpen(true)
  }

  /** 打开编辑弹窗(回填当前值) */
  const openEdit = (record: WeiboAccount) => {
    setEditing(record)
    form.setFieldsValue({ name: record.name, uid: record.uid })
    setModalOpen(true)
  }

  /** 提交表单(antd Form 已做必填校验) */
  const handleOk = async () => {
    const values = await form.validateFields()
    saveMutation.mutate(values)
  }

  const columns: ColumnsType<WeiboAccount> = [
    {
      title: '名称',
      dataIndex: 'name',
      width: 240,
      render: (name: string, record) => (
        <a href={record.url} target="_blank" rel="noreferrer">
          {name}
        </a>
      ),
    },
    { title: 'uid', dataIndex: 'uid', width: 200 },
    {
      title: '最近抓取时间',
      dataIndex: 'crawl_time',
      width: 180,
      sorter: (a, b) => a.crawl_time.localeCompare(b.crawl_time),
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_, record) => (
        <Space>
          <Button size="small" icon={<EditOutlined />} onClick={() => openEdit(record)}>
            编辑
          </Button>
          <Popconfirm
            title="此操作将永久删除该配置, 是否继续?"
            okText="确定"
            cancelText="取消"
            onConfirm={() => deleteMutation.mutate(record.id)}
          >
            <Button size="small" danger icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <Card size="small" style={{ borderRadius: 12 }} styles={{ body: { paddingTop: 16 } }}>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
          添加微博号
        </Button>
      </div>
      <Table<WeiboAccount>
        rowKey="id"
        columns={columns}
        dataSource={rows}
        loading={isFetching}
        pagination={false}
        size="middle"
      />

      <Modal
        title="微博账号信息"
        open={modalOpen}
        onOk={handleOk}
        onCancel={() => setModalOpen(false)}
        confirmLoading={saveMutation.isPending}
        okText="确定"
        cancelText="取消"
        destroyOnClose
      >
        <Form form={form} layout="vertical" style={{ marginTop: 16 }}>
          <Form.Item
            name="name"
            label="名称"
            rules={[{ required: true, message: '请输入名称' }]}
          >
            <Input placeholder="请输入名称" />
          </Form.Item>
          <Form.Item
            name="uid"
            label="uid"
            rules={[{ required: true, message: '请输入uid' }]}
          >
            <Input placeholder="请输入uid" />
          </Form.Item>
        </Form>
      </Modal>
    </Card>
  )
}

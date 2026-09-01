import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Card, Table } from 'antd'
import type { ColumnsType, TablePaginationConfig } from 'antd/es/table'
import type { FilterValue, SorterResult } from 'antd/es/table/interface'
import { fetchMessages } from '../api'
import type { MsgQueryParams } from '../api'
import type { WeiboMsg } from '../types'
import SearchFilter, { toQueryParams } from '../components/SearchFilter'
import type { FilterValue as FilterState } from '../components/SearchFilter'

/** 服务端排序字段:后端列名 <-> 前端 dataIndex 映射(后端已做白名单校验) */
const sortableFields: Record<string, string> = {
  id: 'id',
  forward: 'forward',
  comment: 'comment',
  like: 'like',
  pubtime: 'pubtime',
}

/**
 * 详细数据页:微博文列表。
 * 服务端分页 + 服务端排序(表头点击),搜索条件与统计页一致。
 */
export default function MessagesPage() {
  const [filter, setFilter] = useState<FilterState>({})
  const [query, setQuery] = useState<MsgQueryParams>({ page: 1, page_size: 20 })

  const { data, isFetching } = useQuery({
    queryKey: ['messages', query],
    queryFn: () => fetchMessages(query),
    placeholderData: (prev) => prev,
  })

  /** 点击查询:合并搜索条件并回到第一页 */
  const handleSearch = () => {
    setQuery({ ...toQueryParams(filter), page: 1, page_size: 20 })
  }

  /** 点击清空:重置搜索条件与排序,回到第一页 */
  const handleReset = () => {
    setFilter({})
    setQuery({ page: 1, page_size: 20 })
  }

  /** 分页/排序变化统一处理:排序变更时回到第一页 */
  const handleTableChange = (
    pagination: TablePaginationConfig,
    _filters: Record<string, FilterValue | null>,
    sorter: SorterResult<WeiboMsg> | SorterResult<WeiboMsg>[],
  ) => {
    const s = Array.isArray(sorter) ? sorter[0] : sorter
    const field = s?.order && s.field ? sortableFields[String(s.field)] : undefined
    setQuery((prev) => ({
      ...prev,
      sort_field: field,
      sort_order: s?.order === 'descend' ? 'desc' : 'asc',
      page: pagination.current ?? 1,
      page_size: pagination.pageSize ?? 20,
    }))
  }

  const columns: ColumnsType<WeiboMsg> = [
    { title: 'id', dataIndex: 'id', width: 90, sorter: true },
    { title: '账号', dataIndex: 'account_name', width: 140 },
    {
      title: 'url',
      dataIndex: 'url',
      width: 280,
      render: (url: string) => (
        <a href={url} target="_blank" rel="noreferrer" style={{ wordBreak: 'break-all' }}>
          {url}
        </a>
      ),
    },
    { title: '转发', dataIndex: 'forward', width: 100, align: 'right', sorter: true },
    { title: '评论', dataIndex: 'comment', width: 100, align: 'right', sorter: true },
    { title: '点赞', dataIndex: 'like', width: 100, align: 'right', sorter: true },
    { title: '发布时间', dataIndex: 'pubtime', width: 130, sorter: true },
    { title: '抓取时间', dataIndex: 'crawl_time', width: 180 },
  ]

  const pagination: TablePaginationConfig = {
    current: query.page,
    pageSize: query.page_size,
    total: data?.total ?? 0,
    showTotal: (t) => `共 ${t} 条`,
    showSizeChanger: false,
  }

  return (
    <Card size="small" style={{ borderRadius: 12 }} styles={{ body: { paddingTop: 16 } }}>
      <SearchFilter
        value={filter}
        onChange={setFilter}
        onSearch={handleSearch}
        onReset={handleReset}
        loading={isFetching}
      />
      <Table<WeiboMsg>
        rowKey="id"
        columns={columns}
        dataSource={data?.msg_list ?? []}
        loading={isFetching}
        pagination={pagination}
        onChange={handleTableChange}
        // 不设 scroll.x:表格宽度 100%,各列按声明的 width 比例自然分配,宽屏无右侧留白
        size="middle"
      />
    </Card>
  )
}

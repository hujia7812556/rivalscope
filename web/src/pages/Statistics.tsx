import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Button, Table } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { DownloadOutlined } from '@ant-design/icons'
import { fetchStatistics } from '../api'
import type { WeiboStatistics } from '../types'
import SearchFilter, { toQueryParams } from '../components/SearchFilter'
import type { FilterValue } from '../components/SearchFilter'
import { exportToExcel } from '../utils/exportExcel'

/** Excel 导出的中文表头与字段映射(顺序即列顺序) */
const excelHeaders: Record<string, string> = {
  name: '名称',
  fans: '粉丝数',
  feed: '总微博数',
  count: '微博数',
  forward_sum: '总转发数',
  comment_sum: '总评论数',
  like_sum: '总点赞数',
  forward_avg: '平均转发数',
  comment_avg: '平均评论数',
  like_avg: '平均点赞数',
  forward_max: '最大转发数',
  comment_max: '最大评论数',
  like_max: '最大点赞数',
  crawl_time: '抓取时间',
}

/** 数值列排序(Sum/Avg 是字符串,统一转数字比较) */
const numSorter = (key: keyof WeiboStatistics) => (a: WeiboStatistics, b: WeiboStatistics) =>
  Number(a[key]) - Number(b[key])

/**
 * 统计结果页:按账号聚合的转发/评论/点赞统计。
 * 全量拉取(不分页),排序全部在前端完成,支持导出 Excel。
 */
export default function StatisticsPage() {
  const [filter, setFilter] = useState<FilterValue>({})
  const [query, setQuery] = useState<ReturnType<typeof toQueryParams>>({})

  const { data: rows = [], isFetching } = useQuery({
    queryKey: ['statistics', query],
    queryFn: () => fetchStatistics(query),
  })

  /** 点击查询:按当前搜索条件拉取 */
  const handleSearch = () => setQuery(toQueryParams(filter))

  /** 点击清空:重置搜索条件并重新查询全量 */
  const handleReset = () => {
    setFilter({})
    setQuery({})
  }

  const columns: ColumnsType<WeiboStatistics> = [
    { title: 'id', dataIndex: 'id', width: 70, fixed: 'left', sorter: numSorter('id') },
    { title: '名称', dataIndex: 'name', width: 150, fixed: 'left' },
    { title: '粉丝数', dataIndex: 'fans', width: 110, sorter: numSorter('fans') },
    { title: '总微博数', dataIndex: 'feed', width: 110, sorter: numSorter('feed') },
    { title: '微博数', dataIndex: 'count', width: 100, sorter: numSorter('count') },
    { title: '总转发数', dataIndex: 'forward_sum', width: 110, sorter: numSorter('forward_sum') },
    { title: '总评论数', dataIndex: 'comment_sum', width: 110, sorter: numSorter('comment_sum') },
    { title: '总点赞数', dataIndex: 'like_sum', width: 110, sorter: numSorter('like_sum') },
    { title: '平均转发数', dataIndex: 'forward_avg', width: 120, sorter: numSorter('forward_avg') },
    { title: '平均评论数', dataIndex: 'comment_avg', width: 120, sorter: numSorter('comment_avg') },
    { title: '平均点赞数', dataIndex: 'like_avg', width: 120, sorter: numSorter('like_avg') },
    { title: '最大转发数', dataIndex: 'forward_max', width: 120, sorter: numSorter('forward_max') },
    { title: '最大评论数', dataIndex: 'comment_max', width: 120, sorter: numSorter('comment_max') },
    { title: '最大点赞数', dataIndex: 'like_max', width: 120, sorter: numSorter('like_max') },
    { title: '抓取时间', dataIndex: 'crawl_time', width: 170, sorter: numSorter('crawl_time') },
  ]

  return (
    <div>
      <SearchFilter
        value={filter}
        onChange={setFilter}
        onSearch={handleSearch}
        onReset={handleReset}
        loading={isFetching}
        accountLabel="名称"
        extra={
          <Button
            type="primary"
            icon={<DownloadOutlined />}
            disabled={rows.length === 0}
            onClick={() => exportToExcel('weibo.xlsx', excelHeaders, rows)}
          >
            导出
          </Button>
        }
      />
      <Table<WeiboStatistics>
        rowKey="id"
        columns={columns}
        dataSource={rows}
        loading={isFetching}
        pagination={false}
        scroll={{ x: 1700, y: 650 }}
        size="middle"
      />
    </div>
  )
}

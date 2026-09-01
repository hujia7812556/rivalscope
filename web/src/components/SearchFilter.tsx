import { useQuery } from '@tanstack/react-query'
import { Button, DatePicker, Select, Space } from 'antd'
import type { Dayjs } from 'dayjs'
import type { ReactNode } from 'react'
import dayjs from 'dayjs'
import { ClearOutlined, SearchOutlined } from '@ant-design/icons'
import { fetchAccounts } from '../api'

const { RangePicker } = DatePicker

/** 日期快捷选项(结构与 antd RangePicker presets 一致) */
interface RangePreset {
  label: string
  value: [Dayjs, Dayjs]
}

/** 搜索条件值 */
export interface FilterValue {
  accountId?: number
  dateRange?: [Dayjs | null, Dayjs | null] | null
}

/** 日期快捷选项(禁选未来日期) */
const rangePresets: RangePreset[] = [
  { label: '最近一周', value: [dayjs().subtract(7, 'day'), dayjs()] },
  { label: '最近一个月', value: [dayjs().subtract(1, 'month'), dayjs()] },
  { label: '最近三个月', value: [dayjs().subtract(3, 'month'), dayjs()] },
  { label: '最近半年', value: [dayjs().subtract(6, 'month'), dayjs()] },
  { label: '最近一年', value: [dayjs().subtract(1, 'year'), dayjs()] },
]

interface SearchFilterProps {
  value: FilterValue
  onChange: (value: FilterValue) => void
  /** 点击「查询」 */
  onSearch: () => void
  /** 点击「清空」(组件只负责回调,由页面重置状态) */
  onReset: () => void
  loading?: boolean
  /** 账号下拉的标签文案 */
  accountLabel?: string
  /** 追加在查询/清空之后的控件(如导出按钮),与搜索区同一行对齐 */
  extra?: ReactNode
}

/**
 * 业务页公用搜索区:账号下拉(可搜索)+ 日期范围(快捷选项)+ 查询/清空。
 * 受控组件,由页面维护 FilterValue 状态。
 */
export default function SearchFilter({
  value,
  onChange,
  onSearch,
  onReset,
  loading,
  accountLabel = '账号',
  extra,
}: SearchFilterProps) {
  const { data: accounts = [] } = useQuery({ queryKey: ['accounts'], queryFn: fetchAccounts })

  return (
    <Space wrap style={{ marginBottom: 16 }}>
      <Select
        style={{ width: 220 }}
        showSearch
        allowClear
        placeholder={`请选择${accountLabel}`}
        optionFilterProp="label"
        value={value.accountId}
        onChange={(v) => onChange({ ...value, accountId: v })}
        options={accounts.map((a) => ({ value: a.id, label: a.name }))}
      />
      <RangePicker
        value={value.dateRange}
        onChange={(v) => onChange({ ...value, dateRange: v })}
        disabledDate={(d) => d.isAfter(dayjs(), 'day')}
        presets={rangePresets}
        placeholder={['开始日期', '结束日期']}
      />
      <Button type="primary" icon={<SearchOutlined />} loading={loading} onClick={onSearch}>
        查询
      </Button>
      <Button icon={<ClearOutlined />} onClick={onReset}>
        清空
      </Button>
      {extra}
    </Space>
  )
}

/**
 * 把搜索条件转换为接口查询参数(dateRange → start_date/end_date)。
 * 返回新对象,只包含有值的字段。
 */
export function toQueryParams(filter: FilterValue): {
  account_id?: number
  start_date?: string
  end_date?: string
} {
  const params: { account_id?: number; start_date?: string; end_date?: string } = {}
  if (filter.accountId) {
    params.account_id = filter.accountId
  }
  const [start, end] = filter.dateRange ?? []
  if (start) {
    params.start_date = start.format('YYYY-MM-DD')
  }
  if (end) {
    params.end_date = end.format('YYYY-MM-DD')
  }
  return params
}

import * as XLSX from 'xlsx'

/**
 * 把行数据导出为 xlsx 文件。
 * @param filename 文件名(如 weibo.xlsx)
 * @param headers  表头映射:列 key -> 中文列名(顺序即列顺序)
 * @param rows     数据行
 */
export function exportToExcel<T>(filename: string, headers: Record<string, string>, rows: T[]): void {
  const keys = Object.keys(headers)
  const aoa: unknown[][] = [
    Object.values(headers),
    ...rows.map((row) => keys.map((key) => (row as Record<string, unknown>)[key] ?? '')),
  ]
  const worksheet = XLSX.utils.aoa_to_sheet(aoa)
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, 'Sheet1')
  XLSX.writeFile(workbook, filename)
}
